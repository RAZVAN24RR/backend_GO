package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User reprezintă un utilizator în baza de date.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}

// Credentials este folosit pentru autentificare și înregistrare.
type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}