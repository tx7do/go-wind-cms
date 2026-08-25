package client

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"entgo.io/ent/dialect/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	entCrud "github.com/tx7do/go-crud/entgo"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	entBootstrap "github.com/tx7do/kratos-bootstrap/database/ent"

	"go-wind-cms/app/core/service/internal/data/ent"
	"go-wind-cms/app/core/service/internal/data/ent/migrate"
	_ "go-wind-cms/app/core/service/internal/data/ent/runtime"
)

// NewEntClient 创建Ent ORM数据库客户端
func NewEntClient(ctx *bootstrap.Context) (*entCrud.EntClient[*ent.Client], func(), error) {
	l := ctx.NewLoggerHelper("ent/data/core-service")

	cfg := ctx.GetConfig()
	if cfg == nil || cfg.Data == nil {
		l.Fatalf("[ENT] failed getting config")
		return nil, func() {}, nil
	}

	cli, err := entBootstrap.NewEntClient(cfg, func(drv *sql.Driver) *ent.Client {
		client := ent.NewClient(
			ent.Driver(drv),
			ent.Log(func(a ...any) {
				l.Debug(a...)
			}),
		)
		if client == nil {
			l.Fatalf("[ENT] failed creating ent client")
			return nil
		}

		ApplyTimeDefaultHooks(client)

		// run the auto migration tool
		if cfg.Data.Database.GetMigrate() {
			if err := client.Schema.Create(ctx.Context(), migrate.WithForeignKeys(true)); err != nil {
				l.Fatalf("[ENT] failed creating schema resources: %v", err)
			}
		}

		return client
	})
	if err != nil {
		log.Fatalf("[ENT] failed creating ent client: %v", err)
		return nil, func() {}, err
	}

	return cli, func() {
		if cleanErr := cli.Close(); cleanErr != nil {
			log.Errorf("[ENT] failed closing ent client: %v", cleanErr)
		}
	}, nil
}

// 带 created_at/updated_at setter 的 mutation(所有使用 go-crud 时间 mixin 的实体均实现)。
type createdAtSetter interface{ SetCreatedAt(time.Time) }
type updatedAtSetter interface{ SetUpdatedAt(time.Time) }

// ApplyTimeDefaultHooks 注册全局 mutation hook,为时间列兜底:
// Create 时补 created_at/updated_at,Update 时补 updated_at。
//
// 背景:项目所用 go-crud 的 TimeAt/CreatedAt mixin 只声明了
// Optional().Nillable(),不带任何默认值,也不依赖数据库列 DEFAULT;
// 而大量写入路径(user/post/comment/interaction ledger 等)未显式设置,
// 导致这些表 created_at 为 NULL,无法支撑任何按时间维度的统计与审计。
// 已显式设置的时间不被覆盖。enttest 测试客户端同样需要调用本函数。
func ApplyTimeDefaultHooks(client *ent.Client) {
	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			op := m.Op()
			if op.Is(ent.OpCreate) {
				setTimeIfMissing(m, "created_at", func(v time.Time) {
					if s, ok := m.(createdAtSetter); ok {
						s.SetCreatedAt(v)
					}
				})
				setTimeIfMissing(m, "updated_at", func(v time.Time) {
					if s, ok := m.(updatedAtSetter); ok {
						s.SetUpdatedAt(v)
					}
				})
			} else if op.Is(ent.OpUpdate) || op.Is(ent.OpUpdateOne) {
				setTimeIfMissing(m, "updated_at", func(v time.Time) {
					if s, ok := m.(updatedAtSetter); ok {
						s.SetUpdatedAt(v)
					}
				})
			}
			return next.Mutate(ctx, m)
		})
	})
}

// setTimeIfMissing 当字段未被显式设置时以 now 填充(通过 setter 接口 duck typing)。
func setTimeIfMissing(m ent.Mutation, field string, set func(time.Time)) {
	for _, f := range m.Fields() {
		if f == field {
			return
		}
	}
	set(time.Now())
}
