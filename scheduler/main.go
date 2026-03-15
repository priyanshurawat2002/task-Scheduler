package main

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"distributed_scheduler/shared"
)

func main() {

	rdb := shared.NewRedisClient()

	router := gin.Default()

	router.POST("/task", func(c *gin.Context) {

		var payload struct {
			Data string `json:"data"`
		}

		if err := c.BindJSON(&payload); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		task := shared.Task{
			ID:      uuid.New().String(),
			Payload: payload.Data,
		}

		taskJSON, _ := json.Marshal(task)

		// select best worker
		workerID, err := selectBestWorker(rdb)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		queue := "worker_queue:" + workerID

		err = rdb.Client.LPush(shared.Ctx, queue, taskJSON).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"status":  "task assigned",
			"worker":  workerID,
			"task_id": task.ID,
		})
	})

	router.Run(":8080")
}
