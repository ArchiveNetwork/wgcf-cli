package utils

import (
	"encoding/json"
	"os"
)

func UpdateConfigFile(filePath string) error {
	var ReadedFile, response Response
	var file *os.File
	var err error
	content := ReadConfig(filePath)
	var body []byte
	var updatedContent []byte

	if err = json.Unmarshal(content, &ReadedFile); err != nil {
		panic(err)
	}
	defer file.Close()
	if ReadedFile.Config.ReservedDec == nil || ReadedFile.Config.ReservedHex == "" {
		response.Config.ReservedDec, response.Config.ReservedHex = clientIDtoReserved(ReadedFile.Config.ClientID)
	} else {
		response.Config.ReservedDec = ReadedFile.Config.ReservedDec
		response.Config.ReservedHex = ReadedFile.Config.ReservedHex
	}

	if body, err = request([]byte(``), ReadedFile.Token, ReadedFile.ID, "update"); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	if ReadedFile.Account.PrivateKey != "" {
		response.Config.PrivateKey = ReadedFile.Account.PrivateKey
	} else {
		response.Config.PrivateKey = ReadedFile.Config.PrivateKey
	}
	response.Token = ReadedFile.Token
	if updatedContent, err = json.MarshalIndent(response, "", "    "); err != nil {
		panic(err)
	}

	if err = os.WriteFile(filePath, updatedContent, 0600); err != nil {
		panic(err)
	}
	return nil
}
