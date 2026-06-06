// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_seo_meta.dart';

part 'content_service_v1_tag_translation.g.dart';

/// 标签翻译
@JsonSerializable()
class ContentServiceV1TagTranslation {
  const ContentServiceV1TagTranslation({
    this.id,
    this.tagId,
    this.languageCode,
    this.name,
    this.slug,
    this.description,
    this.coverImage,
    this.fullPath,
    this.seo,
    this.createdBy,
    this.updatedBy,
    this.deletedBy,
    this.createdAt,
    this.updatedAt,
    this.deletedAt,
  });
  
  factory ContentServiceV1TagTranslation.fromJson(Map<String, Object?> json) => _$ContentServiceV1TagTranslationFromJson(json);
  
  /// 翻译记录ID
  final int? id;

  /// 关联的标签ID
  final int? tagId;

  /// 语言代码
  final String? languageCode;

  /// 标签名称
  final String? name;

  /// 标签别名（同语言下唯一，不同语言可重复，如 zh-CN 的 'go' 与 en-US 的 'go' 允许共存）
  final String? slug;

  /// 标签描述
  final String? description;

  /// 封面图
  final String? coverImage;

  /// 完整路径（格式：/{locale}/tag/{slug}，如 /zh-CN/tag/golang）
  final String? fullPath;
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

  Map<String, Object?> toJson() => _$ContentServiceV1TagTranslationToJson(this);
}
