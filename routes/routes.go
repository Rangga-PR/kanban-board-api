package routes

import (
	"context"
	dbconfig "kanban-app-api/config"
	"kanban-app-api/controllers/usercontroller"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx            context.Context
	cancel         context.CancelFunc
	db             *mongo.Database
	userController usercontroller.Controller
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db = dbconfig.ConnectDB(ctx, os.Getenv("MONGO_URI"))
	userController = usercontroller.Controller{Collection: db.Collection("user")}
}

//Routes : define server available routes
func Routes(router *gin.Engine) {

	router.POST("/signup", userController.SignUpHandler())
	router.POST("/signin", userController.SignInHandler())
}
