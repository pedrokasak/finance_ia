package payment

// PaymentGateway is the port (adapter interface) for payment providers.
// To switch from Stripe to PayPal, implement this interface with a new adapter.
type PaymentGateway interface {
	// CreateCheckoutSession creates a hosted checkout page for subscription upgrade
	CreateCheckoutSession(opts CreateCheckoutOptions) (*CheckoutResult, error)

	// CreatePortalSession creates a billing management portal for the customer
	CreatePortalSession(customerID string, returnURL string) (*PortalResult, error)

	// CreateCustomer creates a customer in the payment provider
	CreateCustomer(email, name string) (customerID string, err error)

	// GetSubscription fetches subscription details by external ID
	GetSubscription(externalID string) (*SubscriptionInfo, error)

	// ListInvoices lists invoices for a customer
	ListInvoices(customerID string) ([]*InvoiceInfo, error)

	// ListActiveSubscriptions lists active subscriptions for a customer
	ListActiveSubscriptions(customerID string) ([]*SubscriptionInfo, error)

	// CancelSubscription cancels a subscription
	CancelSubscription(externalID string) error

	// ValidateWebhook validates the webhook payload and returns parsed events
	ValidateWebhook(payload []byte, signature string) ([]*WebhookEvent, error)
}

// CreateCheckoutOptions holds parameters for creating a checkout session
type CreateCheckoutOptions struct {
	CustomerID     string
	Email          string
	PriceID        string
	SuccessURL     string
	CancelURL      string
	Plan           string
	PaymentMethods []string // e.g. ["card"], ["card","pix"], ["boleto"]
	Metadata       map[string]string
}

// CheckoutResult contains the checkout session details
type CheckoutResult struct {
	SessionID string
	URL       string
}

// PortalResult contains the billing portal URL
type PortalResult struct {
	URL string
}

// InvoiceInfo holds normalized invoice data
type InvoiceInfo struct {
	ID               string  `json:"id"`
	Date             string  `json:"date"`
	Amount           float64 `json:"amount"`
	Status           string  `json:"status"`
	Currency         string  `json:"currency"`
	HostedInvoiceURL string  `json:"hosted_invoice_url,omitempty"`
}

// SubscriptionInfo holds normalized subscription data from the gateway
type SubscriptionInfo struct {
	ExternalID  string
	CustomerID  string
	PriceID     string
	Plan        string
	Status      string
	PeriodStart int64
	PeriodEnd   int64
}

// WebhookEventType defines the type of payment event
type WebhookEventType string

const (
	EventCheckoutCompleted       WebhookEventType = "checkout.completed"
	EventSubscriptionUpdated     WebhookEventType = "subscription.updated"
	EventSubscriptionDeleted     WebhookEventType = "subscription.deleted"
	EventInvoicePaymentSucceeded WebhookEventType = "invoice.payment_succeeded"
	EventInvoicePaymentFailed    WebhookEventType = "invoice.payment_failed"
)

// WebhookEvent is a normalized payment event from any gateway
type WebhookEvent struct {
	Type           WebhookEventType
	SubscriptionID string
	CustomerID     string
	UserID         string // from metadata
	Plan           string
	Status         string
	Metadata       map[string]string
}
