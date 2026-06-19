package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticate() gin.HandlerFunc{
	return func(c *gin.Context){
      token:=c.Request.Header.Get("Authorization")
	  if token == ""{
		c.JSON(http.StatusBadRequest,gin.H{"message":"No Authorization"})
		c.Abort()
		return 
	  }
	  extracted:=strings.Split(token," ")
	  if len(extracted)!=2{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Invalid token Format"})
		 c.Abort()
		 return
	  }
	  tokenString:=extracted[1]
      finalToken,err:=jwt.Parse(tokenString,func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")),nil
	  })	  		
	
	if err!=nil || !finalToken.Valid{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Expired Token"})
			c.Abort()
			return
	    }
		claims,ok:=finalToken.Claims.(jwt.MapClaims)
		if ok && finalToken.Valid{
			c.Set("email", claims["email"])
			c.Set("username", claims["username"])
		}
		c.Next()
	}
} 
func AuthenticateReserToken() gin.HandlerFunc{
	return func(c *gin.Context){
		token:=c.Request.Header.Get("Authorization")
		if token==""{
			c.JSON(http.StatusBadRequest,gin.H{"Message":"No token Available"})
			c.Abort()
			return 
		}
		extracted:=strings.Split(token," ")
		if len(extracted)!=2{
			c.JSON(400,gin.H{"Message":"Token is not in format"})
			c.Abort()
			return
		}
		tokenstr:=extracted[1]
		finalToken,err:=jwt.Parse(tokenstr,func(t *jwt.Token) (any, error) {
			return []byte(os.Getenv("RESET_KEY")),nil
		})
		if err!=nil || !finalToken.Valid{
			c.JSON(http.StatusBadRequest,gin.H{"Message":"Invalid token"})
			c.Abort()
			return
		}
		claims,ok:=finalToken.Claims.(jwt.MapClaims)
		if ok && finalToken.Valid{
			c.Set("email", claims["email"])
			c.Set("username", claims["username"])
		}
		c.Next()
	}

}