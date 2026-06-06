// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_post.dart';

part 'content_service_v1_create_post_request.g.dart';

/// 请求 - 创建帖子
@JsonSerializable()
class ContentServiceV1CreatePostRequest {
  const ContentServiceV1CreatePostRequest({
    this.data,
  });
  
  factory ContentServiceV1CreatePostRequest.fromJson(Map<String, Object?> json) => _$ContentServiceV1CreatePostRequestFromJson(json);
  
  final ContentServiceV1Post? data;

  Map<String, Object?> toJson() => _$ContentServiceV1CreatePostRequestToJson(this);
}
