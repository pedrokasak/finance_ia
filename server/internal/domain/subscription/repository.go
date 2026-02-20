package subscription

import "github.com/google/uuid"

// SubscriptionRepository defines persistence operations for subscriptions
type SubscriptionRepository interface {
	Upsert(sub *Subscription) error
	FindByUserID(userID uuid.UUID) (*Subscription, error)
	FindByExternalID(externalID string) (*Subscription, error)
}
