// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_tag.dart';

part 'content_service_v1_list_tag_response.g.dart';

/// 回应 - 标签列表
@JsonSerializable()
class ContentServiceV1ListTagResponse {
  const ContentServiceV1ListTagResponse({
    this.items,
    this.total,
  });
  
  factory ContentServiceV1ListTagResponse.fromJson(Map<String, Object?> json) => _$ContentServiceV1ListTagResponseFromJson(json);
  
  final List<ContentServiceV1Tag>? items;
  final String? total;

  Map<String, Object?> toJson() => _$ContentServiceV1ListTagResponseToJson(this);
}
