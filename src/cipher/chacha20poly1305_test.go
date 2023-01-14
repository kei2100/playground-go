package cipher

import (
	"crypto/rand"
	"testing"
)

func TestXChaCha20Poly1305EncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	encrypted, err := XChaCha20Poly1305Encrypt(key, []byte("foo bar baz"))
	if err != nil {
		t.Fatal(err)
	}
	decrypted, err := XChaCha20Poly1305Decrypt(key, encrypted)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := string(decrypted), "foo bar baz"; g != w {
		t.Errorf("\ngot :%v\nwant:%v", string(decrypted), "foo bar baz")
	}
}
