package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

var defaultPassenv = "ENVY_PASSWORD"

func createHash(key string) []byte {
	h := sha256.Sum256([]byte(key))
	return h[:]
}

func encryptData(data []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)
	if n, err := rand.Read(nonce); err != nil {
		return nil, err
	} else if n != nonceSize {
		return nil, errors.New("nonce")
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decryptData(data []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher(createHash(passphrase))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
