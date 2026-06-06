// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'site_service_v1_create_navigation_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SiteServiceV1CreateNavigationRequest
_$SiteServiceV1CreateNavigationRequestFromJson(Map<String, dynamic> json) =>
    SiteServiceV1CreateNavigationRequest(
      data: json['data'] == null
          ? null
          : SiteServiceV1Navigation.fromJson(
              json['data'] as Map<String, dynamic>,
            ),
    );

Map<String, dynamic> _$SiteServiceV1CreateNavigationRequestToJson(
  SiteServiceV1CreateNavigationRequest instance,
) => <String, dynamic>{'data': instance.data};
