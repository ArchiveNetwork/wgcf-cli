package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
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

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	if action.Register {
		store, output, err := register()
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
		if action.FileName != "" {
			err = os.WriteFile(action.FileName, store, 0600)
			if err != nil {
				panic(err)
			}
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
		}
	}

	if action.FileName == "" {
		action.FileName = "wgcf.json"
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
		fmt.Println(output)
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
		fmt.Println(output)
	}

	if action.Cancle {
		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		err = cancleAccount(token, id)
		if err != nil {
			panic(err)
		}
		fmt.Println("Cancled")
	}

	if action.License != "" {

		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := changeLicense(token, id, action.License)
		if err != nil {
			panic(err)
		}
		fmt.Println(output)

		err = editFile(action.FileName, action.License)
		if err != nil {
			panic(err)
		}
	}

	if action.Name != "" {
		token, id, err := readConfigFile(action.FileName)
		if err != nil {
			panic(err)
		}

		output, err := changeName(token, id, action.Name)
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
	}
}

type Actions struct {
	Register bool
	Bind     bool
	UnBind   bool
	Cancle   bool
	License  string
	FileName string
	Name     string
}

func ParseCommandLine() Actions {
	var action Actions
	flag.BoolVar(&action.Register, "r", false, "")
	flag.BoolVar(&action.Register, "register", false, "")

	flag.BoolVar(&action.Bind, "b", false, "")
	flag.BoolVar(&action.Bind, "bind", false, "")

	flag.BoolVar(&action.UnBind, "u", false, "")
	flag.BoolVar(&action.UnBind, "unbind", false, "")

	flag.BoolVar(&action.Cancle, "c", false, "")
	flag.BoolVar(&action.Cancle, "cancle", false, "")

	flag.StringVar(&action.FileName, "f", "", "")
	flag.StringVar(&action.FileName, "file", "", "")

	flag.StringVar(&action.License, "l", "", "")
	flag.StringVar(&action.License, "license", "", "")

	flag.StringVar(&action.Name, "n", "", "")
	flag.StringVar(&action.Name, "name", "", "")

	flag.Visit(
		func(f *flag.Flag) {
			if f.Name == "r" || f.Name == "register" {
				action.Register = true
			}
			if f.Name == "b" || f.Name == "bind" {
				action.Bind = true
			}
			if f.Name == "u" || f.Name == "unbind" {
				action.UnBind = true
			}
			if f.Name == "c" || f.Name == "cancle" {
				action.Cancle = true
			}
			if f.Name == "f" || f.Name == "file" {
				action.FileName = f.Value.String()
			}
			if f.Name == "l" || f.Name == "license" {
				action.License = f.Value.String()
			}
			if f.Name == "n" || f.Name == "name" {
				action.Name = f.Value.String()
			}
		},
	)

	flag.Usage = func() {
		help()
	}

	flag.Parse()

	if action.License != "" {
		expectedPattern := `^[0-9A-Za-z]{8}-[0-9A-Za-z]{8}-[0-9A-Za-z]{8}$`

		match, err := regexp.MatchString(expectedPattern, action.License)
		if err != nil {
			panic(err)
		}

		if !match {
			panic("License should be something matchs: ^[0-9A-Za-z]{8}-[0-9A-Za-z]{8}-[0-9A-Za-z]{8}$")
		}
	}

	return action
}

func help() {
	fmt.Fprintf(os.Stderr,
		`Usage:	%s [Options]
Options:		-h/--help			help
			-f/--file [string]		Configuration file (default "wgcf.json")
			-r/--register			Register an account
			-b/--bind			Get the account binding devices
			-n/--name [string]		Change the device name
			-l/--license [string]		Change the license
			-u/--unbind			Unbind a device from the account
			-c/--cancle			Cancle the account
`, os.Args[0])
}
