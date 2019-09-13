package cafebazaar

// ConsumptionState is the data type for consumption states.
type ConsumptionState int

// List of consumption states.
const (
	Consumed    = 0
	NotConsumed = 1
)

// PurchaseState is the data type for purchase states.
type PurchaseState int

// List of purchase states.
const (
	PurchaseDone     = 0
	PurchaseRefunded = 1
)

// Product indicates the status of an in-app product purchase.
type Product struct {
	Kind               string           `json:"kind"`
	PurchaseTimeMillis int64            `json:"purchaseTime"`
	PurchaseState      PurchaseState    `json:"purchaseState"`
	ConsumptionState   ConsumptionState `json:"consumptionState"`
	DeveloperPayload   string           `json:"developerPayload"`
}

// Subscription indicates the status of a subscription purchase.
type Subscription struct {
	Kind                 string `json:"kind"`
	InitiationTimeMillis int64  `json:"initiationTimestampMsec"`
	ValidUntilTimeMillis int64  `json:"validUntilTimestampMsec"`
	AutoRenewing         bool   `json:"autoRenewing"`
}
