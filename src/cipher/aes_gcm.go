package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// AESGCMEncrypt encrypts the text.
// The key argument should be the AES key,
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
func AESGCMEncrypt(key []byte, text []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("cipher: creates an aes cipher block: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("cipher: creates an aesgcm: %w", err)
	}
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	// 12 byte is cipher.gcmStandardNonceSize
	nonce = make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("cipher: creates a nonece: %w", err)
	}
	ciphertext = aesgcm.Seal(nil, nonce, text, nil)
	return ciphertext, nonce, nil
}

func AESGCMDecrypt(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher: creates an aes cipher block: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher: creates an aesgcm: %w", err)
	}
	text, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("cipher: decrypt ciphertext: %w", err)
	}
	return text, nil
}
