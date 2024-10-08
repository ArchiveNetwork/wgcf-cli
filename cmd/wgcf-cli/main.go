package main

import (
	"os"

	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

const ConfigPathDefault string = "wgcf.json"

var rootCmd = &cobra.Command{
	Use:  os.Args[0],
	Long: "A command-line tool for Cloudflare-WARP API, built using Cobra.",
}

var (
	client     utils.HTTPClient
	configPath string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", ConfigPathDefault, "set configuration file path")
}
func main() {
	rootCmd.Execute()
}
