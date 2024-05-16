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
	Short:     "Generate a xray/sing-box config",
	Run:       generate,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"--xray", "--sing-box", "--wg"},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().Bool("xray", false, "generate a xray config")
	generateCmd.Flags().Bool("sing-box", false, "generate a sing-box config")
	generateCmd.Flags().Bool("wg", false, "generate a wg-quick config")
}

func generate(cmd *cobra.Command, args []string) {
	var err error
	xray, _ := cmd.Flags().GetBool("xray")
	sing, _ := cmd.Flags().GetBool("sing-box")
	wg, _ := cmd.Flags().GetBool("wg")

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
	} else if wg {
		body, err = utils.GenWgQuick(resStruct)
		gen_type = "ini"
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	store := strings.TrimSuffix(configPath, "json") + gen_type
	if _, err := os.Stat(store); !os.IsNotExist(err) {
		var input string
		fmt.Fprintf(os.Stderr, "Warn: File %s exist, are you sure to continue? [y/N]: ", store)
		fmt.Scanln(&input)
		input = strings.ToLower(input)
		if input != "y" {
			os.Exit(1)
		}
	}
	if err = os.WriteFile(store, body, 0600); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	if strings.HasSuffix(gen_type, "ini") {
		gen_type = "wg-quick"
	} else {
		gen_type = strings.TrimSuffix(gen_type, ".json")
	}
	fmt.Printf("Generate configuration file (ID: %s) for %s successfully\n", resStruct.ID, gen_type)
}
