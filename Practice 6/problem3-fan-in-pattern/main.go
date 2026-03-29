package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Metric struct {
	Name   string
	Value  int
	Server string
	Time   time.Time
}

func FanIn(ctx context.Context, channels ...<-chan Metric) <-chan Metric {
	out := make(chan Metric)
	var wg sync.WaitGroup

	// Helper function (collect metrics from a single channel)
	collect := func(ch <-chan Metric) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case metric, ok := <-ch:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- metric:
				}
			}
		}
	}

	for _, ch := range channels {
		wg.Add(1)
		go collect(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// startServer simulates a server sending metrics
func startServer(ctx context.Context, serverName string, metricsChan chan<- Metric) {
	defer close(metricsChan)

	metrics := []string{"cpu", "memory", "disk"}

	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			metric := Metric{
				Name:   metrics[rand.Intn(len(metrics))],
				Value:  rand.Intn(100),
				Server: serverName,
				Time:   time.Now(),
			}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
			metricsChan <- metric
		}
	}
}

func main() {
	fmt.Println("=== Fan-In Pattern: Metrics Aggregation System ===\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alphaChan := make(chan Metric)
	betaChan := make(chan Metric)
	gammaChan := make(chan Metric)

	fmt.Println("Starting 3 servers: Alpha, Beta, Gamma...\n")
	go startServer(ctx, "Alpha", alphaChan)
	go startServer(ctx, "Beta", betaChan)
	go startServer(ctx, "Gamma", gammaChan)

	mergedMetrics := FanIn(ctx, alphaChan, betaChan, gammaChan)

	fmt.Println("Aggregated Metrics:\n")
	metricsCount := 0
	serverMetrics := make(map[string]int)

	for metric := range mergedMetrics {
		fmt.Printf("[%s] %s: %d\n", metric.Server, metric.Name, metric.Value)
		metricsCount++
		serverMetrics[metric.Server]++
	}

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total metrics received: %d\n", metricsCount)
	for server, count := range serverMetrics {
		fmt.Printf("%s: %d metrics\n", server, count)
	}
}
