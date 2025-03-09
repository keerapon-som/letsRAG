package main

import (
	"encoding/json"
	"fmt"
	"letsrag/api"
	"letsrag/entities"
	"letsrag/ollama"
	"letsrag/postgresql"
	"letsrag/repository"
	"log"
)

type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func init() {
	connStr := "postgres://yourusername:yourpassword@localhost:5432/yourdatabase?sslmode=disable"
	if err := postgresql.InitDB(connStr); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// err := repository.NewDocumentRepoPostgresql().CreateDocumentTable()
	// if err != nil {
	// 	log.Fatalf("Failed to create document table: %v", err)
	// }

}

func main() {

	ollama := ollama.NewOllama("http://localhost:11434")

	lrag := api.NewLetsRag(
		api.NewTextToVector(ollama),
		repository.NewDocumentRepoPostgresql(),
		ollama,
	)

	resCh := make(chan []byte)

	err := lrag.GenerateCompletionRAG("Do you know Doggo ?", api.MODEL_RAG_LLAMA, api.MODEL_ALL_MINILM, 2).Stream(resCh)
	if err != nil {
		log.Fatalf("Failed to generate completion: %v", err)
	}

	result := ""
	for res := range resCh {
		// fmt.Println(string(res))
		var respJsonStruct entities.GenerateACompletionResponse
		if err := json.Unmarshal(res, &respJsonStruct); err != nil {
			log.Fatalf("Failed to unmarshal response: %v", err)
		}
		if respJsonStruct.Done {
			fmt.Println(result)
			return
		}
		result += respJsonStruct.Response
	}

}
