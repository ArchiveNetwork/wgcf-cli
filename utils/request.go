package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	U "net/url"
	"os"
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
	} else if action == "update" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + id
		method = "GET"
	}
	var body []byte
	var request *http.Request
	var response *http.Response
	var err error
	var proxy string
	var proxyURL *U.URL

	httpProxy := os.Getenv("http_proxy")
	httpsProxy := os.Getenv("https_proxy")
	if httpProxy != "" {
		proxy = httpProxy
	} else if httpsProxy != "" {
		proxy = httpsProxy
	}
	if proxyURL, err = U.Parse(proxy); err != nil {
		panic("Error parsing proxy URL: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			},
		},
	}
	if proxy != "" {
		client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
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
