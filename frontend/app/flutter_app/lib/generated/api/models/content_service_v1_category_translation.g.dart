// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_category_translation.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1CategoryTranslation
_$ContentServiceV1CategoryTranslationFromJson(Map<String, dynamic> json) =>
    ContentServiceV1CategoryTranslation(
      id: (json['id'] as num?)?.toInt(),
      categoryId: (json['categoryId'] as num?)?.toInt(),
      languageCode: json['languageCode'] as String?,
      name: json['name'] as String?,
      slug: json['slug'] as String?,
      description: json['description'] as String?,
      thumbnail: json['thumbnail'] as String?,
      coverImage: json['coverImage'] as String?,
      fullPath: json['fullPath'] as String?,
      seo: json['seo'] == null
          ? null
          : ContentServiceV1SeoMeta.fromJson(
              json['seo'] as Map<String, dynamic>,
            ),
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

Map<String, dynamic> _$ContentServiceV1CategoryTranslationToJson(
  ContentServiceV1CategoryTranslation instance,
) => <String, dynamic>{
  'id': instance.id,
  'categoryId': instance.categoryId,
  'languageCode': instance.languageCode,
  'name': instance.name,
  'slug': instance.slug,
  'description': instance.description,
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
