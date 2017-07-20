package main

import (
	"crypto/rsa"
	"log"
	"crypto/rand"
	"math/big"
)

func genRSAPrivateKey() *rsa.PrivateKey {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to gen key: %v", err)
	}
	return priv
}

func randSerialNumber() *big.Int {
	// generate a random serial number (a real Certificate authority would have some logic behind this)
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	sn, err := rand.Int(rand.Reader, limit)
	if err != nil {
		log.Fatalf("failed to gen serial number: %v", err)
	}
	return sn
}