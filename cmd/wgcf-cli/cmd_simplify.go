package main

import (
	"encoding/json"
	"fmt"
	"os"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var simplifyCmd = &cobra.Command{
	Use:   "simplify",
	Short: "Simplify a config from original one",
	Run:   simplify,
}

func init() {
	rootCmd.AddCommand(simplifyCmd)
}

func simplify(cmd *cobra.Command, args []string) {
	var resStruct C.Response
	body := utils.ReadConfig(configPath)
	if err := json.Unmarshal(body, &resStruct); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	utils.SimplifyOutput(resStruct)
}
