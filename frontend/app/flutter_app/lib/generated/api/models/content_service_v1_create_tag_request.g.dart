// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_create_tag_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1CreateTagRequest _$ContentServiceV1CreateTagRequestFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1CreateTagRequest(
  data: json['data'] == null
      ? null
      : ContentServiceV1Tag.fromJson(json['data'] as Map<String, dynamic>),
);

Map<String, dynamic> _$ContentServiceV1CreateTagRequestToJson(
  ContentServiceV1CreateTagRequest instance,
) => <String, dynamic>{'data': instance.data};
