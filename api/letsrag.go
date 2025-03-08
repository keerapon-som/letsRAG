package api

import (
	"fmt"
	"letsrag/entities"
	"letsrag/entities/db"
	"letsrag/ollama"
)

const (
	MODEL_RAG_QWEN  = "qwen2.5:3b"
	MODEL_RAG_LLAMA = "llama3.2:latest"
)

type DocumentRepo interface {
	SaveDocument(document string, vector []float64) error
	GetRelatedDocumentsByVector(vector []float64, limit int) ([]db.Document, error)
}

type TextToVectorService interface {
	ConvertTextToVector(text string, model string) ([][]float64, error)
}

type AIService interface {
	GenerateACompletion(req entities.GenerateACompletionRequest) *ollama.GenerateACompletion
}

type LetRags struct {
	textToVector TextToVectorService
	repo         DocumentRepo
	ollama       *ollama.Ollama
}

func NewLetsRag(textToVector TextToVectorService, repo DocumentRepo, ollama *ollama.Ollama) *LetRags {
	return &LetRags{
		textToVector: textToVector,
		repo:         repo,
		ollama:       ollama,
	}
}

func (l *LetRags) SaveDocumentToDB(text string, modelName string) error {
	vector, err := l.textToVector.ConvertTextToVector(text, modelName)
	if err != nil {
		return err
	}

	return l.repo.SaveDocument(text, vector[0])
}

func (l *LetRags) GetRelatedDocuments(text string, modelName string, limit int) ([]string, error) {
	vector, err := l.textToVector.ConvertTextToVector(text, modelName)
	if err != nil {
		return []string{}, err
	}

	docs, err := l.repo.GetRelatedDocumentsByVector(vector[0], limit)
	if err != nil {
		return []string{}, err
	}

	document := []string{}
	for _, doc := range docs {
		document = append(document, doc.Document)
	}

	return document, nil
}

type GenerateCompletionRAG struct {
	generateCompletionRequest *entities.GenerateACompletionRequest
	GenCompletion             *ollama.GenerateACompletion
}

func (l *LetRags) GenerateCompletionRAG(ask string, modelName string, vectorModel string) *GenerateCompletionRAG {
	// findDatabasefirst
	// turn promth to vector
	vectors, _ := l.textToVector.ConvertTextToVector(ask, vectorModel)

	docs, err := l.repo.GetRelatedDocumentsByVector(vectors[0], 2)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	documentList := ""
	for _, doc := range docs {
		documentList += doc.Document + "\n"
	}

	promth := "This is Document that you can get for answer \n" + documentList + "So please answer this queation " + ask
	fmt.Println("Promth : ", promth)

	genRequest := entities.GenerateACompletionRequest{
		Model:  modelName,
		Prompt: promth,
	}

	return &GenerateCompletionRAG{
		generateCompletionRequest: &genRequest,
		GenCompletion:             l.ollama.GenerateACompletion(genRequest),
	}
}

func (g *GenerateCompletionRAG) Stream(ch chan []byte) error {
	g.generateCompletionRequest.Stream = true
	g.GenCompletion.Stream(ch)
	return nil
}

// func (g *GenerateCompletionRAG) Normall() (entities.GenerateACompletionResponse, error) {
// 	return g.GenCompletion.Normall()
// }
