package handlers

import (
	"context"
	"edubuddy/pkg/database"
	"edubuddy/pkg/models"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
var Upgrader =websocket.Upgrader {
   CheckOrigin: func(r *http.Request) bool {
       return true
    },
}
func Wshandler(c *gin.Context){
   roomId:=c.Param("roomId")
   val, exists := c.Get("id")
   if !exists {
      log.Println("User ID not found in context")
      return
   }
   userIDStr, ok := val.(string)
   if !ok {
      log.Println("User ID is not a string")
      return
   }
   senderObjectID, err := primitive.ObjectIDFromHex(userIDStr)
   if err != nil {
      log.Println("Invalid User ID format:", err)
      return
    }
   conn,err:=Upgrader.Upgrade(c.Writer,c.Request,nil)
   if err!=nil{
	return
   }
   defer conn.Close()
   for{
       _, messageByte,err := conn.ReadMessage()
       if err != nil {
          break
       }
     var msg models.Message
     err=json.Unmarshal(messageByte,&msg) 
     if err!=nil{
      continue
     }
     msg.RoomId=roomId
     msg.TimeStamp=time.Now()
     msg.SenderID=senderObjectID
     collection:=database.Client.Database("Light").Collection("Messages")
     ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
     _,err=collection.InsertOne(ctx,msg,)
       cancel()
     if err!=nil{
      log.Println("Not able to print the data")
      continue
     }
     err = conn.WriteJSON(msg)
		if err != nil {
			break
		}
   }
}