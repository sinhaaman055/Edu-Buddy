package routes

import (
	"edubuddy/pkg/handlers"
	"edubuddy/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RoomRoutes(r *gin.Engine){
	r.POST("/createRoom",middleware.Authenticate(),handlers.CreateRoom)
	r.POST("/joinRoom",middleware.Authenticate(),handlers.JoinRoom)
}