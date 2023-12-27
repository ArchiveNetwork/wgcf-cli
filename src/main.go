package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	action := ParseCommandLine()
	var err error
	var store []byte
	var config, token, id, reserved, output string

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if action.Help {
		help()
		return
	}
	if action.Version {
		version()
		return
	}

	if action.Register {
		if strings.HasPrefix(action.TeamToken, "-") {
			err := fmt.Sprintln(`The parameter must not start with "-"`)
			panic(err)
		}
		if action.TeamToken != "" {
			if store, output, err = register(action.TeamToken); err != nil {
				panic(err)
			}
		} else {
			if store, output, err = register(""); err != nil {
				panic(err)
			}
		}
		fmt.Println(output)
		if action.FileName != "" {
			if err = os.WriteFile(action.FileName, store, 0600); err != nil {
				panic(err)
			}
			return
		} else {
			fileName := "wgcf.json"
			editedFileName := "wgcf.json"
			i := 0

			for {
				if _, err := os.Stat(fileName); err == nil {
					fileName = fmt.Sprintf("%s-%d.json", editedFileName[:len(editedFileName)-5], i)
					i++
				} else {
					break
				}
			}

			err := os.WriteFile(fileName, store, 0600)
			if err != nil {
				panic(err)
			}
			return
		}
	}

	if strings.HasPrefix(action.TeamToken, "-") {
		panic(`The parameter must not start with "-"`)
	}
	if action.TeamToken != "" {
		panic(`You need to use this parameter with "-r/--register"`)
	}

	if action.FileName == "" {
		action.FileName = "wgcf.json"
	} else if strings.HasPrefix(action.FileName, "-") {
		panic(`The parameter must not start with "-"`)
	}

	if !action.Bind && !action.UnBind && !action.Cancel && action.License == "" && action.Name == "" && action.Generate == "" && !action.Plus && !action.Version && !action.Update {
		panic("You need to specify an action")
	}

	if action.Bind {
		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := getBindingDevices(token, id)
		if err != nil {
			panic(err)
		}

		if err = updateConfigFile(action.FileName); err != nil {
			panic(err)
		}
		fmt.Println(output)

		return
	}

	if action.UnBind {
		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := unBind(token, id)
		if err != nil {
			panic(err)
		}

		if err = updateConfigFile(action.FileName); err != nil {
			panic(err)
		}

		fmt.Println(output)
		return
	}

	if action.Cancel {
		if token, id, err = readConfigFile(action.FileName); err != nil {
			panic(err)
		}

		if err = cancleAccount(token, id); err != nil {
			panic(err)
		}

		if err = os.Remove(action.FileName); err != nil {
			panic(err)
		}
		fmt.Println("Canceled")
		return
	}

	if action.License != "" {
		if token, id, err = readConfigFile(action.FileName); err != nil {
			panic(err)
		}

		if output, err = changeLicense(token, id, action.License); err != nil {
			panic(err)
		}

		if err = updateConfigFile(action.FileName); err != nil {
			panic(err)
		}

		fmt.Println(output)
		return
	}

	if strings.HasPrefix(action.Name, "-") {
		panic(`The parameter must not start with "-"`)
	}
	if action.Name != "" {

		if token, id, err = readConfigFile(action.FileName); err != nil {
			panic(err)
		}

		if output, err = changeName(token, id, action.Name); err != nil {
			panic(err)
		}

		if err = updateConfigFile(action.FileName); err != nil {
			panic(err)
		}
		fmt.Println(output)
		return
	}

	if strings.HasPrefix(action.Generate, "-") {
		panic(`The parameter must not start with "-"`)
	}
	if action.Generate != "" && action.Generate == "wg" {
		if config, reserved, err = configGenerate("wireguard", action.FileName); err != nil {
			panic(err)
		}

		if err = os.WriteFile(action.FileName+".wgcf.conf", []byte(config), 0600); err != nil {
			panic(err)
		}

		if output, err = nftConfigGenerate(reserved); err != nil {
			panic(err)
		}

		if err = os.WriteFile("wgcf.nft.conf", []byte(output), 0600); err != nil {
			panic(err)
		}
	} else if action.Generate != "" && action.Generate == "xray" {
		if config, _, err = configGenerate("xray", action.FileName); err != nil {
			panic(err)
		}

		if err = os.WriteFile(action.FileName+".xray.json", []byte(config), 0600); err != nil {
			panic(err)
		}
		return
	} else if action.Generate != "" && action.Generate == "sing-box" {
		if config, _, err = configGenerate("sing-box", action.FileName); err != nil {
			panic(err)
		}

		if err = os.WriteFile(action.FileName+".sing-box.json", []byte(config), 0600); err != nil {
			panic(err)
		}
		return
	}
	if action.Plus {
		if err = plus(action.FileName, 1); err != nil {
			panic(err)
		}
		return
	}
	if action.Update {
		if err = updateConfigFile(action.FileName); err != nil {
			panic(err)
		}
		println("Updated config file successfully")
		return
	}
}
