package main

import (
	"context"
	"encoding/json"
	"fmt"

	"distributed_scheduler/shared"
)

func selectBestWorker(rdb *shared.RedisClient) (string, error) {

	ctx := context.Background()

	workers, err := rdb.Client.HGetAll(ctx, "workers").Result()
	if err != nil {
		return "", err
	}

	minLoad := 999999
	bestWorker := ""

	for id, data := range workers {

		var w shared.Worker
		json.Unmarshal([]byte(data), &w)

		if w.Load < minLoad && w.Alive {
			minLoad = w.Load
			bestWorker = id
		}
	}

	if bestWorker == "" {
		return "", fmt.Errorf("no workers available")
	}

	return bestWorker, nil
}
