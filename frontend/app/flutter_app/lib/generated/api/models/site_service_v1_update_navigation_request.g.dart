// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'site_service_v1_update_navigation_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SiteServiceV1UpdateNavigationRequest
_$SiteServiceV1UpdateNavigationRequestFromJson(Map<String, dynamic> json) =>
    SiteServiceV1UpdateNavigationRequest(
      id: (json['id'] as num?)?.toInt(),
      data: json['data'] == null
          ? null
          : SiteServiceV1Navigation.fromJson(
              json['data'] as Map<String, dynamic>,
            ),
      updateMask: json['updateMask'] as String?,
      allowMissing: json['allowMissing'] as bool?,
    );

Map<String, dynamic> _$SiteServiceV1UpdateNavigationRequestToJson(
  SiteServiceV1UpdateNavigationRequest instance,
) => <String, dynamic>{
  'id': instance.id,
  'data': instance.data,
  'updateMask': instance.updateMask,
  'allowMissing': instance.allowMissing,
};
