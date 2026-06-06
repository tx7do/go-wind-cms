// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'site_service_v1_navigation.dart';

part 'site_service_v1_create_navigation_request.g.dart';

/// 请求 - 创建导航
@JsonSerializable()
class SiteServiceV1CreateNavigationRequest {
  const SiteServiceV1CreateNavigationRequest({
    this.data,
  });
  
  factory SiteServiceV1CreateNavigationRequest.fromJson(Map<String, Object?> json) => _$SiteServiceV1CreateNavigationRequestFromJson(json);
  
  final SiteServiceV1Navigation? data;

  Map<String, Object?> toJson() => _$SiteServiceV1CreateNavigationRequestToJson(this);
}
