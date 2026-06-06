// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 性别
@JsonEnum()
enum IdentityServiceV1UserGender {
  @JsonValue('SECRET')
  secret('SECRET'),
  @JsonValue('MALE')
  male('MALE'),
  @JsonValue('FEMALE')
  female('FEMALE'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const IdentityServiceV1UserGender(this.json);

  factory IdentityServiceV1UserGender.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<IdentityServiceV1UserGender> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
