package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of wgcf-cli",
	Run:   printVersoin,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersoin(cmd *cobra.Command, args []string) {
	var revision string
	debugInfo, loaded := debug.ReadBuildInfo()

	if loaded {
		for _, setting := range debugInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			}
		}
	}
	if revision == "" {
		revision = "unknow"
	}
	fmt.Fprintf(os.Stderr, `wgcf-cli version `+"\033[1;36m"+C.Version+"\033[0m"+` 
`+"Revision: \033[1;35m"+revision+"\033[0m"+`
Environment: %s %s/%s
`, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
