package main

import (
	"31_5/pkg/user"
	"encoding/json"
	"log"
	"os"
)

type Storage struct {
	Users []*user.User
}

func main() {
	rawDataIn, err := os.ReadFile("users.json")
	if err != nil {
		log.Fatal("Cannot load storage:", err)
	}

	var storage Storage
	err = json.Unmarshal(rawDataIn, &storage)
	if err != nil {
		log.Fatal("Invalid storage format:", err)
	}

	newUser := user.User{
		Name:    "Маша",
		Age:     18,
		Friends: []string{"Игорь", "Толя"},
	}

	storage.Users = append(storage.Users, &newUser)

	rawDataOut, err := json.MarshalIndent(&storage, "", "  ")
	if err != nil {
		log.Fatal("JSON marshaling failed:", err)
	}

	err = os.WriteFile("users.json", rawDataOut, 0)
	if err != nil {
		log.Fatal("Cannot write updated storage file:", err)
	}
}
