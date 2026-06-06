// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_page.dart';

part 'content_service_v1_list_page_response.g.dart';

/// 回应 - 页面列表
@JsonSerializable()
class ContentServiceV1ListPageResponse {
  const ContentServiceV1ListPageResponse({
    this.items,
    this.total,
  });
  
  factory ContentServiceV1ListPageResponse.fromJson(Map<String, Object?> json) => _$ContentServiceV1ListPageResponseFromJson(json);
  
  final List<ContentServiceV1Page>? items;
  final String? total;

  Map<String, Object?> toJson() => _$ContentServiceV1ListPageResponseToJson(this);
}
