// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_post_translation.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1PostTranslation _$ContentServiceV1PostTranslationFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1PostTranslation(
  id: (json['id'] as num?)?.toInt(),
  postId: (json['postId'] as num?)?.toInt(),
  languageCode: json['languageCode'] as String?,
  title: json['title'] as String?,
  slug: json['slug'] as String?,
  summary: json['summary'] as String?,
  content: json['content'] as String?,
  originalContent: json['originalContent'] as String?,
  thumbnail: json['thumbnail'] as String?,
  fullPath: json['fullPath'] as String?,
  wordCount: (json['wordCount'] as num?)?.toInt(),
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

Map<String, dynamic> _$ContentServiceV1PostTranslationToJson(
  ContentServiceV1PostTranslation instance,
) => <String, dynamic>{
  'id': instance.id,
  'postId': instance.postId,
  'languageCode': instance.languageCode,
  'title': instance.title,
  'slug': instance.slug,
  'summary': instance.summary,
  'content': instance.content,
  'originalContent': instance.originalContent,
  'thumbnail': instance.thumbnail,
  'fullPath': instance.fullPath,
  'wordCount': instance.wordCount,
  'seo': instance.seo,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
};
