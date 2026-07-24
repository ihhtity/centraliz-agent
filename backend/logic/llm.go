package logic

import (
	"bufio"
	"bytes"
	"centraliz-backend/pkg/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LLMClient struct {
	provider    string
	apiKey      string
	model       string
	baseURL     string
	maxTokens   int
	temperature float64
	timeout     int
	maxRetries  int
}

type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type LLMResponse struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Thought   string     `json:"thought,omitempty"`
}

type StreamHandler func(content string, toolCalls []ToolCall, done bool) error

var llmClient *LLMClient

func InitLLMClient() {
	cfg := config.AppConfig.LLM
	if cfg == nil {
		cfg = &config.LLMConfig{
			Provider:    "openai",
			APIKey:      "",
			Model:       "gpt-4o-mini",
			BaseURL:     "https://api.openai.com/v1",
			MaxTokens:   4096,
			Temperature: 0.7,
			Timeout:     60,
			MaxRetries:  3,
		}
	}

	llmClient = &LLMClient{
		provider:    cfg.Provider,
		apiKey:      cfg.APIKey,
		model:       cfg.Model,
		baseURL:     cfg.BaseURL,
		maxTokens:   cfg.MaxTokens,
		temperature: cfg.Temperature,
		timeout:     cfg.Timeout,
		maxRetries:  cfg.MaxRetries,
	}
}

func GetLLMClient() *LLMClient {
	if llmClient == nil {
		InitLLMClient()
	}
	return llmClient
}

func (c *LLMClient) ChatStream(messages []ChatMessage, tools []ToolDefinition, handler StreamHandler) error {
	if c.apiKey == "" {
		handler("抱歉，LLM 服务未配置。请联系管理员设置 API Key。", nil, true)
		return nil
	}

	type RequestMessage struct {
		Role       string                   `json:"role"`
		Content    string                   `json:"content"`
		ToolCalls  []map[string]interface{} `json:"tool_calls,omitempty"`
		ToolCallID string                   `json:"tool_call_id,omitempty"`
	}

	type Tool struct {
		Type     string                 `json:"type"`
		Function map[string]interface{} `json:"function"`
	}

	type Request struct {
		Model       string           `json:"model"`
		Messages    []RequestMessage `json:"messages"`
		Tools       []Tool           `json:"tools,omitempty"`
		ToolChoice  string           `json:"tool_choice,omitempty"`
		Stream      bool             `json:"stream"`
		Temperature float64          `json:"temperature"`
		MaxTokens   int              `json:"max_tokens"`
	}

	reqMessages := make([]RequestMessage, len(messages))
	for i, msg := range messages {
		reqMessages[i] = RequestMessage{
			Role:       msg.Role,
			Content:    msg.Content,
			ToolCallID: msg.ToolCallID,
		}
		if len(msg.ToolCalls) > 0 {
			reqMessages[i].ToolCalls = make([]map[string]interface{}, len(msg.ToolCalls))
			for j, tc := range msg.ToolCalls {
				reqMessages[i].ToolCalls[j] = map[string]interface{}{
					"id":   tc.ID,
					"type": "function",
					"function": map[string]interface{}{
						"name":      tc.ToolName,
						"arguments": tc.Parameters,
					},
				}
			}
		}
	}

	req := Request{
		Model:       c.model,
		Messages:    reqMessages,
		Stream:      true,
		Temperature: c.temperature,
		MaxTokens:   c.maxTokens,
	}

	if len(tools) > 0 {
		req.Tools = make([]Tool, len(tools))
		for i, tool := range tools {
			req.Tools[i] = Tool{
				Type: "function",
				Function: map[string]interface{}{
					"name":        tool.Name,
					"description": tool.Description,
					"parameters":  tool.Parameters,
				},
			}
		}
		req.ToolChoice = "auto"
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request error: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", c.baseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{
		Timeout: time.Duration(c.timeout) * time.Second,
	}

	var resp *http.Response
	for retry := 0; retry < c.maxRetries; retry++ {
		resp, err = client.Do(httpReq)
		if err != nil {
			if retry < c.maxRetries-1 {
				time.Sleep(time.Duration(retry+1) * 2 * time.Second)
				continue
			}
			return fmt.Errorf("http request error: %w", err)
		}

		if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			if retry < c.maxRetries-1 {
				time.Sleep(time.Duration(retry+1) * 3 * time.Second)
				continue
			}
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("LLM API error: %s, body: %s", resp.Status, string(body))
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("LLM API error: %s, body: %s", resp.Status, string(body))
		}
		break
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	var accumulatedToolCalls []ToolCall
	var currentToolCallIndex int

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read stream error: %w", err)
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		type ToolCallDelta struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Function struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments"`
			} `json:"function"`
		}

		type Delta struct {
			Content   string          `json:"content"`
			ToolCalls []ToolCallDelta `json:"tool_calls,omitempty"`
		}

		type Choice struct {
			Delta Delta `json:"delta"`
		}

		type StreamResponse struct {
			Choices []Choice `json:"choices"`
		}

		var streamResp StreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		if len(streamResp.Choices) > 0 {
			delta := streamResp.Choices[0].Delta

			if delta.Content != "" {
				if err := handler(delta.Content, nil, false); err != nil {
					return err
				}
			}

			if len(delta.ToolCalls) > 0 {
				for _, tcDelta := range delta.ToolCalls {
					if tcDelta.ID != "" {
						currentToolCallIndex = len(accumulatedToolCalls)
						accumulatedToolCalls = append(accumulatedToolCalls, ToolCall{
							ID:         tcDelta.ID,
							ToolName:   tcDelta.Function.Name,
							Parameters: make(map[string]interface{}),
						})
					}

					if currentToolCallIndex < len(accumulatedToolCalls) {
						if tcDelta.Function.Name != "" {
							accumulatedToolCalls[currentToolCallIndex].ToolName = tcDelta.Function.Name
						}
						if tcDelta.Function.Arguments != nil {
							for k, v := range tcDelta.Function.Arguments {
								accumulatedToolCalls[currentToolCallIndex].Parameters[k] = v
							}
						}
					}
				}
			}
		}
	}

	if len(accumulatedToolCalls) > 0 {
		if err := handler("", accumulatedToolCalls, false); err != nil {
			return err
		}
	}

	handler("", nil, true)
	return nil
}

func (c *LLMClient) Chat(messages []ChatMessage, tools []ToolDefinition) (*LLMResponse, error) {
	if c.apiKey == "" {
		return &LLMResponse{
			Content: "抱歉，LLM 服务未配置。请联系管理员设置 API Key。",
		}, nil
	}

	type RequestMessage struct {
		Role       string                   `json:"role"`
		Content    string                   `json:"content"`
		ToolCalls  []map[string]interface{} `json:"tool_calls,omitempty"`
		ToolCallID string                   `json:"tool_call_id,omitempty"`
	}

	type Tool struct {
		Type     string                 `json:"type"`
		Function map[string]interface{} `json:"function"`
	}

	type Request struct {
		Model       string           `json:"model"`
		Messages    []RequestMessage `json:"messages"`
		Tools       []Tool           `json:"tools,omitempty"`
		ToolChoice  string           `json:"tool_choice,omitempty"`
		Stream      bool             `json:"stream"`
		Temperature float64          `json:"temperature"`
		MaxTokens   int              `json:"max_tokens"`
	}

	reqMessages := make([]RequestMessage, len(messages))
	for i, msg := range messages {
		reqMessages[i] = RequestMessage{
			Role:       msg.Role,
			Content:    msg.Content,
			ToolCallID: msg.ToolCallID,
		}
		if len(msg.ToolCalls) > 0 {
			reqMessages[i].ToolCalls = make([]map[string]interface{}, len(msg.ToolCalls))
			for j, tc := range msg.ToolCalls {
				reqMessages[i].ToolCalls[j] = map[string]interface{}{
					"id":   tc.ID,
					"type": "function",
					"function": map[string]interface{}{
						"name":      tc.ToolName,
						"arguments": tc.Parameters,
					},
				}
			}
		}
	}

	req := Request{
		Model:       c.model,
		Messages:    reqMessages,
		Stream:      false,
		Temperature: c.temperature,
		MaxTokens:   c.maxTokens,
	}

	if len(tools) > 0 {
		req.Tools = make([]Tool, len(tools))
		for i, tool := range tools {
			req.Tools[i] = Tool{
				Type: "function",
				Function: map[string]interface{}{
					"name":        tool.Name,
					"description": tool.Description,
					"parameters":  tool.Parameters,
				},
			}
		}
		req.ToolChoice = "auto"
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request error: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", c.baseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{
		Timeout: time.Duration(c.timeout) * time.Second,
	}

	var resp *http.Response
	for retry := 0; retry < c.maxRetries; retry++ {
		resp, err = client.Do(httpReq)
		if err != nil {
			if retry < c.maxRetries-1 {
				time.Sleep(time.Duration(retry+1) * 2 * time.Second)
				continue
			}
			return nil, fmt.Errorf("http request error: %w", err)
		}

		if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			if retry < c.maxRetries-1 {
				time.Sleep(time.Duration(retry+1) * 3 * time.Second)
				continue
			}
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("LLM API error: %s, body: %s", resp.Status, string(body))
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("LLM API error: %s, body: %s", resp.Status, string(body))
		}
		break
	}

	defer resp.Body.Close()

	type APIToolCall struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Function struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		} `json:"function"`
	}

	type APIMessage struct {
		Content   string        `json:"content"`
		ToolCalls []APIToolCall `json:"tool_calls,omitempty"`
	}

	type APIChoice struct {
		Message APIMessage `json:"message"`
	}

	type APIResponse struct {
		Choices []APIChoice `json:"choices"`
	}

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	result := &LLMResponse{}
	if len(response.Choices) > 0 {
		result.Content = response.Choices[0].Message.Content

		if len(response.Choices[0].Message.ToolCalls) > 0 {
			result.ToolCalls = make([]ToolCall, len(response.Choices[0].Message.ToolCalls))
			for i, tc := range response.Choices[0].Message.ToolCalls {
				result.ToolCalls[i] = ToolCall{
					ID:         tc.ID,
					ToolName:   tc.Function.Name,
					Parameters: tc.Function.Arguments,
				}
			}
		}
	}

	return result, nil
}
