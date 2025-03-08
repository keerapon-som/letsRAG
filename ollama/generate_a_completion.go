package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"letsrag/entities"
	"net/http"
)

type GenerateACompletion struct {
	payload entities.GenerateACompletionRequest
	url     string
}

func (o *Ollama) GenerateACompletion(req entities.GenerateACompletionRequest) *GenerateACompletion {
	return &GenerateACompletion{
		payload: req,
		url:     fmt.Sprintf("%s/api/generate", o.HostURL),
	}
}

func (g *GenerateACompletion) Stream(ch chan []byte) error {
	g.payload.Stream = true

	payloadBytes, err := json.Marshal(g.payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(g.url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("Failed to Generate Completion : %s\nResponse body: %s\n", resp.Status, bodyString)
		return fmt.Errorf("failed to Generate Completion: %s", resp.Status)
	}

	go func() {
		defer resp.Body.Close()
		defer close(ch)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			// var res Response
			ch <- scanner.Bytes()
		}
	}()

	return nil
}

func (g *GenerateACompletion) Normall() ([]entities.GenerateACompletionResponse, error) {
	g.payload.Stream = false

	testpayload := map[string]interface{}{
		"model":  g.payload.Model,
		"stream": false,
		"prompt": g.payload.Prompt,
		"format": "json",
	}

	payloadBytes, err := json.Marshal(testpayload)
	if err != nil {
		return []entities.GenerateACompletionResponse{}, err
	}

	resp, err := http.Post(g.url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return []entities.GenerateACompletionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []entities.GenerateACompletionResponse{}, fmt.Errorf("failed to pull model: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []entities.GenerateACompletionResponse{}, err
	}

	var response []entities.GenerateACompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return []entities.GenerateACompletionResponse{}, err
	}

	return response, nil
}
