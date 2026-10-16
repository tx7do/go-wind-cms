package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"google.golang.org/protobuf/types/known/timestamppb"

	entCrud "github.com/tx7do/go-crud/entgo"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/collabdocstate"
	contentV1 "go-wind-cms/api/gen/go/content/service/v1"
)

// CollabDocStateRepo 协同文档状态仓储（房间 → 最新 CRDT update）。
// collab ws 引擎（admin 内嵌 ygo）的持久化后端；系统管道存储，
// 不走租户视图（房间名全局唯一，键含内容 id）。
type CollabDocStateRepo struct {
	entClient *entCrud.EntClient[*ent.Client]
	log       *log.Helper
}

func NewCollabDocStateRepo(ctx *bootstrap.Context, entClient *entCrud.EntClient[*ent.Client]) *CollabDocStateRepo {
	return &CollabDocStateRepo{
		entClient: entClient,
		log:       ctx.NewLoggerHelper("collab-doc-state/repo/core-service"),
	}
}

// Get 取房间状态；无状态返回 (nil, nil)。
func (r *CollabDocStateRepo) Get(ctx context.Context, room string) (*contentV1.CollabDocState, error) {
	row, err := r.entClient.Client().CollabDocState.Query().
		Where(collabdocstate.RoomEQ(room)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return &contentV1.CollabDocState{Room: &room}, nil
		}
		r.log.Errorf("get collab doc state failed (room=%s): %s", room, err.Error())
		return nil, err
	}
	return r.toDTO(row), nil
}

// Store upsert 房间状态（整量替换最新合并 update）。
func (r *CollabDocStateRepo) Store(ctx context.Context, room string, update []byte) error {
	exist, err := r.entClient.Client().CollabDocState.Query().
		Where(collabdocstate.RoomEQ(room)).
		Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		r.log.Errorf("query collab doc state failed (room=%s): %s", room, err.Error())
		return err
	}

	if exist == nil {
		_, err = r.entClient.Client().CollabDocState.Create().
			SetRoom(room).
			SetDocUpdate(update).
			Save(ctx)
	} else {
		_, err = r.entClient.Client().CollabDocState.UpdateOneID(exist.ID).
			SetDocUpdate(update).
			Save(ctx)
	}
	if err != nil {
		r.log.Errorf("store collab doc state failed (room=%s): %s", room, err.Error())
		return err
	}
	return nil
}

// DeleteByPrefix 按房间名前缀清理（内容删除联动）。
func (r *CollabDocStateRepo) DeleteByPrefix(ctx context.Context, prefix string) error {
	count, err := r.entClient.Client().CollabDocState.Delete().
		Where(collabdocstate.RoomHasPrefix(prefix)).
		Exec(ctx)
	if err != nil {
		r.log.Errorf("delete collab doc states by prefix failed (prefix=%s): %s", prefix, err.Error())
		return err
	}
	if count > 0 {
		r.log.Infof("deleted %d collab doc states (prefix=%s)", count, prefix)
	}
	return nil
}

func (r *CollabDocStateRepo) toDTO(row *ent.CollabDocState) *contentV1.CollabDocState {
	if row == nil {
		return nil
	}
	dto := &contentV1.CollabDocState{
		Room:      &row.Room,
		DocUpdate: row.DocUpdate,
	}
	if row.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*row.UpdatedAt)
	}
	return dto
}
