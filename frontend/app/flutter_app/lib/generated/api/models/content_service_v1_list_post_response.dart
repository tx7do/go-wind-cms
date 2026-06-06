// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_post.dart';

part 'content_service_v1_list_post_response.g.dart';

/// 回应 - 帖子列表
@JsonSerializable()
class ContentServiceV1ListPostResponse {
  const ContentServiceV1ListPostResponse({
    this.items,
    this.total,
  });
  
  factory ContentServiceV1ListPostResponse.fromJson(Map<String, Object?> json) => _$ContentServiceV1ListPostResponseFromJson(json);
  
  final List<ContentServiceV1Post>? items;
  final String? total;

  Map<String, Object?> toJson() => _$ContentServiceV1ListPostResponseToJson(this);
}
