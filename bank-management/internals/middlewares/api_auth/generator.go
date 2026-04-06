package apiAuth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	KeyLength  = 32
	KeyPrefix  = "npci"
	KeyVersion = "v1"
)

type APIKeyGenerator struct {
	prefix  string
	version string
}

func NewAPIKeyGenerator() *APIKeyGenerator {
	return &APIKeyGenerator{
		prefix:  KeyPrefix,
		version: KeyVersion,
	}
}

func (g *APIKeyGenerator) Generate() (fullKey string, keyID string, err error) {
	randomBytes := make([]byte, KeyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	randomPart := base64.RawURLEncoding.EncodeToString(randomBytes)

	keyID = randomPart[:8]
	fullKey = fmt.Sprintf("%s_%s_%s", g.prefix, g.version, randomPart)

	return fullKey, keyID, nil
}


func (g *APIKeyGenerator) ParseKey(fullKey string) (prefix, version, randomPart string, err error){
	parts := strings.SplitN(fullKey, "_", 3)
	if len(parts)!= 3{
		return "", "","", fmt.Errorf("invalid key format: expected 3 parts, got %d", len(parts))
	}

	return parts[0], parts[1], parts[2], nil
}


func (g *APIKeyGenerator) ValidateFormat(fullKey string) bool {
	prefix, version, randomPart, err := g.ParseKey(fullKey)

	if err != nil {
		return false 
	}

	if prefix != g.prefix{
		return	false 
	}

	if version != g.version {
		return  false
	}

	expectedLen := base64.RawURLEncoding.EncodedLen(KeyLength)
    if len(randomPart) != expectedLen {
        return false
    }

    return true
}

