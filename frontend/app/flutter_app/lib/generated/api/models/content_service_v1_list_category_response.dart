// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_category.dart';

part 'content_service_v1_list_category_response.g.dart';

/// 回应 - 类别列表
@JsonSerializable()
class ContentServiceV1ListCategoryResponse {
  const ContentServiceV1ListCategoryResponse({
    this.items,
    this.total,
  });
  
  factory ContentServiceV1ListCategoryResponse.fromJson(Map<String, Object?> json) => _$ContentServiceV1ListCategoryResponseFromJson(json);
  
  final List<ContentServiceV1Category>? items;
  final String? total;

  Map<String, Object?> toJson() => _$ContentServiceV1ListCategoryResponseToJson(this);
}
