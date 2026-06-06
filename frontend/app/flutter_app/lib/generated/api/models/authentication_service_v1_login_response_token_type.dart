// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 令牌的类型，该值大小写不敏感，必选项，可以是bearer类型或mac类型，通常只是字符串“Bearer”。
@JsonEnum()
enum AuthenticationServiceV1LoginResponseTokenType {
  @JsonValue('bearer')
  bearer('bearer'),
  @JsonValue('mac')
  mac('mac'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const AuthenticationServiceV1LoginResponseTokenType(this.json);

  factory AuthenticationServiceV1LoginResponseTokenType.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<AuthenticationServiceV1LoginResponseTokenType> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
