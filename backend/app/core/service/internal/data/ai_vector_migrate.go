package data

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
)

// AiEmbeddingDimensions embedding 输出向量维度。
// 与语义检索/知识库向量化请求的 text-embedding-3-small（1536 维）一致；
// 换 embedding 模型必须同步改这里并触发全量重索引，否则写入按维度不符报错。
const AiEmbeddingDimensions = 1536

// HnswIndexRowThreshold HNSW 索引自动门槛（精确表行数）。实测基线（2026-09-28，
// pgvector 0.8.6、1536 维、1.1 万行探针）：无过滤 ANN 查询 41.5ms→1.0ms（~40×），
// 但带过滤查询在低选择性下 planner 仍保守选顺序扫描——索引先建好，待规模/选择性
// 使迭代扫描更优时由 planner 自主启用；且建索引代价随表规模超线性增长（1.1 万行
// 2.7s），越线即建是最便宜的建设时机，拖到更大表反而把秒级构建拖成分钟级阻塞启动。
const HnswIndexRowThreshold = 10000

var vectorTableIdent = regexp.MustCompile(`^[a-z0-9_]+$`)

var dbNameIdent = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// EnsureVectorColumnDim 幂等确保 <table>.embedding 存在且为定维
// vector(AiEmbeddingDimensions)：列缺失则补建；维度不匹配（含历史非定维列）
// 则整列重建——旧向量随列丢弃，调用方须提示重索引。定维是硬约束：
// 混维向量会让 <=> 比较在查询期整体报错，宁可写入期 fail-fast。
func EnsureVectorColumnDim(ctx context.Context, db *sql.DB, table string) (rebuilt bool, err error) {
	if !vectorTableIdent.MatchString(table) {
		return false, fmt.Errorf("invalid table name: %s", table)
	}
	if _, err = db.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS vector`); err != nil {
		return false, fmt.Errorf("create pgvector extension: %w", err)
	}
	if _, err = db.ExecContext(ctx, fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN IF NOT EXISTS embedding vector(%d)", table, AiEmbeddingDimensions),
	); err != nil {
		return false, fmt.Errorf("add embedding column: %w", err)
	}

	var typmod int
	if err = db.QueryRowContext(ctx,
		`SELECT a.atttypmod FROM pg_attribute a
		 JOIN pg_type t ON t.oid = a.atttypid
		 WHERE a.attrelid = $1::regclass AND a.attname = 'embedding'
		   AND a.attnum > 0 AND NOT a.attisdropped AND t.typname = 'vector'`,
		table,
	).Scan(&typmod); err != nil {
		return false, fmt.Errorf("read embedding column dimension: %w", err)
	}
	if typmod == AiEmbeddingDimensions {
		return false, nil
	}

	if _, err = db.ExecContext(ctx, fmt.Sprintf(
		"ALTER TABLE %s DROP COLUMN embedding", table)); err != nil {
		return false, fmt.Errorf("drop mismatched embedding column: %w", err)
	}
	if _, err = db.ExecContext(ctx, fmt.Sprintf(
		"ALTER TABLE %s ADD COLUMN embedding vector(%d)", table, AiEmbeddingDimensions),
	); err != nil {
		return false, fmt.Errorf("re-add embedding column: %w", err)
	}
	return true, nil
}

// CountTableRows 精确表行数（HNSW 门槛的触发信号）。
// 不用 pg_class.reltuples：实测未 ANALYZE 前恒为 -1，会在启动迁移里错过越线；
// 精确 count 实测 1.1 万行仅 4ms（向量列 TOAST 化后堆内是指针行），启动期可承受。
func CountTableRows(ctx context.Context, db *sql.DB, table string) (int64, error) {
	if !vectorTableIdent.MatchString(table) {
		return 0, fmt.Errorf("invalid table name: %s", table)
	}
	var n int64
	if err := db.QueryRowContext(ctx, fmt.Sprintf("SELECT count(*) FROM %s", table)).Scan(&n); err != nil {
		return 0, fmt.Errorf("count table rows: %w", err)
	}
	return n, nil
}

// EnsureHnswIndexIfLarge 表行数越过 HnswIndexRowThreshold 时建 HNSW cosine 索引。
// 索引存在或未越线均返回 created=false；仅新建成功返回 true。
// 副作用（由调用方在 created=true 时一并告警并执行 SetHnswIterativeScanDefault）：
// ANN 查询变为近似召回（实测对抗性近重复数据上 top-8 尾部错位）、写入承担索引
// 维护代价（实测 ~3×/行）。只建索引不放开 iterative_scan 的话，带过滤的检索
// 永远用不上它，索引沦为纯写入负担。
func EnsureHnswIndexIfLarge(ctx context.Context, db *sql.DB, table string, rowCount int64) (created bool, err error) {
	if !vectorTableIdent.MatchString(table) {
		return false, fmt.Errorf("invalid table name: %s", table)
	}
	if rowCount < HnswIndexRowThreshold {
		return false, nil
	}
	indexName := table + "_embedding_hnsw"
	var exists int
	err = db.QueryRowContext(ctx,
		`SELECT 1 FROM pg_indexes WHERE tablename = $1 AND indexname = $2`,
		table, indexName,
	).Scan(&exists)
	if err == nil {
		return false, nil
	}
	if err != sql.ErrNoRows {
		return false, fmt.Errorf("check existing hnsw index: %w", err)
	}
	if _, err = db.ExecContext(ctx, fmt.Sprintf(
		"CREATE INDEX %s ON %s USING hnsw (embedding vector_cosine_ops)", indexName, table),
	); err != nil {
		return false, fmt.Errorf("create hnsw index: %w", err)
	}
	return true, nil
}

// SetHnswIterativeScanDefault 把当前库的 hnsw.iterative_scan 库级默认设为
// relaxed_order：允许 planner 对带过滤（租户/知识库/文档限定）的 ANN 查询走
// 迭代索引扫描，代价是结果为近似召回且距离序不严格。只对新会话生效（实测），
// 迁移期连接池尚未铺开，启动迁移是设置它的正确时机。
func SetHnswIterativeScanDefault(ctx context.Context, db *sql.DB) (err error) {
	var dbName string
	if err = db.QueryRowContext(ctx, `SELECT current_database()`).Scan(&dbName); err != nil {
		return fmt.Errorf("read current database: %w", err)
	}
	if !dbNameIdent.MatchString(dbName) {
		return fmt.Errorf("invalid database name: %s", dbName)
	}
	if _, err = db.ExecContext(ctx, fmt.Sprintf(
		`ALTER DATABASE "%s" SET hnsw.iterative_scan = 'relaxed_order'`, dbName)); err != nil {
		return fmt.Errorf("set hnsw.iterative_scan default: %w", err)
	}
	return nil
}
