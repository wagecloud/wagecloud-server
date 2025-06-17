package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/client/pgxpool"
	instancestorage "github.com/wagecloud/wagecloud-server/internal/modules/instance/storage"
	"github.com/wagecloud/wagecloud-server/internal/shared/pagination"
)

var storage *instancestorage.Storage

func setUpConfig() {
	config.SetConfig("../../config/config.dev.yml")
}

func init() {
	setUpConfig()

	read, _ := pgxpool.NewPgxpool(pgxpool.PgxpoolOptions{
		Url:             config.GetConfig().Postgres.Url,
		Host:            config.GetConfig().Postgres.Host,
		Port:            config.GetConfig().Postgres.Port,
		Username:        config.GetConfig().Postgres.Username,
		Password:        config.GetConfig().Postgres.Password,
		Database:        config.GetConfig().Postgres.Database,
		MaxConnections:  config.GetConfig().Postgres.MaxConnections,
		MaxConnIdleTime: config.GetConfig().Postgres.MaxConnIdleTime,
	})
	write, _ := pgxpool.NewPgxpool(pgxpool.PgxpoolOptions{
		Url:             config.GetConfig().PostgresWrite.Url,
		Host:            config.GetConfig().PostgresWrite.Host,
		Port:            config.GetConfig().PostgresWrite.Port,
		Username:        config.GetConfig().PostgresWrite.Username,
		Password:        config.GetConfig().PostgresWrite.Password,
		Database:        config.GetConfig().PostgresWrite.Database,
		MaxConnections:  config.GetConfig().PostgresWrite.MaxConnections,
		MaxConnIdleTime: config.GetConfig().PostgresWrite.MaxConnIdleTime,
	})

	db := pgxpool.NewDBRouter(read, write)
	_ = db
	storage = instancestorage.NewStorage(write)
}

func main() {
	fmt.Println("Starting benchmark")
	storage.ListInstances(context.Background(), instancestorage.ListInstancesParams{
		PaginationParams: pagination.PaginationParams{
			Page:  1,
			Limit: 10,
		},
	})
	totalRequests := 1000
	concurrency := 10
	ctx := context.Background()
	store := storage
	var wg sync.WaitGroup

	requests := make(chan struct{}, totalRequests)
	latencies := make(chan time.Duration, totalRequests)

	start := time.Now()

	// Spawn workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requests {
				startTime := time.Now()

				// _, err := store.CreateInstance(ctx, instancemodel.Instance{
				// 	ID:        fmt.Sprintf("test-%d", time.Now().UnixNano()),
				// 	AccountID: 1,
				// 	OSID:      "ubuntu",
				// 	ArchID:    "x86_64",
				// 	Name:      "test",
				// 	CPU:       2,
				// 	RAM:       2048,
				// 	Storage:   20,
				// })
				_, err := store.ListInstances(ctx, instancestorage.ListInstancesParams{
					PaginationParams: pagination.PaginationParams{
						Page:  int32(i) + 1,
						Limit: 10,
					},
				})
				if err != nil {
					log.Printf("Error: %v", err)
					continue
				}

				latency := time.Since(startTime)
				latencies <- latency
			}
		}()
	}

	// Feed requests
	for i := 0; i < totalRequests; i++ {
		requests <- struct{}{}
	}
	close(requests)

	wg.Wait()
	elapsed := time.Since(start)

	close(latencies)

	// Analyze latency
	var totalLatency time.Duration
	var min, max time.Duration
	min = time.Hour // arbitrarily large

	for l := range latencies {
		totalLatency += l
		if l < min {
			min = l
		}
		if l > max {
			max = l
		}
	}
	avgLatency := totalLatency / time.Duration(totalRequests)

	fmt.Println("Benchmark results:")
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("QPS: %.2f\n", float64(totalRequests)/elapsed.Seconds())
	fmt.Printf("Latency - Avg: %v | Min: %v | Max: %v\n", avgLatency, min, max)
}
