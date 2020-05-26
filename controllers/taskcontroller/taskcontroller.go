package taskcontroller

import (
	"context"
	"kanban-app-api/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//Controller : act as constructor
type Controller struct {
	Collection *mongo.Collection
}

//PostTaskHandler : handle post task routes logic
func (con *Controller) PostTaskHandler() gin.HandlerFunc {
	taskCol := con.Collection

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var t model.Task
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  "something went wrong",
			})
			return
		}

		taskData := model.Task{
			UserID:    t.UserID,
			Title:     t.Title,
			Content:   t.Content,
			Status:    t.Status,
			Icon:      t.Icon,
			CreatedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
		}

		newTask, err := taskCol.InsertOne(ctx, taskData)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "success",
			"data": gin.H{
				"newTaskID": newTask.InsertedID,
			},
		})
	}
}

//GetTaskHandler : handle get task route  logic
func (con *Controller) GetTaskHandler() gin.HandlerFunc {
	taskCol := con.Collection

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userid, err := primitive.ObjectIDFromHex(c.Param("id"))

		tasks := []model.Task{}
		cursor, err := taskCol.Find(ctx, bson.M{"user_id": userid})
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  "something went wrong",
			})
			return
		}

		for cursor.Next(ctx) {
			var t model.Task
			cursor.Decode(&t)
			tasks = append(tasks, t)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   tasks,
		})

	}
}
