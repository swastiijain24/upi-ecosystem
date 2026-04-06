package apiAuth

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swastiijain24/bank-management/internals/audit"
	"github.com/swastiijain24/bank-management/internals/services"
)

const APIKeyContextKey = "api_key"

type APIMiddleware struct {
	generator     *APIKeyGenerator
	hasher        *APIKeyHasher
	auditLogger   *audit.Logger
	apiKeyService services.ApiKeyService
}

func NewApiAuthMiddleware(generator *APIKeyGenerator, hasher *APIKeyHasher, auditLogger *audit.Logger, apiKeyService services.ApiKeyService) *APIMiddleware {
	return &APIMiddleware{
		generator:     generator,
		hasher:        hasher,
		auditLogger:   auditLogger,
		apiKeyService: apiKeyService,
	}
}

func (m *APIMiddleware) ApiAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		apiKey := m.extractAPIKey(c)
		if apiKey == "" {
			m.auditLogger.LogAuthFailure(ctx, "", "missing_api_key", c.Request.RemoteAddr)
			c.AbortWithStatusJSON(401, gin.H{"error": "API Key required"})
			return

		}

		if !m.generator.ValidateFormat(apiKey) {
			m.auditLogger.LogAuthFailure(ctx, "", "invalid_key_format", c.Request.RemoteAddr)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key format"})
			return
		}

		_, _, randomPart, err := m.generator.ParseKey(apiKey)
		if err != nil || len(randomPart) < 8 {
			m.auditLogger.LogAuthFailure(ctx, "", "incorrect_key_length", c.Request.RemoteAddr)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
			return
		}
		keyID := randomPart[:8]

		key, err := m.apiKeyService.GetAPIKeyByKeyID(c, keyID)
		if err != nil {
			m.auditLogger.LogAuthFailure(ctx, keyID, "key_not_found", c.Request.RemoteAddr)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
			return
		}

		valid, err := m.hasher.Verify(apiKey, key.KeyHash)
		if err != nil || !valid {
			m.auditLogger.LogAuthFailure(ctx, keyID, "hash_mismatch", c.Request.RemoteAddr)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
			return
		}

		IsActive, err := m.apiKeyService.IsValid(c, key.KeyID)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if !IsActive || (key.ExpiresAt.Valid && time.Now().After(key.ExpiresAt.Time)) {
			m.auditLogger.LogAuthFailure(ctx, keyID, "key_inactive", c.Request.RemoteAddr)
			c.AbortWithStatusJSON(401, gin.H{"error": "API key is no longer valid"})
			return

		}

		if len(key.AllowedIps) > 0 && !m.isIPAllowed(c, key.AllowedIps) {
			m.auditLogger.LogAuthFailure(ctx, keyID, "ip_not_allowed", c.Request.RemoteAddr)
			c.AbortWithStatusJSON(403, gin.H{"error": "IP not allowed"})
			return
		}

		go func() {
			_ = m.apiKeyService.UpdateAPIKeyLastUsed(context.Background(), key.KeyID)
		}()


		m.auditLogger.LogAuthSuccess(ctx, keyID, c.Request.RemoteAddr, c.Request.URL.RawPath)

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
