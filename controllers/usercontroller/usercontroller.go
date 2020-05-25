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

//Controller : act as constructor
type Controller struct {
	Collection *mongo.Collection
}

//SignUpHandler : handle signup routes logic
func (con *Controller) SignUpHandler() gin.HandlerFunc {
	userCol := con.Collection

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var u model.User
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  "something went wrong",
			})
			return
		}

		var existingUser model.User
		err := userCol.FindOne(ctx, gin.H{"username": u.Username}).Decode(&existingUser)
		if existingUser.Username != "" {
			c.JSON(http.StatusConflict, gin.H{
				"status": "failed",
				"error":  "username already taken",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 11)
		if err != nil {
			log.Fatal(err)
		}

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

		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"data": gin.H{
				"newUserID": newUser,
			},
		})
	}
}

//SignInHandler : handle signin routes logic
func (con *Controller) SignInHandler() gin.HandlerFunc {
	userCol := con.Collection

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
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
				"error":  "user not found",
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
