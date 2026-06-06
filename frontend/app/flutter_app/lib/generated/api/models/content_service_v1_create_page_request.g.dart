// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_create_page_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1CreatePageRequest _$ContentServiceV1CreatePageRequestFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1CreatePageRequest(
  data: json['data'] == null
      ? null
      : ContentServiceV1Page.fromJson(json['data'] as Map<String, dynamic>),
);

Map<String, dynamic> _$ContentServiceV1CreatePageRequestToJson(
  ContentServiceV1CreatePageRequest instance,
) => <String, dynamic>{'data': instance.data};
