// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_post_editor_type.dart';
import 'content_service_v1_post_status.dart';
import 'content_service_v1_post_translation.dart';

part 'content_service_v1_post.g.dart';

/// 帖子
@JsonSerializable()
class ContentServiceV1Post {
  const ContentServiceV1Post({
    this.id,
    this.status,
    this.editorType,
    this.code,
    this.disallowComment,
    this.inProgress,
    this.autoSummary,
    this.isFeatured,
    this.sortOrder,
    this.visits,
    this.likes,
    this.commentCount,
    this.authorId,
    this.authorName,
    this.customFields,
    this.translations,
    this.availableLanguages,
    this.categoryIds,
    this.tagIds,
    this.passwordHash,
    this.createdBy,
    this.updatedBy,
    this.deletedBy,
    this.createdAt,
    this.updatedAt,
    this.deletedAt,
    this.publishTime,
  });
  
  factory ContentServiceV1Post.fromJson(Map<String, Object?> json) => _$ContentServiceV1PostFromJson(json);
  
  /// 帖子ID
  final int? id;

  /// 帖子状态
  final ContentServiceV1PostStatus? status;

  /// 编辑器类型
  final ContentServiceV1PostEditorType? editorType;

  /// 唯一编码（必须唯一，建议使用英文）
  final String? code;

  /// 帖子是否禁止评论
  final bool? disallowComment;

  /// 帖子是否正在编辑中（如正在编辑则不允许发布，避免内容不完整）
  final bool? inProgress;

  /// 是否自动生成摘要
  final bool? autoSummary;

  /// 是否推荐
  final bool? isFeatured;

  /// 排序优先级（数值越小越靠前，同组内排序）
  final int? sortOrder;

  /// 帖子访问次数
  final int? visits;

  /// 帖子点赞次数
  final int? likes;

  /// 帖子评论数
  final int? commentCount;

  /// 评论作者ID，0表示游客
  final int? authorId;

  /// 帖子作者名称（游客填写）
  final String? authorName;

  /// 自定义字段，键值对形式，便于扩展
  final Map<String, String>? customFields;

  /// 多语言翻译列表
  final List<ContentServiceV1PostTranslation>? translations;

  /// 可用的语言代码列表（快速查询，避免遍历 translations）
  final List<String>? availableLanguages;

  /// 关联的分类ID列表（多选）
  final List<int>? categoryIds;

  /// 关联的标签ID列表（多选）
  final List<int>? tagIds;

  /// 密码哈希
  final String? passwordHash;

  /// 创建者用户ID
  final int? createdBy;

  /// 更新者用户ID
  final int? updatedBy;

  /// 删除者用户ID
  final int? deletedBy;

  /// 创建时间
  final DateTime? createdAt;

  /// 更新时间
  final DateTime? updatedAt;

  /// 删除时间
  final DateTime? deletedAt;

  /// 发布时间
  final DateTime? publishTime;

  Map<String, Object?> toJson() => _$ContentServiceV1PostToJson(this);
}
