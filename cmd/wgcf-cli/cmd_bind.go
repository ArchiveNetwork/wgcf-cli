package main

import (
	"fmt"
	"os"

	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var bindCmd = &cobra.Command{
	Use:     "bind",
	Short:   "Check current bind devices",
	Run:     bind,
	PostRun: update,
}

func init() {
	rootCmd.AddCommand(bindCmd)
}

func bind(cmd *cobra.Command, args []string) {
	token, id := utils.GetTokenID(configPath)
	r := utils.Request{
		Action: "bind",
		Token:  token,
		ID:     id,
	}
	request, err := r.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	var client utils.HTTPClient
	if _, err := client.Do(request); err != nil {
		client.HandleBody()
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	client.HandleBody()
}
