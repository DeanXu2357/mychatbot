package llm

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/goccy/go-json"
)

const (
	Byte          = 1
	maxBufferSize = 512 * Byte * 1000
)

type Ollama struct {
	httpClient *http.Client

	base     *url.URL
	model    string
	system   string
	template string
	wordWrap bool
	format   string
}

func (o *Ollama) Generate(ctx context.Context, generateContext []int, prompt string) (string, []int, error) {
	req := GenerateRequest{
		Model:    o.model,
		Prompt:   prompt,
		System:   o.system,
		Template: o.template,
		Context:  generateContext,
		Options:  map[string]interface{}{},
	}

	path := "/api/generate"
	method := http.MethodPost

	var buf *bytes.Buffer
	bts, errB := json.Marshal(req)
	if errB != nil {
		return "", []int{}, fmt.Errorf("marshal data failed: %w", errB)
	}
	buf = bytes.NewBuffer(bts)

	requestURL := o.base.JoinPath(path)
	request, errR := http.NewRequestWithContext(ctx, method, requestURL.String(), buf)
	if errR != nil {
		return "", []int{}, fmt.Errorf("create request failed: %w", errR)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/x-ndjson")
	request.Header.Set("User-Agent", "api integration")

	response, errD := o.httpClient.Do(request)
	if errD != nil {
		return "", []int{}, fmt.Errorf("do request failed: %w", errD)
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	// increase the buffer size to avoid running out of space
	scanBuf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(scanBuf, maxBufferSize)

	var resp GenerateResponse
	var respString string
	// keep scanner parts for word wrap mode
	for scanner.Scan() {
		var errorResponse struct {
			Error string `json:"error,omitempty"`
		}

		bts := scanner.Bytes()
		if err := json.Unmarshal(bts, &errorResponse); err != nil {
			return "", []int{}, fmt.Errorf("unmarshal: %w", err)
		}

		if errorResponse.Error != "" {
			return "", []int{}, fmt.Errorf(errorResponse.Error)
		}

		if response.StatusCode >= http.StatusBadRequest {
			return "", []int{}, fmt.Errorf("status code: %d, status: %s, error message: %s", response.StatusCode, response.Status, errorResponse.Error)
		}

		if err := json.Unmarshal(bts, &resp); err != nil {
			return "", []int{}, fmt.Errorf("unmarshal response failed: %w", err)
		}

		respString += resp.Response
	}

	generateCTX := resp.Context

	return respString, generateCTX, nil
}

func NewOllama(baseURL string, system string) (*Ollama, error) {
	u, errU := url.Parse(baseURL)
	if errU != nil {
		return nil, fmt.Errorf("parse url failed: %w", errU)
	}

	return &Ollama{
		base:       u,
		model:      "llama2:13b",
		system:     system,
		httpClient: &http.Client{},
	}, nil
}

type GenerateRequest struct {
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	System   string `json:"system"`
	Template string `json:"template"`
	Context  []int  `json:"context,omitempty"`
	Stream   *bool  `json:"stream,omitempty"`
	Raw      bool   `json:"raw,omitempty"`
	Format   string `json:"format"`

	Options map[string]interface{} `json:"options"`
}

type GenerateResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`

	Done    bool  `json:"done"`
	Context []int `json:"context,omitempty"`

	TotalDuration      time.Duration `json:"total_duration,omitempty"`
	LoadDuration       time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration time.Duration `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       time.Duration `json:"eval_duration,omitempty"`
}
