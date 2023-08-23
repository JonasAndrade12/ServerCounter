package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Run setup code before all tests
	setup()

	// Run the tests
	exitCode := m.Run()

	// Run teardown code after all tests
	teardown()

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func setup() {
	// Initialize/reset your shared resources here
	requestsDate = []time.Time{}
}

func teardown() {
	// Perform cleanup here if needed
}

func TestLoadRequestDataFromFile(t *testing.T) {

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_request_data.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write test data to the temporary file
	testData := "1630040000\n1630040010\n1630040020\n" // Example UNIX timestamps
	tmpFile.WriteString(testData)

	// Set the file name to the temporary file
	dataFilePath = tmpFile.Name()

	// Call the function being tested
	loadRequestDataFromFile()

	// Check if the loaded data matches the expected values
	expectedCount := 3
	if len(requestsDate) != expectedCount {
		t.Errorf("Expected request count: %d, but got: %d", expectedCount, len(requestsDate))
	}

	// Check if the loaded timestamps match the expected values
	expectedTimestamps := []int64{1630040000, 1630040010, 1630040020}
	for i, expected := range expectedTimestamps {
		actual := requestsDate[i].Unix()
		if actual != expected {
			t.Errorf("Expected timestamp: %d, but got: %d", expected, actual)
		}
	}
}

func TestUpdateRequestCount(t *testing.T) {

	// Clean data
	requestsDate = []time.Time{}

	// Start the updateRequestCount function in a goroutine
	go updateRequestCount()

	// Simulate time passing and request events
	for i := 0; i < 5; i++ {
		requestsMutex.Lock()
		requestsDate = append(requestsDate, time.Now())
		requestsMutex.Unlock()

		// Wait for a short time to simulate the passage of time
		time.Sleep(time.Second)
	}

	// Let the updateRequestCount goroutine run for a while
	time.Sleep(5 * time.Second)

	// Check if data is correct
	requestsMutex.Lock()
	if len(requestsDate) != 5 {
		t.Errorf("Expected requestsDate to be 5, but got length %d", len(requestsDate))
	}
	requestsMutex.Unlock()
}

func TestUpdateRequestCount_OutdatedData(t *testing.T) {

	// Clean data
	requestsDate = []time.Time{}

	// Start the updateRequestCount function in a goroutine
	go updateRequestCount()

	outdatedData := time.Now().Add(-time.Minute * 3)

	// Add request events
	requestsMutex.Lock()
	requestsDate = append(requestsDate, outdatedData)
	requestsMutex.Unlock()

	// Let the updateRequestCount goroutine run for a while
	time.Sleep(2 * time.Second)

	// Check if the outdated requests are removed
	requestsMutex.Lock()
	if len(requestsDate) != 0 {
		t.Errorf("Expected requestsDate to be empty, but got length %d", len(requestsDate))
	}
	requestsMutex.Unlock()
}

func TestSaveRequestDataToFile(t *testing.T) {

	// Clean data
	requestsDate = []time.Time{}

	// Set the file name to the temporary file
	dataFilePath = "test_request_data.txt"
	timeNow := time.Now()

	// Simulate time passing and request events
	requestsMutex.Lock()
	requestsDate = append(requestsDate, timeNow)
	requestsMutex.Unlock()

	// Save on temporary file
	saveRequestDataToFile()

	// Read the content of the temporary file
	savedData, err := os.OpenFile(dataFilePath, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("Error reading saved file: %v", err)
	}

	// Check if saved data matches the expected
	var timestamp int64
	var data = []time.Time{}
	for {
		_, err := fmt.Fscanf(savedData, "%d\n", &timestamp)
		if err != nil {
			break
		}
		data = append(data, time.Unix(timestamp, 0))
	}

	if len(data) != 1 {
		t.Errorf("Expected 1 lines in saved file, but got %d", len(data))
	}

	if data[0].Second() != timeNow.Second() {
		t.Errorf("Expected that the data in saved file are equal than the request")
	}
}
