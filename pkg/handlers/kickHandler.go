package handlers

import (
	"context"
	"edubuddy/pkg/database"
	"edubuddy/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)
func Kick(c *gin.Context) {
    id:=c.Param("id")
	if id==""{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Id feild is empty in URL"})
		return
	}
	user,exists:=c.Get("username")
	if !exists{
		c.JSON(http.StatusUnauthorized,gin.H{"Message":"No usermae in token"})
		return
	}
	username:=user.(string)
	var activeRoom models.Room
	collection:=database.Client.Database("Light").Collection("StudyRooms")
    ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	err:=collection.FindOne(ctx,bson.M{"_id":id}).Decode(&activeRoom)
	if err!=nil{
		c.JSON(http.StatusBadGateway,gin.H{"Message":"room id not present in Active room "})
      return
	}
	if username!=activeRoom.CreatedBy{
       c.JSON(http.StatusUnauthorized,gin.H{"Message":"User is Not the Admin so not able to kick"})
	   return
	}
	var Input struct{
	KickUsername string   `json:"kickusername"`
    } 
	err=c.BindJSON(&Input)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to Bind the user"})
		return
	}
	if Input.KickUsername==""{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Username Not availble"})
		return
	}
	filter:=bson.M{"_id":id}
	update:=bson.M{"$pull":bson.M{"members":Input.KickUsername}}
	result,err:=collection.UpdateOne(ctx,filter,update)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to remove the member"})
		return
	}
    if result.ModifiedCount==0{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"User Not present in the Room"})
		return
	}
   c.JSON(http.StatusOK,gin.H{"Message":"Member Removed Successfully"})
}