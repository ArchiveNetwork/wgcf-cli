package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-ini/ini"
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

func ConfigGenerate(generateType string, filePath string) (string, string, error) {
	var err error
	var content []byte
	var config []byte
	var fileType string
	var ReadedFile Response
	if fileType, err = GetFileType(filePath); err != nil {
		panic(err)
	}
	if fileType == "json" {
		content = ReadConfig(filePath)
		if err = json.Unmarshal(content, &ReadedFile); err != nil {
			panic(err)
		}
	} else {
		var cfg *ini.File
		if cfg, err = ini.Load(filePath); err != nil {
			panic(err)
		}
		section_Config := cfg.Section("Config")
		ReadedFile.Config.PrivateKey = section_Config.Key("PrivateKey").String()
		ReadedFile.Config.Interface.Addresses.V4 = section_Config.Key("IPv4").String()
		ReadedFile.Config.Interface.Addresses.V6 = section_Config.Key("IPv6").String()
		ReadedFile.Config.ClientID = section_Config.Key("ClientID").String()
		ReadedFile.Config.ReservedHex = section_Config.Key("ReservedHex").String()
		var intSlice []int
		re := regexp.MustCompile(`\d+`)
		matches := re.FindAllString(section_Config.Key("ReservedDec").String(), -1)
		for _, str := range matches {
			var i int
			if i, err = strconv.Atoi(str); err != nil {
				panic(err)
			}
			intSlice = append(intSlice, i)
		}
		ReadedFile.Config.ReservedDec = intSlice
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
				Address:   []string{ReadedFile.Config.Interface.Addresses.V4 + "/32", ReadedFile.Config.Interface.Addresses.V6 + "/128"},
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
		config := fmt.Sprintf(`
[Interface]
PrivateKey = `+ReadedFile.Config.PrivateKey+`
Address = `+ReadedFile.Config.Interface.Addresses.V4+`/32
Address = `+ReadedFile.Config.Interface.Addresses.V6+`/128
MTU = 1280

PreUp = ip rule add oif %s lookup 300
PostDown = ip rule del oif %s lookup 300
PreUp = ip -6 rule add oif %s lookup 300
PostDown = ip -6 rule del oif %s lookup 300

PreUp = ip rule add fwmark 32975 lookup 300
PostDown = ip rule del fwmark 32975 lookup 300
PreUp = ip -6 rule add fwmark 32975 lookup 300
PostDown = ip -6 rule del fwmark 32975 lookup 300

#PreUp = ip rule add from `+ReadedFile.Config.Interface.Addresses.V4+`/32 lookup 300
#PostDown = ip rule del from `+ReadedFile.Config.Interface.Addresses.V4+`/32 lookup 300
#PreUp = ip -6 rule add from `+ReadedFile.Config.Interface.Addresses.V6+`/128 lookup 300
#PostDown = ip -6 rule del from `+ReadedFile.Config.Interface.Addresses.V6+`/128 lookup 300
# Alternative

PostUp = iptables -t mangle -A OUTPUT -s `+ReadedFile.Config.Interface.Addresses.V4+` -j MARK --set-mark 32975
PreDown = iptables -t mangle -D OUTPUT -s `+ReadedFile.Config.Interface.Addresses.V4+` -j MARK --set-mark 32975
PostUp = ip6tables -t mangle -A OUTPUT -s `+ReadedFile.Config.Interface.Addresses.V6+` -j MARK --set-mark 32975
PreDown = ip6tables -t mangle -D OUTPUT -s `+ReadedFile.Config.Interface.Addresses.V6+` -j MARK --set-mark 32975

[Peer]
PublicKey = `+"bmXOC+F1FxEMF9dyiK2H5/1SUtzH0JuVo51h2wPfgyo="+`
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = 162.159.193.199:2408
`, "%i", "%i", "%i", "%i")
		return config, ReadedFile.Config.ReservedHex, nil
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

func NftConfigGenerate(reserved string) (string, error) {
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
