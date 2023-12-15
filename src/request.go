package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func request(payload []byte, token string, id string, action string) ([]byte, error) {
	var url, method string
	if action == "register" {
		url = "https://api.cloudflareclient.com/v0a2158/reg"
		method = "POST"
	} else if action == "license" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + id + "/account"
		method = "PUT"
	} else if action == "bind" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + id + "/account/devices"
		method = "GET"
	} else if action == "name" || action == "unbind" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + id + "/account/reg/" + id
		method = "PATCH"
	} else if action == "cancle" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + id
		method = "DELETE"
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			},
		},
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		panic(err)
	}
	request.Header.Add("CF-Client-Version", "a-7.21-0721")
	request.Header.Add("User-Agent", "okhttp/0.7.21")
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	statusCode := response.StatusCode
	if statusCode != 204 && statusCode != 200 {
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, body, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(prettyJSON.String())
		panic("REST API returned " + fmt.Sprint(statusCode) + " " + http.StatusText(statusCode))
	}

	return body, nil
}
