package api

import (
	"encoding/json"
	"letsrag/ollama"
	"letsrag/postgresql"
	"letsrag/repository"
	"letsrag/utils"
	"log"
	"testing"
)

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
