package box

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/nacl/box"
)

// https://pkg.go.dev/golang.org/x/crypto/nacl/box
// Package box authenticates and encrypts small messages using public-key cryptography.
// Box uses Curve25519, XSalsa20 and Poly1305 to encrypt and authenticate messages.

func TestBox(t *testing.T) {
	senderPublicKey, senderPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	recipientPublicKey, recipientPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	// The shared key can be used to speed up processing when using the same
	// pair of keys repeatedly.
	sharedEncryptKey := new([32]byte)
	box.Precompute(sharedEncryptKey, recipientPublicKey, senderPrivateKey)

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	msg := []byte("A fellow of infinite jest, of most excellent fancy")
	// This encrypts msg and appends the result to the nonce.
	encrypted := box.SealAfterPrecomputation(nonce[:], msg, &nonce, sharedEncryptKey)

	// The shared key can be used to speed up processing when using the same
	// pair of keys repeatedly.
	var sharedDecryptKey [32]byte
	box.Precompute(&sharedDecryptKey, senderPublicKey, recipientPrivateKey)

	// The recipient can decrypt the message using the shared key. When you
	// decrypt, you must use the same nonce you used to encrypt the message.
	// One way to achieve this is to store the nonce alongside the encrypted
	// message. Above, we stored the nonce in the first 24 bytes of the
	// encrypted text.
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	decrypted, ok := box.OpenAfterPrecomputation(nil, encrypted[24:], &decryptNonce, &sharedDecryptKey)
	if !ok {
		panic("decryption error")
	}

	assert.Equal(t, string(msg), string(decrypted))
}

func TestBoxAnonymous(t *testing.T) {
	recipientPublicKey, recipientPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	// SealAnonymous appends an encrypted and authenticated copy of message to out,
	// which will be AnonymousOverhead bytes longer than the original and must not
	// overlap it. This differs from Seal in that the sender is not required to
	// provide a private key.j
	message := []byte("test message")
	encrypted, err := box.SealAnonymous(nil, message, recipientPublicKey, nil)
	if err != nil {
		panic(err)
	}

	decrypted, ok := box.OpenAnonymous(nil, encrypted, recipientPublicKey, recipientPrivateKey)
	if !ok {
		t.Fatalf("failed to open box")
	}

	assert.Equal(t, string(message), string(decrypted))
}
