// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 页面类型（首页/404/普通页面等，影响路由和渲染逻辑）
@JsonEnum()
enum ContentServiceV1PageType {
  @JsonValue('PAGE_TYPE_UNSPECIFIED')
  pageTypeUnspecified('PAGE_TYPE_UNSPECIFIED'),
  @JsonValue('PAGE_TYPE_DEFAULT')
  pageTypeDefault('PAGE_TYPE_DEFAULT'),
  @JsonValue('PAGE_TYPE_HOME')
  pageTypeHome('PAGE_TYPE_HOME'),
  @JsonValue('PAGE_TYPE_ERROR_404')
  pageTypeError404('PAGE_TYPE_ERROR_404'),
  @JsonValue('PAGE_TYPE_ERROR_500')
  pageTypeError500('PAGE_TYPE_ERROR_500'),
  @JsonValue('PAGE_TYPE_CUSTOM')
  pageTypeCustom('PAGE_TYPE_CUSTOM'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const ContentServiceV1PageType(this.json);

  factory ContentServiceV1PageType.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<ContentServiceV1PageType> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
