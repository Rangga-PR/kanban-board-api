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
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  "something went wrong",
			})
			return
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
		if err != nil {
			log.Println(err)
		}

		tasks := []model.Task{}
		cursor, err := taskCol.Find(ctx, bson.M{"user_id": userid})
		if err != nil {
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

//DeleteTaskHandler : handle delete task route  logic
func (con *Controller) DeleteTaskHandler() gin.HandlerFunc {
	taskCol := con.Collection

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		taskid, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			log.Println()
		}

		deletedTask, err := taskCol.DeleteOne(ctx, bson.M{"_id": taskid})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  "something went wrong",
			})
			return
		}

		if deletedTask.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "failed",
				"error":  "no document with specified id",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"task_deleted": deletedTask.DeletedCount,
		})

	}
}

//UpdateTaskHandler : handle patch task route  logic
func (con *Controller) UpdateTaskHandler() gin.HandlerFunc {
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

		taskid, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			log.Println()
		}

		update := bson.D{
			primitive.E{
				Key: "$set", Value: bson.D{
					primitive.E{Key: "title", Value: t.Title},
					primitive.E{Key: "content", Value: t.Content},
				},
			},
		}

		updatedTask, err := taskCol.UpdateOne(
			ctx,
			bson.M{"_id": taskid},
			update,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  "something went wrong",
			})
			return
		}

		if updatedTask.MatchedCount < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "failed",
				"error":  "no document found with the specified id",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"task_updated": updatedTask.ModifiedCount,
		})
	}
}
