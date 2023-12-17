package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func updateConfigFile(filePath string) error {
	var ReadedFile Response
	var responseBody Response
	var response *http.Response
	var request *http.Request
	var body []byte
	var client *http.Client
	var file *os.File
	var content []byte
	var err error
	var updatedContent []byte

	if file, err = os.Open(filePath); err != nil {
		panic(err)
	}
	if content, err = io.ReadAll(file); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(content, &ReadedFile); err != nil {
		panic(err)
	}
	defer file.Close()

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			},
		},
	}

	if request, err = http.NewRequest("GET", "https://api.cloudflareclient.com/v0a2158/reg/"+ReadedFile.ID, nil); err != nil {
		panic(err)
	}
	request.Header.Add("CF-Client-Version", "a-7.21-0721")
	request.Header.Add("User-Agent", "okhttp/0.7.21")
	request.Header.Add("Authorization", "Bearer "+ReadedFile.Token)

	if response, err = client.Do(request); err != nil {
		panic(err)
	}
	if body, err = io.ReadAll(response.Body); err != nil {
		panic(err)
	}
	if response.StatusCode != 204 && response.StatusCode != 200 {
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, body, "", "  "); err != nil {
			fmt.Println(string(body))
		} else {
			fmt.Println(prettyJSON.String())
		}
		panic("REST API returned " + fmt.Sprint(response.StatusCode) + " " + http.StatusText(response.StatusCode))
	}
	if err = json.Unmarshal(body, &responseBody); err != nil {
		panic(err)
	}
	responseBody.Account.ReservedDec = ReadedFile.Account.ReservedDec
	responseBody.Account.ReservedHex = ReadedFile.Account.ReservedHex
	responseBody.Account.PrivateKey = ReadedFile.Account.PrivateKey
	responseBody.Token = ReadedFile.Token

	if updatedContent, err = json.MarshalIndent(responseBody, "", "  "); err != nil {
		panic(err)
	}

	if err = os.WriteFile(filePath, updatedContent, 0600); err != nil {
		panic(err)
	}
	return nil
}
