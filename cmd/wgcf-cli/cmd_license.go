package main

import (
	"fmt"
	"os"

	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:     "license",
	Short:   "Change to a new license",
	Run:     change_license,
	PostRun: update,
}

var license string

func init() {
	rootCmd.AddCommand(licenseCmd)
	licenseCmd.PersistentFlags().StringVarP(&license, "license", "l", "", "set the license that change to")
	licenseCmd.MarkPersistentFlagRequired("license")
}

func change_license(cmd *cobra.Command, args []string) {
	token, id := utils.GetTokenID(configPath)

	r := utils.Request{
		Action: "license",
		Payload: []byte(
			`{
				"license":"` + license + `"
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
	fmt.Printf("License changed to %s\n", license)
}
