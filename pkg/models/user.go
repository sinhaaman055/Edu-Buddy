package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Name string             `bson:"name" json:"name" binding:"required"`
	Id  primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Email string            `bson:"email" json:"email" binding:"required" validate:"required,email"`
	Username string         `bson:"username" json:"username"`
	Gender string           `bson:"gender" json:"gender" binding:"required"`
	Password string         `bson:"password" json:"password" binding:"required" validate:"required,min=7"`
	ConfirmPass string      `bson:"-" json:"confirmPass" validate:"required"`
	Country string          `bson:"country" json:"country" binding:"required"`
	DOB string              `bson:"dob" json:"dob" binding:"required"`
}