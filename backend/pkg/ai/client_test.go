package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx7do/go-wind-plugins/ai/openai"
)

// ── ModelType 常量对齐 proto 值域 ──

func TestModelTypeConstants(t *testing.T) {
	// proto ai.service.v1.AiProvider.ModelType：1=LOCAL，2=CLOUD。
	// 本包是本地别名，值漂移会让 switch 走 default 分支报 unsupported。
	assert.Equal(t, int32(1), int32(openai.ModelTypeLocal))
	assert.Equal(t, int32(2), int32(openai.ModelTypeCloud))
	assert.Equal(t, openai.ModelTypeLocal, ModelTypeLocal)
	assert.Equal(t, openai.ModelTypeCloud, ModelTypeCloud)
}

// ── NewClient 配置翻译 ──

func TestNewClientNilConfig(t *testing.T) {
	client, err := NewClient(nil, "sk-xxx")
	assert.Nil(t, client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestNewClientUnsupportedModelType(t *testing.T) {
	client, err := NewClient(&ClientConfig{ModelType: 99}, "")
	assert.Nil(t, client)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported ai model type")
}

func TestNewClientZeroModelType(t *testing.T) {
	// 0 既非 LOCAL 也非 CLOUD，必须显式拒绝而不是静默构造出不可用客户端。
	client, err := NewClient(&ClientConfig{ModelType: 0}, "")
	assert.Nil(t, client)
	require.Error(t, err)
}

func TestNewClientCloud(t *testing.T) {
	client, err := NewClient(&ClientConfig{
		ModelType:      int32(openai.ModelTypeCloud),
		ModelName:      "gpt-4o-mini",
		BaseUrl:        "https://api.openai.com/v1",
		Organization:   "org-1",
		TimeoutSeconds: 60,
	}, "sk-secret")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewClientLocal(t *testing.T) {
	// 本地 Ollama 无需密钥（上游以 "none" 占位），key 传空即可。
	client, err := NewClient(&ClientConfig{
		ModelType: int32(openai.ModelTypeLocal),
		ModelName: "qwen2.5:7b",
		LocalHost: "127.0.0.1",
		LocalPort: 11434,
	}, "")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewClientCloudEmptyApiKey(t *testing.T) {
	// 客户端构造本身允许空 key（错误推迟到调用时 401）——语义钉测：
	// 调用方须在进入本工厂前完成密钥解密与归属校验。
	client, err := NewClient(&ClientConfig{
		ModelType: int32(openai.ModelTypeCloud),
		BaseUrl:   "https://api.openai.com/v1",
	}, "")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
