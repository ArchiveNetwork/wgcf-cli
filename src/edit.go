package main

import (
	"encoding/json"
	"io"
	"os"
)

func editFile(filePath string, license string) error {
	var file *os.File
	var content []byte
	var err error
	var updatedContent []byte
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
	ReadedFile.Account.License = license

	if updatedContent, err = json.MarshalIndent(ReadedFile, "", "  "); err != nil {
		panic(err)
	}

	if err = os.WriteFile(filePath, updatedContent, 0600); err != nil {
		panic(err)
	}

	return nil
}
