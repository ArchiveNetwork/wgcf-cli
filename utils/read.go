package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
)

func ReadConfig(filePath string) []byte {
	var file *os.File
	var err error
	var body []byte
	if file, err = os.Open(filePath); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	if body, err = io.ReadAll(file); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	return body
}

func GetTokenID(filePath string) (string, string) {
	var response C.Response

	body := ReadConfig(filePath)

	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	return response.Token, response.ID
}
