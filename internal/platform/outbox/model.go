
package outbox

import "time"

type Event struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventType     string     `gorm:"type:varchar(100);not null"`
	AggregateType string     `gorm:"type:varchar(50);not null"`
	AggregateID   string     `gorm:"type:uuid;not null"`
	Payload       string     `gorm:"type:jsonb;not null"`
	Status        string     `gorm:"type:varchar(20);not null;default:pending"`
	Attempts      int        `gorm:"not null;default:0"`
	AvailableAt   time.Time  `gorm:"not null"`
	LockedUntil   *time.Time
	PublishedAt   *time.Time
	LastError     *string
	CreatedAt     time.Time  `gorm:"not null"`
}

func (Event) TableName() string {
	return "outbox_events"
}
