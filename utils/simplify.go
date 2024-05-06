package utils

import (
	"encoding/json"
	"fmt"
	"os"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
)

func SimplifyOutput(resStruct C.Response) {
	jsonStruct := C.SimpleOutput{
		Endpoint: struct {
			V4 string `json:"v4"`
			V6 string `json:"v6"`
		}{
			V4: resStruct.Config.Peers[0].Endpoint.V4,
			V6: resStruct.Config.Peers[0].Endpoint.V6,
		},
		ReservedStr: resStruct.Config.ClientID,
		ReservedHex: resStruct.Config.ReservedHex,
		ReservedDec: resStruct.Config.ReservedDec,
		PrivateKey:  resStruct.Config.PrivateKey,
		PublicKey:   resStruct.Config.Peers[0].PublicKey,
		Addresses:   resStruct.Config.Interface.Addresses,
	}

	output, err := json.MarshalIndent(jsonStruct, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, string(output))
}
