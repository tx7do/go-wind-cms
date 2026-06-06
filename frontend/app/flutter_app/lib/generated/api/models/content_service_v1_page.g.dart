// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_page.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1Page _$ContentServiceV1PageFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1Page(
  id: (json['id'] as num?)?.toInt(),
  status: json['status'] == null
      ? null
      : ContentServiceV1PageStatus.fromJson(json['status'] as String),
  type: json['type'] == null
      ? null
      : ContentServiceV1PageType.fromJson(json['type'] as String),
  editorType: json['editorType'] == null
      ? null
      : ContentServiceV1PageEditorType.fromJson(json['editorType'] as String),
  slug: json['slug'] as String?,
  authorId: (json['authorId'] as num?)?.toInt(),
  authorName: json['authorName'] as String?,
  disallowComment: json['disallowComment'] as bool?,
  redirectUrl: json['redirectUrl'] as String?,
  showInNavigation: json['showInNavigation'] as bool?,
  sortOrder: (json['sortOrder'] as num?)?.toInt(),
  template: json['template'] as String?,
  isCustomTemplate: json['isCustomTemplate'] as bool?,
  visits: (json['visits'] as num?)?.toInt(),
  customFields: (json['customFields'] as Map<String, dynamic>?)?.map(
    (k, e) => MapEntry(k, e as String),
  ),
  translations: (json['translations'] as List<dynamic>?)
      ?.map(
        (e) =>
            ContentServiceV1PageTranslation.fromJson(e as Map<String, dynamic>),
      )
      .toList(),
  availableLanguages: (json['availableLanguages'] as List<dynamic>?)
      ?.map((e) => e as String)
      .toList(),
  parentId: (json['parentId'] as num?)?.toInt(),
  children: (json['children'] as List<dynamic>?)
      ?.map((e) => ContentServiceV1Page.fromJson(e as Map<String, dynamic>))
      .toList(),
  depth: (json['depth'] as num?)?.toInt(),
  path: json['path'] as String?,
  createdBy: (json['createdBy'] as num?)?.toInt(),
  updatedBy: (json['updatedBy'] as num?)?.toInt(),
  deletedBy: (json['deletedBy'] as num?)?.toInt(),
  createdAt: json['createdAt'] == null
      ? null
      : DateTime.parse(json['createdAt'] as String),
  updatedAt: json['updatedAt'] == null
      ? null
      : DateTime.parse(json['updatedAt'] as String),
  deletedAt: json['deletedAt'] == null
      ? null
      : DateTime.parse(json['deletedAt'] as String),
);

Map<String, dynamic> _$ContentServiceV1PageToJson(
  ContentServiceV1Page instance,
) => <String, dynamic>{
  'id': instance.id,
  'status': _$ContentServiceV1PageStatusEnumMap[instance.status],
  'type': _$ContentServiceV1PageTypeEnumMap[instance.type],
  'editorType': _$ContentServiceV1PageEditorTypeEnumMap[instance.editorType],
  'slug': instance.slug,
  'authorId': instance.authorId,
  'authorName': instance.authorName,
  'disallowComment': instance.disallowComment,
  'redirectUrl': instance.redirectUrl,
  'showInNavigation': instance.showInNavigation,
  'sortOrder': instance.sortOrder,
  'template': instance.template,
  'isCustomTemplate': instance.isCustomTemplate,
  'visits': instance.visits,
  'customFields': instance.customFields,
  'translations': instance.translations,
  'availableLanguages': instance.availableLanguages,
  'parentId': instance.parentId,
  'children': instance.children,
  'depth': instance.depth,
  'path': instance.path,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
};

const _$ContentServiceV1PageStatusEnumMap = {
  ContentServiceV1PageStatus.pageStatusUnspecified: 'PAGE_STATUS_UNSPECIFIED',
  ContentServiceV1PageStatus.pageStatusDraft: 'PAGE_STATUS_DRAFT',
  ContentServiceV1PageStatus.pageStatusPublished: 'PAGE_STATUS_PUBLISHED',
  ContentServiceV1PageStatus.pageStatusArchived: 'PAGE_STATUS_ARCHIVED',
  ContentServiceV1PageStatus.$unknown: r'$unknown',
};

const _$ContentServiceV1PageTypeEnumMap = {
  ContentServiceV1PageType.pageTypeUnspecified: 'PAGE_TYPE_UNSPECIFIED',
  ContentServiceV1PageType.pageTypeDefault: 'PAGE_TYPE_DEFAULT',
  ContentServiceV1PageType.pageTypeHome: 'PAGE_TYPE_HOME',
  ContentServiceV1PageType.pageTypeError404: 'PAGE_TYPE_ERROR_404',
  ContentServiceV1PageType.pageTypeError500: 'PAGE_TYPE_ERROR_500',
  ContentServiceV1PageType.pageTypeCustom: 'PAGE_TYPE_CUSTOM',
  ContentServiceV1PageType.$unknown: r'$unknown',
};

const _$ContentServiceV1PageEditorTypeEnumMap = {
  ContentServiceV1PageEditorType.editorTypeUnspecified:
      'EDITOR_TYPE_UNSPECIFIED',
  ContentServiceV1PageEditorType.editorTypeMarkdown: 'EDITOR_TYPE_MARKDOWN',
  ContentServiceV1PageEditorType.editorTypeRichText: 'EDITOR_TYPE_RICH_TEXT',
  ContentServiceV1PageEditorType.editorTypePlainText: 'EDITOR_TYPE_PLAIN_TEXT',
  ContentServiceV1PageEditorType.editorTypeCode: 'EDITOR_TYPE_CODE',
  ContentServiceV1PageEditorType.editorTypeJsonBlock: 'EDITOR_TYPE_JSON_BLOCK',
  ContentServiceV1PageEditorType.editorTypeVisualBuilder:
      'EDITOR_TYPE_VISUAL_BUILDER',
  ContentServiceV1PageEditorType.$unknown: r'$unknown',
};
