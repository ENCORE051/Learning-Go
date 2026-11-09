package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

func main() {
	person := Person{
		UserId:   "12345",
		Username: "JohnDoe",
		Age:      30,
		Address:  "123 Main St",
	}

	// jsonData, err := json.Marshal(person)
	jsonData, err := json.MarshalIndent(person, "", "\t")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("JSON Data:", string(jsonData))
}
