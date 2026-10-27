package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"

	"go-wind-cms/app/core/service/internal/data"
	"go-wind-cms/pkg/constants"
)

// 评论 AI 反垃圾（v1，系统参数开关）：
//   - 开关：sys_configs 键 comment.ai_moderation.enabled（GetConfigBool，默认 false）；
//   - 判定：默认启用 provider（与 chat/agent 同源）零轮工具 chat completion，
//     要求模型输出 JSON {"spam":bool,"reason":"..."}，temperature=0；
//   - 语义：spam → 拒绝创建（400）；AI 任何错误 → fail-open 放行（可用性优先，
//     漏网垃圾由既有审核队列兜底）；provider 未配置同样放行。

// configKeyCommentAIModeration 反垃圾开关的系统参数键（内置种子见 constants.DefaultConfigs）。
const configKeyCommentAIModeration = constants.ConfigKeyCommentAIModeration

// commentAIModerationMaxChars 送审内容截断上限（长文截断省 token，判定语义足够）。
const commentAIModerationMaxChars = 2000

// CommentAIModerationDeps AI 反垃圾的外呼依赖（wiring 注入真实现；单测可注入桩）。
type CommentAIModerationDeps struct {
	ProviderRepo *data.AiProviderRepo
	UsageLogRepo *data.AiUsageLogRepo
}

// commentVerdict 模型输出的结构化判定。
type commentVerdict struct {
	Spam   bool   `json:"spam"`
	Reason string `json:"reason"`
}

// moderateCommentAI 反垃圾闸门。返回 error 即拒绝创建。
func (s *CommentService) moderateCommentAI(ctx context.Context, content string) error {
	if s.configRepo == nil || s.aiDeps == nil || s.aiDeps.ProviderRepo == nil {
		return nil
	}
	if !s.configRepo.GetConfigBool(ctx, configKeyCommentAIModeration, false) {
		return nil
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	if len(content) > commentAIModerationMaxChars {
		content = content[:commentAIModerationMaxChars]
	}

	verdict, err := s.runCommentModeration(ctx, content)
	if err != nil {
		// fail-open：AI 链路不可用不放死真人评论
		s.log.Errorf("comment ai moderation failed (fail-open): %s", err.Error())
		return nil
	}
	if verdict.Spam {
		s.log.Warnf("comment rejected by ai moderation: %s", verdict.Reason)
		return fmt.Errorf("comment rejected by moderation")
	}
	return nil
}

// runCommentModeration 调默认 provider 做零轮判定（温度 0、限输出长度），落用量流水。
func (s *CommentService) runCommentModeration(ctx context.Context, content string) (*commentVerdict, error) {
	provider, err := s.aiDeps.ProviderRepo.GetEnabledDefault(ctx)
	if err != nil {
		return nil, fmt.Errorf("no enabled ai provider: %w", err)
	}
	client, err := newOpenAIClientForProvider(provider)
	if err != nil {
		return nil, err
	}
	model := derefAgentStr(provider.ModelName)

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You are a comment spam classifier. Reply with ONLY a JSON object " +
					`{"spam": true|false, "reason": "<short reason>"}. ` +
					"Spam means: advertising, scam, SEO link dumping, gibberish flooding, or abusive harassment. " +
					"Normal on-topic discussion, questions and criticism are NOT spam.",
			},
			{Role: openai.ChatMessageRoleUser, Content: content},
		},
		Temperature: 0,
		MaxTokens:   100,
	})
	if err != nil {
		return nil, fmt.Errorf("moderation chat failed: %w", err)
	}

	// 用量流水（conversation_id=0 的非对话调用，与 agent/审查同口径；失败仅留日志）
	if s.aiDeps.UsageLogRepo != nil && resp.Usage.TotalTokens > 0 {
		if uerr := s.aiDeps.UsageLogRepo.Create(ctx, &data.UsageRecord{
			ProviderID:       provider.ID,
			ModelName:        model,
			PromptTokens:     uint32(resp.Usage.PromptTokens),
			CompletionTokens: uint32(resp.Usage.CompletionTokens),
			TotalTokens:      uint32(resp.Usage.TotalTokens),
		}); uerr != nil {
			_ = uerr // 记账失败不影响判定结论
		}
	}

	return parseCommentVerdict(resp.Choices[0].Message.Content)
}

func ptrFloat32(v float32) *float32 { return &v }

// parseCommentVerdict 从模型输出提取 JSON 判定（容忍 ```json 围栏与前后杂文）。
func parseCommentVerdict(text string) (*commentVerdict, error) {
	text = strings.TrimSpace(text)
	if begin := strings.Index(text, "{"); begin >= 0 {
		if end := strings.LastIndex(text, "}"); end > begin {
			text = text[begin : end+1]
		}
	}
	var v commentVerdict
	if err := json.Unmarshal([]byte(text), &v); err != nil {
		return nil, fmt.Errorf("moderation verdict is not valid JSON: %w", err)
	}
	return &v, nil
}
