package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"letsrag/entities"
	"net/http"
)

type Ollama struct {
	HostURL string
}

type PullAModel struct {
	url       string
	modelName string
	insecure  bool
}

func NewOllama(hosturl string) *Ollama {
	return &Ollama{
		HostURL: hosturl,
	}
}

func (o *Ollama) PullModel(modelName string, insecure bool) *PullAModel {
	return &PullAModel{
		modelName: modelName,
		insecure:  insecure,
		url:       fmt.Sprintf("%s/api/pull", o.HostURL),
	}
}
func (p *PullAModel) Stream(ch chan entities.PullAModelStatus, errorChan chan struct{}) error {
	payload := map[string]interface{}{
		"model":    p.modelName,
		"stream":   true,
		"insecure": p.insecure,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(p.url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("Failed to pull model: %s\nResponse body: %s\n", resp.Status, bodyString)
		return fmt.Errorf("failed to pull model: %s", resp.Status)
	}

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		decoder := json.NewDecoder(resp.Body)
		for {
			var status entities.PullAModelStatus
			if err := decoder.Decode(&status); err != nil {
				fmt.Println("It's getting error ? ", err)
				if err == io.EOF {
					break
				}
				// Handle error (you might want to send it through the channel or log it)
				errorChan <- struct{}{}
				return
			}
			ch <- status
			if status.Status == "success" {
				break
			}
		}
	}()

	return nil
}

func (p *PullAModel) Normall(modelName string, insecure bool) (string, error) {

	payload := map[string]interface{}{
		"model":    modelName,
		"stream":   false,
		"insecure": insecure,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(p.url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to pull model: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
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
