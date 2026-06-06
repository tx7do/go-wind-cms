// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_category.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1Category _$ContentServiceV1CategoryFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1Category(
  id: (json['id'] as num?)?.toInt(),
  status: json['status'] == null
      ? null
      : ContentServiceV1CategoryStatus.fromJson(json['status'] as String),
  sortOrder: (json['sortOrder'] as num?)?.toInt(),
  isNav: json['isNav'] as bool?,
  icon: json['icon'] as String?,
  code: json['code'] as String?,
  postCount: (json['postCount'] as num?)?.toInt(),
  directPostCount: (json['directPostCount'] as num?)?.toInt(),
  translations: (json['translations'] as List<dynamic>?)
      ?.map(
        (e) => ContentServiceV1CategoryTranslation.fromJson(
          e as Map<String, dynamic>,
        ),
      )
      .toList(),
  availableLanguages: (json['availableLanguages'] as List<dynamic>?)
      ?.map((e) => e as String)
      .toList(),
  customFields: (json['customFields'] as Map<String, dynamic>?)?.map(
    (k, e) => MapEntry(k, e as String),
  ),
  parentId: (json['parentId'] as num?)?.toInt(),
  children: (json['children'] as List<dynamic>?)
      ?.map((e) => ContentServiceV1Category.fromJson(e as Map<String, dynamic>))
      .toList(),
  depth: (json['depth'] as num?)?.toInt(),
  path: json['path'] as String?,
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

Map<String, dynamic> _$ContentServiceV1CategoryToJson(
  ContentServiceV1Category instance,
) => <String, dynamic>{
  'id': instance.id,
  'status': _$ContentServiceV1CategoryStatusEnumMap[instance.status],
  'sortOrder': instance.sortOrder,
  'isNav': instance.isNav,
  'icon': instance.icon,
  'code': instance.code,
  'postCount': instance.postCount,
  'directPostCount': instance.directPostCount,
  'translations': instance.translations,
  'availableLanguages': instance.availableLanguages,
  'customFields': instance.customFields,
  'parentId': instance.parentId,
  'children': instance.children,
  'depth': instance.depth,
  'path': instance.path,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
};

const _$ContentServiceV1CategoryStatusEnumMap = {
  ContentServiceV1CategoryStatus.categoryStatusUnspecified:
      'CATEGORY_STATUS_UNSPECIFIED',
  ContentServiceV1CategoryStatus.categoryStatusActive: 'CATEGORY_STATUS_ACTIVE',
  ContentServiceV1CategoryStatus.categoryStatusHidden: 'CATEGORY_STATUS_HIDDEN',
  ContentServiceV1CategoryStatus.categoryStatusArchived:
      'CATEGORY_STATUS_ARCHIVED',
  ContentServiceV1CategoryStatus.$unknown: r'$unknown',
};
