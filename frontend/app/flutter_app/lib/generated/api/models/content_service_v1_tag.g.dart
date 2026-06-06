// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_tag.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1Tag _$ContentServiceV1TagFromJson(Map<String, dynamic> json) =>
    ContentServiceV1Tag(
      id: (json['id'] as num?)?.toInt(),
      status: json['status'] == null
          ? null
          : ContentServiceV1TagStatus.fromJson(json['status'] as String),
      color: json['color'] as String?,
      icon: json['icon'] as String?,
      group: json['group'] as String?,
      sortOrder: (json['sortOrder'] as num?)?.toInt(),
      isFeatured: json['isFeatured'] as bool?,
      code: json['code'] as String?,
      postCount: (json['postCount'] as num?)?.toInt(),
      translations: (json['translations'] as List<dynamic>?)
          ?.map(
            (e) => ContentServiceV1TagTranslation.fromJson(
              e as Map<String, dynamic>,
            ),
          )
          .toList(),
      availableLanguages: (json['availableLanguages'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
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

Map<String, dynamic> _$ContentServiceV1TagToJson(
  ContentServiceV1Tag instance,
) => <String, dynamic>{
  'id': instance.id,
  'status': _$ContentServiceV1TagStatusEnumMap[instance.status],
  'color': instance.color,
  'icon': instance.icon,
  'group': instance.group,
  'sortOrder': instance.sortOrder,
  'isFeatured': instance.isFeatured,
  'code': instance.code,
  'postCount': instance.postCount,
  'translations': instance.translations,
  'availableLanguages': instance.availableLanguages,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
};

const _$ContentServiceV1TagStatusEnumMap = {
  ContentServiceV1TagStatus.tagStatusUnspecified: 'TAG_STATUS_UNSPECIFIED',
  ContentServiceV1TagStatus.tagStatusActive: 'TAG_STATUS_ACTIVE',
  ContentServiceV1TagStatus.tagStatusHidden: 'TAG_STATUS_HIDDEN',
  ContentServiceV1TagStatus.tagStatusArchived: 'TAG_STATUS_ARCHIVED',
  ContentServiceV1TagStatus.$unknown: r'$unknown',
};
