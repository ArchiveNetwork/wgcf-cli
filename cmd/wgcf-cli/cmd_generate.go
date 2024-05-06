package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:       "generate",
	Short:     "generate a config from original one",
	Run:       generate,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"xray", "sing-box"},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().Bool("xray", false, "generate a xray config")
	generateCmd.Flags().Bool("sing-box", false, "generate a sing-box config")
}
func generate(cmd *cobra.Command, args []string) {
	var err error
	xray, _ := cmd.Flags().GetBool("xray")
	sing, _ := cmd.Flags().GetBool("sing-box")

	var resStruct C.Response
	body := utils.ReadConfig(configPath)
	if err := json.Unmarshal(body, &resStruct); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	var gen_type string
	if xray {
		body, err = utils.GenXray(resStruct)
		gen_type = "xray.json"
	} else if sing {
		body, err = utils.GenSing(resStruct)
		gen_type = "sing-box.json"
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if err = os.WriteFile(strings.TrimSuffix(configPath, "json")+gen_type, body, 0600); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
