package api

import (
	"encoding/json"
	"fmt"
	"letsrag/entities"
	"letsrag/ollama"
	"letsrag/postgresql"
	"letsrag/repository"
	"letsrag/utils"
	"log"
	"testing"
)

func TestListAllModel(t *testing.T) {
	ollama := ollama.NewOllama("http://localhost:11434")
	models, err := ollama.ListLocalModels()
	if err != nil {
		log.Fatalf("Failed to list models: %v", err)
	}

	for _, model := range models {
		fmt.Println(model.Name)
	}
}

func TestDeleteAModel(t *testing.T) {
	ollama := ollama.NewOllama("http://localhost:11434")
	err := ollama.DeleteModel("qqq")
	if err != nil {
		log.Fatalf("Failed to delete model: %v", err)
	}
}

func TestPullAModel(t *testing.T) {
	ollama := ollama.NewOllama("http://localhost:11434")

	resCh := make(chan entities.PullAModelStatus)
	errorCh := make(chan struct{})
	err := ollama.PullModel("all-minilm", false).Stream(resCh, errorCh)
	if err != nil {
		log.Fatalf("Failed to list models: %v", err)
	}

	for res := range resCh {
		fmt.Println(res)
	}
}

func TestCreateTable(t *testing.T) {
	connStr := "postgres://yourusername:yourpassword@localhost:5432/yourdatabase?sslmode=disable"
	if err := postgresql.InitDB(connStr); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	err := repository.NewDocumentRepoPostgresql().CreateDocumentTable()
	if err != nil {
		log.Fatalf("Failed to create document table: %v", err)
	}
}

func TestImportData(t *testing.T) {
	connStr := "postgres://yourusername:yourpassword@localhost:5432/yourdatabase?sslmode=disable"
	if err := postgresql.InitDB(connStr); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	ollama := ollama.NewOllama("http://localhost:11434")

	lrag := NewLetsRag(
		NewTextToVector(ollama),
		repository.NewDocumentRepoPostgresql(),
		ollama,
	)

	bytes, err := utils.ReadFile("mockupDocument.json")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	var resp []string
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	for _, text := range resp {
		lrag.SaveDocumentToDB(text, MODEL_ALL_MINILM)
	}

}
