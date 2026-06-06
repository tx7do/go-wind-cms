// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'site_service_v1_navigation.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SiteServiceV1Navigation _$SiteServiceV1NavigationFromJson(
  Map<String, dynamic> json,
) => SiteServiceV1Navigation(
  id: (json['id'] as num?)?.toInt(),
  name: json['name'] as String?,
  location: json['location'] == null
      ? null
      : SiteServiceV1NavigationLocation.fromJson(json['location'] as String),
  locale: json['locale'] as String?,
  isActive: json['isActive'] as bool?,
  items: (json['items'] as List<dynamic>?)
      ?.map(
        (e) => SiteServiceV1NavigationItem.fromJson(e as Map<String, dynamic>),
      )
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

Map<String, dynamic> _$SiteServiceV1NavigationToJson(
  SiteServiceV1Navigation instance,
) => <String, dynamic>{
  'id': instance.id,
  'name': instance.name,
  'location': _$SiteServiceV1NavigationLocationEnumMap[instance.location],
  'locale': instance.locale,
  'isActive': instance.isActive,
  'items': instance.items,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
};

const _$SiteServiceV1NavigationLocationEnumMap = {
  SiteServiceV1NavigationLocation.navigationLocationUnspecified:
      'NAVIGATION_LOCATION_UNSPECIFIED',
  SiteServiceV1NavigationLocation.header: 'HEADER',
  SiteServiceV1NavigationLocation.footer: 'FOOTER',
  SiteServiceV1NavigationLocation.sidebar: 'SIDEBAR',
  SiteServiceV1NavigationLocation.mobile: 'MOBILE',
  SiteServiceV1NavigationLocation.topBar: 'TOP_BAR',
  SiteServiceV1NavigationLocation.offcanvas: 'OFFCANVAS',
  SiteServiceV1NavigationLocation.$unknown: r'$unknown',
};
