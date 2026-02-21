package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/tx7do/go-crud/entgo/mixin"
)

// MediaAsset holds the schema definition for the MediaAsset entity.
type MediaAsset struct {
	ent.Schema
}

func (MediaAsset) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "media_assets",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
		schema.Comment("媒体资源库表"),
	}
}

// Fields of the MediaAsset.
func (MediaAsset) Fields() []ent.Field {
	return []ent.Field{
		field.String("filename").
			Comment("原始文件名").
			Optional().
			Nillable(),

		field.Enum("type").
			Comment("媒体类型").
			NamedValues(
				"AssetTypeImage", "ASSET_TYPE_IMAGE",
				"AssetTypeVideo", "ASSET_TYPE_VIDEO",
				"AssetTypeDocument", "ASSET_TYPE_DOCUMENT",
				"AssetTypeAudio", "ASSET_TYPE_AUDIO",
				"AssetTypeArchive", "ASSET_TYPE_ARCHIVE",
				"AssetTypeOther", "ASSET_TYPE_OTHER",
			).
			Optional().
			Nillable(),

		field.String("mime_type").
			Comment("MIME 类型").
			Optional().
			Nillable(),

		field.Uint64("size").
			Comment("文件大小（字节）").
			Default(0).
			Optional().
			Nillable(),

		field.String("storage_path").
			Comment("存储路径").
			Optional().
			Nillable(),

		field.String("url").
			Comment("CDN 访问 URL").
			Optional().
			Nillable(),

		field.Uint32("width").
			Comment("宽度（像素）").
			Default(0).
			Optional().
			Nillable(),

		field.Uint32("height").
			Comment("高度（像素）").
			Default(0).
			Optional().
			Nillable(),

		field.Uint32("duration").
			Comment("时长（秒）").
			Default(0).
			Optional().
			Nillable(),

		field.String("alt_text").
			Comment("ALT 文本").
			Optional().
			Nillable(),

		field.String("title").
			Comment("标题").
			Optional().
			Nillable(),

		field.String("caption").
			Comment("说明文字").
			Optional().
			Nillable(),

		field.Enum("processing_status").
			Comment("媒体处理状态").
			NamedValues(
				"ProcessingStatusUploading", "PROCESSING_STATUS_UPLOADING",
				"ProcessingStatusProcessing", "PROCESSING_STATUS_PROCESSING",
				"ProcessingStatusCompleted", "PROCESSING_STATUS_COMPLETED",
				"ProcessingStatusFailed", "PROCESSING_STATUS_FAILED",
			).
			Optional().
			Nillable(),

		field.String("processing_error").
			Comment("处理失败时的错误信息").
			Optional().
			Nillable(),

		field.String("file_hash").
			Comment("文件哈希值").
			Optional().
			Nillable(),

		field.Uint32("folder_id").
			Comment("所属文件夹ID").
			Default(0).
			Optional().
			Nillable(),

		field.Uint32("file_id").
			Comment("存储文件ID").
			Optional().
			Nillable(),

		field.Uint32("reference_count").
			Comment("被引用次数").
			Default(0).
			Optional().
			Nillable(),

		field.Bool("is_private").
			Comment("是否私密").
			Default(false).
			Optional().
			Nillable(),

		//field.JSON("variant_file_ids", &map[string]uint32{}).
		//	Comment("变体文件URL").
		//	Optional(),
	}
}

// Mixin of the MediaAsset.
func (MediaAsset) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.TimeAt{},
		mixin.OperatorID{},
	}
}

func (MediaAsset) Indexes() []ent.Index {
	return []ent.Index{
		// 单字段索引，用于按媒体类型查询
		index.Fields("type"),
		// 单字段索引，用于按处理状态查询
		index.Fields("processing_status"),
		// 单字段索引，用于按所属文件夹查询媒体资源
		index.Fields("folder_id"),
		// 单字段索引，用于按存储文件ID查询
		index.Fields("file_id"),
		// 单字段索引，用于按私密状态过滤
		index.Fields("is_private"),
		// 单字段索引，优化文件哈希值的去重和查询
		index.Fields("file_hash"),
		// 复合索引，优化按文件夹和私密状态查询
		index.Fields("folder_id", "is_private"),
		// 复合索引，优化按媒体类型和处理状态查询
		index.Fields("type", "processing_status"),
	}
}
