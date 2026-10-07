// Package ai 是 AI 模型客户端工厂：把 sys_ai_providers 表行提炼成可调用的
// OpenAI 兼容客户端。实现层复用上游 go-wind-plugins/ai/openai
// （sashabaranov/go-openai 的薄封装，云=OpenAI 兼容 API，本地=Ollama）。
//
// 这里只做"表行 → 客户端"的翻译与超时语义修正，不含任何业务逻辑；
// provider 的归属校验、密钥解密在 service 层完成后再进来。
package ai

import (
	"fmt"
	"net/http"

	"github.com/sashabaranov/go-openai"
	openaiPlugin "github.com/tx7do/go-wind-plugins/ai/openai"
)

// ClientConfig 从 AiProvider 表行提炼的连接配置。
// ModelType 取值与 proto ai.service.v1.AiProvider.ModelType 一致（1=LOCAL，2=CLOUD）。
type ClientConfig struct {
	ModelType      int32
	ModelName      string
	BaseUrl        string
	Organization   string
	LocalHost      string
	LocalPort      int32
	TimeoutSeconds int32
}

// ModelType 及两值常量：对上游插件类型的本地别名，调用方按 proto 值域（1/2）比较用。
type ModelType = openaiPlugin.ModelType

const (
	ModelTypeLocal = openaiPlugin.ModelTypeLocal
	ModelTypeCloud = openaiPlugin.ModelTypeCloud
)

// NewClient 构造 OpenAI 兼容客户端。
//
// 云端模型必须传解密后的 apiKey（上游按 DefaultConfig(apiKey) 构造，空 key
// 构造出的客户端调用必 401）；本地 Ollama 由上游填 "none" 占位，key 传空即可。
//
// 显式注入无超时的 http.Client 覆盖上游默认 30s Timeout——上游把超时挂在
// http.Client 上，对 CreateChatCompletionStream 意味着"整条流 30 秒必须收完"，
// 长回复会被腰斩。单轮流式对话的时长由调用方以 context deadline 控制。
func NewClient(cfg *ClientConfig, decryptedApiKey string) (*openai.Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("ai client config is nil")
	}

	pluginCfg := &openaiPlugin.Config{
		Type:      openaiPlugin.ModelType(cfg.ModelType),
		ModelName: cfg.ModelName,
	}
	if cfg.TimeoutSeconds > 0 {
		pluginCfg.TimeoutSeconds = cfg.TimeoutSeconds
	}
	switch openaiPlugin.ModelType(cfg.ModelType) {
	case openaiPlugin.ModelTypeCloud:
		pluginCfg.Cloud = &openaiPlugin.CloudConfig{
			ApiKey:       decryptedApiKey,
			BaseUrl:      cfg.BaseUrl,
			Organization: cfg.Organization,
		}
	case openaiPlugin.ModelTypeLocal:
		pluginCfg.Local = &openaiPlugin.LocalConfig{
			Host: cfg.LocalHost,
			Port: cfg.LocalPort,
		}
	default:
		return nil, fmt.Errorf("unsupported ai model type: %d", cfg.ModelType)
	}

	return openaiPlugin.NewClient(pluginCfg, openaiPlugin.WithHTTPClient(&http.Client{}))
}
