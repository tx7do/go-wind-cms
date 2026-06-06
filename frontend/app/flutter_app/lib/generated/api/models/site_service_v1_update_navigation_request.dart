// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'site_service_v1_navigation.dart';

part 'site_service_v1_update_navigation_request.g.dart';

/// 请求 - 更新导航
@JsonSerializable()
class SiteServiceV1UpdateNavigationRequest {
  const SiteServiceV1UpdateNavigationRequest({
    this.id,
    this.data,
    this.updateMask,
    this.allowMissing,
  });
  
  factory SiteServiceV1UpdateNavigationRequest.fromJson(Map<String, Object?> json) => _$SiteServiceV1UpdateNavigationRequestFromJson(json);
  
  final int? id;
  final SiteServiceV1Navigation? data;

  /// 要更新的字段列表
  final String? updateMask;

  /// 如果设置为true的时候，资源不存在则会新增(插入)，并且在这种情况下`updateMask`字段将会被忽略。
  final bool? allowMissing;

  Map<String, Object?> toJson() => _$SiteServiceV1UpdateNavigationRequestToJson(this);
}
