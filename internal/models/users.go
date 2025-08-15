// internal/models/user.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Definisikan tipe Role untuk kejelasan
type Role string

// Definisikan konstanta untuk setiap peran yang valid
const (
	RoleAdmin    Role = "admin"
	RoleManajer  Role = "manajer"
	RoleDesigner Role = "designer"
	RoleOperator Role = "operator"
	RoleKasir    Role = "kasir"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name" json:"name" binding:"required"`
	Email     string             `bson:"email" json:"email" binding:"required"`
	Password  string             `bson:"password" json:"password" binding:"required"`
	Role      Role               `bson:"role" json:"role"` // <-- Gunakan tipe Role kita
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
