// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 页面状态
@JsonEnum()
enum ContentServiceV1PageStatus {
  @JsonValue('PAGE_STATUS_UNSPECIFIED')
  pageStatusUnspecified('PAGE_STATUS_UNSPECIFIED'),
  @JsonValue('PAGE_STATUS_DRAFT')
  pageStatusDraft('PAGE_STATUS_DRAFT'),
  @JsonValue('PAGE_STATUS_PUBLISHED')
  pageStatusPublished('PAGE_STATUS_PUBLISHED'),
  @JsonValue('PAGE_STATUS_ARCHIVED')
  pageStatusArchived('PAGE_STATUS_ARCHIVED'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const ContentServiceV1PageStatus(this.json);

  factory ContentServiceV1PageStatus.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<ContentServiceV1PageStatus> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
