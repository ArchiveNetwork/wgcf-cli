package main_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ArchiveNetwork/wgcf-cli/cmd"
	"github.com/ArchiveNetwork/wgcf-cli/utils"
)

func readLicense(filePath string) (string, error) {
	var ReadedFile utils.Response

	content := utils.ReadConfig(filePath)

	if err := json.Unmarshal(content, &ReadedFile); err != nil {
		panic(err)
	}

	return ReadedFile.Account.License, nil
}

func TestMain(m *testing.M) {
	var action cmd.Actions
	var store []byte
	var config, token, id, reserved, output string
	var err error

	action.FileName = "wgcf.json"
	if store, output, err = utils.Register(""); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if err = os.WriteFile(action.FileName, store, 0600); err != nil {
		panic(err)
	}

	if store, output, err = utils.Register(""); err != nil {
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

	if config, reserved, err = utils.ConfigGenerate("wireguard", "test.json"); err != nil {
		panic(err)
	}

	if err = os.WriteFile("test.conf", []byte(config), 0600); err != nil {
		panic(err)
	}

	if output, err = utils.NftConfigGenerate(reserved); err != nil {
		panic(err)
	}

	if err = os.WriteFile("test.nft.conf", []byte(output), 0600); err != nil {
		panic(err)
	}

	if config, _, err = utils.ConfigGenerate("xray", "test.json"); err != nil {
		panic(err)
	}

	if err = os.WriteFile("wgcf.xray.json", []byte(config), 0600); err != nil {
		panic(err)
	}

	if config, _, err = utils.ConfigGenerate("sing-box", "test.json"); err != nil {
		panic(err)
	}

	if err = os.WriteFile("wgcf.sing-box.json", []byte(config), 0600); err != nil {
		panic(err)
	}

	if token, id, err = utils.GetTokenID(action.FileName); err != nil {
		panic(err)
	}

	if output, err = utils.GetBindingDevices(token, id); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if output, err = utils.UnBind(token, id); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if output, err = utils.ChangeLicense(token, id, action.License); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if output, err = utils.ChangeName(token, id, action.Name); err != nil {
		panic(err)
	}
	fmt.Println(output)

	if token, id, err = utils.GetTokenID("test.json"); err != nil {
		panic(err)
	}

	if err = utils.CancleAccount(token, id); err != nil {
		panic(err)
	}
	fmt.Println("Cancled")

	if token, id, err = utils.GetTokenID(action.FileName); err != nil {
		panic(err)
	}

	if err = utils.CancleAccount(token, id); err != nil {
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
