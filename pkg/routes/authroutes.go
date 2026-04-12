package routes

import (
	"edubuddy/pkg/handlers"
	"edubuddy/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.Engine) {
   r.POST("/signup",handlers.Signup)
   r.POST("/login",handlers.Login)
   r.POST("/forgetpassword",handlers.ForgetPass)
   r.PATCH("/resetpassword",middleware.Authenticate(),handlers.ResetPass)
   r.POST("/home/deleteaccount",middleware.Authenticate(),handlers.DeleteAccount)
}