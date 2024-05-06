package main

import (
	"fmt"
	"os"

	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "cancel a config from original one",
	Run:   cancel,
}

func init() {
	rootCmd.AddCommand(cancelCmd)
	cancelCmd.PersistentFlags().Bool("yes", false, "confirm that you want to cancel the account")
	cancelCmd.MarkPersistentFlagRequired("yes")
}

func cancel(cmd *cobra.Command, args []string) {
	token, id := utils.GetTokenID(configPath)

	r := utils.Request{
		Action: "cancel",
		ID:     id,
		Token:  token,
	}
	requset, err := r.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	var client utils.HTTPClient
	if _, err = client.Do(requset); err != nil {
		client.HandleBody()
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	client.HandleBody()

	if err = os.Remove(configPath); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Printf("Canceled account (ID: %s) successfully\n", id)
}
