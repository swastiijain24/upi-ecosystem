package apiAuth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	KeyLength = 32
	KeyPrefix = "npci" 
)

type KeyAuth struct{
	prefix string 
}

func NewKeyAuth() *KeyAuth {
	return &KeyAuth{
		prefix: KeyPrefix,
	}
}

func (h *KeyAuth) Hash(key string) string {
	key = strings.TrimSpace(key)

	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (h *KeyAuth) Verify(providedKey, storedHash string) (bool, error) {
	hashed := h.Hash(providedKey)

	return subtle.ConstantTimeCompare(
		[]byte(hashed),
		[]byte(storedHash),
	) == 1, nil
}

func (h* KeyAuth) ParseKey(fullKey string) (prefix string, randomPart string, err error){
	parts := strings.Split(fullKey, "_")
	if len(parts)!= 2{
		return "", "", fmt.Errorf("invalid key format: expected 2 parts, got %d", len(parts))
	}

	return parts[0], parts[1], nil
}

func (h *KeyAuth) ValidateFormat(fullKey string) bool {
	prefix, randomPart, err := h.ParseKey(fullKey)

	if err != nil {
		return false 
	}

	if prefix != h.prefix{
		return	false 
	}

	expectedLen := base64.RawURLEncoding.EncodedLen(KeyLength)
    if len(randomPart) != expectedLen {
        return false
    }

    return true
}