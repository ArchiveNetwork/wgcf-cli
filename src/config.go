package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Xray struct {
	Protocol string `json:"protocol"`
	Settings struct {
		SecretKey string   `json:"secretKey"`
		Address   []string `json:"address"`
		Peers     []struct {
			PublicKey  string   `json:"publicKey"`
			AllowedIPs []string `json:"allowedIPs"`
			Endpoint   string   `json:"endpoint"`
		} `json:"peers"`
		Reserved []int `json:"reserved"`
		MTU      int   `json:"mtu"`
	} `json:"settings"`
	Tag string `json:"tag"`
}

type singBox struct {
	Type          string   `json:"type"`
	Tag           string   `json:"tag"`
	Server        string   `json:"server"`
	ServerPort    int      `json:"server_port"`
	LocalAddress  []string `json:"local_address"`
	PrivateKey    string   `json:"private_key"`
	PeerPublicKey string   `json:"peer_public_key"`
	Reserved      string   `json:"reserved"`
	MTU           int      `json:"mtu"`
}

type Wireguard struct {
	Reserved  string `conf:"Reserved,omitempty"`
	Interface struct {
		PrivateKey string   `conf:"PrivateKey"`
		Address    []string `conf:"Address"`
		MTU        int      `conf:"MTU"`
		PostUp     string   `conf:"PostUp"`
		PostDown   []string `conf:"PostDown"`
		Table      int      `conf:"Table"`
		PreUp      []string `conf:"PreUp"`
	} `conf:"Interface"`
	Peer struct {
		PublicKey  string `conf:"PublicKey"`
		AllowedIPs string `conf:"AllowedIPs"`
		Endpoint   string `conf:"Endpoint"`
	} `conf:"Peer"`
}

func configGenerate(generateType string, filePath string) (string, string, error) {
	var err error
	var content []byte
	var config []byte
	var file *os.File
	if file, err = os.Open(filePath); err != nil {
		panic(err)
	}
	defer file.Close()

	if content, err = io.ReadAll(file); err != nil {
		panic(err)
	}

	var ReadedFile NormalResponse

	if err = json.Unmarshal(content, &ReadedFile); err != nil {
		panic(err)
	}

	if generateType == "xray" {
		input := Xray{
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
				SecretKey: ReadedFile.Config.PrivateKey,
				Address:   []string{"172.16.0.2/32", ReadedFile.Config.Interface.Addresses.V6 + "/128"},
				Peers: []struct {
					PublicKey  string   `json:"publicKey"`
					AllowedIPs []string `json:"allowedIPs"`
					Endpoint   string   `json:"endpoint"`
				}{
					{
						PublicKey:  "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=",
						AllowedIPs: []string{"0.0.0.0/0", "::/0"},
						Endpoint:   "162.159.192.1:2408",
					},
				},
				Reserved: ReadedFile.Config.ReservedDec,
				MTU:      1280,
			},
			Tag: "wireguard",
		}

		if config, err = json.MarshalIndent(input, "", "    "); err != nil {
			panic(err)
		}
		return fmt.Sprintf((string(config) + "\n")), "", nil

	} else if generateType == "wireguard" {
		input := Wireguard{
			Reserved: ReadedFile.Config.ReservedHex,
			Interface: struct {
				PrivateKey string   `conf:"PrivateKey"`
				Address    []string `conf:"Address"`
				MTU        int      `conf:"MTU"`
				PostUp     string   `conf:"PostUp"`
				PostDown   []string `conf:"PostDown"`
				Table      int      `conf:"Table"`
				PreUp      []string `conf:"PreUp"`
			}{
				PrivateKey: ReadedFile.Config.PrivateKey,
				Address:    []string{ReadedFile.Config.Interface.Addresses.V6 + "/128", "172.16.0.2/32"},
				MTU:        1280,
				PostUp:     "nft -f /etc/wireguard/wgcf.nft.conf",
				PostDown:   []string{"nft delete table inet wgcf", "ip rule del oif %i lookup 300", "ip -6 rule del oif %i lookup 300"},
				Table:      300,
				PreUp:      []string{"ip rule add oif %i lookup 300", "ip -6 rule add oif %i lookup 300"},
			},
			Peer: struct {
				PublicKey  string `conf:"PublicKey"`
				AllowedIPs string `conf:"AllowedIPs"`
				Endpoint   string `conf:"Endpoint"`
			}{
				PublicKey:  "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=",
				AllowedIPs: "::/0, 0.0.0.0/0",
				Endpoint:   "162.159.193.199:2408",
			},
		}
		config := fmt.Sprintf(`
[Interface]
PrivateKey = %s
Address = %s
Address = %s
MTU = %d

Table = %d

PreUp = %s
PostDown = %s
PreUp = %s
PostDown = %s
PostUp = %s
PostDown = %s
	
[Peer]
PublicKey = %s
AllowedIPs = %v
Endpoint = %s`,
			input.Interface.PrivateKey,
			input.Interface.Address[0],
			input.Interface.Address[1],
			input.Interface.MTU,
			input.Interface.Table,
			input.Interface.PreUp[0],
			input.Interface.PostDown[1],
			input.Interface.PreUp[1],
			input.Interface.PostDown[2],
			input.Interface.PostUp,
			input.Interface.PostDown[0],
			input.Peer.PublicKey,
			input.Peer.AllowedIPs,
			input.Peer.Endpoint,
		)
		return config, input.Reserved, nil
	} else if generateType == "sing-box" {
		input := singBox{
			Type:          "wireguard",
			Tag:           "wireguard-out",
			Server:        "engage.cloudflareclient.com",
			ServerPort:    2408,
			LocalAddress:  []string{ReadedFile.Config.Interface.Addresses.V4 + "/32", ReadedFile.Config.Interface.Addresses.V6 + "/128"},
			PrivateKey:    ReadedFile.Config.PrivateKey,
			PeerPublicKey: "bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo=",
			Reserved:      ReadedFile.Config.ClientID,
			MTU:           1280,
		}

		if config, err = json.MarshalIndent(input, "", "    "); err != nil {
			panic(err)
		}
		return fmt.Sprintf((string(config) + "\n")), "", nil
	}
	panic("unsupported generateType")
}

func nftConfigGenerate(reserved string) (string, error) {
	config := fmt.Sprintf(`
#!/usr/bin/nft -f
# vim:set ts=4 sw=4 et:

define routing_id = ` + reserved + `
	
table inet wgcf
delete table inet wgcf
table inet wgcf {
	chain output {
		type filter hook output priority mangle - 1; policy accept;
		ip daddr 162.159.193.199 udp dport 2408 @th,72,24 set $routing_id;
	}
	chain input {
		type filter hook input priority mangle; policy accept;
		ip saddr 162.159.193.199 udp sport 2408 @th,72,24 $routing_id @th,72,24 set 0;
	}
}
`)
	return config, nil
}
