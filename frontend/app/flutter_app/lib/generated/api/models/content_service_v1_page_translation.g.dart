// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_page_translation.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1PageTranslation _$ContentServiceV1PageTranslationFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1PageTranslation(
  id: (json['id'] as num?)?.toInt(),
  pageId: (json['pageId'] as num?)?.toInt(),
  languageCode: json['languageCode'] as String?,
  title: json['title'] as String?,
  slug: json['slug'] as String?,
  sections: (json['sections'] as List<dynamic>?)
      ?.map((e) => ContentServiceV1Section.fromJson(e as Map<String, dynamic>))
      .toList(),
  thumbnail: json['thumbnail'] as String?,
  coverImage: json['coverImage'] as String?,
  fullPath: json['fullPath'] as String?,
  seo: json['seo'] == null
      ? null
      : ContentServiceV1SeoMeta.fromJson(json['seo'] as Map<String, dynamic>),
  createdBy: (json['createdBy'] as num?)?.toInt(),
  updatedBy: (json['updatedBy'] as num?)?.toInt(),
  deletedBy: (json['deletedBy'] as num?)?.toInt(),
  createdAt: json['createdAt'] == null
      ? null
      : DateTime.parse(json['createdAt'] as String),
  updatedAt: json['updatedAt'] == null
      ? null
      : DateTime.parse(json['updatedAt'] as String),
  deletedAt: json['deletedAt'] == null
      ? null
      : DateTime.parse(json['deletedAt'] as String),
);

Map<String, dynamic> _$ContentServiceV1PageTranslationToJson(
  ContentServiceV1PageTranslation instance,
) => <String, dynamic>{
  'id': instance.id,
  'pageId': instance.pageId,
  'languageCode': instance.languageCode,
  'title': instance.title,
  'slug': instance.slug,
  'sections': instance.sections,
  'thumbnail': instance.thumbnail,
  'coverImage': instance.coverImage,
  'fullPath': instance.fullPath,
  'seo': instance.seo,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
};
