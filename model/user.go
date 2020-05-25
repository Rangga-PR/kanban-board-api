package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User : define user document in database
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username" binding:"required"`
	Password  string             `json:"password" bson:"password" binding:"required"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at,omitempty"`
}
