package jwtgo

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
)

func TestHmac(t *testing.T) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		t.Fatal(err)
	}
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		t.Fatal(err)
	}
	exp := strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo":   "bar",
		"exp":   exp,
		"nonce": fmt.Sprintf("%x", nonce),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// NOTICE: Don't forget validate the alg is what you expect.
		// https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries//
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %T", token.Method)
		}
		return secret, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Errorf("unexpected type of claims: %T", parsedToken.Claims)
	}
	if g, w := claims["foo"], "bar"; g != w {
		t.Errorf("foo got %v, want %v", g, w)
	}
	if g, w := claims["exp"], exp; g != w {
		t.Errorf("exp got %v, want %v", g, w)
	}
	if g, w := claims["nonce"], fmt.Sprintf("%x", nonce); g != w {
		t.Errorf("nonce got %v, want %v", g, w)
	}
}
