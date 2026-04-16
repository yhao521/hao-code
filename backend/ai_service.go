package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// AIConfig AI 服务配置
type AIConfig struct {
	APIKey    string `json:"apiKey"`
	BaseURL   string `json:"baseURL"`
	Model     string `json:"model"`
	MaxTokens int    `json:"maxTokens"`
}

// GhostTextRequest 行内建议请求
type GhostTextRequest struct {
	Prefix   string `json:"prefix"`
	Suffix   string `json:"suffix"`
	Language string `json:"language"`
	FilePath string `json:"filePath"`
}

// GhostTextResponse 行内建议响应
type GhostTextResponse struct {
	Text string `json:"text"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
	Context  string        `json:"context,omitempty"`
}

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Reply string `json:"reply"`
}

// AIService AI 服务
type AIService struct {
	config AIConfig
}

// NewAIService 创建 AI 服务
func NewAIService(config AIConfig) *AIService {
	return &AIService{config: config}
}

// UpdateConfig 更新配置
func (s *AIService) UpdateConfig(config AIConfig) {
	s.config = config
}

// GetGhostText 获取行内代码建议
func (s *AIService) GetGhostText(req GhostTextRequest) (*GhostTextResponse, error) {
	if s.config.APIKey == "" {
		return nil, fmt.Errorf("AI API Key not configured")
	}

	// 构建提示词
	prompt := fmt.Sprintf("Complete the following %s code:\n\n%s", req.Language, req.Prefix)

	// 构建请求体（以 OpenAI 兼容格式为例）
	payload := map[string]interface{}{
		"model": s.config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a professional code completion assistant."},
			{"role": "user", "content": prompt},
		},
		"max_tokens":  s.config.MaxTokens,
		"temperature": 0.2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// 发送请求
	url := fmt.Sprintf("%s/chat/completions", s.config.BaseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.APIKey))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI API error: %s", string(respBody))
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return &GhostTextResponse{Text: ""}, nil
	}

	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content := message["content"].(string)

	return &GhostTextResponse{Text: content}, nil
}

// ChatWithAI 与 AI 进行对话
func (s *AIService) ChatWithAI(req ChatRequest) (*ChatResponse, error) {
	if s.config.APIKey == "" {
		return nil, fmt.Errorf("AI API Key not configured")
	}

	// 构建提示词，加入上下文代码
	var promptBuilder strings.Builder
	if req.Context != "" {
		promptBuilder.WriteString(fmt.Sprintf("Code Context:\n```\n%s\n```\n\n", req.Context))
	}
	promptBuilder.WriteString("Conversation:\n")
	for _, msg := range req.Messages {
		promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}

	// 构建请求体
	payload := map[string]interface{}{
		"model": s.config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful coding assistant in Hao-Code editor."},
			{"role": "user", "content": promptBuilder.String()},
		},
		"max_tokens":  1024,
		"temperature": 0.7,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/chat/completions", s.config.BaseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.APIKey))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI API error: %s", string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return &ChatResponse{Reply: "No response from AI"}, nil
	}

	choice := choices[0].(map[string]interface{})
	message := choice["message"].(map[string]interface{})
	content := message["content"].(string)

	return &ChatResponse{Reply: content}, nil
}

// ChatWithAIStream 流式聊天（模拟实现，实际需配合 SSE）
func (s *AIService) ChatWithAIStream(req ChatRequest) (*ChatResponse, error) {
	// 目前 Wails v3 对原生 SSE 支持有限，我们先通过常规接口返回
	// 后续可以通过 WebSocket 实现真正的打字机效果
	return s.ChatWithAI(req)
}

// GetAIConfig 获取当前 AI 配置（脱敏）
func (s *AIService) GetAIConfig() AIConfig {
	cfg := s.config
	if cfg.APIKey != "" && len(cfg.APIKey) > 8 {
		cfg.APIKey = cfg.APIKey[:4] + "****" + cfg.APIKey[len(cfg.APIKey)-4:]
	}
	return cfg
}

// BuildContextFromReferences 根据文件引用构建上下文
func (s *AIService) BuildContextFromReferences(rootPath string, references []string) string {
	var contextBuilder strings.Builder
	for _, ref := range references {
		// 简单的安全检查，防止路径遍历
		if strings.Contains(ref, "..") {
			continue
		}
		fullPath := filepath.Join(rootPath, ref)
		content, err := os.ReadFile(fullPath)
		if err == nil {
			contextBuilder.WriteString(fmt.Sprintf("\n--- File: %s ---\n%s\n", ref, string(content)))
		}
	}
	return contextBuilder.String()
}
