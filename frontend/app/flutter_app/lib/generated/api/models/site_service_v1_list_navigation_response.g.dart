// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'site_service_v1_list_navigation_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SiteServiceV1ListNavigationResponse
_$SiteServiceV1ListNavigationResponseFromJson(Map<String, dynamic> json) =>
    SiteServiceV1ListNavigationResponse(
      items: (json['items'] as List<dynamic>?)
          ?.map(
            (e) => SiteServiceV1Navigation.fromJson(e as Map<String, dynamic>),
          )
          .toList(),
      total: json['total'] as String?,
    );

Map<String, dynamic> _$SiteServiceV1ListNavigationResponseToJson(
  SiteServiceV1ListNavigationResponse instance,
) => <String, dynamic>{'items': instance.items, 'total': instance.total};
