// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_page.dart';

part 'content_service_v1_update_page_request.g.dart';

/// 请求 - 更新页面
@JsonSerializable()
class ContentServiceV1UpdatePageRequest {
  const ContentServiceV1UpdatePageRequest({
    this.id,
    this.data,
    this.updateMask,
    this.allowMissing,
  });
  
  factory ContentServiceV1UpdatePageRequest.fromJson(Map<String, Object?> json) => _$ContentServiceV1UpdatePageRequestFromJson(json);
  
  final int? id;
  final ContentServiceV1Page? data;

  /// 要更新的字段列表
  final String? updateMask;

  /// 如果设置为true的时候，资源不存在则会新增(插入)，并且在这种情况下`updateMask`字段将会被忽略。
  final bool? allowMissing;

  Map<String, Object?> toJson() => _$ContentServiceV1UpdatePageRequestToJson(this);
}
