package main

import (
	"edubuddy/pkg/database"
	"edubuddy/pkg/hub"
	"edubuddy/pkg/routes"
	"os"

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
	routes.AuthRoute(router)
	routes.RoomRoutes(router)
    router.Run(":"+Port)
}