// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_tag.dart';

part 'content_service_v1_create_tag_request.g.dart';

/// 请求 - 创建标签
@JsonSerializable()
class ContentServiceV1CreateTagRequest {
  const ContentServiceV1CreateTagRequest({
    this.data,
  });
  
  factory ContentServiceV1CreateTagRequest.fromJson(Map<String, Object?> json) => _$ContentServiceV1CreateTagRequestFromJson(json);
  
  final ContentServiceV1Tag? data;

  Map<String, Object?> toJson() => _$ContentServiceV1CreateTagRequestToJson(this);
}
