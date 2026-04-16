package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// VSTheme VSCode 主题结构 (简化版)
type VSTheme struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"` // dark or light
	Colors      map[string]string `json:"colors"`
	TokenColors []TokenColor      `json:"tokenColors"`
}

// TokenColor VSCode Token 颜色定义
type TokenColor struct {
	Scope    interface{} `json:"scope"` // string or []string
	Settings struct {
		Foreground string `json:"foreground"`
		Background string `json:"background,omitempty"`
		FontStyle  string `json:"fontStyle,omitempty"`
	} `json:"settings"`
}

// ImportTheme 导入并解析 VSCode 主题文件
func ImportTheme(themePath string) (string, error) {
	data, err := os.ReadFile(themePath)
	if err != nil {
		return "", err
	}

	var vsTheme VSTheme
	if err := json.Unmarshal(data, &vsTheme); err != nil {
		return "", fmt.Errorf("invalid theme format: %v", err)
	}

	// 生成一个唯一的主题 ID
	themeID := strings.ToLower(strings.ReplaceAll(vsTheme.Name, " ", "-"))

	// 这里可以将主题保存到本地存储，为了演示，我们直接返回解析后的数据
	// 实际生产中应存入数据库或配置文件
	return themeID, nil
}

// GetThemeDefinition 获取 Monaco 兼容的主题定义
func GetThemeDefinition(themePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(themePath)
	if err != nil {
		return nil, err
	}

	var vsTheme VSTheme
	if err := json.Unmarshal(data, &vsTheme); err != nil {
		return nil, err
	}

	// 转换为 Monaco IStandaloneThemeData 格式
	rules := []map[string]interface{}{}
	for _, tc := range vsTheme.TokenColors {
		var scopes []string
		switch v := tc.Scope.(type) {
		case string:
			scopes = []string{v}
		case []interface{}:
			for _, s := range v {
				if str, ok := s.(string); ok {
					scopes = append(scopes, str)
				}
			}
		}

		for _, scope := range scopes {
			rule := map[string]interface{}{
				"token": scope,
			}
			if tc.Settings.Foreground != "" {
				rule["foreground"] = tc.Settings.Foreground
			}
			if tc.Settings.Background != "" {
				rule["background"] = tc.Settings.Background
			}
			if tc.Settings.FontStyle != "" {
				rule["fontStyle"] = tc.Settings.FontStyle
			}
			rules = append(rules, rule)
		}
	}

	// 提取基础颜色
	baseColors := make(map[string]interface{})
	if bg, ok := vsTheme.Colors["editor.background"]; ok {
		baseColors["editor.background"] = bg
	}
	if fg, ok := vsTheme.Colors["editor.foreground"]; ok {
		baseColors["editor.foreground"] = fg
	}

	return map[string]interface{}{
		"base":    vsTheme.Type, // "vs-dark" or "vs"
		"inherit": true,
		"rules":   rules,
		"colors":  baseColors,
	}, nil
}
