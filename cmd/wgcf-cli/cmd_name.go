package main

import (
	"fmt"
	"os"

	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var nameCmd = &cobra.Command{
	Use:     "name",
	Short:   "change the device name",
	Run:     change_name,
	PostRun: update,
}

var name string

func init() {
	rootCmd.AddCommand(nameCmd)
	nameCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "change the device name to")
	nameCmd.MarkPersistentFlagRequired("name")
}

func change_name(cmd *cobra.Command, args []string) {
	token, id := utils.GetTokenID(configPath)

	r := utils.Request{
		Action: "name",
		Payload: []byte(
			`{
				"name":"` + name + `"
			 }`,
		),
		ID:    id,
		Token: token,
	}
	requset, err := r.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	var client utils.HTTPClient

	if _, err := client.Do(requset); err != nil {
		client.HandleBody()
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	client.HandleBody()
	fmt.Printf("Name changed to %s\n", name)
}
