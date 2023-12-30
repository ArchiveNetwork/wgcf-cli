package utils

import (
	"bytes"
	"encoding/json"
)

func GetBindingDevices(token string, id string) (string, error) {
	var body, output []byte
	var prettyJSON bytes.Buffer
	var err error

	if body, err = request([]byte(``), token, id, "bind"); err != nil {
		panic(err)
	}

	if err = json.Indent(&prettyJSON, body, "", "    "); err == nil {
		output = prettyJSON.Bytes()
	} else {
		panic(err)
	}
	return string(output), nil
}
