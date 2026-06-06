// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 分类状态
@JsonEnum()
enum ContentServiceV1CategoryStatus {
  @JsonValue('CATEGORY_STATUS_UNSPECIFIED')
  categoryStatusUnspecified('CATEGORY_STATUS_UNSPECIFIED'),
  @JsonValue('CATEGORY_STATUS_ACTIVE')
  categoryStatusActive('CATEGORY_STATUS_ACTIVE'),
  @JsonValue('CATEGORY_STATUS_HIDDEN')
  categoryStatusHidden('CATEGORY_STATUS_HIDDEN'),
  @JsonValue('CATEGORY_STATUS_ARCHIVED')
  categoryStatusArchived('CATEGORY_STATUS_ARCHIVED'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const ContentServiceV1CategoryStatus(this.json);

  factory ContentServiceV1CategoryStatus.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<ContentServiceV1CategoryStatus> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
