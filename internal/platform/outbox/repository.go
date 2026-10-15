
package outbox

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var ErrInvalidEvent = errors.New("outbox: invalid event")

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// Create inserts an event using the caller's transaction.
// It does not commit the transaction.
func (r *Repository) Create(ctx context.Context,tx *gorm.DB,event *Event,) error {
    if tx == nil || event == nil {
        return ErrInvalidEvent
    }

    if strings.TrimSpace(event.EventType) == "" ||
        strings.TrimSpace(event.AggregateType) == "" ||
        strings.TrimSpace(event.AggregateID) == "" ||
        strings.TrimSpace(event.Payload) == "" {
        return ErrInvalidEvent
    }

    return tx.WithContext(ctx).Create(event).Error
}
