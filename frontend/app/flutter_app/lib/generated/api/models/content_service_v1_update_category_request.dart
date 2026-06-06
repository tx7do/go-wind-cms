// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_category.dart';

part 'content_service_v1_update_category_request.g.dart';

/// 请求 - 更新类别
@JsonSerializable()
class ContentServiceV1UpdateCategoryRequest {
  const ContentServiceV1UpdateCategoryRequest({
    this.id,
    this.data,
    this.updateMask,
    this.allowMissing,
  });
  
  factory ContentServiceV1UpdateCategoryRequest.fromJson(Map<String, Object?> json) => _$ContentServiceV1UpdateCategoryRequestFromJson(json);
  
  final int? id;
  final ContentServiceV1Category? data;

  /// 要更新的字段列表
  final String? updateMask;

  /// 如果设置为true的时候，资源不存在则会新增(插入)，并且在这种情况下`updateMask`字段将会被忽略。
  final bool? allowMissing;

  Map<String, Object?> toJson() => _$ContentServiceV1UpdateCategoryRequestToJson(this);
}
