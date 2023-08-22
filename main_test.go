package main

import (
	"os"
	"testing"
)

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
