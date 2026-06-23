package handlers

import (
	"context"
	"edubuddy/pkg/database"
	"edubuddy/pkg/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)
func RoomIdGeneration(username string) string{
	suffix:=time.Now().Unix() % 10000
	return fmt.Sprintf("%s%d","roomOK"+username, suffix)
}
func CreateRoom(c* gin.Context){

	var room models.Room
    err:=c.BindJSON(&room)
    if err!=nil{
       c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to read Input infromation "})
       return
    }
	username,exists:=c.Get("username")
   if !exists{
	c.JSON(http.StatusUnauthorized,gin.H{"Message":"Not able to find user in context"})
	return
   }
   usernamestr:=username.(string)
   roomCollection := database.Client.Database("Light").Collection("StudyRooms")
   ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
   defer cancel()
   var existingRoom models.Room
	err=roomCollection.FindOne(ctx, bson.M{"created_by":usernamestr}).Decode(&existingRoom)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "You already have created one room",
			"room_id": existingRoom.ID,
		})
		return
	}
   room.CreatedBy=usernamestr
   room.ID=RoomIdGeneration(usernamestr)
   room.CreatedAt=time.Now()
   room.Members=[]string{usernamestr}
    _, err = roomCollection.InsertOne(ctx, room)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to create room in database"})
      return
    }
	c.JSON(http.StatusAccepted,gin.H{
		"message":"Room Created Successfully",
		"data":room,
	})
}

