// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 帖子状态
@JsonEnum()
enum ContentServiceV1PostStatus {
  @JsonValue('POST_STATUS_UNSPECIFIED')
  postStatusUnspecified('POST_STATUS_UNSPECIFIED'),
  @JsonValue('POST_STATUS_DRAFT')
  postStatusDraft('POST_STATUS_DRAFT'),
  @JsonValue('POST_STATUS_PUBLISHED')
  postStatusPublished('POST_STATUS_PUBLISHED'),
  @JsonValue('POST_STATUS_SCHEDULED')
  postStatusScheduled('POST_STATUS_SCHEDULED'),
  @JsonValue('POST_STATUS_TRASHED')
  postStatusTrashed('POST_STATUS_TRASHED'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const ContentServiceV1PostStatus(this.json);

  factory ContentServiceV1PostStatus.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<ContentServiceV1PostStatus> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
