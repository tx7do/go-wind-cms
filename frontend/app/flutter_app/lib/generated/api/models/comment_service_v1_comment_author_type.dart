// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 作者类型（用于前端高亮管理员回复）
@JsonEnum()
enum CommentServiceV1CommentAuthorType {
  @JsonValue('AUTHOR_TYPE_UNSPECIFIED')
  authorTypeUnspecified('AUTHOR_TYPE_UNSPECIFIED'),
  @JsonValue('AUTHOR_TYPE_GUEST')
  authorTypeGuest('AUTHOR_TYPE_GUEST'),
  @JsonValue('AUTHOR_TYPE_USER')
  authorTypeUser('AUTHOR_TYPE_USER'),
  @JsonValue('AUTHOR_TYPE_ADMIN')
  authorTypeAdmin('AUTHOR_TYPE_ADMIN'),
  @JsonValue('AUTHOR_TYPE_MODERATOR')
  authorTypeModerator('AUTHOR_TYPE_MODERATOR'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const CommentServiceV1CommentAuthorType(this.json);

  factory CommentServiceV1CommentAuthorType.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<CommentServiceV1CommentAuthorType> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
