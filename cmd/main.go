package main

import (
    "edubuddy/pkg/routes"
	"edubuddy/pkg/database"
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
	router:=gin.Default()
	routes.AuthRoute(router)
  router.Run(":"+Port)

}