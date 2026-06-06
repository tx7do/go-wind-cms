// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_seo_meta.dart';

part 'content_service_v1_post_translation.g.dart';

/// 帖子翻译
@JsonSerializable()
class ContentServiceV1PostTranslation {
  const ContentServiceV1PostTranslation({
    this.id,
    this.postId,
    this.languageCode,
    this.title,
    this.slug,
    this.summary,
    this.content,
    this.originalContent,
    this.thumbnail,
    this.fullPath,
    this.wordCount,
    this.seo,
    this.createdBy,
    this.updatedBy,
    this.deletedBy,
    this.createdAt,
    this.updatedAt,
    this.deletedAt,
  });
  
  factory ContentServiceV1PostTranslation.fromJson(Map<String, Object?> json) => _$ContentServiceV1PostTranslationFromJson(json);
  
  /// 帖子翻译ID
  final int? id;

  /// 关联的帖子ID
  final int? postId;

  /// 语言代码
  final String? languageCode;

  /// 帖子标题
  final String? title;

  /// 语言特定 slug（覆盖主表 slug，如中文用拼音 /post/hou-duan）
  final String? slug;

  /// 帖子摘要
  final String? summary;

  /// 帖子内容
  final String? content;

  /// 原始内容
  final String? originalContent;

  /// 缩略图（某些语言可能需要不同图片）
  final String? thumbnail;

  /// 完整路径（如 /zh-CN/blog/post-slug）
  final String? fullPath;

  /// 当前语言版本的字数（中文按字符数，英文按单词数）
  final int? wordCount;
  final ContentServiceV1SeoMeta? seo;

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

  Map<String, Object?> toJson() => _$ContentServiceV1PostTranslationToJson(this);
}
