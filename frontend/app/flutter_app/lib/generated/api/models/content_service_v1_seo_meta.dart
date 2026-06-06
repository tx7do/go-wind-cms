// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'content_service_v1_seo_meta.g.dart';

/// SEO 元信息结构设计
@JsonSerializable()
class ContentServiceV1SeoMeta {
  const ContentServiceV1SeoMeta({
    this.seoTitle,
    this.metaKeywords,
    this.metaDescription,
    this.ogTitle,
    this.ogDescription,
    this.ogImage,
    this.canonicalUrl,
  });
  
  factory ContentServiceV1SeoMeta.fromJson(Map<String, Object?> json) => _$ContentServiceV1SeoMetaFromJson(json);
  
  /// SEO 标题（覆盖页面默认标题）
  final String? seoTitle;

  /// SEO 关键词
  final String? metaKeywords;

  /// SEO 描述
  final String? metaDescription;

  /// OG 标题（覆盖页面默认标题）
  final String? ogTitle;

  /// OG 描述
  final String? ogDescription;

  /// OG 图片 URL
  final String? ogImage;

  /// 规范 URL（SEO 用，格式：/{locale}/tag/{slug}，如 /zh-CN/tag/golang，允许自定义以覆盖默认生成的 URL）
  final String? canonicalUrl;

  Map<String, Object?> toJson() => _$ContentServiceV1SeoMetaToJson(this);
}
