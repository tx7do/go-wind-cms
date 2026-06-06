// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 评论状态
@JsonEnum()
enum CommentServiceV1CommentStatus {
  @JsonValue('STATUS_UNSPECIFIED')
  statusUnspecified('STATUS_UNSPECIFIED'),
  @JsonValue('STATUS_PENDING')
  statusPending('STATUS_PENDING'),
  @JsonValue('STATUS_APPROVED')
  statusApproved('STATUS_APPROVED'),
  @JsonValue('STATUS_REJECTED')
  statusRejected('STATUS_REJECTED'),
  @JsonValue('STATUS_SPAM')
  statusSpam('STATUS_SPAM'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const CommentServiceV1CommentStatus(this.json);

  factory CommentServiceV1CommentStatus.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<CommentServiceV1CommentStatus> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
