package utils

import (
	"bytes"
	"encoding/json"
)

func ChangeLicense(token string, id string, license string) (string, error) {
	var body []byte
	var err error
	var output []byte
	payload := []byte(
		`{
			"license":"` + license + `"
		 }`,
	)

	if body, err = request(payload, token, id, "license"); err != nil {
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
