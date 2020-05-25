package usercontroller

import (
	"context"
	"kanban-app-api/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//SignUpHandler : handle signup routes logic
func SignUpHandler(ctx context.Context, cancel context.CancelFunc, db *mongo.Database) gin.HandlerFunc {
	userCol := db.Collection("user")
	return func(c *gin.Context) {
		defer cancel()
		var u model.User
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 11)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(hashedPassword)

		userData := model.User{
			Username:  u.Username,
			Password:  string(hashedPassword),
			CreatedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
		}

		newUser, err := userCol.InsertOne(ctx, userData)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   newUser,
		})
	}
}

//SignInHandler : handle signin routes logic
func SignInHandler(ctx context.Context, cancel context.CancelFunc, db *mongo.Database) gin.HandlerFunc {
	userCol := db.Collection("user")
	return func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")

		if username == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  "please fill username and password",
			})
			return
		}

		var loginUser model.User
		err := userCol.FindOne(ctx, gin.H{"username": username}).Decode(&loginUser)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "failed",
				"error":  "user does not exist",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(password))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "failed",
				"error":  "incorrect password",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"id":       loginUser.ID,
				"username": loginUser.Username,
			},
		})
	}
}
