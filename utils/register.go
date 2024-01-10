package utils

import (
	"encoding/json"
	"time"
)

func Register(teamToken string) ([]byte, string, error) {
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
	if teamToken != "" {
		if body, err = request(payload, teamToken, "", "registerTeam"); err != nil {
			panic(err)
		}
	} else {
		if body, err = request(payload, "", "", "register"); err != nil {
			panic(err)
		}
	}

	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	response.Config.ReservedDec, response.Config.ReservedHex = clientIDtoReserved(response.Config.ClientID)
	response.Config.PrivateKey = privateKey
	jsonIn := RegisterOutput{
		Endpoint: struct {
			V4 string `json:"v4"`
			V6 string `json:"v6"`
		}{
			V4: response.Config.Peers[0].Endpoint.V4,
			V6: response.Config.Peers[0].Endpoint.V6,
		},
		ReservedStr: response.Config.ClientID,
		ReservedHex: response.Config.ReservedHex,
		ReservedDec: response.Config.ReservedDec,
		PrivateKey:  response.Config.PrivateKey,
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
