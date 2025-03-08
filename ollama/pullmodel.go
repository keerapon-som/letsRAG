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

type PullAModel struct {
	url       string
	modelName string
	insecure  bool
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
				if err == io.EOF {
					break
				}
				// Handle error (you might want to send it through the channel or log it)
				errorChan <- struct{}{}
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
