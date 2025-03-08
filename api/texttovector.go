package api

import (
	"fmt"
	"letsrag/ollama"
)

// mxbai-embed-large	334M	View model
// nomic-embed-text	137M	View model
// all-minilm	23M	View model

type TextToVector struct {
	ollamaService *ollama.Ollama
}

const (
	MODEL_MXBAI_EMBED_LARGE = "mxbai-embed-large"
	MODEL_NOMIC_EMBED_TEXT  = "nomic-embed-text"
	MODEL_ALL_MINILM        = "all-minilm"
)

func NewTextToVector(ollamaService *ollama.Ollama) *TextToVector {
	return &TextToVector{
		ollamaService: ollamaService,
	}
}

func (t *TextToVector) ConvertTextToVector(text string, model string) ([][]float64, error) {
	resp, err := t.ollamaService.Embed().TextToVector(text, model)
	if err != nil {
		return [][]float64{}, err
	}

	dimention := len(resp.Embeddings[0])
	fmt.Println("Dimention : ", dimention)
	return resp.Embeddings, nil
}
