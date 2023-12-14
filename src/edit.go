package main

import (
	"encoding/json"
	"io"
	"os"
)

func editFile(filePath string, license string) error {
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
	ReadedFile.Account.License = license
	updatedContent, err := json.MarshalIndent(ReadedFile, "", "  ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filePath, updatedContent, 0600)
	if err != nil {
		panic(err)
	}

	return nil
}
