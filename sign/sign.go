package sign

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/tula3and/me-sign/utils"
)

const (
	fileName string = "server.key"
)

var k *ecdsa.PrivateKey

func hasServerKey() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func CreatePrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

func restoreKey() *ecdsa.PrivateKey {
	bytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	key, err := x509.ParseECPrivateKey(bytes)
	utils.HandleErr(err)
	return key
}

func Sign(payload string, k *ecdsa.PrivateKey) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, k, payloadAsB)
	utils.HandleErr(err)
	signedText := append(r.Bytes(), s.Bytes()...)
	return fmt.Sprintf("%x", signedText)
}

func RestorePublicKey(k *ecdsa.PrivateKey) string {
	pK := append(k.X.Bytes(), k.Y.Bytes()...)
	return fmt.Sprintf("%x", pK)
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	Bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := Bytes[:len(Bytes)/2]
	secondHalfBytes := Bytes[len(Bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func Verify(signedText, payload, address string) bool {
	r, s, err := restoreBigInts(signedText)
	utils.HandleErr(err)
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	// create the public key from address
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

func Key() *ecdsa.PrivateKey {
	if k == nil {
		if hasServerKey() {
			k = restoreKey()
		} else {
			k = CreatePrivKey()
			persistKey(k)
		}
	}
	return k
}
