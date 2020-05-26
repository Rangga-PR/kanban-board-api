package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Task : define task document in database
type Task struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	Title     string             `json:"title" bson:"title,omitempty"`
	Content   string             `json:"content" bson:"content,omitempty"`
	Status    string             `json:"status" bson:"status"`
	Icon      string             `json:"icon" bson:"icon"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at,omitempty"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at,omitempty"`
}
