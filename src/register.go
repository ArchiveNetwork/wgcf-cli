package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"
)

func register() ([]byte, string, error) {
	type RegisterOutput struct {
		Endpoint struct {
			V4 string `json:"v4"`
			V6 string `json:"v6"`
		} `json:"endpoint"`
		ReservedStr string `json:"reserved_str"`
		ReservedHex string `json:"reserved_hex"`
		ReservedDec []int  `json:"reserved_dec"`
		PrivateKey  string `json:"private_key"`
		PublicKey   string `json:"public_key"`
		Addresses   struct {
			V4 string `json:"v4"`
			V6 string `json:"v6"`
		} `json:"addresses"`
	}
	var err error
	var output, body, store []byte
	var response Response
	var publicKey, privateKey string

	if privateKey, publicKey, err = GenerateKey(); err != nil {
		panic(err)
	}

	installID := RandStringRunes(22, nil)
	fcmtoken := RandStringRunes(134, nil)

	payload := []byte(
		`{
			"key":"` + publicKey + `",
			"install_id":"` + installID + `",
			"fcm_token":"` + installID + `:APA91b` + fcmtoken + `",
			"tos":"` + time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + `",
			"model":"Android",
			"serial_number":"` + installID + `"
		}`,
	)

	if body, err = request(payload, "", "", "register"); err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	clientID := response.Config.ClientID
	response.Account.PrivateKey = privateKey
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

	response.Account.ReservedDec = reserved
	response.Account.ReservedHex = "0x" + hexString
	jsonIn := RegisterOutput{
		Endpoint: struct {
			V4 string `json:"v4"`
			V6 string `json:"v6"`
		}{
			V4: response.Config.Peers[0].Endpoint.V4,
			V6: response.Config.Peers[0].Endpoint.V6,
		},
		ReservedStr: response.Config.ClientID,
		ReservedHex: response.Account.ReservedHex,
		ReservedDec: response.Account.ReservedDec,
		PrivateKey:  privateKey,
		PublicKey:   response.Config.Peers[0].PublicKey,
		Addresses:   response.Config.Interface.Addresses,
	}

	if output, err = json.MarshalIndent(jsonIn, "", "    "); err != nil {
		panic(err)
	}

	if store, err = json.MarshalIndent(response, "", "    "); err != nil {
		panic(err)
	}

	return store, string(output), nil
}
