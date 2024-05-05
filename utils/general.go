package utils

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/curve25519"
)

func ClientIDtoReserved(clientID string) ([]int, string) {
	decoded, err := base64.StdEncoding.DecodeString(clientID)
	if err != nil {
		panic(err)
	}
	hexString := hex.EncodeToString(decoded)

	reserved := []int{}
	for i := 0; i < len(hexString); i += 2 {
		hexByte := hexString[i : i+2]
		decValue, _ := strconv.ParseInt(hexByte, 16, 64)
		reserved = append(reserved, int(decValue))
	}
	hexString = "0x" + hexString
	return reserved, hexString
}

func RandStringRunes(n int, letterRunes []rune) string {
	if n <= 0 {
		return ""
	}

	if len(letterRunes) == 0 {
		letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	randomRunes := make([]rune, n)
	for i := range randomRunes {
		randomRunes[i] = letterRunes[random.Intn(len(letterRunes))]
	}
	return string(randomRunes)
}

func GenerateKey() (string, string, error) {
	var priv, pub []byte
	var err error

	priv = make([]byte, curve25519.ScalarSize)
	if _, err = crand.Read(priv); err != nil {
		panic(err)
	}

	priv[0] &= 248
	priv[31] &= 127 | 64

	if pub, err = curve25519.X25519(priv, curve25519.Basepoint); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(priv[:]), base64.StdEncoding.EncodeToString(pub[:]), nil
}
