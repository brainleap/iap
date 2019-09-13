package playstore

// AcknowledgementState is the data type for acknowledgement states.
type AcknowledgementState int

// List of acknowledgement states.
const (
	NotAcknowledged = 0
	Acknowledged    = 1
)

// CancelReason is the data type for cancel reasons.
type CancelReason int

// List of cancel reasons.
const (
	CRUserCanceled      = 0
	CRSystemCanceled    = 1
	CRReplaced          = 2
	CRDeveloperCanceled = 3
)

// CancelSurveyReason is the data type for cancel survey reasons.
type CancelSurveyReason int

// List of cancel survey reasons.
const (
	CSROther          = 0
	CSRNoUse          = 1
	CSRTechnicalIssue = 2
	CSRCostRelated    = 3
	CSRFoundBetter    = 4
)

// ConsumptionState is the data type for consumption states.
type ConsumptionState int

// List of consumption states.
const (
	NotConsumed = 0
	Consumed    = 1
)

// PaymentState is the data type for payment states.
type PaymentState int

// List of payment states.
const (
	PaymentPending   = 0
	PaymentReceived  = 1
	PaymentFreeTrial = 2
	PaymentDeferred  = 3
)

// PriceChangeState is the data type for price change states.
type PriceChangeState int

// List of price change states.
const (
	PCOutstanding = 0
	PCAccepted    = 1
)

// PurchaseState is the data type for purchase states.
type PurchaseState int

// List of purchase states.
const (
	PurchaseDone     = 0
	PurchaseCanceled = 1
	PurchasePending  = 2
)

// PurchaseType is the data type for purchase types.
type PurchaseType int

// List of purchase types.
const (
	PTTest     = 0
	PTPromo    = 1
	PTRewarded = 2
)

// Product indicates the status of an in-app product purchase.
type Product struct {
	Kind                 string               `json:"kind"`
	PurchaseTimeMillis   int64                `json:"purchaseTimeMillis"`
	PurchaseState        PurchaseState        `json:"purchaseState"`
	ConsumptionState     ConsumptionState     `json:"consumptionState"`
	DeveloperPayload     string               `json:"developerPayload"`
	OrderID              string               `json:"orderId"`
	PurchaseType         PurchaseType         `json:"purchaseType"`
	AcknowledgementState AcknowledgementState `json:"acknowledgementState"`
}

// Subscription indicates the status of a subscription purchase.
type Subscription struct {
	Kind                       string                 `json:"kind"`
	StartTimeMillis            int64                  `json:"startTimeMillis"`
	ExpiryTimeMillis           int64                  `json:"expiryTimeMillis"`
	AutoResumeTimeMillis       int64                  `json:"autoResumeTimeMillis"`
	AutoRenewing               bool                   `json:"autoRenewing"`
	PriceCurrencyCode          string                 `json:"priceCurrencyCode"`
	PriceAmountMicros          int64                  `json:"priceAmountMicros"`
	IntroductoryPriceInfo      *IntroductoryPriceInfo `json:"introductoryPriceInfo"`
	CountryCode                string                 `json:"countryCode"`
	DeveloperPayload           string                 `json:"developerPayload"`
	PaymentState               PaymentState           `json:"paymentState"`
	CancelReason               CancelReason           `json:"cancelReason"`
	UserCancellationTimeMillis int64                  `json:"userCancellationTimeMillis"`
	CancelSurveyResult         *CancelSurveyResult    `json:"cancelSurveyResult"`
	OrderID                    string                 `json:"orderId"`
	LinkedPurchaseToken        string                 `json:"linkedPurchaseToken"`
	PurchaseType               PurchaseType           `json:"purchaseType"`
	PriceChange                *PriceChange           `json:"priceChange"`
	ProfileName                string                 `json:"profileName"`
	EmailAddress               string                 `json:"emailAddress"`
	GivenName                  string                 `json:"givenName"`
	FamilyName                 string                 `json:"familyName"`
	ProfileID                  string                 `json:"profileId"`
	AcknowledgementState       AcknowledgementState   `json:"acknowledgementState"`
}

// IntroductoryPriceInfo is the introductory price info of a subscription.
type IntroductoryPriceInfo struct {
	CurrencyCode string `json:"introductoryPriceCurrencyCode"`
	AmountMicros int64  `json:"introductoryPriceAmountMicros"`
	Period       string `json:"introductoryPricePeriod"`
	Cycles       int    `json:"introductoryPriceCycles"`
}

// CancelSurveyResult is the info provided by the user when they complete the
// subscription cancellation flow.
type CancelSurveyResult struct {
	Reason    CancelSurveyReason `json:"cancelSurveyReason"`
	UserInput string             `json:"userInputCancelReason"`
}

// PriceChange is the latest price change info available.
type PriceChange struct {
	NewPrice *NewPrice        `json:"newPrice"`
	State    PriceChangeState `json:"state"`
}

// NewPrice is the new price info in a price change.
type NewPrice struct {
	PriceMicros int64  `json:"priceMicros"`
	Currency    string `json:"currency"`
}

// DeferralInfo is the info about the new desired expiry time for
// the subscription.
type DeferralInfo struct {
	ExpectedTimeMillis int64 `json:"expectedExpiryTimeMillis"`
	DesiredTimeMillis  int64 `json:"desiredExpiryTimeMillis"`
}
