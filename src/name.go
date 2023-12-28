package main

import (
	"bytes"
	"encoding/json"
)

func changeName(token string, id string, name string) (string, error) {
	var err error
	var body, output []byte
	var prettyJSON bytes.Buffer
	payload := []byte(
		`{
			"name":"` + name + `"
		 }`,
	)
	if body, err = request(payload, token, id, "name"); err != nil {
		panic(err)
	}

	if err = json.Indent(&prettyJSON, body, "", "    "); err == nil {
		output = prettyJSON.Bytes()
	} else {
		panic(err)
	}

	return string(output), nil
}
