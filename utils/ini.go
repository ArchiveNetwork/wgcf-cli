package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-ini/ini"
)

func ConvertJsonToIni(filePath string) error {
	var err error
	if err = WriteIniConfig(filePath, nil); err != nil {
		panic(err)
	}
	return err
}

func IniGetTokenID(filePath string) (string, string, error) {
	var err error
	var cfg *ini.File

	if cfg, err = ini.Load(filePath); err != nil {
		panic(err)
	}
	section_Account := cfg.Section("Account")
	id := section_Account.Key("ID").String()
	token := section_Account.Key("Token").String()
	return token, id, nil
}

func UpdateIniConfig(filePath string) error {
	var token, id string
	var err error
	var body []byte
	var ReadedFile Response
	if token, id, err = IniGetTokenID(filePath); err != nil {
		panic(err)
	}
	if body, err = request([]byte(``), token, id, "update"); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &ReadedFile); err != nil {
		panic(err)
	}
	ReadedFile.Token = token
	if err = WriteIniConfig(filePath, &ReadedFile); err != nil {
		panic(err)
	}

	return nil
}

func WriteIniConfig(filePath string, ReadedFile *Response) error {
	var err error
	var section_Account, section_Usage, section_Config *ini.Section
	var body []byte
	var cfg *ini.File
	var fileType string
	if fileType, err = GetFileType(filePath); err != nil {
		panic(err)
	}
	if fileType == "ini" {
		if cfg, err = ini.Load(filePath); err != nil {
			panic(err)
		}
		section_Config = cfg.Section("Config")
		ReadedFile.Config.ReservedDec, ReadedFile.Config.ReservedHex = clientIDtoReserved(section_Config.Key("ClientID").String())
		ReadedFile.Config.PrivateKey = section_Config.Key("PrivateKey").String()
	} else {
		body = ReadConfig(filePath)
		if err = json.Unmarshal(body, &ReadedFile); err != nil {
			panic(err)
		}
	}

	cfg = ini.Empty()
	if section_Account, err = cfg.NewSection("Account"); err != nil {
		panic(err)
	}
	section_Account.NewKey("ID", ReadedFile.ID)
	section_Account.NewKey("Token", ReadedFile.Token)
	if section_Usage, err = cfg.NewSection("Usage"); err != nil {
		panic(err)
	}
	section_Usage.NewKey("PremiumData", fmt.Sprint(float32(ReadedFile.Account.PremiumData)/1024/1024/1024, "GB"))
	section_Usage.NewKey("License", fmt.Sprint(ReadedFile.Account.License))
	if section_Config, err = cfg.NewSection("Config"); err != nil {
		panic(err)
	}
	section_Config.NewKey("PrivateKey", ReadedFile.Config.PrivateKey)
	section_Config.NewKey("IPv4", fmt.Sprint(ReadedFile.Config.Interface.Addresses.V4))
	section_Config.NewKey("IPv6", fmt.Sprint(ReadedFile.Config.Interface.Addresses.V6))
	section_Config.NewKey("ClientID", fmt.Sprint(ReadedFile.Config.ClientID))
	section_Config.NewKey("ReservedHex", fmt.Sprint(ReadedFile.Config.ReservedHex))
	section_Config.NewKey("ReservedDec", strings.Join(strings.Fields(fmt.Sprint(ReadedFile.Config.ReservedDec)), ", "))
	filePath = strings.Replace(filePath, ".json", ".ini", 1)
	cfg.SaveTo(filePath)
	os.Chmod(filePath, 0600)
	return nil
}

func GetFileType(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	var jsonData interface{}
	err = json.Unmarshal(content, &jsonData)
	if err == nil {
		return "json", nil
	}
	_, err = ini.LoadSources(ini.LoadOptions{AllowShadows: true}, strings.NewReader(string(content)))
	if err == nil {
		return "ini", nil
	}
	panic("File type not supported")
}
