package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var action Actions
	action.FileName = "wgcf.json"
	action.License = "m48b9e0W-3N9eR1T7-78I0d1Qh"
	store, output, err := register()
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
	err = os.WriteFile("test.json", store, 0600)
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

	err = os.Remove(action.FileName)
	if err != nil {
		panic(err)
	}

	err = os.Remove("test.json")
	if err != nil {
		panic(err)
	}

}
