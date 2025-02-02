package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Read the JSON file
	file, err := os.ReadFile("dataset.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal JSON into a slice (for arrays)
	var data []interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Print the JSON data
	fmt.Println(data)
}
