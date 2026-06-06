// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 授权类型，此处的值固定为"password"，必选项。
@JsonEnum()
enum AuthenticationServiceV1LoginRequestGrantType {
  @JsonValue('password')
  password('password'),
  @JsonValue('client_credentials')
  clientCredentials('client_credentials'),
  @JsonValue('authorization_code')
  authorizationCode('authorization_code'),
  @JsonValue('refresh_token')
  refreshToken('refresh_token'),
  @JsonValue('implicit')
  implicit('implicit'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const AuthenticationServiceV1LoginRequestGrantType(this.json);

  factory AuthenticationServiceV1LoginRequestGrantType.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<AuthenticationServiceV1LoginRequestGrantType> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
