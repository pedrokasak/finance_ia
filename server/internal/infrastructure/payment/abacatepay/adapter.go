package abacatepay

import (
	"bytes"
	"encoding/json"
	"errors"
	"finance-ia/internal/domain/payment"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Adapter struct {
	baseURL       string
	apiKey        string
	webhookSecret string
	client        *http.Client
}

func NewAdapter() *Adapter {
	baseURL := strings.TrimRight(os.Getenv("ABACATEPAY_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = "https://api.abacatepay.com"
	}
	return &Adapter{
		baseURL:       baseURL,
		apiKey:        os.Getenv("ABACATEPAY_API_KEY"),
		webhookSecret: os.Getenv("ABACATEPAY_WEBHOOK_SECRET"),
		client:        &http.Client{Timeout: 20 * time.Second},
	}
}

type envelope[T any] struct {
	Data    T       `json:"data"`
	Success bool    `json:"success"`
	Error   *string `json:"error"`
}

type customerResp struct {
	ID string `json:"id"`
}

type checkoutResp struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type productResp struct {
	ID         string `json:"id"`
	ExternalID string `json:"externalId"`
	Cycle      string `json:"cycle"`
}

func (a *Adapter) doJSON(method, path string, payload any, out any) error {
	if a.apiKey == "" {
		return errors.New("abacatepay: ABACATEPAY_API_KEY is not set")
	}
	var body []byte
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = b
	}
	req, err := http.NewRequest(method, a.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("abacatepay: http %d at %s", resp.StatusCode, path)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (a *Adapter) CreateCheckoutSession(opts payment.CreateCheckoutOptions) (*payment.CheckoutResult, error) {
	if opts.PriceID == "" {
		return nil, errors.New("abacatepay: missing product id for subscription checkout")
	}

	externalID := ""
	if uid, ok := opts.Metadata["user_id"]; ok && uid != "" && opts.Plan != "" {
		externalID = "uid:" + uid + "|plan:" + opts.Plan
	}
	body := map[string]any{
		"items": []map[string]any{
			{"id": opts.PriceID, "quantity": 1},
		},
		"methods":       []string{"CARD"},
		"customerId":    opts.CustomerID,
		"returnUrl":     opts.CancelURL,
		"completionUrl": opts.SuccessURL,
	}
	if externalID != "" {
		body["externalId"] = externalID
	}
	if len(opts.Metadata) > 0 {
		body["metadata"] = opts.Metadata
	}

	var resp envelope[checkoutResp]
	if err := a.doJSON(http.MethodPost, "/v2/subscriptions/create", body, &resp); err != nil {
		return nil, err
	}
	return &payment.CheckoutResult{
		SessionID: resp.Data.ID,
		URL:       resp.Data.URL,
	}, nil
}

// EnsurePlanCatalog creates/gets recurring products for monthly and yearly cycles.
// Returns product IDs to be used at subscription checkout.
func (a *Adapter) EnsurePlanCatalog(slug, name, description string, monthly, yearly float64, currency string) (string, string, string, error) {
	monthlyExternalID := fmt.Sprintf("finance-ia:%s:monthly", slug)
	yearlyExternalID := fmt.Sprintf("finance-ia:%s:yearly", slug)

	monthlyID, err := a.findOrCreateRecurringProduct(monthlyExternalID, name+" Mensal", description, monthly, currency, "MONTHLY")
	if err != nil {
		return "", "", "", err
	}
	yearlyID, err := a.findOrCreateRecurringProduct(yearlyExternalID, name+" Anual", description, yearly, currency, "ANNUALLY")
	if err != nil {
		return "", "", "", err
	}
	// Keep compatibility with existing sync signature:
	// productID (unused) + monthly item ID + yearly item ID
	return monthlyID, monthlyID, yearlyID, nil
}

func (a *Adapter) findOrCreateRecurringProduct(externalID, name, description string, amount float64, currency, cycle string) (string, error) {
	products, err := a.listProducts()
	if err != nil {
		return "", err
	}
	for _, p := range products {
		if p.ExternalID == externalID {
			return p.ID, nil
		}
	}

	body := map[string]any{
		"externalId":  externalID,
		"name":        name,
		"description": description,
		"price":       int64(amount * 100),
		"currency":    strings.ToUpper(currency),
		"cycle":       cycle,
	}
	var resp envelope[productResp]
	if err := a.doJSON(http.MethodPost, "/v2/products/create", body, &resp); err != nil {
		return "", err
	}
	return resp.Data.ID, nil
}

func (a *Adapter) listProducts() ([]productResp, error) {
	var resp envelope[[]productResp]
	if err := a.doJSON(http.MethodGet, "/v2/products/list", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (a *Adapter) CreatePortalSession(customerID string, returnURL string) (*payment.PortalResult, error) {
	return nil, errors.New("abacatepay: billing portal is not supported by this adapter")
}

func (a *Adapter) CreateCustomer(email, name string) (string, error) {
	body := map[string]any{
		"email": email,
	}
	if strings.TrimSpace(name) != "" {
		body["name"] = name
	}
	var resp envelope[customerResp]
	if err := a.doJSON(http.MethodPost, "/v2/customers/create", body, &resp); err != nil {
		return "", err
	}
	return resp.Data.ID, nil
}

func (a *Adapter) GetSubscription(externalID string) (*payment.SubscriptionInfo, error) {
	return nil, errors.New("abacatepay: GetSubscription not implemented")
}

func (a *Adapter) ListInvoices(customerID string) ([]*payment.InvoiceInfo, error) {
	// Avoid breaking UI sections that expect invoice array.
	return []*payment.InvoiceInfo{}, nil
}

func (a *Adapter) ListActiveSubscriptions(customerID string) ([]*payment.SubscriptionInfo, error) {
	return []*payment.SubscriptionInfo{}, nil
}

func (a *Adapter) CancelSubscription(externalID string) error {
	return errors.New("abacatepay: CancelSubscription not implemented")
}

func (a *Adapter) CreateProduct(name, description string) (string, error) {
	return "", errors.New("abacatepay: product provisioning not implemented in admin handler")
}

func (a *Adapter) CreatePrice(productID string, amount float64, currency string, interval string) (string, error) {
	return "", errors.New("abacatepay: pricing is managed by product cycle")
}

func (a *Adapter) ArchiveProduct(productID string) error {
	return errors.New("abacatepay: archive product not implemented")
}

type abacateWebhook struct {
	Event string `json:"event"`
	Data  struct {
		Subscription struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"subscription"`
		Payment struct {
			ExternalID string `json:"externalId"`
		} `json:"payment"`
	} `json:"data"`
}

func (a *Adapter) ValidateWebhook(payload []byte, signature string) ([]*payment.WebhookEvent, error) {
	if a.webhookSecret != "" && signature != a.webhookSecret {
		return nil, errors.New("abacatepay: invalid webhook signature")
	}
	var p abacateWebhook
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if p.Event == "" {
		return nil, errors.New("abacatepay: missing event")
	}

	ev := &payment.WebhookEvent{
		SubscriptionID: p.Data.Subscription.ID,
	}
	switch p.Event {
	case "subscription.completed":
		ev.Type = payment.EventCheckoutCompleted
		ev.Status = "active"
	case "subscription.renewed":
		ev.Type = payment.EventSubscriptionUpdated
		ev.Status = "active"
	case "subscription.cancelled":
		ev.Type = payment.EventSubscriptionDeleted
		ev.Status = "canceled"
	default:
		return []*payment.WebhookEvent{}, nil
	}

	// externalId format created in checkout: uid:<uuid>|plan:<slug>
	parts := strings.Split(p.Data.Payment.ExternalID, "|")
	for _, part := range parts {
		if strings.HasPrefix(part, "uid:") {
			ev.UserID = strings.TrimPrefix(part, "uid:")
		}
		if strings.HasPrefix(part, "plan:") {
			ev.PriceID = strings.TrimPrefix(part, "plan:")
		}
	}

	return []*payment.WebhookEvent{ev}, nil
}
