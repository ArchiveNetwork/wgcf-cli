package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

func ConvertJsonToIni(filePath string, fileType string) error {
	var err error
	if err = WriteIniConfig(filePath, nil, fileType); err != nil {
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
	if err = WriteIniConfig(filePath, &ReadedFile, "ini"); err != nil {
		panic(err)
	}

	return nil
}

func WriteIniConfig(filePath string, ReadedFile *Response, fileType string) error {
	var err error
	var section_Account, section_Usage, section_Config *ini.Section
	var body []byte
	var cfg *ini.File
	if fileType == "ini" {
		if cfg, err = ini.Load(filePath); err != nil {
			panic(err)
		}
		section_Config = cfg.Section("Config")
		ReadedFile.Config.ReservedDec, ReadedFile.Config.ReservedHex = clientIDtoReserved(section_Config.Key("ClientID").String())
		ReadedFile.Config.PrivateKey = section_Config.Key("PrivateKey").String()
	} else if fileType == "json" {
		body = ReadConfig(filePath)
		if err = json.Unmarshal(body, &ReadedFile); err != nil {
			panic(err)
		}
	} else {
		panic("No file type detected")
	}

	cfg = ini.Empty()
	if section_Account, err = cfg.NewSection("Account"); err != nil {
		panic(err)
	}
	section_Account.NewKey("ID", ReadedFile.ID)
	section_Account.NewKey("Token", ReadedFile.Token)
	section_Account.NewKey("Type", ReadedFile.Account.AccountType)
	if section_Account.Key("Type").String() != "team" {
		if section_Usage, err = cfg.NewSection("Usage"); err != nil {
			panic(err)
		}
		section_Usage.NewKey("PremiumData", fmt.Sprint(float32(ReadedFile.Account.PremiumData)/1000/1000/1000, "GB"))
		section_Usage.NewKey("License", fmt.Sprint(ReadedFile.Account.License))
	}
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
		panic(err)
	}
	var jsonData Response
	err = json.Unmarshal(content, &jsonData)
	if err == nil {
		return "json", nil
	}
	var cfg *ini.File
	if cfg, err = ini.Load(filePath); err != nil {
		panic(err)
	}
	section := cfg.Section("Account")
	if id := section.Key("ID").String(); id != "" {
		return "ini", nil
	}
	panic("File type not supported")
}

func ConvertIniToJson(filePath string) error {
	var err error
	var token, id string
	if token, id, err = IniGetTokenID(filePath); err != nil {
		panic(err)
	}
	jsonFilePath := strings.Replace(filePath, ".ini", ".json", 1)
	tmpPath := os.TempDir() + "/tmp." + filePath
	if runtime.GOOS == "windows" {
		tmpPath = os.TempDir() + "\\tmp." + filePath
	}
	defer os.Remove(tmpPath)
	if err = os.WriteFile(tmpPath, []byte(`{"id":"`+id+`","token":"`+token+`"}`), 0600); err != nil {
		panic(err)
	}
	if err = UpdateConfigFile(tmpPath); err != nil {
		panic(err)
	}
	var ReadedFile Response
	if err = json.Unmarshal(ReadConfig(tmpPath), &ReadedFile); err != nil {
		panic(err)
	}
	var cfg *ini.File
	if cfg, err = ini.Load(filePath); err != nil {
		panic(err)
	}
	section_Config := cfg.Section("Config")
	ReadedFile.Config.PrivateKey = section_Config.Key("PrivateKey").String()
	ReadedFile.Config.ClientID = section_Config.Key("ClientID").String()
	ReadedFile.Config.ReservedHex = section_Config.Key("ReservedHex").String()
	str := strings.Trim(section_Config.Key("ReservedDec").String(), "[]")
	strSlice := strings.Split(str, ",")
	intSlice := make([]int, len(strSlice))
	for i, s := range strSlice {
		s = strings.TrimSpace(s)
		num, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		intSlice[i] = num
	}
	ReadedFile.Config.ReservedDec = intSlice
	var updatedContent []byte
	if updatedContent, err = json.MarshalIndent(ReadedFile, "", "    "); err != nil {
		panic(err)
	}
	if err = os.WriteFile(jsonFilePath, updatedContent, 0600); err != nil {
		panic(err)
	}
	return err
}
