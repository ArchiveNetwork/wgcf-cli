package main

import (
	"encoding/json"
	"time"
)

func changeName(token string, id string, name string) (string, error) {
	type Device struct {
		ID        string    `json:"id"`
		Type      string    `json:"type"`
		Model     string    `json:"model"`
		Name      string    `json:"name,omitempty"`
		Created   time.Time `json:"created"`
		Activated time.Time `json:"activated"`
		Active    bool      `json:"active"`
		Role      string    `json:"role"`
	}
	var response []Device
	payload := []byte(
		`{
			"name":"` + name + `"
		 }`,
	)
	body, err := request(payload, token, id, "name")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}
	output, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(output), nil
}
