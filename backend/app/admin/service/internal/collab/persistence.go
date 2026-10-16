package collab

import (
	"context"
	"time"

	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// grpcPersistence ygo PersistenceAdapter 的 gRPC 实现：
// 状态存 core（collab_doc_states 表），admin 服务自身不触库。
// ygo 侧已对 update 做防抖合并，StoreUpdate 收到的是合并后的整量状态。
type grpcPersistence struct {
	client contentV1.CollabDocStateServiceClient
}

func newGRPCPersistence(client contentV1.CollabDocStateServiceClient) *grpcPersistence {
	return &grpcPersistence{client: client}
}

func (p *grpcPersistence) LoadDoc(room string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	st, err := p.client.GetDocState(ctx, &contentV1.GetCollabDocStateRequest{Room: room})
	if err != nil {
		return nil, err
	}
	return st.GetDocUpdate(), nil
}

func (p *grpcPersistence) StoreUpdate(room string, update []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := p.client.StoreDocState(ctx, &contentV1.StoreCollabDocStateRequest{
		Room:      room,
		DocUpdate: update,
	})
	return err
}
