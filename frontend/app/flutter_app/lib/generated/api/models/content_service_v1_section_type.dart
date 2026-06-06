// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

/// 区块类型
@JsonEnum()
enum ContentServiceV1SectionType {
  @JsonValue('SECTION_TYPE_UNSPECIFIED')
  sectionTypeUnspecified('SECTION_TYPE_UNSPECIFIED'),
  @JsonValue('SECTION_TYPE_RICH_TEXT')
  sectionTypeRichText('SECTION_TYPE_RICH_TEXT'),
  @JsonValue('SECTION_TYPE_MARKDOWN')
  sectionTypeMarkdown('SECTION_TYPE_MARKDOWN'),
  @JsonValue('SECTION_TYPE_TITLE')
  sectionTypeTitle('SECTION_TYPE_TITLE'),
  @JsonValue('SECTION_TYPE_IMAGE')
  sectionTypeImage('SECTION_TYPE_IMAGE'),
  @JsonValue('SECTION_TYPE_GALLERY')
  sectionTypeGallery('SECTION_TYPE_GALLERY'),
  @JsonValue('SECTION_TYPE_VIDEO')
  sectionTypeVideo('SECTION_TYPE_VIDEO'),
  @JsonValue('SECTION_TYPE_BUTTON')
  sectionTypeButton('SECTION_TYPE_BUTTON'),
  @JsonValue('SECTION_TYPE_DIVIDER')
  sectionTypeDivider('SECTION_TYPE_DIVIDER'),
  @JsonValue('SECTION_TYPE_SPACER')
  sectionTypeSpacer('SECTION_TYPE_SPACER'),
  @JsonValue('SECTION_TYPE_CODE')
  sectionTypeCode('SECTION_TYPE_CODE'),
  @JsonValue('SECTION_TYPE_HTML')
  sectionTypeHtml('SECTION_TYPE_HTML'),
  @JsonValue('SECTION_TYPE_FORM')
  sectionTypeForm('SECTION_TYPE_FORM'),
  @JsonValue('SECTION_TYPE_CAROUSEL')
  sectionTypeCarousel('SECTION_TYPE_CAROUSEL'),
  @JsonValue('SECTION_TYPE_CUSTOM')
  sectionTypeCustom('SECTION_TYPE_CUSTOM'),
  /// Default value for all unparsed values, allows backward compatibility when adding new values on the backend.
  $unknown(null);

  const ContentServiceV1SectionType(this.json);

  factory ContentServiceV1SectionType.fromJson(String json) => values.firstWhere(
        (e) => e.json == json,
        orElse: () => $unknown,
      );

  final String? json;

  @override
  String toString() => json?.toString() ?? super.toString();
  /// Returns all defined enum values excluding the $unknown value.
  static List<ContentServiceV1SectionType> get $valuesDefined => values.where((value) => value != $unknown).toList();
}
