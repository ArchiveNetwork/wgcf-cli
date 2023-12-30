package utils

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"
)

func clientIDtoReserved(clientID string) ([]int, string) {
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
