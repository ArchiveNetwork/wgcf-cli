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

type HTTPClient struct {
	client *http.Client
	body   []byte
}

func (h *HTTPClient) New() {
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

	h.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			},
		},
	}
	if proxy != "" {
		h.client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
	}
}
func (h *HTTPClient) Do(request *http.Request) (body []byte, err error) {
	response, err := h.client.Do(request)
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
		err = fmt.Errorf("REST API returned %d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}
	h.body = body
	return
}

func (h *HTTPClient) HandleBody() {
	var buffer bytes.Buffer
	if err := json.Indent(&buffer, h.body, "", "    "); err != nil {
		fmt.Fprint(os.Stderr, string(h.body))
	} else {
		output := buffer.Bytes()
		if output[len(output)-1] != '\n' {
			output = append(output, '\n')
		}
		fmt.Fprint(os.Stderr, string(output))
	}
}
