package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	RoomId    string	               `json:"roomId" bson:"roomId"`
	SenderID  primitive.ObjectID       `json:"senderid" bson:"senderid"` 
	TimeStamp time.Time                `json:"timeStamp" bson:"timeStamp"`
	Content   string                   `json:"content" bson:"content" binding:"required"`
}