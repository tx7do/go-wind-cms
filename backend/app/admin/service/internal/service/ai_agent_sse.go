package service

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/tx7do/go-utils/id"
	sse "github.com/tx7do/kratos-transport/transport/sse"
	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
	"go-wind-cms/pkg/sseevent"
)

// StartAgentRunSubscriber 订阅 core 服务经 Redis Pub/Sub 发出的 Agent 试运行
// 轨迹事件，按发起用户的全部在线管理端会话（accessToken 流）转 SSE 实时推送。
// 与系统通知的 StartNotifySubscriber 同模式：落库为准，频道丢失只影响实时观感。
func (s *AiAgentService) StartAgentRunSubscriber(redisClient *redis.Client, authenticationServiceClient authenticationV1.AuthenticationServiceClient) {
	if redisClient == nil {
		return
	}
	go func() {
		ctx := context.Background()
		sub := redisClient.Subscribe(ctx, sseevent.RedisChannelAgentRun)
		defer func() {
			_ = sub.Close()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-sub.Channel():
				var ev sseevent.AgentRunEvent
				if err := json.Unmarshal([]byte(m.Payload), &ev); err != nil {
					log.Errorf("agent run subscriber: unmarshal event failed: %v", err)
					continue
				}
				s.pushAgentRunFrames(ctx, authenticationServiceClient, ev)
			}
		}
	}()
}

// pushAgentRunFrames 把一帧轨迹推给目标用户的全部在线管理端会话。
func (s *AiAgentService) pushAgentRunFrames(ctx context.Context, authenticationServiceClient authenticationV1.AuthenticationServiceClient, ev sseevent.AgentRunEvent) {
	if s.agentRunPublisher == nil {
		return
	}
	streamIds, err := authenticationServiceClient.GetAccessTokens(ctx, &authenticationV1.GetAccessTokensRequest{
		UserId:     ev.UserID,
		ClientType: authenticationV1.ClientType_admin,
	})
	if err != nil {
		log.Errorf("agent run subscriber: get access tokens failed (user=%d): %v", ev.UserID, err)
		return
	}
	for _, streamId := range streamIds.AccessTokens {
		s.agentRunPublisher.Publish(ctx, sse.StreamID(streamId), &sse.Event{
			ID:    []byte(id.NewGUIDv4(false)),
			Data:  ev.Frame,
			Event: []byte(sseevent.EventAgentRun),
		})
	}
}
