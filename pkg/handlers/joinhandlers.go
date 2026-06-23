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

type join struct{
	Id  string             `json:"id"`
	Password string        `json:"password"`
}
func JoinRoom(c *gin.Context){
   var joined join
   var user models.User
   var created models.Room
   err:=c.BindJSON(&joined)
   if err!=nil{
	c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to read the input"})
	return
   }
	if(joined.Id==""){
	c.JSON(http.StatusBadRequest,gin.H{"Message":"Id not present"})
	return
   }
   collection:=database.Client.Database("Light").Collection("StudyRooms")
   ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
   defer cancel()
   err=collection.FindOne(ctx,bson.M{"_id":joined.Id}).Decode(&created)
   if err!=nil{
	c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to find the room"})
	return
   }
   if(created.Password!=joined.Password){
	c.JSON(http.StatusBadRequest,gin.H{"Message":"Password dont match"})
	return;
   }
   email,exists:=c.Get("email")
   if !exists{
	c.JSON(http.StatusUnauthorized,gin.H{"Message":"email not present in context"})
     return 
   }
     emailstr:=email.(string)
   collection=database.Client.Database("Light").Collection("Users")
   err=collection.FindOne(ctx,bson.M{"email":emailstr}).Decode(&user)
   if err!=nil{
	c.JSON(http.StatusBadRequest,gin.H{"Message":"NOt able to decode the user"})
    return
   }
   roomCollection:=database.Client.Database("Light").Collection("StudyRooms")
   filter:=bson.M{"_id": joined.Id}
update:=bson.M{
    "$addToSet":bson.M{
        "members":user.Username, 
    },
}
_, err=roomCollection.UpdateOne(ctx, filter, update)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"Message":"Failed to join room"})
    return
}
c.JSON(http.StatusOK, gin.H{"Message":"Successfully joined the room!"})
}