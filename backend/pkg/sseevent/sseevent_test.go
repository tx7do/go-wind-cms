package sseevent

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── 跨服务契约钉测 ──
//
// 事件名与 Redis 频道是后端发布方与前端 SSE 注册表的双边约定：
// 改名需两侧同步，已发布的事件名不可复用为其他语义。

func TestEventNamesPinned(t *testing.T) {
	assert.Equal(t, "notification", EventNotification)
	assert.Equal(t, "ai_agent_run", EventAgentRun)
}

func TestRedisChannelsPinned(t *testing.T) {
	// 频道名带版本号 v1；换载荷结构应升版本而非复用旧频道。
	assert.Equal(t, "go-wind-cms:notify:v1", RedisChannelNotify)
	assert.Equal(t, "go-wind-cms:agent-run:v1", RedisChannelAgentRun)
}

func TestNotifyEventJSONRoundTrip(t *testing.T) {
	in := NotifyEvent{UserID: 42, Title: "标题", Content: "正文"}
	b, err := json.Marshal(in)
	require.NoError(t, err)

	// 键名是跨服务 JSON 契约（snake_case），前端/服务端双侧依赖。
	assert.JSONEq(t, `{"user_id":42,"title":"标题","content":"正文"}`, string(b))

	var out NotifyEvent
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, in, out)
}

func TestAgentRunEventJSONRoundTrip(t *testing.T) {
	// Frame 是 ai.service.v1.AgentRunEvent 的 protojson 字节（camelCase 键），
	// 载荷对 frame 必须原样透传、不得转义破坏。
	frame := json.RawMessage(`{"runId":"r-1","nodeId":"node-2","phase":"PHASE_END","seq":3,"status":"STATUS_OK"}`)
	in := AgentRunEvent{UserID: 7, Frame: frame}
	b, err := json.Marshal(in)
	require.NoError(t, err)

	var out AgentRunEvent
	require.NoError(t, json.Unmarshal(b, &out))
	assert.Equal(t, int64(7), out.UserID)
	require.NotNil(t, out.Frame)
	assert.True(t, json.Valid(out.Frame))
	assert.JSONEq(t, string(frame), string(out.Frame))
}
