// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// Incorrect name has been replaced. Original name: `filterExpr.type`.
@JsonEnum()
enum Enum3 {
  @JsonValue('EXPR_TYPE_UNSPECIFIED')
  exprTypeUnspecified('EXPR_TYPE_UNSPECIFIED'),
  @JsonValue('AND')
  and('AND'),
  @JsonValue('OR')
  or('OR'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const Enum3(this.json);

  factory Enum3.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<Enum3> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
