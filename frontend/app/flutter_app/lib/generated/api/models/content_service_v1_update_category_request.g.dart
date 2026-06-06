// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_update_category_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1UpdateCategoryRequest
_$ContentServiceV1UpdateCategoryRequestFromJson(Map<String, dynamic> json) =>
    ContentServiceV1UpdateCategoryRequest(
      id: (json['id'] as num?)?.toInt(),
      data: json['data'] == null
          ? null
          : ContentServiceV1Category.fromJson(
              json['data'] as Map<String, dynamic>,
            ),
      updateMask: json['updateMask'] as String?,
      allowMissing: json['allowMissing'] as bool?,
    );

Map<String, dynamic> _$ContentServiceV1UpdateCategoryRequestToJson(
  ContentServiceV1UpdateCategoryRequest instance,
) => <String, dynamic>{
  'id': instance.id,
  'data': instance.data,
  'updateMask': instance.updateMask,
  'allowMissing': instance.allowMissing,
};
