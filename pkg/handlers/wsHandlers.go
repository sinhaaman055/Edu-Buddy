package handlers

import (
	"context"
	"edubuddy/pkg/database"
	"edubuddy/pkg/hub"
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
    roomCollection := database.Client.Database("Light").Collection("StudyRooms")
	 roomCtx, roomCancel := context.WithTimeout(context.Background(), 5*time.Second)
	 defer roomCancel()
	 count, err := roomCollection.CountDocuments(roomCtx, map[string]string{"_id": roomId})
	 if err != nil || count == 0 {
		log.Println("Room does not exist in Database:", roomId)
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	 }
   conn,err:=Upgrader.Upgrade(c.Writer,c.Request,nil)
   if err!=nil{
	return
   }
   client := &hub.Client{
        Conn:   conn,
        Send:   make(chan []byte, 256),
        RoomId: roomId,
    }
    hub.RoomHub.Register <- client
    go func(cl *hub.Client) {
        defer cl.Conn.Close()
        for msgByte := range cl.Send {
            err := cl.Conn.WriteMessage(websocket.TextMessage, msgByte)
            if err != nil {
                break
            }
		}
	}(client)
   defer func() {
        hub.RoomHub.Unregister <- client
        conn.Close()
    }()
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
    finalMessage, err := json.Marshal(msg)
    if err == nil {
        hub.RoomHub.Broadcast <- finalMessage
    }
   }
}