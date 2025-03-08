package api

import "letsrag/ollama"

type AIService interface{}

type ai struct {
	ollamaService *ollama.Ollama
}

func NewAI(ollamaService *ollama.Ollama) *ai {
	return &ai{
		ollamaService: ollamaService,
	}
}
