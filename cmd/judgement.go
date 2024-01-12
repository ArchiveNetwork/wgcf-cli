package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ArchiveNetwork/wgcf-cli/utils"
)

func Judgement() error {
	action := ParseCommandLine()
	var err error
	var store []byte
	var config, token, id, reserved, output string

	if len(os.Args) == 1 {
		flag.Usage()
		return nil
	}

	if action.Help {
		Help()
		return nil
	}
	if action.Version {
		Version()
		return nil
	}

	if action.Register {
		if strings.HasPrefix(action.TeamToken, "-") {
			err := fmt.Sprintln(`The parameter must not start with "-"`)
			panic(err)
		}
		if action.TeamToken != "" {
			if store, output, err = utils.Register(action.TeamToken); err != nil {
				panic(err)
			}
		} else {
			if store, output, err = utils.Register(""); err != nil {
				panic(err)
			}
		}
		fmt.Println(output)
		if action.FileName != "" {
			if err = os.WriteFile(action.FileName, store, 0600); err != nil {
				panic(err)
			}
			return nil
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
			return nil
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
		token, id, err := utils.GetTokenID(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := utils.GetBindingDevices(token, id)
		if err != nil {
			panic(err)
		}

		if err = utils.UpdateConfigFile(action.FileName); err != nil {
			panic(err)
		}
		fmt.Println(output)

		return nil
	}

	if action.UnBind {
		token, id, err := utils.GetTokenID(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := utils.UnBind(token, id)
		if err != nil {
			panic(err)
		}

		if err = utils.UpdateConfigFile(action.FileName); err != nil {
			panic(err)
		}

		fmt.Println(output)
		return nil
	}

	if action.Cancel {
		if token, id, err = utils.GetTokenID(action.FileName); err != nil {
			panic(err)
		}

		if err = utils.CancleAccount(token, id); err != nil {
			panic(err)
		}

		if err = os.Remove(action.FileName); err != nil {
			panic(err)
		}
		fmt.Println("Canceled")
		return nil
	}

	if action.License != "" {
		if token, id, err = utils.GetTokenID(action.FileName); err != nil {
			panic(err)
		}

		if output, err = utils.ChangeLicense(token, id, action.License); err != nil {
			panic(err)
		}

		if err = utils.UpdateConfigFile(action.FileName); err != nil {
			panic(err)
		}

		fmt.Println(output)
		return nil
	}

	if strings.HasPrefix(action.Name, "-") {
		panic(`The parameter must not start with "-"`)
	}
	if action.Name != "" {

		if token, id, err = utils.GetTokenID(action.FileName); err != nil {
			panic(err)
		}

		if output, err = utils.ChangeName(token, id, action.Name); err != nil {
			panic(err)
		}

		if err = utils.UpdateConfigFile(action.FileName); err != nil {
			panic(err)
		}
		fmt.Println(output)
		return nil
	}

	if strings.HasPrefix(action.Generate, "-") {
		panic(`The parameter must not start with "-"`)
	}
	if action.Generate != "" && action.Generate == "wg" {
		if config, reserved, err = utils.ConfigGenerate("wireguard", action.FileName); err != nil {
			panic(err)
		}

		if err = os.WriteFile(action.FileName+".wgcf.conf", []byte(config), 0600); err != nil {
			panic(err)
		}

		if output, err = utils.NftConfigGenerate(reserved); err != nil {
			panic(err)
		}

		if err = os.WriteFile("wgcf.nft.conf", []byte(output), 0600); err != nil {
			panic(err)
		}
	} else if action.Generate != "" && action.Generate == "xray" {
		if config, _, err = utils.ConfigGenerate("xray", action.FileName); err != nil {
			panic(err)
		}

		if err = os.WriteFile(action.FileName+".xray.json", []byte(config), 0600); err != nil {
			panic(err)
		}
		return nil
	} else if action.Generate != "" && action.Generate == "sing-box" {
		if config, _, err = utils.ConfigGenerate("sing-box", action.FileName); err != nil {
			panic(err)
		}

		if err = os.WriteFile(action.FileName+".sing-box.json", []byte(config), 0600); err != nil {
			panic(err)
		}
		return nil
	}
	if action.Plus {
		if err = utils.Plus(action.FileName, false); err != nil {
			panic(err)
		}
		return nil
	}
	if action.Update {
		if err = utils.UpdateConfigFile(action.FileName); err != nil {
			panic(err)
		}
		println("Updated config file successfully")
		return nil
	}
	return nil
}
