package social

import (
	"context"
	"time"
)

// InboxItem 一条收件箱条目（平台回拉的回复/提及等）。
type InboxItem struct {
	// ExternalID 平台侧条目唯一标识（去重键）。
	ExternalID string
	// Type 条目类型（REPLY/MENTION/LIKE/OTHER）。
	Type string
	// AuthorName 作者显示名。
	AuthorName string
	// AuthorHandle 作者账号标识。
	AuthorHandle string
	// Content 内容摘要（纯文本）。
	Content string
	// URL 平台侧链接。
	URL string
	// OccurredAt 平台侧发生时间。
	OccurredAt time.Time
}

// InboxFetcher 收件箱回拉能力（可选实现）：实现该接口的平台支持把
// 回复/提及回拉到站内收件箱。周期同步任务对未实现的平台静默跳过。
type InboxFetcher interface {
	FetchInbox(ctx context.Context, creds Credentials, since time.Time) ([]InboxItem, error)
}

// AsInboxFetcher 断言平台适配器是否具备收件箱回拉能力。
func AsInboxFetcher(p Provider) (InboxFetcher, bool) {
	f, ok := p.(InboxFetcher)
	return f, ok
}
