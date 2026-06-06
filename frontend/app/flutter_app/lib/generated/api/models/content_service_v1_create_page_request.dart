// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_page.dart';

part 'content_service_v1_create_page_request.g.dart';

/// 请求 - 创建页面
@JsonSerializable()
class ContentServiceV1CreatePageRequest {
  const ContentServiceV1CreatePageRequest({
    this.data,
  });
  
  factory ContentServiceV1CreatePageRequest.fromJson(Map<String, Object?> json) => _$ContentServiceV1CreatePageRequestFromJson(json);
  
  final ContentServiceV1Page? data;

  Map<String, Object?> toJson() => _$ContentServiceV1CreatePageRequestToJson(this);
}
