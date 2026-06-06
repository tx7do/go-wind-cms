// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_update_post_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1UpdatePostRequest _$ContentServiceV1UpdatePostRequestFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1UpdatePostRequest(
  id: (json['id'] as num?)?.toInt(),
  data: json['data'] == null
      ? null
      : ContentServiceV1Post.fromJson(json['data'] as Map<String, dynamic>),
  updateMask: json['updateMask'] as String?,
  allowMissing: json['allowMissing'] as bool?,
);

Map<String, dynamic> _$ContentServiceV1UpdatePostRequestToJson(
  ContentServiceV1UpdatePostRequest instance,
) => <String, dynamic>{
  'id': instance.id,
  'data': instance.data,
  'updateMask': instance.updateMask,
  'allowMissing': instance.allowMissing,
};
