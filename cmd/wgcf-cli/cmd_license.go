package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:     "license",
	Short:   "change to a new license",
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
	body, err := client.Do(requset)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	var buffer bytes.Buffer
	if err = json.Indent(&buffer, body, "", "    "); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Printf("License changed to %s\n", license)
}
