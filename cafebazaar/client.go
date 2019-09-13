package cafebazaar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	paymentBaseURL = "https://pardakht.cafebazaar.ir/devapi/v2/api"
	authorizeURL   = "https://pardakht.cafebazaar.ir/devapi/v2/auth/authorize/"
	tokenURL       = "https://pardakht.cafebazaar.ir/devapi/v2/auth/token/"
	authScope      = "androidpublisher"
)

// NewClient creates a new Cafebazaar client.
func NewClient(clientID, clientSecret string) *Client {
	c := &Client{}

	c.OAuth = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{authScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizeURL,
			TokenURL: tokenURL,
		},
	}

	return c
}

// Client provides Cafebazaar in-app billing API.
type Client struct {
	OAuth  *oauth2.Config
	Client *http.Client
}

// AuthCodeURL returns URL to which user must be redirected to be asked for
// permission.
func (c *Client) AuthCodeURL() string {
	return c.OAuth.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// Setup initializes the client by providing authorization code.
func (c *Client) Setup(authCode string) error {
	tok, err := c.OAuth.Exchange(context.Background(), authCode)
	if err != nil {
		return err
	}

	c.Client = c.OAuth.Client(context.Background(), tok)

	return nil
}

// ValidateProduct checks the purchase and consumption status of an in-app
// product.
func (c *Client) ValidateProduct(pkg, prod, token string) (*Product, error) {
	url := fmt.Sprintf(
		"%s/validate/%s/inapp/%s/purchases/%s/",
		paymentBaseURL,
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

// ValidateSubscription checks the purchase and consumption status of a
// subscription.
func (c *Client) ValidateSubscription(pkg, prod, token string) (*Subscription, error) {
	url := fmt.Sprintf(
		"%s/applications/%s/subscriptions/%s/purchases/%s/",
		paymentBaseURL,
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

	var s Subscription

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

// CancelSubscription cancels a subscription purchase.
func (c *Client) CancelSubscription(pkg, prod, token string) error {
	url := fmt.Sprintf(
		"%s/applications/%s/subscriptions/%s/purchases/%s/cancel/",
		paymentBaseURL,
		url.PathEscape(pkg),
		url.PathEscape(prod),
		url.PathEscape(token),
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
