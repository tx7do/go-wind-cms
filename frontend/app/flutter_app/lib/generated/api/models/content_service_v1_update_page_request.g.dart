// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_update_page_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1UpdatePageRequest _$ContentServiceV1UpdatePageRequestFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1UpdatePageRequest(
  id: (json['id'] as num?)?.toInt(),
  data: json['data'] == null
      ? null
      : ContentServiceV1Page.fromJson(json['data'] as Map<String, dynamic>),
  updateMask: json['updateMask'] as String?,
  allowMissing: json['allowMissing'] as bool?,
);

Map<String, dynamic> _$ContentServiceV1UpdatePageRequestToJson(
  ContentServiceV1UpdatePageRequest instance,
) => <String, dynamic>{
  'id': instance.id,
  'data': instance.data,
  'updateMask': instance.updateMask,
  'allowMissing': instance.allowMissing,
};
