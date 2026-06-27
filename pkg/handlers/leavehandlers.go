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

func Leave(c *gin.Context){
	var activeRoom models.Room
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
	collection:=database.Client.Database("Light").Collection("StudyRooms")
    ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	err:=collection.FindOne(ctx,bson.M{"_id":id}).Decode(&activeRoom)
	if err!=nil{
		c.JSON(http.StatusBadGateway,gin.H{"Message":"room id not present in Active room "})
      return
	}
	if activeRoom.CreatedBy==username{
		_,err=collection.DeleteOne(ctx,bson.M{"_id":id})
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to Delete the Room"})
			return
		}
		c.JSON(http.StatusOK,gin.H{"Message":"Room deleted sucessfully"})
		return
	}
	filter:=bson.M{"_id":id}
	update:=bson.M{"$pull":bson.M{"members":username}}
	_,err=collection.UpdateOne(ctx,filter,update)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to Delete member from the Room"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"Message":"Member deleted sucessfully"})
}