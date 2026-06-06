// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_section.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1Section _$ContentServiceV1SectionFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1Section(
  type: json['type'] == null
      ? null
      : ContentServiceV1SectionType.fromJson(json['type'] as String),
  name: json['name'] as String?,
  sortOrder: (json['sortOrder'] as num?)?.toInt(),
  config: (json['config'] as Map<String, dynamic>?)?.map(
    (k, e) => MapEntry(k, e as String),
  ),
  content: (json['content'] as Map<String, dynamic>?)?.map(
    (k, e) => MapEntry(k, e as String),
  ),
);

Map<String, dynamic> _$ContentServiceV1SectionToJson(
  ContentServiceV1Section instance,
) => <String, dynamic>{
  'type': _$ContentServiceV1SectionTypeEnumMap[instance.type],
  'name': instance.name,
  'sortOrder': instance.sortOrder,
  'config': instance.config,
  'content': instance.content,
};

const _$ContentServiceV1SectionTypeEnumMap = {
  ContentServiceV1SectionType.sectionTypeUnspecified:
      'SECTION_TYPE_UNSPECIFIED',
  ContentServiceV1SectionType.sectionTypeRichText: 'SECTION_TYPE_RICH_TEXT',
  ContentServiceV1SectionType.sectionTypeMarkdown: 'SECTION_TYPE_MARKDOWN',
  ContentServiceV1SectionType.sectionTypeTitle: 'SECTION_TYPE_TITLE',
  ContentServiceV1SectionType.sectionTypeImage: 'SECTION_TYPE_IMAGE',
  ContentServiceV1SectionType.sectionTypeGallery: 'SECTION_TYPE_GALLERY',
  ContentServiceV1SectionType.sectionTypeVideo: 'SECTION_TYPE_VIDEO',
  ContentServiceV1SectionType.sectionTypeButton: 'SECTION_TYPE_BUTTON',
  ContentServiceV1SectionType.sectionTypeDivider: 'SECTION_TYPE_DIVIDER',
  ContentServiceV1SectionType.sectionTypeSpacer: 'SECTION_TYPE_SPACER',
  ContentServiceV1SectionType.sectionTypeCode: 'SECTION_TYPE_CODE',
  ContentServiceV1SectionType.sectionTypeHtml: 'SECTION_TYPE_HTML',
  ContentServiceV1SectionType.sectionTypeForm: 'SECTION_TYPE_FORM',
  ContentServiceV1SectionType.sectionTypeCarousel: 'SECTION_TYPE_CAROUSEL',
  ContentServiceV1SectionType.sectionTypeCustom: 'SECTION_TYPE_CUSTOM',
  ContentServiceV1SectionType.$unknown: r'$unknown',
};
