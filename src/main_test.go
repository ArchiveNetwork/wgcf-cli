package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func readLicense(filePath string) (string, error) {
	var file *os.File
	var content []byte
	var err error

	if file, err = os.Open(filePath); err != nil {
		panic(err)
	}
	defer file.Close()

	if content, err = io.ReadAll(file); err != nil {
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
	var store []byte
	var config, token, id, reserved, output string
	var err error

	action.FileName = "wgcf.json"
	if store, output, err = register(""); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if err = os.WriteFile(action.FileName, store, 0600); err != nil {
		panic(err)
	}

	if store, output, err = register(""); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if err = os.WriteFile("test.json", store, 0600); err != nil {
		panic(err)
	}

	action.License, err = readLicense("test.json")
	if err != nil {
		panic(err)
	}

	if config, reserved, err = configGenerate("wireguard", "test.json"); err != nil {
		panic(err)
	}

	if err = os.WriteFile("test.conf", []byte(config), 0600); err != nil {
		panic(err)
	}

	if output, err = nftConfigGenerate(reserved); err != nil {
		panic(err)
	}

	if err = os.WriteFile("test.nft.conf", []byte(output), 0600); err != nil {
		panic(err)
	}

	if config, _, err = configGenerate("xray", "test.json"); err != nil {
		panic(err)
	}

	if err = os.WriteFile("wgcf.xray.json", []byte(config), 0600); err != nil {
		panic(err)
	}

	if config, _, err = configGenerate("sing-box", "test.json"); err != nil {
		panic(err)
	}

	if err = os.WriteFile("wgcf.sing-box.json", []byte(config), 0600); err != nil {
		panic(err)
	}

	if token, id, err = readConfigFile(action.FileName); err != nil {
		panic(err)
	}

	if output, err = getBindingDevices(token, id); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if output, err = unBind(token, id); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if output, err = changeLicense(token, id, action.License); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if output, err = changeName(token, id, action.Name); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if token, id, err = readConfigFile("test.json"); err != nil {
		panic(err)
	}

	if err = cancleAccount(token, id); err != nil {
		panic(err)
	}
	fmt.Println("Cancled")

	if token, id, err = readConfigFile(action.FileName); err != nil {
		panic(err)
	}

	if err = cancleAccount(token, id); err != nil {
		panic(err)
	}
	fmt.Println("Cancled")

	if err = os.Remove("test.json"); err != nil {
		panic(err)
	}

	if err = os.Remove(action.FileName); err != nil {
		panic(err)
	}

	if err = os.Remove("test.conf"); err != nil {
		panic(err)
	}

	if err = os.Remove("test.nft.conf"); err != nil {
		panic(err)
	}

	if err = os.Remove("wgcf.xray.json"); err != nil {
		panic(err)
	}

	if err = os.Remove("wgcf.sing-box.json"); err != nil {
		panic(err)
	}
}
