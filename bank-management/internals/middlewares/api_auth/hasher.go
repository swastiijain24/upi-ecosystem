package apiAuth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

type APIKeyHasher struct{}

func NewAPIKeyHasher() *APIKeyHasher {
	return &APIKeyHasher{}
}

func (h *APIKeyHasher) Hash(key string) string {
	key = strings.TrimSpace(key)

	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (h *APIKeyHasher) Verify(providedKey, storedHash string) (bool, error) {
	hashed := h.Hash(providedKey)

	return subtle.ConstantTimeCompare(
		[]byte(hashed),
		[]byte(storedHash),
	) == 1, nil
}

