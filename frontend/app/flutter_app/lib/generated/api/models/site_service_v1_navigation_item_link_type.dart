// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 链接类型
@JsonEnum()
enum SiteServiceV1NavigationItemLinkType {
  @JsonValue('LINK_TYPE_UNSPECIFIED')
  linkTypeUnspecified('LINK_TYPE_UNSPECIFIED'),
  @JsonValue('LINK_TYPE_CUSTOM')
  linkTypeCustom('LINK_TYPE_CUSTOM'),
  @JsonValue('LINK_TYPE_POST')
  linkTypePost('LINK_TYPE_POST'),
  @JsonValue('LINK_TYPE_PAGE')
  linkTypePage('LINK_TYPE_PAGE'),
  @JsonValue('LINK_TYPE_CATEGORY')
  linkTypeCategory('LINK_TYPE_CATEGORY'),
  @JsonValue('LINK_TYPE_EXTERNAL')
  linkTypeExternal('LINK_TYPE_EXTERNAL'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const SiteServiceV1NavigationItemLinkType(this.json);

  factory SiteServiceV1NavigationItemLinkType.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<SiteServiceV1NavigationItemLinkType> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
