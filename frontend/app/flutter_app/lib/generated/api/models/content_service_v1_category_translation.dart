// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_seo_meta.dart';

part 'content_service_v1_category_translation.g.dart';

/// 类别翻译
@JsonSerializable()
class ContentServiceV1CategoryTranslation {
  const ContentServiceV1CategoryTranslation({
    this.id,
    this.categoryId,
    this.languageCode,
    this.name,
    this.slug,
    this.description,
    this.thumbnail,
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
  
  factory ContentServiceV1CategoryTranslation.fromJson(Map<String, Object?> json) => _$ContentServiceV1CategoryTranslationFromJson(json);
  
  /// 翻译记录ID
  final int? id;

  /// 关联的分类ID
  final int? categoryId;

  /// 语言代码
  final String? languageCode;

  /// 分类名称
  final String? name;

  /// 分类别名
  final String? slug;

  /// 分类描述
  final String? description;

  /// 缩略图
  final String? thumbnail;

  /// 封面图
  final String? coverImage;

  /// 完整路径（如 /zh-CN/category/tech）
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

  Map<String, Object?> toJson() => _$ContentServiceV1CategoryTranslationToJson(this);
}
