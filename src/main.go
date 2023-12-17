package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Response struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Model   string `json:"model"`
	Name    string `json:"name"`
	Key     string `json:"key"`
	Account struct {
		ID                   string `json:"id"`
		PrivateKey           string `json:"private_key"`
		ReservedHex          string `json:"reserved_hex"`
		ReservedDec          []int  `json:"reserved_dec"`
		AccountType          string `json:"account_type"`
		Created              string `json:"created"`
		Updated              string `json:"updated"`
		PremiumData          int    `json:"premium_data"`
		Quota                int    `json:"quota"`
		Usage                int    `json:"usage"`
		WarpPlus             bool   `json:"warp_plus"`
		ReferralCount        int    `json:"referral_count"`
		ReferralRenewalCount int    `json:"referral_renewal_countdown"`
		Role                 string `json:"role"`
		License              string `json:"license"`
	} `json:"account"`
	Config struct {
		ClientID string `json:"client_id"`
		Peers    []struct {
			PublicKey string `json:"public_key"`
			Endpoint  struct {
				V4   string `json:"v4"`
				V6   string `json:"v6"`
				Host string `json:"host"`
			} `json:"endpoint"`
		} `json:"peers"`
		Interface struct {
			Addresses struct {
				V4 string `json:"v4"`
				V6 string `json:"v6"`
			} `json:"addresses"`
		} `json:"interface"`
		Services struct {
			HTTPProxy string `json:"http_proxy"`
		} `json:"services"`
	} `json:"config"`
	Token     string `json:"token"`
	Warp      bool   `json:"warp_enabled"`
	Waitlist  bool   `json:"waitlist_enabled"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
	TOS       string `json:"tos"`
	Place     int    `json:"place"`
	Locale    string `json:"locale"`
	Enabled   bool   `json:"enabled"`
	InstallID string `json:"install_id"`
	FCMToken  string `json:"fcm_token"`
	SerialNum string `json:"serial_number"`
}

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

	if action.Register {

		if store, output, err = register(); err != nil {
			panic(err)
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

	if action.FileName == "" {
		action.FileName = "wgcf.json"
	} else if strings.HasPrefix(action.FileName, "-") {
		err := fmt.Sprintln("The parameter must not start with '-'")
		panic(err)
	}

	if !action.Bind && !action.UnBind && !action.Cancle && action.License == "" && action.Name == "" && action.Generate == "" {
		err := fmt.Sprintln("You need to specify an action")
		panic(err)
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

	if action.Cancle {
		if token, id, err = readConfigFile(action.FileName); err != nil {
			panic(err)
		}

		if err = cancleAccount(token, id); err != nil {
			panic(err)
		}

		if err = os.Remove(action.FileName); err != nil {
			panic(err)
		}
		fmt.Println("Cancled")
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
		err := fmt.Sprintln("The parameter must not start with '-'")
		panic(err)
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
		err := fmt.Sprintln("The parameter must not start with '-'")
		panic(err)
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

}
