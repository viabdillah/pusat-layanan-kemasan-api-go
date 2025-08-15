// internal/models/notification.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"` // Untuk siapa notifikasi ini
	Message   string             `bson:"message" json:"message"`
	IsRead    bool               `bson:"isRead" json:"isRead"`
	Link      string             `bson:"link,omitempty" json:"link,omitempty"` // Link opsional, misal: /orders/ID_PESANAN
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
