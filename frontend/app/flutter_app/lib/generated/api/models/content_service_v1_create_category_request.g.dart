// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_create_category_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1CreateCategoryRequest
_$ContentServiceV1CreateCategoryRequestFromJson(Map<String, dynamic> json) =>
    ContentServiceV1CreateCategoryRequest(
      data: json['data'] == null
          ? null
          : ContentServiceV1Category.fromJson(
              json['data'] as Map<String, dynamic>,
            ),
    );

Map<String, dynamic> _$ContentServiceV1CreateCategoryRequestToJson(
  ContentServiceV1CreateCategoryRequest instance,
) => <String, dynamic>{'data': instance.data};
