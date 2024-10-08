package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
	"github.com/ArchiveNetwork/wgcf-cli/utils"
)

func expectNoErr(err error, t *testing.T) {
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

type checker func(err error)

func createConfig(check checker) {
	if _, err := os.Open(ConfigPathDefault); errors.Is(err, os.ErrNotExist) {
		rootCmd.SetArgs([]string{"register"})
		check(rootCmd.Execute())
	}
}
func removeConfig(check checker) {
	rootCmd.SetArgs([]string{"cancel", "--yes"})
	check(rootCmd.Execute())
}

func beginGenerateTest(check checker) {
	generateCmd.Flags().Set(asString(Xray), "false")
	generateCmd.Flags().Set("wg", "false")
	generateCmd.Flags().Set(asString(WgQuick), "false")
	generateCmd.Flags().Set(asString(SingBox), "false")
	createConfig(check)
}
func endGenerateTest(check checker, generator GeneratorType) {
	removeConfig(check)
	os.Remove(getDefaultFilePath(generator))
}
func runGenerateTest(check checker, generator GeneratorType, test func()) {
	beginGenerateTest(check)
	test()
	endGenerateTest(check, generator)
}
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

	check := func(err error) { expectNoErr(err, t) }

	rootCmd.SetArgs([]string{"register"})
	check(rootCmd.Execute())

	rootCmd.SetArgs([]string{"simplify"})
	check(rootCmd.Execute())

	rootCmd.SetArgs([]string{"name", "-n", "test"})
	check(rootCmd.Execute())

	rootCmd.SetArgs([]string{"bind"})
	check(rootCmd.Execute())

	rootCmd.SetArgs([]string{"unbind", "--yes"})
	check(rootCmd.Execute())

	rootCmd.SetArgs([]string{"license", "-l", getLicense()})
	check(rootCmd.Execute())

	rootCmd.SetArgs([]string{"cancel", "--yes"})
	check(rootCmd.Execute())
}

func TestGenerateWg(t *testing.T) {

	check := func(err error) { expectNoErr(err, t) }

	runGenerateTest(check, WgQuick, func() {
		rootCmd.SetArgs([]string{"generate", "--wg"})
		check(rootCmd.Execute())

		os.Remove(getDefaultFilePath(WgQuick))
		rootCmd.SetArgs([]string{"generate", "--wg-quick"})
		check(rootCmd.Execute())
	})
}

func TestGenerateSingBox(t *testing.T) {
	check := func(err error) { expectNoErr(err, t) }

	runGenerateTest(check, WgQuick, func() {
		rootCmd.SetArgs([]string{"generate", "--sing-box"})
		check(rootCmd.Execute())
	})
}

func TestGenerateXray(t *testing.T) {
	check := func(err error) { expectNoErr(err, t) }

	runGenerateTest(check, WgQuick, func() {
		rootCmd.SetArgs([]string{"generate", "--xray"})
		check(rootCmd.Execute())
	})
}
