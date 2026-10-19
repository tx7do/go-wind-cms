package channel

// settings 加密/脱敏：schema 感知的 settings 规范化逻辑全部收在这里——
// repo 把 settings 当不透明字符串存取（channel 包 import data 单向，无循环依赖）。

import (
	"encoding/json"
	"fmt"
	"strings"

	pkgcrypto "go-wind-cms/pkg/crypto"
)

// NormalizeSettingsForStore 保存前规范化：JSON 解析 → required 校验 → sensitive 字段加密。
// 返回可直接落 settings 列的 JSON 文本。
func NormalizeSettingsForStore(typeName, plaintextJSON string) (string, error) {
	schema := SettingsSchemaByType(typeName)
	if schema == nil {
		return "", fmt.Errorf("unknown settings channel type %q", typeName)
	}
	if plaintextJSON == "" {
		return "", fmt.Errorf("settings is empty")
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(plaintextJSON), &raw); err != nil {
		return "", fmt.Errorf("settings is not valid JSON: %w", err)
	}
	m := settingsMap(raw)

	for _, f := range schema.Fields {
		if f.Required && strings.TrimSpace(m.str(f.Name)) == "" {
			return "", fmt.Errorf("%s is required", f.Label)
		}
		if f.Sensitive {
			if plain := m.str(f.Name); plain != "" {
				encrypted, encErr := pkgcrypto.EncryptIfNeeded(plain)
				if encErr != nil {
					return "", fmt.Errorf("encrypt %s failed: %w", f.Name, encErr)
				}
				raw[f.Name] = encrypted
			}
		}
	}

	b, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("encode settings failed: %w", err)
	}
	return string(b), nil
}

// MergeSettingsForStore 更新前合并：sensitive 字段新值为空 = 保留已存密文；
// 其余字段以新值为准。旧值为空且新值也为空的 required 字段照常校验报错。
func MergeSettingsForStore(typeName, oldStoreJSON, plaintextJSON string) (string, error) {
	schema := SettingsSchemaByType(typeName)
	if schema == nil {
		return "", fmt.Errorf("unknown settings channel type %q", typeName)
	}
	if plaintextJSON == "" {
		return "", fmt.Errorf("settings is empty")
	}

	var incoming map[string]any
	if err := json.Unmarshal([]byte(plaintextJSON), &incoming); err != nil {
		return "", fmt.Errorf("settings is not valid JSON: %w", err)
	}
	in := settingsMap(incoming)

	var oldRaw map[string]any
	if oldStoreJSON != "" {
		_ = json.Unmarshal([]byte(oldStoreJSON), &oldRaw) // 旧值坏了按缺省处理，新值说了算
	}
	old := settingsMap(oldRaw)

	for _, f := range schema.Fields {
		incomingValue := strings.TrimSpace(in.str(f.Name))
		switch {
		case f.Sensitive && incomingValue == "":
			// 敏感字段留空 = 保留已存密文；连旧值也没有则补空串占位（required 校验兜底）
			if oldPlain := old.str(f.Name); oldPlain != "" {
				incoming[f.Name] = oldPlain
			} else {
				delete(incoming, f.Name)
			}
		case f.Sensitive:
			encrypted, encErr := pkgcrypto.EncryptIfNeeded(incomingValue)
			if encErr != nil {
				return "", fmt.Errorf("encrypt %s failed: %w", f.Name, encErr)
			}
			incoming[f.Name] = encrypted
		}
		if f.Required && strings.TrimSpace(settingsMap(incoming).str(f.Name)) == "" {
			return "", fmt.Errorf("%s is required", f.Label)
		}
	}

	b, err := json.Marshal(incoming)
	if err != nil {
		return "", fmt.Errorf("encode settings failed: %w", err)
	}
	return string(b), nil
}

// DecryptSettings 存储态 JSON → 解密后的 settings map（sender 用）。
func DecryptSettings(typeName, storeJSON string) (map[string]any, error) {
	var raw map[string]any
	if storeJSON == "" {
		return map[string]any{}, nil
	}
	if err := json.Unmarshal([]byte(storeJSON), &raw); err != nil {
		return nil, fmt.Errorf("settings is not valid JSON: %w", err)
	}
	schema := SettingsSchemaByType(typeName)
	if schema == nil {
		return nil, fmt.Errorf("unknown settings channel type %q", typeName)
	}
	m := settingsMap(raw)
	for _, f := range schema.Fields {
		if !f.Sensitive {
			continue
		}
		if stored := m.str(f.Name); stored != "" {
			plain, err := pkgcrypto.DecryptIfNeeded(stored)
			if err != nil {
				return nil, fmt.Errorf("decrypt %s failed: %w", f.Name, err)
			}
			raw[f.Name] = plain
		}
	}
	return raw, nil
}

// SanitizeSettingsForDTO 读回显脱敏：sensitive 字段值恒为空串（前端"留空 = 不修改"）。
func SanitizeSettingsForDTO(typeName, storeJSON string) string {
	if storeJSON == "" {
		return ""
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(storeJSON), &raw); err != nil {
		return "" // 存储态坏了不阻塞列表渲染，M2 可加告警
	}
	if schema := SettingsSchemaByType(typeName); schema != nil {
		for _, f := range schema.Fields {
			if f.Sensitive {
				if _, ok := raw[f.Name]; ok {
					raw[f.Name] = ""
				}
			}
		}
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return ""
	}
	return string(b)
}
