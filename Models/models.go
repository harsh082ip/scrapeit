package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `json:"name" binding:"required"`
	Email    string             `json:"email" binding:"required,email"`
	Password string             `json:"password" binding:"required,min=6"`
	Username string             `json:"username" binding:"min=6,max=12"`
}

type LoginUser struct {
	LoginID  string `json:"login_id" binding:"required"`
	Password string `json:"password" binding:"password"`
}

type AppCredits struct {
	Email        string
	TotalCredits string
}
