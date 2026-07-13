package routes

import (
	"edubuddy/pkg/handlers"
	"edubuddy/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RoomRoutes(r *gin.Engine){
	r.POST("/createRoom",middleware.Authenticate(),handlers.CreateRoom)
	r.POST("/joinRoom",middleware.Authenticate(),handlers.JoinRoom)
	r.POST("/kick/:id",middleware.Authenticate(),handlers.Kick)
	r.POST("/leave/:id",middleware.Authenticate(),handlers.Leave)
	r.GET("/ws/:roomId",middleware.Authenticate(),handlers.Wshandler)
}