package main

import (
	"encoding/json"
	"io"
	"os"
)

func readConfigFile(filePath string) (string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var ReadedFile Response
	err = json.Unmarshal(content, &ReadedFile)
	if err != nil {
		panic(err)
	}

	token := ReadedFile.Token
	id := ReadedFile.ID

	return token, id, nil
}
