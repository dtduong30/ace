package models

import "time"

type Wishlist struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Title       string    `json:"title" validate:"required,min=1,max=100"`
	Type        string    `json:"type" validate:"required,oneof=Read Buy Learn Do"`
	Description string    `json:"description,omitempty" validate:"max=500"`
	Priority    int       `json:"priority,omitempty" validate:"min=1,max=5"`
	Status      string    `json:"status" validate:"required,oneof=pending in-progress completed"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
