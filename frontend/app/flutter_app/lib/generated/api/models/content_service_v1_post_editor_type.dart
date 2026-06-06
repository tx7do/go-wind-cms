// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 编辑器类型
@JsonEnum()
enum ContentServiceV1PostEditorType {
  @JsonValue('EDITOR_TYPE_UNSPECIFIED')
  editorTypeUnspecified('EDITOR_TYPE_UNSPECIFIED'),
  @JsonValue('EDITOR_TYPE_MARKDOWN')
  editorTypeMarkdown('EDITOR_TYPE_MARKDOWN'),
  @JsonValue('EDITOR_TYPE_RICH_TEXT')
  editorTypeRichText('EDITOR_TYPE_RICH_TEXT'),
  @JsonValue('EDITOR_TYPE_PLAIN_TEXT')
  editorTypePlainText('EDITOR_TYPE_PLAIN_TEXT'),
  @JsonValue('EDITOR_TYPE_CODE')
  editorTypeCode('EDITOR_TYPE_CODE'),
  @JsonValue('EDITOR_TYPE_JSON_BLOCK')
  editorTypeJsonBlock('EDITOR_TYPE_JSON_BLOCK'),
  @JsonValue('EDITOR_TYPE_VISUAL_BUILDER')
  editorTypeVisualBuilder('EDITOR_TYPE_VISUAL_BUILDER'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const ContentServiceV1PostEditorType(this.json);

  factory ContentServiceV1PostEditorType.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<ContentServiceV1PostEditorType> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
