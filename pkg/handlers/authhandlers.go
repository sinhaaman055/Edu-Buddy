package handlers

import (
	"context"
	"edubuddy/pkg/database"
	"edubuddy/pkg/models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)
type LoginUser struct{
	Email string     `json:"email" bson:"email"`
	Password string  `json:"password" bson:"password"`
}
type ForgetUser struct{
	Email string     `json:"email" bson:"email"`
	Username string   `json:"username" bson:"username"`
	DOB string        `json:"dob" bson:"dob"`
}
func GenerateUsername(name string)string{
	base:=strings.ToLower(strings.ReplaceAll(name," ",""))
	suffix:=time.Now().Unix() % 10000
	return fmt.Sprintf("%s%d", base, suffix)
}
func Signup(c *gin.Context){
  var User models.User
  err:=c.BindJSON(&User)
    if err!=nil{
	c.JSON(http.StatusBadRequest,gin.H{"error":"Not able to convert the data"})
	return 
    }
	if User.Password!=User.ConfirmPass{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Password should be same",
		})
		return
	}
	collection:=database.Client.Database("Light").Collection("Users")
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	count,_:=collection.CountDocuments(ctx,bson.M{"email":User.Email})
	if count>0{
		c.JSON(http.StatusBadRequest,gin.H{
		  "Message":"Email Already Registered",
		})
		return
	} 
  HashedPass,err:=bcrypt.GenerateFromPassword([]byte(User.Password),bcrypt.DefaultCost)
  if err!=nil{
    	c.JSON(http.StatusBadRequest,gin.H{"Error":"Not able to hash The Password"})
       return
	}
    User.Password=string(HashedPass)
	isUnique:=false
	var FinalUsername string
	for !isUnique{
        FinalUsername=GenerateUsername(User.Name)
		count,_:=collection.CountDocuments(ctx,bson.M{"username":FinalUsername})
		if count==0{
			isUnique=true
		}
	}
	User.Username=FinalUsername
	result,err:=collection.InsertOne(ctx,User)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Error":"Not able to Insert to databse"})
	   return 
	}
     c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully!",
		"username": User.Username,
        "userId":  result.InsertedID,
    })
}
func Login(c *gin.Context){
   var NewUser LoginUser
   var ExistUser models.User
   err:=c.BindJSON(&NewUser)
   if err!=nil{
	c.JSON(http.StatusBadRequest,gin.H{
        "Message":"Not able to bind user",
	})
	return
   }
   collection:=database.Client.Database("Light").Collection("Users")
   ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
   defer cancel()
   err=collection.FindOne(ctx,bson.M{"email":NewUser.Email}).Decode(&ExistUser)
   if err!=nil{
	c.JSON(http.StatusBadRequest,gin.H{"message":"Email is not registered"})
	return 
   }
    err=bcrypt.CompareHashAndPassword([]byte(ExistUser.Password),[]byte(NewUser.Password))
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Password don't match"})
		return
	}
	jwtsecret:=os.Getenv("JWT_SECRET")
	claims:=jwt.MapClaims{
		"username":ExistUser.Username,
		"email":ExistUser.Email,
		"iat":time.Now().Unix(),
		"exp":time.Now().Add(time.Hour*24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	TokenString,err:=token.SignedString([]byte(jwtsecret))
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to Sign the token"})
		return
	}
	c.JSON(http.StatusAccepted,gin.H{
		"Message":"Successfully Log in",
		"JWTtoken":TokenString,
		"User":ExistUser.Username,
	})
}
func ForgetPass(c *gin.Context){
	var olduser ForgetUser
	var CheckedUser models.User
	err:=c.BindJSON(&olduser)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"message":"Not able to load information"})
		return 
	}
	collection:=database.Client.Database("Light").Collection("Users")
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	errr:=collection.FindOne(ctx,bson.M{
		"$or": []bson.M{
        {"email":olduser.Email},
        {"username":olduser.Username},
    },
	}).Decode(&CheckedUser)
	if errr!=nil{
		c.JSON(http.StatusBadGateway,gin.H{"Message":"No User found"})
       return
	}
	if(olduser.DOB!=CheckedUser.DOB){
		c.JSON(http.StatusBadRequest,gin.H{"message":"Credentials donot match"})
		return
	}
	secret_key:=os.Getenv("RESET_KEY")
	claims:=jwt.MapClaims{
		"username":CheckedUser.Username,
		"email":CheckedUser.Email,
		"dob":CheckedUser.DOB,
		"purpose":"password_reset",
		"iat":time.Now().Unix(),
		"exp":time.Now().Add(time.Minute*5).Unix(),
	}
	Newtoken:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	Finaltoken,err:=Newtoken.SignedString([]byte(secret_key))
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to Sign the token"})
		return
	}
	c.JSON(http.StatusAccepted,gin.H{
		"Message":"Successfully verified user for reset password",
		"RESETtoken":Finaltoken,
		"User":CheckedUser.Username,
	})
}
func ResetPass(c *gin.Context){
	type ResetPass struct{
	    Password string   `bson:"password" json:"password"`
		ConfirmPassord string     `json:"confirmpassword"`
	}
	var Credential ResetPass
	err:=c.Bind(&Credential)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Input in wrong format"})
	}
	UserEmail,exists:=c.Get("email")
	if !exists {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Email not found in token"})
	return
     }
	 if Credential.Password=="" || Credential.ConfirmPassord==""{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Fill the correct Input"})
		return
	 }
	 emailstr:=UserEmail.(string)
	 if Credential.Password!=Credential.ConfirmPassord{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Password and Confpass Should be same"})
		return
	 }
	 HashedPAss,err:=bcrypt.GenerateFromPassword([]byte(Credential.Password),bcrypt.DefaultCost)
    if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to hash the Password"})
		return
	}
	collection:=database.Client.Database("Light").Collection("Users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	filter:=bson.M{"email":emailstr}
	update:=bson.M{"$set":bson.M{"password":string(HashedPAss)}}
    result,err:=collection.UpdateOne(ctx,filter,update)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }
	if result.ModifiedCount==0{
        c.JSON(http.StatusOK, gin.H{"message": "Password is same as old one"})
        return
	}
       c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully! Now login."})
    
}
func DeleteAccount(c *gin.Context){
	type info struct{
		Password string     `json:"password"`
		Delete string       `json:"delete"`
	}
	var input info
	err:=c.Bind(&input)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to read the input"})
		return
	}
    Email,Exist:=c.Get("email")
	if !Exist {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Email not found in token"})
	return
     }
	 Emailstr:=Email.(string)
	 Username,exit:=c.Get("username")
	 if !exit{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username not found in token"})
	return
	 }
	 Usernamestr:=Username.(string)
	 collection:=database.Client.Database("Light").Collection("Users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	var ExistUser models.User
	 err=collection.FindOne(ctx,bson.M{"email":Emailstr}).Decode(&ExistUser)
	 if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to Find user with email"})
		return
	 }
	 err=bcrypt.CompareHashAndPassword([]byte(ExistUser.Password),[]byte(input.Password))
	 if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Passwoed Don't Match"})
		return
	 }
	 if input.Delete!="DELETE"{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"String Donot Match"})
		return
	 }
	 filter:=bson.M{"$or":[]bson.M{{"email":Emailstr},{"username":Usernamestr}}}
	result,err:=collection.DeleteOne(ctx,filter)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"Message":"Not able to delete the Acoount"})
		return
	}
	c.JSON(http.StatusAccepted,gin.H{
		"Message":"Deletion successfully",
		"No of deletion":result.DeletedCount,
})
}