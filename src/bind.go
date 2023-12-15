package main

import (
	"encoding/json"
	"time"
)

func getBindingDevices(token string, id string) (string, error) {
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

	body, err := request([]byte(``), token, id, "bind")
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
