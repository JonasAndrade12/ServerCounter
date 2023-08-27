package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	requestsDate  []time.Time
	dataFilePath  = "request_counter.txt"
	requestsMutex sync.Mutex
)

func loadRequestDataFromFile() {
	file, err := os.OpenFile(dataFilePath, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	var timestamp int64
	for {
		_, err := fmt.Fscanf(file, "%d\n", &timestamp)
		if err != nil {
			break
		}
		requestsDate = append(requestsDate, time.Unix(timestamp, 0))
	}
}

func updateRequestCount() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		requestsMutex.Lock()

		now := time.Now()
		// Remove outdated requests
		for len(requestsDate) > 0 && now.Sub(requestsDate[0]) >= time.Minute {
			requestsDate = requestsDate[1:]
		}

		requestsMutex.Unlock()

		saveRequestDataToFile()
	}
}

func saveRequestDataToFile() {
	file, err := os.Create(dataFilePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	requestsMutex.Lock()
	defer requestsMutex.Unlock()

	for _, timestamp := range requestsDate {
		fmt.Fprintf(file, "%d\n", timestamp.Unix())
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	done := make(chan struct{}) // Channel to signal completion

	go func() {

		requestsMutex.Lock()
		defer requestsMutex.Unlock()

		fmt.Fprintf(w, "Total requests in the last 60 seconds: %d\n", len(requestsDate))

		requestsDate = append(requestsDate, time.Now())
		close(done)
	}()

	<-done
}

func main() {

	loadRequestDataFromFile()

	go updateRequestCount()

	http.HandleFunc("/", requestHandler)

	fmt.Println("Server is running on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
