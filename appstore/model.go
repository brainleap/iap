package appstore

import "errors"

// Response contains the status of a receipt validation.
type Response struct {
	Status             int                  `json:"status"`
	Receipt            *Receipt             `json:"receipt"`
	LatestReceipt      string               `json:"latest_receipt"`
	LatestReceiptInfo  []InApp              `json:"latest_receipt_info"`
	PendingRenewalInfo []PendingRenewalInfo `json:"pending_renewal_info"`
	IsRetryable        bool                 `json:"is-retryable"`
}

// Receipt is the receipt that was sent for verification.
type Receipt struct {
	ReceiptType                string  `json:"receipt_type"`
	AdamID                     int64   `json:"adam_id"`
	AppItemID                  string  `json:"app_item_id"`
	BundleID                   string  `json:"bundle_id"`
	ApplicationVersion         string  `json:"application_version"`
	DownloadID                 int64   `json:"download_id"`
	VersionExternalIdentifier  string  `json:"version_external_identifier"`
	InApp                      []InApp `json:"in_app"`
	ReceiptCreationDate        string  `json:"receipt_creation_date"`
	ReceiptCreationDateMS      string  `json:"receipt_creation_date_ms"`
	ReceiptCreationDatePST     string  `json:"receipt_creation_date_pst"`
	RequestDate                string  `json:"request_date"`
	RequestDateMS              string  `json:"request_date_ms"`
	RequestDatePST             string  `json:"request_date_pst"`
	OriginalPurchaseDate       string  `json:"original_purchase_date"`
	OriginalPurchaseDateMS     string  `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePST    string  `json:"original_purchase_date_pst"`
	OriginalApplicationVersion string  `json:"original_application_version"`
}

// InApp is the receipt for an in-app purchase.
type InApp struct {
	Quantity                string `json:"quantity"`
	ProductID               string `json:"product_id"`
	TransactionID           string `json:"transaction_id"`
	OriginalTransactionID   string `json:"original_transaction_id"`
	WebOrderLineItemID      string `json:"web_order_line_item_id"`
	IsTrialPeriod           string `json:"is_trial_period"`
	IsInIntroOfferPeriod    string `json:"is_in_intro_offer_period"`
	ExpiresDate             string `json:"expires_date"`
	ExpiresDateMS           string `json:"expires_date_ms"`
	ExpiresDatePST          string `json:"expires_date_pst"`
	ExpiresDateFormatted    string `json:"expires_date_formatted"`
	ExpiresDateFormattedPST string `json:"expires_date_formatted_pst"`
	PurchaseDate            string `json:"purchase_date"`
	PurchaseDateMS          string `json:"purchase_date_ms"`
	PurchaseDatePST         string `json:"purchase_date_pst"`
	OriginalPurchaseDate    string `json:"original_purchase_date"`
	OriginalPurchaseDateMS  string `json:"original_purchase_date_ms"`
	OriginalPurchaseDatePST string `json:"original_purchase_date_pst"`
	CancellationDate        string `json:"cancellation_date"`
	CancellationDateMS      string `json:"cancellation_date_ms"`
	CancellationDatePST     string `json:"cancellation_date_pst"`
	CancellationReason      string `json:"cancellation_reason"`
}

// PendingRenewalInfo contains the pending renewal info for each auto-renewable
// subscription
type PendingRenewalInfo struct {
	SubscriptionExpirationIntent   string `json:"expiration_intent"`
	SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`
	SubscriptionRetryFlag          string `json:"is_in_billing_retry_period"`
	SubscriptionAutoRenewStatus    string `json:"auto_renew_status"`
	SubscriptionPriceConsentStatus string `json:"price_consent_status"`
	ProductID                      string `json:"product_id"`
	OriginalTransactionID          string `json:"original_transaction_id"`
}

// Err returns receipt validation error.
func (r *Response) Err() error {
	switch r.Status {
	case 0:
		return nil
	case 21000:
		return errors.New("could not read the JSON object you provided")
	case 21002:
		return errors.New("data in the receipt-data property was malformed or missing")
	case 21003:
		return errors.New("receipt could not be authenticated")
	case 21004:
		return errors.New("shared secret you provided does not match the shared secret on file for your account")
	case 21005:
		return errors.New("receipt server is not currently available")
	case 21007:
		return errors.New("receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead")
	case 21008:
		return errors.New("receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead")
	case 21010:
		return errors.New("receipt could not be authorized. Treat this the same as if a purchase was never made")
	default:
		if r.Status >= 21100 && r.Status <= 21199 {
			return errors.New("internal data access error")
		}

		return errors.New("unknown error occurred")
	}
}
