// Package sseevent SSE 事件名与跨服务事件载荷的注册表。
//
// 事件名是后端发布方与前端 sse event 注册表的双边约定：改名需两侧同步，
// 已发布的事件名不可复用为其他语义。
package sseevent

import "encoding/json"

// EventNotification 站内信/系统通知事件（前端通知 Toast 消费）。
const EventNotification = "notification"

// RedisChannelNotify 系统通知事件的 Redis Pub/Sub 频道：
// core 服务（发布 worker 等后台链路）发出，admin 服务订阅后转 SSE 实时推送。
// 落库为准：频道丢失只影响实时性，不影响到达。
const RedisChannelNotify = "go-wind-cms:notify:v1"

// NotifyEvent 跨服务通知事件载荷（Redis 内 JSON 序列化）。
type NotifyEvent struct {
	UserID  uint32 `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// EventAgentRun AI Agent 试运行的逐节点轨迹事件（前端 SSE_EVENT.AgentRun 消费）。
// data 为 AgentRunEvent 的 protojson：camelCase 键（runId/appId/seq/phase/nodeId/summary/status）。
// 一次试运行产生多帧：每节点 node_start + node_end 各一帧，收尾 done 或 error 一帧。
// 尽力而为丢帧可容忍：Run 同步响应携带完整 trace 与终态，以响应为准。
const EventAgentRun = "ai_agent_run"

// RedisChannelAgentRun Agent 试运行事件的 Redis Pub/Sub 频道：
// core 服务（图执行引擎）发出，admin 服务订阅后按发起用户转 SSE 实时推送。
// 落库为准：频道丢失只影响实时观感，不影响同步响应。
const RedisChannelAgentRun = "go-wind-cms:agent-run:v1"

// AgentRunEvent 跨服务试运行事件载荷（Redis 内 JSON 序列化）：
// Frame 是 ai.service.v1.AgentRunEvent 的 protojson 字节。
type AgentRunEvent struct {
	UserID uint32          `json:"user_id"`
	Frame  json.RawMessage `json:"frame"`
}
