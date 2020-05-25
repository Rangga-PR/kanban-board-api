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
)

//Routes : define server available routes
func Routes(router *gin.Engine) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	db, ctx := dbconfig.ConnectDB(ctx, os.Getenv("MONGO_URI"))

	router.POST("/signup", usercontroller.SignUpHandler(ctx, cancel, db))
	router.POST("/signin", usercontroller.SignInHandler(ctx, cancel, db))
}
