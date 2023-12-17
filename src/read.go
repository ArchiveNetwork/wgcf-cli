package main

import (
	"encoding/json"
	"io"
	"os"
)

func readConfigFile(filePath string) (string, string, error) {
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

	var ReadedFile Response

	if err = json.Unmarshal(content, &ReadedFile); err != nil {
		panic(err)
	}

	return ReadedFile.Token, ReadedFile.ID, nil
}
