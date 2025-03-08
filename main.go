package main

import (
	"letsrag/api"
	"letsrag/ollama"
	"letsrag/postgresql"
	"letsrag/repository"
	"log"
)

func main() {
	connStr := "postgres://yourusername:yourpassword@localhost:5432/yourdatabase?sslmode=disable"
	if err := postgresql.InitDB(connStr); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	ollama := ollama.NewOllama("http://localhost:11434")
	lrag := api.NewLetsRag(
		api.NewAI(ollama),
		*api.NewTextToVector(ollama),
		repository.NewDocumentRepoPostgresql(),
	)
	docs, err := lrag.GetRelatedDocuments("Hi brother", api.MODEL_ALL_MINILM, 1)
	if err != nil {
		log.Fatalf("Failed to get related documents: %v", err)
	}
	log.Println(docs)
	postgresql.CloseDB()
}
