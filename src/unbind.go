package main

import (
	"bytes"
	"encoding/json"
)

func unBind(token string, id string) (string, error) {
	var err error
	var payload, body, output []byte

	payload = []byte(
		`{
			"active": false
		 }`,
	)

	if body, err = request(payload, token, id, "unbind"); err != nil {
		panic(err)
	}

	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, body, "", "    "); err == nil {
		output = prettyJSON.Bytes()
	} else {
		panic(err)
	}

	return string(output), nil
}
