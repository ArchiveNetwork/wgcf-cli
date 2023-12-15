package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func readLicense(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var ReadedFile Response
	err = json.Unmarshal(content, &ReadedFile)
	if err != nil {
		panic(err)
	}

	return ReadedFile.Account.License, nil
}

func TestMain(m *testing.M) {
	var action Actions
	var err error

	action.FileName = "wgcf.json"

	store, output, err := register()
	if err != nil {
		panic(err)
	}
	fmt.Println(output)

	err = os.WriteFile("test.json", store, 0600)
	if err != nil {
		panic(err)
	}

	action.License, err = readLicense("test.json")
	if err != nil {
		panic(err)
	}

	config, reserved, err := configGenerate("wireguard", "test.json")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("test.conf", []byte(config), 0600)
	if err != nil {
		panic(err)
	}
	output, err = nftConfigGenerate(reserved)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("test.nft.conf", []byte(output), 0600)
	if err != nil {
		panic(err)
	}

	config, _, err = configGenerate("xray", "test.json")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("wgcf.xray.json", []byte(config), 0600)
	if err != nil {
		panic(err)
	}

	config, _, err = configGenerate("sing-box", "test.json")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("wgcf.sing-box.json", []byte(config), 0600)
	if err != nil {
		panic(err)
	}

	token, id, err := readConfigFile(action.FileName)
	if err != nil {
		panic(err)
	}

	output, err = getBindingDevices(token, id)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)

	output, err = unBind(token, id)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)

	output, err = changeLicense(token, id, action.License)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)

	err = editFile(action.FileName, action.License)
	if err != nil {
		panic(err)
	}

	output, err = changeName(token, id, action.Name)
	if err != nil {
		panic(err)
	}
	fmt.Println(output)

	token, id, err = readConfigFile("test.json")
	if err != nil {
		panic(err)
	}

	err = cancleAccount(token, id)
	if err != nil {
		panic(err)
	}
	fmt.Println("Cancled")

	err = os.Remove("test.json")
	if err != nil {
		panic(err)
	}

	err = os.Remove(action.FileName)
	if err != nil {
		panic(err)
	}

	err = os.Remove("test.conf")
	if err != nil {
		panic(err)
	}

	err = os.Remove("test.nft.conf")
	if err != nil {
		panic(err)
	}

	err = os.Remove("wgcf.xray.json")
	if err != nil {
		panic(err)
	}

	err = os.Remove("wgcf.sing-box.json")
	if err != nil {
		panic(err)
	}
}
