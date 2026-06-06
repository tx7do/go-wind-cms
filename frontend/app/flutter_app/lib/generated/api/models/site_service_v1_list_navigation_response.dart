// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'site_service_v1_navigation.dart';

part 'site_service_v1_list_navigation_response.g.dart';

/// 回应 - 导航列表
@JsonSerializable()
class SiteServiceV1ListNavigationResponse {
  const SiteServiceV1ListNavigationResponse({
    this.items,
    this.total,
  });
  
  factory SiteServiceV1ListNavigationResponse.fromJson(Map<String, Object?> json) => _$SiteServiceV1ListNavigationResponseFromJson(json);
  
  final List<SiteServiceV1Navigation>? items;
  final String? total;

  Map<String, Object?> toJson() => _$SiteServiceV1ListNavigationResponseToJson(this);
}
