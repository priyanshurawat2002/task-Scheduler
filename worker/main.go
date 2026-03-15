package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"distributed_scheduler/shared"
)

func main() {

	rdb := shared.NewRedisClient()

	workerID := uuid.New().String()

	fmt.Println("Worker started:", workerID)

	registerWorker(rdb, workerID)

	queue := "worker_queue:" + workerID

	for {

		result, err := rdb.Client.BRPop(shared.Ctx, 0, queue).Result()
		if err != nil {
			fmt.Println("Redis error:", err)
			continue
		}

		var task shared.Task

		err = json.Unmarshal([]byte(result[1]), &task)
		if err != nil {
			fmt.Println("JSON error:", err)
			continue
		}

		updateLoad(rdb, workerID, 1)

		fmt.Println("Processing task:", task.ID)

		time.Sleep(3 * time.Second)

		fmt.Println("Task completed:", task.ID)

		updateLoad(rdb, workerID, -1)
	}
}

func registerWorker(rdb *shared.RedisClient, id string) {

	worker := shared.Worker{
		ID:    id,
		Load:  0,
		Alive: true,
	}

	data, _ := json.Marshal(worker)

	rdb.Client.HSet(shared.Ctx, "workers", id, data)
}

func updateLoad(rdb *shared.RedisClient, id string, delta int) {

	data, err := rdb.Client.HGet(shared.Ctx, "workers", id).Result()
	if err != nil {
		fmt.Println("Redis error:", err)
		return
	}

	var worker shared.Worker

	json.Unmarshal([]byte(data), &worker)

	worker.Load += delta

	newData, _ := json.Marshal(worker)

	rdb.Client.HSet(shared.Ctx, "workers", id, newData)
}
