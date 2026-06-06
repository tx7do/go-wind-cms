// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 状态
@JsonEnum()
enum IdentityServiceV1UserStatus {
  @JsonValue('DISABLED')
  disabled('DISABLED'),
  @JsonValue('NORMAL')
  normal('NORMAL'),
  @JsonValue('PENDING')
  pending('PENDING'),
  @JsonValue('LOCKED')
  locked('LOCKED'),
  @JsonValue('EXPIRED')
  expired('EXPIRED'),
  @JsonValue('CLOSED')
  closed('CLOSED'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const IdentityServiceV1UserStatus(this.json);

  factory IdentityServiceV1UserStatus.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<IdentityServiceV1UserStatus> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
