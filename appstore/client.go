package appstore

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	productionBaseURL = "https://buy.itunes.apple.com/verifyReceipt"
	sandboxBaseURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
)

// Mode is the data type for verification mode.
type Mode int

// List of verification modes.
const (
	ProductionMode Mode = 0
	SandboxMode    Mode = 1
)

// NewClient creates a new AppStore client.
func NewClient(mode Mode) (*Client, error) {
	return NewClientWithProxy(mode, "")
}

// NewClientWithProxy creates a new AppStore client with a proxy.
func NewClientWithProxy(mode Mode, proxy string) (*Client, error) {
	c := &http.Client{}

	if proxy == "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}

		c.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	return &Client{Client: c, Mode: mode}, nil
}

// Client provides AppStore in-app billing API.
type Client struct {
	Client *http.Client
	Mode   Mode
}

// Verify validates a purchase receipt.
func (c *Client) Verify(receipt, password string) (*Response, error) {
	body := struct {
		ReceiptData            string `json:"receipt-data"`
		Password               string `json:"password"`
		ExcludeOldTransactions bool   `json:"exclude-old-transactions"`
	}{
		ReceiptData:            receipt,
		Password:               password,
		ExcludeOldTransactions: true,
	}

	reqBody, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	var url string
	if c.Mode == ProductionMode {
		url = productionBaseURL
	} else {
		url = sandboxBaseURL
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r Response

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}
