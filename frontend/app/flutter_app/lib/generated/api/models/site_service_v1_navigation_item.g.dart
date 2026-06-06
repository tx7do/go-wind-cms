// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'site_service_v1_navigation_item.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SiteServiceV1NavigationItem _$SiteServiceV1NavigationItemFromJson(
  Map<String, dynamic> json,
) => SiteServiceV1NavigationItem(
  id: (json['id'] as num?)?.toInt(),
  navigationId: (json['navigationId'] as num?)?.toInt(),
  title: json['title'] as String?,
  description: json['description'] as String?,
  icon: json['icon'] as String?,
  url: json['url'] as String?,
  linkType: json['linkType'] == null
      ? null
      : SiteServiceV1NavigationItemLinkType.fromJson(
          json['linkType'] as String,
        ),
  objectId: (json['objectId'] as num?)?.toInt(),
  sortOrder: (json['sortOrder'] as num?)?.toInt(),
  isOpenNewTab: json['isOpenNewTab'] as bool?,
  isInvalid: json['isInvalid'] as bool?,
  requiredPermission: json['requiredPermission'] as String?,
  parentId: (json['parentId'] as num?)?.toInt(),
  children: (json['children'] as List<dynamic>?)
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

Map<String, dynamic> _$SiteServiceV1NavigationItemToJson(
  SiteServiceV1NavigationItem instance,
) => <String, dynamic>{
  'id': instance.id,
  'navigationId': instance.navigationId,
  'title': instance.title,
  'description': instance.description,
  'icon': instance.icon,
  'url': instance.url,
  'linkType': _$SiteServiceV1NavigationItemLinkTypeEnumMap[instance.linkType],
  'objectId': instance.objectId,
  'sortOrder': instance.sortOrder,
  'isOpenNewTab': instance.isOpenNewTab,
  'isInvalid': instance.isInvalid,
  'requiredPermission': instance.requiredPermission,
  'parentId': instance.parentId,
  'children': instance.children,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
};

const _$SiteServiceV1NavigationItemLinkTypeEnumMap = {
  SiteServiceV1NavigationItemLinkType.linkTypeUnspecified:
      'LINK_TYPE_UNSPECIFIED',
  SiteServiceV1NavigationItemLinkType.linkTypeCustom: 'LINK_TYPE_CUSTOM',
  SiteServiceV1NavigationItemLinkType.linkTypePost: 'LINK_TYPE_POST',
  SiteServiceV1NavigationItemLinkType.linkTypePage: 'LINK_TYPE_PAGE',
  SiteServiceV1NavigationItemLinkType.linkTypeCategory: 'LINK_TYPE_CATEGORY',
  SiteServiceV1NavigationItemLinkType.linkTypeExternal: 'LINK_TYPE_EXTERNAL',
  SiteServiceV1NavigationItemLinkType.$unknown: r'$unknown',
};
