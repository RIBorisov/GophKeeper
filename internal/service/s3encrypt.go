package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// Encrypt encrypts data and returns it as a byte slice.
func Encrypt(secretKey string, data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to create nonce: %w", err)
	}

	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

// Decrypt decrypts the encrypted data and returns it as a byte slice.
func Decrypt(secretKey string, data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, cText := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	decrypted, err := gcm.Open(nil, nonce, cText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return decrypted, nil
}
