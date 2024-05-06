package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new WARP account",
	Run:   register,
}

var (
	teamToken string
)

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.PersistentFlags().StringVarP(&teamToken, "token", "t", "", "set register ZeroTrust Token")
}

func register(cmd *cobra.Command, args []string) {
	privateKey, publicKey, err := utils.GenerateKey()
	if err != nil {
		panic(err)
	}

	installID := utils.RandStringRunes(22, nil)
	fcmtoken := utils.RandStringRunes(134, nil)

	r := utils.Request{
		Payload: []byte(
			`{
				"key":"` + publicKey + `",
				"install_id":"` + installID + `",
				"fcm_token":"` + installID + `:APA91b` + fcmtoken + `",
				"tos":"` + time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + `",
				"model":"Android",
				"serial_number":"` + installID + `"
			}`,
		),
		Action:    "register",
		TeamToken: teamToken,
	}

	request, err := r.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	var client utils.HTTPClient
	body, err := client.Do(request)
	if err != nil {
		client.HandleBody()
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	var resStruct C.Response
	if err = json.Unmarshal(body, &resStruct); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	resStruct.Config.ReservedDec, resStruct.Config.ReservedHex = utils.ClientIDtoReserved(resStruct.Config.ClientID)
	resStruct.Config.PrivateKey = privateKey
	utils.SimplifyOutput(resStruct)

	store, err := json.MarshalIndent(resStruct, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	if err = os.WriteFile(configPath, store, 0600); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

}
