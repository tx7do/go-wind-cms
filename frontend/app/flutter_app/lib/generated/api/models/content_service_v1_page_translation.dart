// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_section.dart';
import 'content_service_v1_seo_meta.dart';

part 'content_service_v1_page_translation.g.dart';

/// 页面翻译
@JsonSerializable()
class ContentServiceV1PageTranslation {
  const ContentServiceV1PageTranslation({
    this.id,
    this.pageId,
    this.languageCode,
    this.title,
    this.slug,
    this.sections,
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
  
  factory ContentServiceV1PageTranslation.fromJson(Map<String, Object?> json) => _$ContentServiceV1PageTranslationFromJson(json);
  
  /// 翻译记录ID
  final int? id;

  /// 关联的页面ID
  final int? pageId;

  /// 语言代码（ISO 639-1 标准）
  final String? languageCode;

  /// 页面标题
  final String? title;

  /// 语言特定的 slug（覆盖主表 slug，用于多语言路由，如中文用拼音）
  final String? slug;

  /// 页面内容区块（模块化）
  final List<ContentServiceV1Section>? sections;

  /// 缩略图
  final String? thumbnail;

  /// 封面图
  final String? coverImage;

  /// 完整路径（如 /zh-CN/about-us），首页应为 '/' 或 '/zh-CN/'
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

  Map<String, Object?> toJson() => _$ContentServiceV1PageTranslationToJson(this);
}
