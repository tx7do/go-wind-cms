// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'comment_service_v1_comment.dart';

part 'comment_service_v1_list_comment_response.g.dart';

/// 回应 - 评论列表
@JsonSerializable()
class CommentServiceV1ListCommentResponse {
  const CommentServiceV1ListCommentResponse({
    this.items,
    this.total,
  });
  
  factory CommentServiceV1ListCommentResponse.fromJson(Map<String, Object?> json) => _$CommentServiceV1ListCommentResponseFromJson(json);
  
  final List<CommentServiceV1Comment>? items;
  final String? total;

  Map<String, Object?> toJson() => _$CommentServiceV1ListCommentResponseToJson(this);
}
