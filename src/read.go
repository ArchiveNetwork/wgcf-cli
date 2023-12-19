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
	var genericResponse map[string]interface{}
	if file, err = os.Open(filePath); err != nil {
		panic(err)
	}
	defer file.Close()

	if content, err = io.ReadAll(file); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(content, &genericResponse); err == nil {
		if account, ok := genericResponse["account"].(map[string]interface{}); ok {
			if _, ok := account["license"]; ok {
				var ReadedFile NormalResponse
				if err = json.Unmarshal(content, &ReadedFile); err != nil {
					panic(err)
				}
				return ReadedFile.Token, ReadedFile.ID, nil
			}
		}
		if _, ok := genericResponse["version"]; ok {
			var ReadedFile TeamResponse
			if err = json.Unmarshal(content, &ReadedFile); err != nil {
				panic(err)
			}
			return ReadedFile.Token, ReadedFile.ID, nil
		}
	}
	panic(nil)
}
