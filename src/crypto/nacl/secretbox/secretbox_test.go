package secretbox

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/nacl/secretbox"
)

func TestSecretBox(t *testing.T) {
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		t.Fatal(err)
	}
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		t.Fatal(err)
	}
	// This encrypts "hello world" and appends the result to the nonce.
	encrypted := secretbox.Seal(nonce[:], []byte("hello world"), &nonce, &secret)

	// decrypt
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	got, ok := secretbox.Open(nil, encrypted[24:], &decryptNonce, &secret)
	assert.True(t, ok)
	assert.Equal(t, "hello world", string(got))
}
