package cipher

import (
	"crypto/rand"
	"testing"
)

func TestAES256GCM(t *testing.T) {
	// The key argument should be the AES key, either 16, 24 or 32 bytes
	// to select AES-128, AES-192 or AES-256.
	key256 := make([]byte, 32)
	if _, err := rand.Read(key256); err != nil {
		t.Fatal(err)
	}
	text := []byte("test test test")

	ciphertext, err := AESGCMEncrypt(key256, text)
	if err != nil {
		t.Fatal(err)
	}
	plaintext, err := AESGCMDecrypt(key256, ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := string(text), string(plaintext); g != w {
		t.Errorf("plaintext got %v, want %v", g, w)
	}
}
