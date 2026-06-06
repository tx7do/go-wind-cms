// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 渲染位置
@JsonEnum()
enum SiteServiceV1NavigationLocation {
  @JsonValue('NAVIGATION_LOCATION_UNSPECIFIED')
  navigationLocationUnspecified('NAVIGATION_LOCATION_UNSPECIFIED'),
  @JsonValue('HEADER')
  header('HEADER'),
  @JsonValue('FOOTER')
  footer('FOOTER'),
  @JsonValue('SIDEBAR')
  sidebar('SIDEBAR'),
  @JsonValue('MOBILE')
  mobile('MOBILE'),
  @JsonValue('TOP_BAR')
  topBar('TOP_BAR'),
  @JsonValue('OFFCANVAS')
  offcanvas('OFFCANVAS'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const SiteServiceV1NavigationLocation(this.json);

  factory SiteServiceV1NavigationLocation.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<SiteServiceV1NavigationLocation> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
