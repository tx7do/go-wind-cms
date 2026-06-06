// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_create_post_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1CreatePostRequest _$ContentServiceV1CreatePostRequestFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1CreatePostRequest(
  data: json['data'] == null
      ? null
      : ContentServiceV1Post.fromJson(json['data'] as Map<String, dynamic>),
);

Map<String, dynamic> _$ContentServiceV1CreatePostRequestToJson(
  ContentServiceV1CreatePostRequest instance,
) => <String, dynamic>{'data': instance.data};
