package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func updateConfigFile(filePath string) error {
	var genericResponse map[string]interface{}
	var file *os.File
	var content []byte
	var err error
	if file, err = os.Open(filePath); err != nil {
		panic(err)
	}
	if content, err = io.ReadAll(file); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(content, &genericResponse); err == nil {
		if account, ok := genericResponse["account"].(map[string]interface{}); ok {
			if _, ok := account["license"]; ok {
				var ReadedFile NormalResponse
				var responseBody NormalResponse
				var response *http.Response
				var request *http.Request
				var body []byte
				var client *http.Client
				var updatedContent []byte

				if err = json.Unmarshal(content, &ReadedFile); err != nil {
					panic(err)
				}
				defer file.Close()
				if ReadedFile.Config.ReservedDec == nil || ReadedFile.Config.ReservedHex == "" {
					clientID := ReadedFile.Config.ClientID
					decoded, err := base64.StdEncoding.DecodeString(clientID)
					if err != nil {
						panic(err)
					}
					hexString := hex.EncodeToString(decoded)

					reserved := []int{}
					for i := 0; i < len(hexString); i += 2 {
						hexByte := hexString[i : i+2]
						decValue, _ := strconv.ParseInt(hexByte, 16, 64)
						reserved = append(reserved, int(decValue))
					}

					ReadedFile.Config.ReservedDec = reserved
					ReadedFile.Config.ReservedHex = "0x" + hexString
				}

				client = &http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							MinVersion: tls.VersionTLS12,
							MaxVersion: tls.VersionTLS13,
						},
					},
				}

				if request, err = http.NewRequest("GET", "https://api.cloudflareclient.com/v0a2158/reg/"+ReadedFile.ID, nil); err != nil {
					panic(err)
				}
				request.Header.Add("CF-Client-Version", "a-7.21-0721")
				request.Header.Add("User-Agent", "okhttp/0.7.21")
				request.Header.Add("Authorization", "Bearer "+ReadedFile.Token)

				if response, err = client.Do(request); err != nil {
					panic(err)
				}
				if body, err = io.ReadAll(response.Body); err != nil {
					panic(err)
				}
				if response.StatusCode != 204 && response.StatusCode != 200 {
					var prettyJSON bytes.Buffer
					if err = json.Indent(&prettyJSON, body, "", "  "); err != nil {
						fmt.Println(string(body))
					} else {
						fmt.Println(prettyJSON.String())
					}
					panic("REST API returned " + fmt.Sprint(response.StatusCode) + " " + http.StatusText(response.StatusCode))
				}
				if err = json.Unmarshal(body, &responseBody); err != nil {
					panic(err)
				}
				if ReadedFile.Account.PrivateKey != "" {
					responseBody.Config.PrivateKey = ReadedFile.Account.PrivateKey
				} else {
					responseBody.Config.PrivateKey = ReadedFile.Config.PrivateKey
				}
				responseBody.Config.ReservedDec = ReadedFile.Config.ReservedDec
				responseBody.Config.ReservedHex = ReadedFile.Config.ReservedHex
				responseBody.Token = ReadedFile.Token
				if updatedContent, err = json.MarshalIndent(responseBody, "", "    "); err != nil {
					panic(err)
				}

				if err = os.WriteFile(filePath, updatedContent, 0600); err != nil {
					panic(err)
				}
				return nil
			}
		}
		if _, ok := genericResponse["version"]; ok {
			var ReadedFile TeamResponse
			var responseBody TeamResponse
			var response *http.Response
			var request *http.Request
			var body []byte
			var client *http.Client
			var updatedContent []byte

			if err = json.Unmarshal(content, &ReadedFile); err != nil {
				panic(err)
			}
			defer file.Close()
			if ReadedFile.Config.ReservedDec == nil || ReadedFile.Config.ReservedHex == "" {
				clientID := ReadedFile.Config.ClientID
				decoded, err := base64.StdEncoding.DecodeString(clientID)
				if err != nil {
					panic(err)
				}
				hexString := hex.EncodeToString(decoded)

				reserved := []int{}
				for i := 0; i < len(hexString); i += 2 {
					hexByte := hexString[i : i+2]
					decValue, _ := strconv.ParseInt(hexByte, 16, 64)
					reserved = append(reserved, int(decValue))
				}

				ReadedFile.Config.ReservedDec = reserved
				ReadedFile.Config.ReservedHex = "0x" + hexString
			}

			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						MinVersion: tls.VersionTLS12,
						MaxVersion: tls.VersionTLS13,
					},
				},
			}

			if request, err = http.NewRequest("GET", "https://api.cloudflareclient.com/v0a2158/reg/"+ReadedFile.ID, nil); err != nil {
				panic(err)
			}
			request.Header.Add("CF-Client-Version", "a-7.21-0721")
			request.Header.Add("User-Agent", "okhttp/0.7.21")
			request.Header.Add("Authorization", "Bearer "+ReadedFile.Token)

			if response, err = client.Do(request); err != nil {
				panic(err)
			}
			if body, err = io.ReadAll(response.Body); err != nil {
				panic(err)
			}
			if response.StatusCode != 204 && response.StatusCode != 200 {
				var prettyJSON bytes.Buffer
				if err = json.Indent(&prettyJSON, body, "", "  "); err != nil {
					fmt.Println(string(body))
				} else {
					fmt.Println(prettyJSON.String())
				}
				panic("REST API returned " + fmt.Sprint(response.StatusCode) + " " + http.StatusText(response.StatusCode))
			}
			if err = json.Unmarshal(body, &responseBody); err != nil {
				panic(err)
			}

			if ReadedFile.Account.PrivateKey != "" {
				responseBody.Config.PrivateKey = ReadedFile.Account.PrivateKey
			} else {
				responseBody.Config.PrivateKey = ReadedFile.Config.PrivateKey
			}
			responseBody.Config.ReservedDec = ReadedFile.Config.ReservedDec
			responseBody.Config.ReservedHex = ReadedFile.Config.ReservedHex
			responseBody.InstallID = ReadedFile.InstallID
			responseBody.Token = ReadedFile.Token

			if updatedContent, err = json.MarshalIndent(responseBody, "", "    "); err != nil {
				panic(err)
			}

			if err = os.WriteFile(filePath, updatedContent, 0600); err != nil {
				panic(err)
			}
			return nil
		}
	}
	panic(nil)
}
