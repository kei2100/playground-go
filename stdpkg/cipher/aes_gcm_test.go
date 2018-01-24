package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"testing"
)

func TestAES256GCM(t *testing.T) {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	key256 := make([]byte, 32)
	if _, err := rand.Read(key256); err != nil {
		t.Fatal(err)
	}
	text := []byte("test test test")

	ciphertext, nonce := aes256GCMEncrypt(t, key256, text)
	plaintext := aes256GCMDecrypt(t, key256, nonce, ciphertext)

	if g, w := string(text), string(plaintext); g != w {
		t.Errorf("plaintext got %v, want %v", g, w)
	}
}

func aes256GCMEncrypt(t *testing.T, key, text []byte) (ciphertext, nonce []byte) {
	t.Helper()

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatal(err)
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	// 12 byte is cipher.gcmStandardNonceSize
	nonce = make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		t.Fatal(err)
	}
	ciphertext = aesgcm.Seal(nil, nonce, text, nil)
	return ciphertext, nonce
}

func aes256GCMDecrypt(t *testing.T, key, nonce, ciphertext []byte) []byte {
	t.Helper()

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatal(err)
	}

	text, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		t.Fatal(err)
	}
	return text
}
