package utils

import (
	"bytes"
	"fmt"
	"net/http"
)

type Request struct {
	Payload   []byte
	Token     string
	TeamToken string
	ID        string
	Action    string
}

func (r Request) New() (request *http.Request, err error) {
	var url, method string
	if r.Action == "register" {
		url = "https://api.cloudflareclient.com/v0a2158/reg"
		method = "POST"
	} else if r.Action == "license" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + r.ID + "/account"
		method = "PUT"
	} else if r.Action == "bind" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + r.ID + "/account/devices"
		method = "GET"
	} else if r.Action == "name" || r.Action == "unbind" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + r.ID + "/account/reg/" + r.ID
		method = "PATCH"
	} else if r.Action == "cancel" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + r.ID
		method = "DELETE"
	} else if r.Action == "update" {
		url = "https://api.cloudflareclient.com/v0a2158/reg/" + r.ID
		method = "GET"
	} else {
		err = fmt.Errorf("no action specified")
		return
	}

	if request, err = http.NewRequest(method, url, bytes.NewBuffer(r.Payload)); err != nil {
		return nil, err
	}
	request.Header.Add("CF-Client-Version", "a-7.21-0721")
	request.Header.Add("User-Agent", "okhttp/0.7.21")
	if r.Token != "" {
		request.Header.Add("Authorization", "Bearer "+r.Token)
	}
	if r.TeamToken != "" {
		request.Header.Add("Cf-Access-Jwt-Assertion", r.TeamToken)
	}

	return request, nil
}
