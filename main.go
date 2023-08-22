package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	requestsDate    []time.Time
	windowStartTime time.Time
	dataFilePath    = "request_counter.txt"
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

func requestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Total requests in the last 60 seconds: %d\n", len(requestsDate))
}

func main() {

	loadRequestDataFromFile()

	windowStartTime = time.Now()

	http.HandleFunc("/", requestHandler)

	fmt.Println("Server is running on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
