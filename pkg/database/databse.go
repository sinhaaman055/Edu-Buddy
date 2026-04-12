package database
import (
	"context"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var Client *mongo.Client
func LoadEnv(){
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal(err)
	}
}
func Connectdb(){
	uri:=os.Getenv("MONGO")
	if uri==""{
		log.Fatal("Uri is empty in .env")
	}
	ServerAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts:=options.Client().ApplyURI(uri).SetServerAPIOptions(ServerAPI)
	ctx,cncl:=context.WithTimeout(context.Background(),10*time.Second)
	defer cncl()
	var err error
	Client,err=mongo.Connect(ctx,opts)
	if err!=nil{
		log.Fatal("Not able to connect",err)
	}
    err=Client.Ping(ctx,nil)
	if err!=nil{
		log.Fatal("Ping disconnected",err)
	}
     log.Println("Connection succeeded")
}



