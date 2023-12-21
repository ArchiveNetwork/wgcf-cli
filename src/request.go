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
	var teamToken string
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
	} else if action == "registerTeam" {
		url = "https://api.cloudflareclient.com/v0a2158/reg"
		method = "POST"
		teamToken = token
		token = ""
	}
	var body []byte
	var request *http.Request
	var response *http.Response
	var err error

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			},
		},
	}

	if request, err = http.NewRequest(method, url, bytes.NewBuffer(payload)); err != nil {
		panic(err)
	}
	request.Header.Add("CF-Client-Version", "a-7.21-0721")
	request.Header.Add("User-Agent", "okhttp/0.7.21")
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}
	if teamToken != "" {
		request.Header.Add("Cf-Access-Jwt-Assertion", teamToken)
	}

	if response, err = client.Do(request); err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if body, err = io.ReadAll(response.Body); err != nil {
		panic(err)
	}

	if response.StatusCode != 204 && response.StatusCode != 200 {
		var prettyJSON bytes.Buffer
		if err = json.Indent(&prettyJSON, body, "", "    "); err != nil {
			fmt.Println(string(body))
		} else {
			fmt.Println(prettyJSON.String())
		}
		err = fmt.Errorf("REST API returned " + fmt.Sprint(response.StatusCode) + " " + http.StatusText(response.StatusCode))
	}

	return body, err
}
