// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 内容类型
@JsonEnum()
enum CommentServiceV1CommentContentType {
  @JsonValue('CONTENT_TYPE_UNSPECIFIED')
  contentTypeUnspecified('CONTENT_TYPE_UNSPECIFIED'),
  @JsonValue('CONTENT_TYPE_POST')
  contentTypePost('CONTENT_TYPE_POST'),
  @JsonValue('CONTENT_TYPE_PAGE')
  contentTypePage('CONTENT_TYPE_PAGE'),
  @JsonValue('CONTENT_TYPE_PRODUCT')
  contentTypeProduct('CONTENT_TYPE_PRODUCT'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const CommentServiceV1CommentContentType(this.json);

  factory CommentServiceV1CommentContentType.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<CommentServiceV1CommentContentType> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
