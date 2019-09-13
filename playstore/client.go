package playstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	baseURL   = "https://www.googleapis.com/androidpublisher/v3"
	authScope = "https://www.googleapis.com/auth/androidpublisher"
)

// Option is the common type for optional arguments.
type Option interface {
	isOption()
}

// DeveloperPayload is the optional developer payload argument.
type DeveloperPayload string

func (DeveloperPayload) isOption() {}

func getDeveloperPayload(opts []Option) string {
	for _, o := range opts {
		p, ok := o.(DeveloperPayload)
		if ok {
			return string(p)
		}
	}
	return ""
}

// NewClient creates a new PlayStore client.
func NewClient(jsonKey []byte) (*Client, error) {
	return NewClientWithProxy(jsonKey, "")
}

// NewClientWithProxy creates a new PlayStore client with a proxy.
func NewClientWithProxy(jsonKey []byte, proxy string) (*Client, error) {
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

	conf, err := google.JWTConfigFromJSON(jsonKey, authScope)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)

	return &Client{Client: conf.Client(ctx)}, nil
}

// Client provides PlayStore in-app billing API.
type Client struct {
	Client *http.Client
}

// AcknowledgeProduct acknowledges purchase of an in-app product.
func (c *Client) AcknowledgeProduct(pkg, prod, token string, opts ...Option) error {
	body := struct {
		DeveloperPayload string `json:"developerPayload"`
	}{
		DeveloperPayload: getDeveloperPayload(opts),
	}

	reqBody, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/applications/%s/purchases/products/%s/tokens/%s:acknowledge",
		baseURL,
		url.PathEscape(pkg),
		url.PathEscape(prod),
		url.PathEscape(token),
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status: %d", res.StatusCode)
	}

	return nil
}

// GetProduct checks the purchase and consumption status of an in-app product.
func (c *Client) GetProduct(pkg, prod, token string) (*Product, error) {
	url := fmt.Sprintf(
		"%s/applications/%s/purchases/products/%s/tokens/%s",
		baseURL,
		url.PathEscape(pkg),
		url.PathEscape(prod),
		url.PathEscape(token),
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status: %d", res.StatusCode)
	}

	var p Product

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

// AcknowledgeSubscription acknowledges a subscription purchase.
func (c *Client) AcknowledgeSubscription(pkg, sub, token string, opts ...Option) error {
	body := struct {
		DeveloperPayload string `json:"developerPayload"`
	}{
		DeveloperPayload: getDeveloperPayload(opts),
	}

	reqBody, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/applications/%s/purchases/subscriptions/%s/tokens/%s:acknowledge",
		baseURL,
		url.PathEscape(pkg),
		url.PathEscape(sub),
		url.PathEscape(token),
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status: %d", res.StatusCode)
	}

	return nil
}

// CancelSubscription cancels a subscription purchase.
func (c *Client) CancelSubscription(pkg, sub, token string) error {
	url := fmt.Sprintf(
		"%s/applications/%s/purchases/subscriptions/%s/tokens/%s:cancel",
		baseURL,
		url.PathEscape(pkg),
		url.PathEscape(sub),
		url.PathEscape(token),
	)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status: %d", res.StatusCode)
	}

	return nil
}

// DeferSubscription defers a subscription purchase.
func (c *Client) DeferSubscription(pkg, sub, token string, expected, desired int64) (int64, error) {
	body := struct {
		Info *DeferralInfo `json:"deferralInfo"`
	}{
		Info: &DeferralInfo{
			ExpectedTimeMillis: expected,
			DesiredTimeMillis:  desired,
		},
	}

	reqBody, err := json.Marshal(&body)
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf(
		"%s/applications/%s/purchases/subscriptions/%s/tokens/%s:defer",
		baseURL,
		url.PathEscape(pkg),
		url.PathEscape(sub),
		url.PathEscape(token),
	)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return 0, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed with status: %d", res.StatusCode)
	}

	var respBody struct {
		NewTimeMillis int64 `json:"newExpiryTimeMillis"`
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&respBody); err != nil {
		return 0, err
	}

	return respBody.NewTimeMillis, nil
}

// RefundSubscription refunds a subscription purchase.
func (c *Client) RefundSubscription(pkg, sub, token string) error {
	url := fmt.Sprintf(
		"%s/applications/%s/purchases/subscriptions/%s/tokens/%s:refund",
		baseURL,
		url.PathEscape(pkg),
		url.PathEscape(sub),
		url.PathEscape(token),
	)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status: %d", res.StatusCode)
	}

	return nil
}

// RevokeSubscription revokes a subscription purchase.
func (c *Client) RevokeSubscription(pkg, sub, token string) error {
	url := fmt.Sprintf(
		"%s/applications/%s/purchases/subscriptions/%s/tokens/%s:revoke",
		baseURL,
		url.PathEscape(pkg),
		url.PathEscape(sub),
		url.PathEscape(token),
	)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status: %d", res.StatusCode)
	}

	return nil
}

// GetSubscription checks the purchase and consumption status of a subscription.
func (c *Client) GetSubscription(pkg, sub, token string) (*Subscription, error) {
	url := fmt.Sprintf(
		"%s/applications/%s/purchases/subscriptions/%s/tokens/%s",
		baseURL,
		url.PathEscape(pkg),
		url.PathEscape(sub),
		url.PathEscape(token),
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status: %d", res.StatusCode)
	}

	var s Subscription

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}

	return &s, nil
}
