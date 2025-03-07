package main

import (
	"fmt"
	"letsrag/entities"
	"letsrag/ollama"
	"letsrag/postgresql"
	"log"
)

func main() {
	connStr := "user=new_user dbname=postgres sslmode=disable password=new_password host=localhost port=5432"
	if err := postgresql.InitDB(connStr); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := postgresql.DB().Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	ollama := ollama.NewOllama("http://localhost:11434")
	list, err := ollama.ListLocalModels()
	if err != nil {
		log.Fatalf("Failed to list local models: %v", err)
	}
	fmt.Println("List is ", list)
	// err = ollama.DeleteModel(list[0].Name)
	// if err != nil {
	// 	log.Fatalf("Failed to delete model: %v", err)
	// }

	receiveCh := make(chan entities.PullAModelStatus)
	closeCh := make(chan struct{})

	err = ollama.PullModel("llama3.2", false).Stream(receiveCh, closeCh)
	if err != nil {
		log.Fatalf("Failed to stream model: %v", err)
	}
	defer close(receiveCh)
	defer close(closeCh)

	for {
		select {
		case status := <-receiveCh:
			fmt.Println(status)
			if status.Status == "success" {
				return
			}
		case <-closeCh:
			fmt.Println("Stream closed")
			return

		}
	}

}
