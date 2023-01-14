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
func AESGCMEncrypt(key []byte, text []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher: creates an aes cipher block: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher: creates an aesgcm: %w", err)
	}
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	// 12 byte is cipher.gcmStandardNonceSize
	// AES-GCM では 同一キーにおいて nonce の再利用は決してされるべきではない。本来は乱数よりは単調増加するカウンターなどの利用が推奨される。
	// （nonce 自体は推測困難である必要はない。）
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("cipher: creates a nonece: %w", err)
	}
	// Encrypt the text and append the ciphertext to the nonce.
	encrypted = aesgcm.Seal(nonce, nonce, text, nil)
	return encrypted, nil
}

func AESGCMDecrypt(key, encrypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher: creates an aes cipher block: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher: creates an aesgcm: %w", err)
	}
	nonceSize := aesgcm.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, fmt.Errorf("cipher: invalid encrypted value")
	}
	nonce := encrypted[:nonceSize]
	ciphertext := encrypted[nonceSize:]
	text, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("cipher: decrypt encrypted: %w", err)
	}
	return text, nil
}
