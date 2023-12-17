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
	var err error
	var body, output []byte
	payload := []byte(
		`{
			"name":"` + name + `"
		 }`,
	)

	if body, err = request(payload, token, id, "name"); err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	if output, err = json.MarshalIndent(response, "", "    "); err != nil {
		panic(err)
	}
	return string(output), nil
}
