package main

import (
	"edubuddy/pkg/database"
	"edubuddy/pkg/hub"
	"edubuddy/pkg/routes"
	"os"
     "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)
func main(){
	database.LoadEnv()
	database.Connectdb()
	Port:=os.Getenv("PORT")
	if Port==""{
		Port="8090"
	}
	go hub.RoomHub.Run()
	router:=gin.Default()
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"POST", "GET", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
	routes.AuthRoute(router)
	routes.RoomRoutes(router)
    router.Run(":"+Port)
}