package main

import (
	"encoding/json"
	"io"
	"os"
)

func readConfig(filePath string) []byte {
	var file *os.File
	var err error
	var content []byte
	if file, err = os.Open(filePath); err != nil {
		panic(err)
	}
	defer file.Close()

	if content, err = io.ReadAll(file); err != nil {
		panic(err)
	}
	return content
}

func getTokenID(filePath string) (string, string, error) {
	var err error
	var content []byte
	var ReadedFile Response

	content = readConfig(filePath)

	if err = json.Unmarshal(content, &ReadedFile); err != nil {
		panic(err)
	}
	return ReadedFile.Token, ReadedFile.ID, nil
}
