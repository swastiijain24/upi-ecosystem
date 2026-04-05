package apiAuth

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/services"
)

type contextKey string

const (
	APIKeyContextKey contextKey = "api_key"
)

type APIMiddleware struct {
	keyAuth *KeyAuth
	apiKeyService services.ApiKeyService
}

func NewApiAuthMiddleware(keyAuth *KeyAuth, apiKeyService services.ApiKeyService) *APIMiddleware {
	return &APIMiddleware{
		keyAuth: keyAuth,
		apiKeyService: apiKeyService,
	}
}

func (m *APIMiddleware) ApiAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		apiKey := m.extractAPIKey(c)
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "API Key required"})
			return

		}

		if !m.keyAuth.ValidateFormat(apiKey) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key format"})
			return
		}

		_, randomPart, _ := m.keyAuth.ParseKey(apiKey)
		keyID := randomPart[:8]

		key, err := m.apiKeyService.GetAPIKeyByKeyID(c, keyID)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
			return
		}

		valid, err := m.keyAuth.Verify(apiKey, key.KeyHash)
		if err != nil || !valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
			return
		}

		IsActive, err := m.apiKeyService.IsValid(c, key.KeyID)
		if !IsActive || time.Now().After(key.ExpiresAt.Time) {

			c.AbortWithStatusJSON(401, gin.H{"error": "API key is no longer valid"})
			return

		}

		if len(key.AllowedIps) > 0 && !m.isIPAllowed(c, key.AllowedIps) {
			c.AbortWithStatusJSON(403, gin.H{"error": "IP not allowed"})
			return
		}

		go func() {
			_ = m.apiKeyService.UpdateAPIKeyLastUsed(context.Background(), key.KeyID)
		}()

		c.Set(APIKeyContextKey, key)

		c.Next()
	}
}

func (m *APIMiddleware) extractAPIKey(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	if key := c.GetHeader("X-API-Key"); key != "" {
		return key
	}

	return ""

}

func (m *APIMiddleware) isIPAllowed(c *gin.Context, allowedIPs []string) bool {

	clientIP := c.ClientIP()

	for _, allowed := range allowedIPs {
		if strings.Contains(allowed, "/") {
			_, network, err := net.ParseCIDR(allowed)
			if err == nil && network.Contains(net.ParseIP(clientIP)) {
				return true
			}

		} else if clientIP == allowed {
			return true
		}
	}
	return false
}
