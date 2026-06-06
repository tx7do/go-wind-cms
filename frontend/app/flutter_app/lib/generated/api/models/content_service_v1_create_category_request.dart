// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_category.dart';

part 'content_service_v1_create_category_request.g.dart';

/// 请求 - 创建类别
@JsonSerializable()
class ContentServiceV1CreateCategoryRequest {
  const ContentServiceV1CreateCategoryRequest({
    this.data,
  });
  
  factory ContentServiceV1CreateCategoryRequest.fromJson(Map<String, Object?> json) => _$ContentServiceV1CreateCategoryRequestFromJson(json);
  
  final ContentServiceV1Category? data;

  Map<String, Object?> toJson() => _$ContentServiceV1CreateCategoryRequestToJson(this);
}
