package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var ErrMismatchedHashAndPassword = errors.New("argon2: mismatched hash and password")

func GenerateFromPassword(password string) (string, error) {
	// RFC recommended Argon2id parameters
	// https://datatracker.ietf.org/doc/html/rfc9106#section-4
	const iterations = 3
	const threads = 4
	const memory = 64 * 1024
	// salt and key length
	const saltLen = 16
	const keyLen = 32
	// generate salt
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("argon2: generate salt %w", err)
	}
	// generate key
	key := argon2.IDKey([]byte(password), salt, iterations, memory, threads, keyLen)
	//　encode like https://github.com/P-H-C/phc-winner-argon2#command-line-utility
	var b strings.Builder
	b.WriteString("$argon2id")
	b.WriteString(fmt.Sprintf("$v=%d", argon2.Version))
	b.WriteString(fmt.Sprintf("$m=%d,t=%d,p=%d", memory, iterations, threads))
	b.WriteString("$" + base64.RawStdEncoding.EncodeToString(salt))
	b.WriteString("$" + base64.RawStdEncoding.EncodeToString(key))
	return b.String(), nil
}

func CompareHashAndPassword(encodedHash, password string) error {
	// endocodedHash e.g
	// $argon2id$v=19$m=65536,t=3,p=4$dDfMhYJIkUq8fzLMM7+tiw$AOWjpr5Psw3HxDMGSdPJdzEktQ/d3OIzZ/wuQWtUBvk
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return errors.New("argon2: unexpected hash format")
	}
	if vals[1] != "argon2id" {
		return fmt.Errorf("argon2: unsupported %s", vals[1])
	}
	var version int
	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return fmt.Errorf("argon2: scan version: %w", err)
	}
	if version != argon2.Version {
		return fmt.Errorf("argon2: unsupported version %d", version)
	}
	var iterations, memory uint32
	var threads uint8
	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &threads); err != nil {
		return fmt.Errorf("argon2: scan m t p: %w", err)
	}
	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return fmt.Errorf("argon2: decode salt %s", vals[4])
	}
	key, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return fmt.Errorf("argon2: decode salt %s", vals[5])
	}
	const keyLen = 32
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, threads, keyLen)
	if subtle.ConstantTimeCompare(key, hash) != 1 {
		return ErrMismatchedHashAndPassword
	}
	return nil
}
