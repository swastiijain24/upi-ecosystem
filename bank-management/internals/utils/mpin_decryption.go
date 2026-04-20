package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func DecryptAES(cryptoText string, key []byte) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(cryptoText)
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	nonce := data[:gcm.NonceSize()]
	plainText, err := gcm.Open(nil, nonce, data[gcm.NonceSize():], nil)
	return string(plainText), err
}