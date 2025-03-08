package entities

import "time"

type PullAModelStatus struct {
	Status    string `json:"status"`
	Digest    string `json:"digest"`
	Total     int    `json:"total"`
	Completed int    `json:"completed"`
}

type ModelDetails struct {
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

type Model struct {
	Name       string       `json:"name"`
	ModifiedAt string       `json:"modified_at"`
	Size       int64        `json:"size"`
	Digest     string       `json:"digest"`
	Details    ModelDetails `json:"details"`
}

type ListLocalModelsResponse struct {
	Models []Model `json:"models"`
}

type TestToVectorResponse struct {
	Model           string      `json:"model"`
	Embeddings      [][]float64 `json:"embeddings"`
	TotalDuration   int         `json:"total_duration"`
	LoadDuration    int         `json:"load_duration"`
	PromptEvalCount int         `json:"prompt_eval_count"`
}

// type GenerateACompletionResponse struct {
// 	Model              string    `json:"model"`
// 	CreatedAt          time.Time `json:"created_at"`
// 	Response           string    `json:"response"`
// 	Done               bool      `json:"done"`
// 	Context            []int     `json:"context"`
// 	TotalDuration      int64     `json:"total_duration"`
// 	LoadDuration       int64     `json:"load_duration"`
// 	PromptEvalCount    int       `json:"prompt_eval_count"`
// 	PromptEvalDuration int       `json:"prompt_eval_duration"`
// 	EvalCount          int       `json:"eval_count"`
// 	EvalDuration       int64     `json:"eval_duration"`
// }

type GenerateACompletionResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         *string   `json:"done_reason"`
	Context            *[]int    `json:"context"`
	TotalDuration      *int      `json:"total_duration"`
	LoadDuration       *int      `json:"load_duration"`
	PromptEvalCount    *int      `json:"prompt_eval_count"`
	PromptEvalDuration *int      `json:"prompt_eval_duration"`
	EvalCount          *int      `json:"eval_count"`
	EvalDuration       *int      `json:"eval_duration"`
}

type GenerateACompletionRequest struct {
	Model     string                  `json:"model"`      // required
	Prompt    string                  `json:"prompt"`     // optional
	Suffix    *string                 `json:"suffix"`     // optional
	Images    *[]string               `json:"images"`     // optional, base64-encoded images
	Format    *interface{}            `json:"format"`     // optional, can be "json" or a JSON schema
	Options   *map[string]interface{} `json:"options"`    // optional, additional model parameters
	System    *string                 `json:"system"`     // optional, system message
	Template  *string                 `json:"template"`   // optional, prompt template
	Stream    bool                    `json:"stream"`     // optional, default is false
	Raw       *bool                   `json:"raw"`        // optional, default is false
	KeepAlive *string                 `json:"keep_alive"` // optional, default is "5m"
	Context   *string                 `json:"context"`    // deprecated
}
