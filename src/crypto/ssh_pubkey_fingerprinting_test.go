package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestFingerPrints(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	sshpub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	// FingerprintLegacyMD5 returns the user presentation of the key's fingerprint as described by RFC 4716 section 4.
	log.Println(ssh.FingerprintLegacyMD5(sshpub))
	// => e.g. c7:9a:03:86:86:8d:c7:06:99:32:88:8c:ac:99:ee:0a

	// FingerprintSHA256 returns the user presentation of the key's fingerprint as unpadded base64 encoded sha256 hash.
	// This format was introduced from OpenSSH 6.8. https://www.openssh.com/txt/release-6.8
	// https://tools.ietf.org/html/rfc4648#section-3.2  (unpadded base64 encoding)
	log.Println(ssh.FingerprintSHA256(sshpub))
	// => e.g. SHA256:Snf3igD5pr/RYc0kV8eYYJAhTa08b69seeR76Id5684
}
