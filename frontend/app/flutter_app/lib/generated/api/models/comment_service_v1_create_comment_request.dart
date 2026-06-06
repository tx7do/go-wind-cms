// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'comment_service_v1_comment.dart';

part 'comment_service_v1_create_comment_request.g.dart';

/// 请求 - 创建评论
@JsonSerializable()
class CommentServiceV1CreateCommentRequest {
  const CommentServiceV1CreateCommentRequest({
    this.data,
  });
  
  factory CommentServiceV1CreateCommentRequest.fromJson(Map<String, Object?> json) => _$CommentServiceV1CreateCommentRequestFromJson(json);
  
  final CommentServiceV1Comment? data;

  Map<String, Object?> toJson() => _$CommentServiceV1CreateCommentRequestToJson(this);
}
