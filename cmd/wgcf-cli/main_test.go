package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
	"github.com/ArchiveNetwork/wgcf-cli/utils"
)

func TestRootCmd(t *testing.T) {
	var output bytes.Buffer
	rootCmd.SetOutput(&output)

	getLicense := func() string {
		var response C.Response
		body := utils.ReadConfig("wgcf.json")
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		return response.Account.License
	}
	must := func(err error) {
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
	}

	rootCmd.SetArgs([]string{"register"})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"simplify"})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"name", "-n", "test"})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"bind"})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"unbind", "--yes"})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"license", "-l", getLicense()})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"generate", "--wg"})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"generate", "--wg-quick"})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"generate", "--sing-box"})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"generate", "--xray"})
	must(rootCmd.Execute())

	rootCmd.SetArgs([]string{"cancel", "--yes"})
	must(rootCmd.Execute())
}
