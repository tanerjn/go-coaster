package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runEntrance(ctx, "North", 200, 3000)
	go runEntrance(ctx, "East", 100, 2000)
	go runEntrance(ctx, "South", 300, 4000)
	go runEntrance(ctx, "West", 400, 6000)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
}

func runEntrance(ctx context.Context, name string, minWaitTime, maxWaitTime int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(rand.Intn(maxWaitTime-minWaitTime)+minWaitTime) * time.Millisecond)

			data := map[string]interface{}{"entrance": name} // Use a map for flexibility

			d, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshaling JSON: %v\n", err)
				continue
			}

			response, err := http.Post("http://server-service:80", "application/json", bytes.NewBuffer(d))
			if err != nil {
				log.Printf("Error sending request: %v\n", err)
				// Consider retry logic here
				continue
			}
			defer response.Body.Close() // Ensure body is closed
		}
	}
}
