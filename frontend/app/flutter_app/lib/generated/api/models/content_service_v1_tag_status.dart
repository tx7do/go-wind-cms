// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 标签状态
@JsonEnum()
enum ContentServiceV1TagStatus {
  @JsonValue('TAG_STATUS_UNSPECIFIED')
  tagStatusUnspecified('TAG_STATUS_UNSPECIFIED'),
  @JsonValue('TAG_STATUS_ACTIVE')
  tagStatusActive('TAG_STATUS_ACTIVE'),
  @JsonValue('TAG_STATUS_HIDDEN')
  tagStatusHidden('TAG_STATUS_HIDDEN'),
  @JsonValue('TAG_STATUS_ARCHIVED')
  tagStatusArchived('TAG_STATUS_ARCHIVED'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const ContentServiceV1TagStatus(this.json);

  factory ContentServiceV1TagStatus.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<ContentServiceV1TagStatus> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
