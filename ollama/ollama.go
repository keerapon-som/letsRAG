package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"letsrag/entities"
	"net/http"
)

type Ollama struct {
	HostURL string
}

func NewOllama(hosturl string) *Ollama {
	return &Ollama{
		HostURL: hosturl,
	}
}

func (p *Ollama) DeleteModel(modelName string) error {
	url := fmt.Sprintf("%s/api/delete", p.HostURL)
	payload := map[string]interface{}{
		"model": modelName,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to pull model: %s", resp.Status)
	}

	return nil
}

func (p *Ollama) ListLocalModels() ([]entities.Model, error) {
	url := fmt.Sprintf("%s/api/tags", p.HostURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list models: %s", resp.Status)
	}

	var response entities.ListLocalModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Models, nil
}

func (p *Ollama) ListRunningModels() ([]entities.PullAModelStatus, error) {
	url := fmt.Sprintf("%s/api/ps", p.HostURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list models: %s", resp.Status)
	}

	var response []entities.PullAModelStatus
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}
