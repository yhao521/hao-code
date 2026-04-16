package backend

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
