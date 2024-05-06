package utils

import (
	"encoding/json"
	"fmt"
	"os"

	C "github.com/ArchiveNetwork/wgcf-cli/constant"
)

func GenXray(resStruct C.Response) (body []byte, err error) {
	in_struct := C.Xray{
		Protocol: "wireguard",
		Settings: struct {
			SecretKey string   `json:"secretKey"`
			Address   []string `json:"address"`
			Peers     []struct {
				PublicKey  string   `json:"publicKey"`
				AllowedIPs []string `json:"allowedIPs"`
				Endpoint   string   `json:"endpoint"`
			} `json:"peers"`
			Reserved []int `json:"reserved"`
			MTU      int   `json:"mtu"`
		}{
			SecretKey: resStruct.Config.PrivateKey,
			Address:   []string{resStruct.Config.Interface.Addresses.V4 + "/32", resStruct.Config.Interface.Addresses.V6 + "/128"},
			Peers: []struct {
				PublicKey  string   `json:"publicKey"`
				AllowedIPs []string `json:"allowedIPs"`
				Endpoint   string   `json:"endpoint"`
			}{
				{
					PublicKey:  "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=",
					AllowedIPs: []string{"0.0.0.0/0", "::/0"},
					Endpoint:   resStruct.Config.Peers[0].Endpoint.Host,
				},
			},
			Reserved: resStruct.Config.ReservedDec,
			MTU:      1280,
		},
		Tag: "wireguard",
	}

	if body, err = json.MarshalIndent(in_struct, "", "    "); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	return
}
func GenSing(resStruct C.Response) (body []byte, err error) {
	in_struct := C.Sing{
		Type:          "wireguard",
		Tag:           "wireguard-out",
		Server:        resStruct.Config.Peers[0].Endpoint.Host,
		ServerPort:    2408,
		LocalAddress:  []string{resStruct.Config.Interface.Addresses.V4 + "/32", resStruct.Config.Interface.Addresses.V6 + "/128"},
		PrivateKey:    resStruct.Config.PrivateKey,
		PeerPublicKey: "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=",
		Reserved:      resStruct.Config.ClientID,
		MTU:           1280,
	}

	if body, err = json.MarshalIndent(in_struct, "", "    "); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	return
}
