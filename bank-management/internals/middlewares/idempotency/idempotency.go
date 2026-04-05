package idempotency

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	//standard header name
	HeaderIdempotencyKey = "Idempotency-Key"

	//cached response
	HeaderIdempotencyReplayed = "Idempotency-Replayed"
)

type IdempotencyMiddleware struct {
	redisStore RedisStore
}

func NewIdempotencyMiddleware(redisStore RedisStore) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		redisStore: redisStore,
	}
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	CreatedAt  time.Time
}

type ginResponseWriter struct {
	gin.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

func (w *ginResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *ginResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (m *IdempotencyMiddleware) IdempotencyCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
			c.Next()
			return
		}

		key := c.GetHeader(HeaderIdempotencyKey)
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Idempotency-Key header is required",
			})
			c.Abort()
			return
		}

		cached, err := m.redisStore.Get(key)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if cached != nil && cached.StatusCode != 0 {
			m.replayGinResponse(c, cached)
			c.Abort()
			return
		}

		acquired, err := m.redisStore.SetProcessing(key)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !acquired {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Request already being processed",
			})
			c.Abort()
			return
		}

		writer := &ginResponseWriter{
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
		}
		c.Writer = writer

		c.Next()

		response := &Response{
			StatusCode: writer.statusCode,
			Body:       writer.body.Bytes(),
			Headers:    map[string]string{},
		}

		for _, header := range []string{"Content-Type", "Location"} {
			if val := c.Writer.Header().Get(header); val != "" {
				response.Headers[header] = val
			}
		}

		_ = m.redisStore.Set(key, response)
	}
}

func (m *IdempotencyMiddleware) replayGinResponse(c *gin.Context, resp *Response) {
	for k, v := range resp.Headers {
		c.Header(k, v)

	}

	c.Header(HeaderIdempotencyReplayed, "true")

	c.Status(resp.StatusCode)
	c.Writer.Write(resp.Body)
}
