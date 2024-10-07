package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/bits"
	"os"
	"path"
	"strings"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
	"github.com/ArchiveNetwork/wgcf-cli/utils"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:       "generate",
	Short:     "Generate a xray/sing-box/wg-quick config",
	Run:       generate,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"--xray", "--sing-box", "--wg", "--wg-quick", "--output-file"},
}

type OutputFileType int8

const (
	Stdout OutputFileType = iota
	Default
	Custom
)

type GeneratorType int8

const (
	Xray GeneratorType = iota
	SingBox
	WgQuick
	None
)

func (t GeneratorType) String() string {
	switch t {
	case Xray:
		return "xray"
	case SingBox:
		return "sing-box"
	case WgQuick:
		return "wg-quick"
	}
	return "unknown"
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().Bool(asString(Xray), false, "generate a xray config")
	generateCmd.Flags().Bool(asString(SingBox), false, "generate a sing-box config")
	generateCmd.Flags().Bool(asString(WgQuick), false, "generate a wg-quick config")
	generateCmd.Flags().Bool("wg", false, "see --"+asString(WgQuick))
	generateCmd.Flags().String("output-file", "default", "output file name. Supported values: 'default'/'stdout'/any file path")
}

func asString[V fmt.Stringer](object V) string {
	return V.String(object)
}

func detectOutputFileType(cmd *cobra.Command) (OutputFileType, error) {
	var err error
	path, err := cmd.Flags().GetString("output-file")
	if err != nil {
		return Stdout, err
	}
	switch path {
	case "stdout":
		return Stdout, nil
	case "default":
		return Default, nil
	}
	return Custom, nil
}

func ternary[V any](condition bool, on_true V, on_false V) V {
	if condition {
		return on_true
	}
	return on_false
}

func detectGeneratorType(cmd *cobra.Command) (GeneratorType, error) {
	xray, _ := cmd.Flags().GetBool(asString(Xray))
	sing, _ := cmd.Flags().GetBool(asString(SingBox))
	wg, _ := cmd.Flags().GetBool(asString(WgQuick))
	if !wg {
		wg, _ = cmd.Flags().GetBool("wg")
	}

	var options uint8 = 0
	options |= ternary(xray, uint8(0b001), 0)
	options |= ternary(sing, uint8(0b010), 0)
	options |= ternary(wg, uint8(0b100), 0)
	if c := bits.OnesCount8(options); c != 1 {
		if c == 0 {
			return None, errors.New("generator not specified")
		} else {
			return None, errors.New("multiple generators not supported")
		}
	}

	if xray {
		return Xray, nil
	} else if sing {
		return SingBox, nil
	} else if wg {
		return WgQuick, nil
	}
	return None, nil
}

func askOutputOverwrite(path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var input string
		fmt.Fprintf(os.Stderr, "Warn: File %s exist, it will be overwritten. Continue? [y/N]: ", path)
		fmt.Scanln(&input)
		input = strings.ToLower(input)
		if input != "y" {
			os.Exit(1)
		}
	}
}

func getDefaultFilePath(generator GeneratorType) string {
	var base_name = strings.TrimSuffix(configPath, path.Ext(configPath))
	switch generator {
	case Xray:
		return base_name + ".xray.json"
	case SingBox:
		return base_name + ".sing-box.json"
	case WgQuick:
		return base_name + ".ini"
	}
	return ""
}

func Exit(err error, exit_code int) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(exit_code)
}
func ExitDefault(err error) {
	Exit(err, 1)
}

func generate(cmd *cobra.Command, args []string) {
	var err error
	var generator GeneratorType
	var output_type OutputFileType

	output_type, err = detectOutputFileType(cmd)
	if err != nil {
		ExitDefault(err)
	}
	generator, err = detectGeneratorType(cmd)
	if err != nil {
		ExitDefault(err)
	}

	var resStruct C.Response
	body := utils.ReadConfig(configPath)
	err = json.Unmarshal(body, &resStruct)
	if err != nil {
		ExitDefault(err)
	}

	switch generator {
	case Xray:
		body, err = utils.GenXray(resStruct)
	case SingBox:
		body, err = utils.GenSing(resStruct)
	case WgQuick:
		body, err = utils.GenWgQuick(resStruct)
	}
	if err != nil {
		ExitDefault(err)
	}

	switch output_type {
	case Stdout:
		_, err = fmt.Print(body)
		if err != nil {
			ExitDefault(err)
		}
	case Default:
		var filepath = getDefaultFilePath(generator)
		askOutputOverwrite(filepath)
		err = os.WriteFile(filepath, body, 0600)
		if err != nil {
			ExitDefault(err)
		}
		fmt.Printf("Generate %s configuration file '%s' (ID: %s) successfully\n", asString(generator), filepath, resStruct.ID)
	case Custom:
		filepath, _ := cmd.Flags().GetString("output-file")
		askOutputOverwrite(filepath)
		err = os.WriteFile(filepath, body, 0600)
		if err != nil {
			ExitDefault(err)
		}
		fmt.Printf("Generate %s configuration file '%s' (ID: %s) successfully\n", asString(generator), filepath, resStruct.ID)
	}
}
