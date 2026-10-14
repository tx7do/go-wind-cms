// Package sseevent SSE 事件名与跨服务事件载荷的注册表。
//
// 事件名是后端发布方与前端 sse event 注册表的双边约定：改名需两侧同步，
// 已发布的事件名不可复用为其他语义。
package sseevent

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
