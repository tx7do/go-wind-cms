// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_seo_meta.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1SeoMeta _$ContentServiceV1SeoMetaFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1SeoMeta(
  seoTitle: json['seoTitle'] as String?,
  metaKeywords: json['metaKeywords'] as String?,
  metaDescription: json['metaDescription'] as String?,
  ogTitle: json['ogTitle'] as String?,
  ogDescription: json['ogDescription'] as String?,
  ogImage: json['ogImage'] as String?,
  canonicalUrl: json['canonicalUrl'] as String?,
);

Map<String, dynamic> _$ContentServiceV1SeoMetaToJson(
  ContentServiceV1SeoMeta instance,
) => <String, dynamic>{
  'seoTitle': instance.seoTitle,
  'metaKeywords': instance.metaKeywords,
  'metaDescription': instance.metaDescription,
  'ogTitle': instance.ogTitle,
  'ogDescription': instance.ogDescription,
  'ogImage': instance.ogImage,
  'canonicalUrl': instance.canonicalUrl,
};
