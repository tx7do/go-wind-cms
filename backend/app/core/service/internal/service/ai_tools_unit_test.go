package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/require"
)

// TestExecAiTool 工具执行器：已知工具返回结果、未知工具返回错误（回传模型自纠）。
func TestExecAiTool(t *testing.T) {
	out, err := execAiTool(context.Background(), "get_current_time", "{}")
	require.NoError(t, err)
	require.NotEmpty(t, out)

	_, err = execAiTool(context.Background(), "no_such_tool", "{}")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no_such_tool")
}

// TestAppendToolCallDelta 流式工具调用分片按 Index 累积：name 与 arguments 可能
// 分属不同 chunk，空分片（仅 index）不破坏已有累积。
func TestAppendToolCallDelta(t *testing.T) {
	var acc []openai.ToolCall
	acc = appendToolCallDelta(acc, openai.ToolCall{Index: transInt(0), ID: "call_1", Type: "function", Function: openai.FunctionCall{Name: "get_time"}})
	acc = appendToolCallDelta(acc, openai.ToolCall{Index: transInt(0), Function: openai.FunctionCall{Arguments: "{\"tz\":"}})
	acc = appendToolCallDelta(acc, openai.ToolCall{Index: transInt(1), ID: "call_2", Function: openai.FunctionCall{Name: "other", Arguments: "{}"}})
	acc = appendToolCallDelta(acc, openai.ToolCall{Index: transInt(0), Function: openai.FunctionCall{Arguments: "\"utc\"}"}})

	require.Len(t, acc, 2)
	require.Equal(t, "get_time", acc[0].Function.Name)
	require.Equal(t, "{\"tz\":\"utc\"}", acc[0].Function.Arguments)
	require.Equal(t, "other", acc[1].Function.Name)
}

// TestAiBuiltinToolDefs 工具清单非空且 schema 合法 JSON。
func TestAiBuiltinToolDefs(t *testing.T) {
	defs := aiBuiltinToolDefs()
	require.NotEmpty(t, defs)
	for _, d := range defs {
		require.Equal(t, openai.ToolTypeFunction, d.Type)
		require.NotEmpty(t, d.Function.Name)
		var parsed map[string]any
		require.NoError(t, json.Unmarshal(d.Function.Parameters.(json.RawMessage), &parsed))
	}
}

func transInt(v int) *int { return &v }
