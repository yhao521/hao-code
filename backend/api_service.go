package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// APIRequest API 请求参数
type APIRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// APIResponse API 响应结果
type APIResponse struct {
	Status     int               `json:"status"`
	StatusText string            `json:"statusText"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Duration   int64             `json:"duration"` // 毫秒
}

// SendHTTPRequest 发送 HTTP 请求
func SendHTTPRequest(req APIRequest) (*APIResponse, error) {
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = bytes.NewBufferString(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	// 设置默认 Header
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}
	if _, ok := req.Headers["Content-Type"]; !ok && req.Body != "" {
		req.Headers["Content-Type"] = "application/json"
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	start := time.Now()

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	duration := time.Since(start).Milliseconds()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 尝试格式化 JSON
	var jsonBody interface{}
	bodyStr := string(bodyBytes)
	if json.Unmarshal(bodyBytes, &jsonBody) == nil {
		formatted, _ := json.MarshalIndent(jsonBody, "", "  ")
		bodyStr = string(formatted)
	}

	headersMap := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headersMap[k] = v[0]
		}
	}

	return &APIResponse{
		Status:     resp.StatusCode,
		StatusText: resp.Status,
		Headers:    headersMap,
		Body:       bodyStr,
		Duration:   duration,
	}, nil
}

// APIHistoryItem 历史记录项
type APIHistoryItem struct {
	ID        string            `json:"id"`
	Timestamp int64             `json:"timestamp"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
}

// SaveApiHistory 保存 API 请求到历史记录
func SaveApiHistory(req APIRequest) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	historyPath := filepath.Join(homeDir, ".hao-code", "api_history.json")

	// 读取现有历史
	var history []APIHistoryItem
	if data, err := os.ReadFile(historyPath); err == nil {
		json.Unmarshal(data, &history)
	}

	// 添加新记录
	newItem := APIHistoryItem{
		ID:        time.Now().Format("20060102150405"),
		Timestamp: time.Now().Unix(),
		Method:    req.Method,
		URL:       req.URL,
		Headers:   req.Headers,
		Body:      req.Body,
	}

	// 限制历史记录数量（例如最多 50 条）
	if len(history) > 50 {
		history = history[:50]
	}

	history = append([]APIHistoryItem{newItem}, history...)

	// 确保目录存在
	os.MkdirAll(filepath.Dir(historyPath), 0755)
	data, _ := json.MarshalIndent(history, "", "  ")
	return os.WriteFile(historyPath, data, 0644)
}

// GetApiHistory 获取 API 历史记录
func GetApiHistory() ([]APIHistoryItem, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	historyPath := filepath.Join(homeDir, ".hao-code", "api_history.json")

	var history []APIHistoryItem
	if data, err := os.ReadFile(historyPath); err == nil {
		json.Unmarshal(data, &history)
	}
	return history, nil
}

// DeleteApiHistory 删除指定的历史记录
func DeleteApiHistory(id string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}
	historyPath := filepath.Join(homeDir, ".hao-code", "api_history.json")

	var history []APIHistoryItem
	if data, err := os.ReadFile(historyPath); err == nil {
		json.Unmarshal(data, &history)
	}

	var newHistory []APIHistoryItem
	for _, item := range history {
		if item.ID != id {
			newHistory = append(newHistory, item)
		}
	}

	data, _ := json.MarshalIndent(newHistory, "", "  ")
	return os.WriteFile(historyPath, data, 0644)
}

// GetEnvVariables 获取环境变量配置
func GetEnvVariables() (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	envPath := filepath.Join(homeDir, ".hao-code", "api_env.json")

	vars := make(map[string]string)
	if data, err := os.ReadFile(envPath); err == nil {
		json.Unmarshal(data, &vars)
	}
	return vars, nil
}

// SaveEnvVariables 保存环境变量配置
func SaveEnvVariables(vars map[string]string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	envPath := filepath.Join(homeDir, ".hao-code", "api_env.json")

	os.MkdirAll(filepath.Dir(envPath), 0755)
	data, _ := json.MarshalIndent(vars, "", "  ")
	return os.WriteFile(envPath, data, 0644)
}
