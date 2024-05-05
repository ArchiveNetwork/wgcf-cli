package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type HTTPClient struct{}

func (h *HTTPClient) Do(request *http.Request) (body []byte, err error) {
	var proxy string
	httpProxy := os.Getenv("http_proxy")
	httpsProxy := os.Getenv("https_proxy")
	if httpProxy != "" {
		proxy = httpProxy
	} else if httpsProxy != "" {
		proxy = httpsProxy
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			},
		},
	}
	if proxy != "" {
		client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	if body, err = io.ReadAll(response.Body); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	if response.StatusCode != 204 && response.StatusCode != 200 {
		var buffer bytes.Buffer
		if err = json.Indent(&buffer, body, "", "    "); err != nil {
			fmt.Println(string(body))
		} else {
			fmt.Println(buffer.String())
		}
		err = fmt.Errorf("REST API returned " + fmt.Sprint(response.StatusCode) + " " + http.StatusText(response.StatusCode))
	}
	return
}
