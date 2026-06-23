package models

import (
	"time"
)

type Room struct {
	ID                 string               `json:"id" bson:"_id"`
	RoomName           string               `json:"roomname" bson:"roomname"`
	Description        string               `json:"description" bson:"description"`
	CreatedBy          string               `json:"createdby" bson:"createdby"`
	CreatedAt          time.Time            `json:"createdat" bson:"createdat"`
	Members            []string             `json:"members" bson:"members"`
	Password           string               `json:"password" bson:"password"` 
}