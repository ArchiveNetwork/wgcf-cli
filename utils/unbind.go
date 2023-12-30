package utils

import (
	"bytes"
	"encoding/json"
)

func UnBind(token string, id string) (string, error) {
	var err error
	var payload, body, output []byte
	var prettyJSON bytes.Buffer
	payload = []byte(
		`{
			"active": false
		 }`,
	)
	if body, err = request(payload, token, id, "unbind"); err != nil {
		panic(err)
	}

	if err = json.Indent(&prettyJSON, body, "", "    "); err == nil {
		output = prettyJSON.Bytes()
	} else {
		panic(err)
	}

	return string(output), nil
}
