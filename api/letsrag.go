package api

import (
	"fmt"
	"letsrag/entities"
	"letsrag/entities/db"
	"letsrag/ollama"
	"strconv"
)

const (
	MODEL_RAG_QWEN  = "qwen2.5:3b"
	MODEL_RAG_LLAMA = "llama3.2:latest"
)

type DocumentRepo interface {
	SaveDocument(document string, vector []float64, vectorModelName string) error
	GetRelatedDocumentsByVector(vector []float64, modelName string, limit int) ([]db.Document, error)
}

type TextToVectorService interface {
	ConvertTextToVector(text string, model string) ([][]float64, error)
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

	return l.repo.SaveDocument(text, vector[0], modelName)
}

func (l *LetRags) GetRelatedDocuments(text string, modelName string, limit int) ([]string, error) {
	vector, err := l.textToVector.ConvertTextToVector(text, modelName)
	if err != nil {
		return []string{}, err
	}

	docs, err := l.repo.GetRelatedDocumentsByVector(vector[0], modelName, limit)
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
	GenCompletion *ollama.GenerateACompletion
}

func (l *LetRags) GenerateCompletionRAG(ask string, modelName string, vectorModel string, numberOfDocumentRef int) *GenerateCompletionRAG {
	// findDatabasefirst
	// turn promth to vector
	vectors, _ := l.textToVector.ConvertTextToVector(ask, vectorModel)

	docs, err := l.repo.GetRelatedDocumentsByVector(vectors[0], vectorModel, numberOfDocumentRef)
	if err != nil {
		fmt.Println("Error : ", err)
	}

	documentList := "\n"
	for i, doc := range docs {
		documentList += "Doc No." + strconv.Itoa(i+1) + " " + doc.Document + "\n"
	}

	promth := "\nThis is Document that you can get for answer \n" + documentList + "\nSo please answer this queation " + ask
	fmt.Println("--------------------------------------------------------------\nคำถาม : ", promth)

	genRequest := entities.GenerateACompletionRequest{
		Model:  modelName,
		Prompt: promth,
	}

	return &GenerateCompletionRAG{
		GenCompletion: l.ollama.GenerateACompletion(genRequest),
	}
}

func (g *GenerateCompletionRAG) Stream(ch chan []byte) error {
	g.GenCompletion.Stream(ch)
	return nil
}

// func (g *GenerateCompletionRAG) Normall() (entities.GenerateACompletionResponse, error) {
// 	return g.GenCompletion.Normall()
// }
