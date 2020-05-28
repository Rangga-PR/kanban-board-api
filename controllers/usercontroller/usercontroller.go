package usercontroller

import (
	"context"
	"kanban-app-api/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//Controller : act as constructor
type Controller struct {
	Collection *mongo.Collection
}

func sendFailedResponse(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, gin.H{
		"status": "failed",
		"error":  msg,
	})
}

func sendSuccessResponse(c *gin.Context, statusCode int, data gin.H) {
	c.JSON(statusCode, gin.H{
		"status": "success",
		"result": data,
	})
}

//SignUpHandler : handle signup routes logic
func (con *Controller) SignUpHandler() gin.HandlerFunc {
	userCol := con.Collection

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var u model.User
		if err := c.ShouldBindJSON(&u); err != nil {
			sendFailedResponse(c, http.StatusBadRequest, "please fill username and password")
			return
		}

		var existingUser model.User
		err := userCol.FindOne(ctx, bson.M{"username": u.Username}).Decode(&existingUser)
		if existingUser.Username != "" {
			sendFailedResponse(c, http.StatusConflict, "username already taken")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 11)
		if err != nil {
			sendFailedResponse(c, http.StatusInternalServerError, "something went wrong")
			return
		}

		userData := model.User{
			Username:  u.Username,
			Password:  string(hashedPassword),
			CreatedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
		}

		newUser, err := userCol.InsertOne(ctx, userData)
		if err != nil {
			sendFailedResponse(c, http.StatusInternalServerError, "something went wrong")
			return
		}

		sendSuccessResponse(c, http.StatusCreated, gin.H{
			"new_user_id": newUser.InsertedID,
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
			sendFailedResponse(c, http.StatusBadRequest, "please fill username and password")
			return
		}

		var loginUser model.User
		err := userCol.FindOne(ctx, bson.M{"username": username}).Decode(&loginUser)
		if err != nil {
			sendFailedResponse(c, http.StatusNotFound, "user not found")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(password))
		if err != nil {
			sendFailedResponse(c, http.StatusNotFound, "incorrect password")
			return
		}

		sendSuccessResponse(c, http.StatusOK, gin.H{
			"id":       loginUser.ID,
			"username": loginUser.Username,
		})
	}
}
