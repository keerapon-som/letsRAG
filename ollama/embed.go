package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"letsrag/entities"
	"net/http"
)

type EmbedService struct {
	url string
}

func (o *Ollama) Embed() *EmbedService {
	return &EmbedService{
		url: fmt.Sprintf("%s/api/embed", o.HostURL),
	}
}

func (e *EmbedService) TextToVector(text string, modelName string) (entities.TestToVectorResponse, error) {
	payload := map[string]interface{}{
		"model": modelName,
		"input": text,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return entities.TestToVectorResponse{}, err
	}

	resp, err := http.Post(e.url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return entities.TestToVectorResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entities.TestToVectorResponse{}, fmt.Errorf("failed to pull model: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return entities.TestToVectorResponse{}, err
	}

	var response entities.TestToVectorResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return entities.TestToVectorResponse{}, err
	}

	return response, nil
}
