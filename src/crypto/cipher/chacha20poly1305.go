package cipher

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

// XChaCha20Poly1305Encrypt は XChaCha20-Poly1305 で text を暗号化する。key は 32byte であること。
// XChaCha20-Poly1305 は以下の特徴がある。
// * nonce が 192bit と長く、カウンターではなく乱数を使って nonce を生成する場合に比較的安全に使うことができる
// * 同じ鍵で実質的に無制限の数のメッセージを安全に暗号化でき、メッセージのサイズに実用的な制限がない（最大2^64バイトまで）
func XChaCha20Poly1305Encrypt(key, text []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("cipher: create aead: %w", err)
	}
	// Select a random nonce, and leave capacity for the ciphertext.
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(text)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("cipher: read rand: %w", err)
	}
	// Encrypt the message and append the ciphertext to the nonce.
	encrypted := aead.Seal(nonce, nonce, text, nil)
	return encrypted, nil
}

func XChaCha20Poly1305Decrypt(key, encrypted []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("cipher: create aead: %w", err)
	}
	// Split nonce and ciphertext.
	nonceSize := aead.NonceSize()
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]
	// Decrypt the message and check it wasn't tampered with.
	text, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("cipher: decrypt: %w", err)
	}
	return text, nil
}
