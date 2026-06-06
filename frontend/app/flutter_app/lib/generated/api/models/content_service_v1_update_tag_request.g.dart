// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_update_tag_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1UpdateTagRequest _$ContentServiceV1UpdateTagRequestFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1UpdateTagRequest(
  id: (json['id'] as num?)?.toInt(),
  data: json['data'] == null
      ? null
      : ContentServiceV1Tag.fromJson(json['data'] as Map<String, dynamic>),
  updateMask: json['updateMask'] as String?,
  allowMissing: json['allowMissing'] as bool?,
);

Map<String, dynamic> _$ContentServiceV1UpdateTagRequestToJson(
  ContentServiceV1UpdateTagRequest instance,
) => <String, dynamic>{
  'id': instance.id,
  'data': instance.data,
  'updateMask': instance.updateMask,
  'allowMissing': instance.allowMissing,
};
