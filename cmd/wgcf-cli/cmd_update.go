package main

import (
	"encoding/json"
	"fmt"
	"os"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a config",
	PreRun: func(cmd *cobra.Command, args []string) {
		client.New()
	},
	Run: update,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update(cmd *cobra.Command, args []string) {
	var resStruct, response C.Response
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	body := utils.ReadConfig(configPath)
	token, id := utils.GetTokenID(configPath)
	var updatedContent []byte

	if err := json.Unmarshal(body, &resStruct); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	if resStruct.Config.ReservedDec == nil || resStruct.Config.ReservedHex == "" {
		response.Config.ReservedDec, response.Config.ReservedHex = utils.ClientIDtoReserved(resStruct.Config.ClientID)
	} else {
		response.Config.ReservedDec = resStruct.Config.ReservedDec
		response.Config.ReservedHex = resStruct.Config.ReservedHex
	}
	r := utils.Request{
		Action: "update",
		Token:  token,
		ID:     id,
	}
	request, err := r.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if body, err = client.Do(request); err != nil {
		client.HandleBody()
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	if err = json.Unmarshal(body, &response); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	if resStruct.Account.PrivateKey != "" {
		response.Config.PrivateKey = resStruct.Account.PrivateKey
	} else {
		response.Config.PrivateKey = resStruct.Config.PrivateKey
	}
	response.Token = resStruct.Token
	if updatedContent, err = json.MarshalIndent(response, "", "    "); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if err = os.WriteFile(configPath, updatedContent, 0600); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Updated configuration file (ID: %s) successfully\n", id)
}
