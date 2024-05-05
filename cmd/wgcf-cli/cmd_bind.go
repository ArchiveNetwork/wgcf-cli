package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var bindCmd = &cobra.Command{
	Use:     "bind",
	Short:   "check current bind devices",
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
	var buffer bytes.Buffer
	var body []byte
	if body, err = client.Do(request); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	if err = json.Indent(&buffer, body, "", "    "); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Println(buffer.String())
}
