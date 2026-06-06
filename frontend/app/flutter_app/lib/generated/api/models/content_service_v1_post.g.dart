// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'content_service_v1_post.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

ContentServiceV1Post _$ContentServiceV1PostFromJson(
  Map<String, dynamic> json,
) => ContentServiceV1Post(
  id: (json['id'] as num?)?.toInt(),
  status: json['status'] == null
      ? null
      : ContentServiceV1PostStatus.fromJson(json['status'] as String),
  editorType: json['editorType'] == null
      ? null
      : ContentServiceV1PostEditorType.fromJson(json['editorType'] as String),
  code: json['code'] as String?,
  disallowComment: json['disallowComment'] as bool?,
  inProgress: json['inProgress'] as bool?,
  autoSummary: json['autoSummary'] as bool?,
  isFeatured: json['isFeatured'] as bool?,
  sortOrder: (json['sortOrder'] as num?)?.toInt(),
  visits: (json['visits'] as num?)?.toInt(),
  likes: (json['likes'] as num?)?.toInt(),
  commentCount: (json['commentCount'] as num?)?.toInt(),
  authorId: (json['authorId'] as num?)?.toInt(),
  authorName: json['authorName'] as String?,
  customFields: (json['customFields'] as Map<String, dynamic>?)?.map(
    (k, e) => MapEntry(k, e as String),
  ),
  translations: (json['translations'] as List<dynamic>?)
      ?.map(
        (e) =>
            ContentServiceV1PostTranslation.fromJson(e as Map<String, dynamic>),
      )
      .toList(),
  availableLanguages: (json['availableLanguages'] as List<dynamic>?)
      ?.map((e) => e as String)
      .toList(),
  categoryIds: (json['categoryIds'] as List<dynamic>?)
      ?.map((e) => (e as num).toInt())
      .toList(),
  tagIds: (json['tagIds'] as List<dynamic>?)
      ?.map((e) => (e as num).toInt())
      .toList(),
  passwordHash: json['passwordHash'] as String?,
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
  publishTime: json['publishTime'] == null
      ? null
      : DateTime.parse(json['publishTime'] as String),
);

Map<String, dynamic> _$ContentServiceV1PostToJson(
  ContentServiceV1Post instance,
) => <String, dynamic>{
  'id': instance.id,
  'status': _$ContentServiceV1PostStatusEnumMap[instance.status],
  'editorType': _$ContentServiceV1PostEditorTypeEnumMap[instance.editorType],
  'code': instance.code,
  'disallowComment': instance.disallowComment,
  'inProgress': instance.inProgress,
  'autoSummary': instance.autoSummary,
  'isFeatured': instance.isFeatured,
  'sortOrder': instance.sortOrder,
  'visits': instance.visits,
  'likes': instance.likes,
  'commentCount': instance.commentCount,
  'authorId': instance.authorId,
  'authorName': instance.authorName,
  'customFields': instance.customFields,
  'translations': instance.translations,
  'availableLanguages': instance.availableLanguages,
  'categoryIds': instance.categoryIds,
  'tagIds': instance.tagIds,
  'passwordHash': instance.passwordHash,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
  'publishTime': instance.publishTime?.toIso8601String(),
};

const _$ContentServiceV1PostStatusEnumMap = {
  ContentServiceV1PostStatus.postStatusUnspecified: 'POST_STATUS_UNSPECIFIED',
  ContentServiceV1PostStatus.postStatusDraft: 'POST_STATUS_DRAFT',
  ContentServiceV1PostStatus.postStatusPublished: 'POST_STATUS_PUBLISHED',
  ContentServiceV1PostStatus.postStatusScheduled: 'POST_STATUS_SCHEDULED',
  ContentServiceV1PostStatus.postStatusTrashed: 'POST_STATUS_TRASHED',
  ContentServiceV1PostStatus.$unknown: r'$unknown',
};

const _$ContentServiceV1PostEditorTypeEnumMap = {
  ContentServiceV1PostEditorType.editorTypeUnspecified:
      'EDITOR_TYPE_UNSPECIFIED',
  ContentServiceV1PostEditorType.editorTypeMarkdown: 'EDITOR_TYPE_MARKDOWN',
  ContentServiceV1PostEditorType.editorTypeRichText: 'EDITOR_TYPE_RICH_TEXT',
  ContentServiceV1PostEditorType.editorTypePlainText: 'EDITOR_TYPE_PLAIN_TEXT',
  ContentServiceV1PostEditorType.editorTypeCode: 'EDITOR_TYPE_CODE',
  ContentServiceV1PostEditorType.editorTypeJsonBlock: 'EDITOR_TYPE_JSON_BLOCK',
  ContentServiceV1PostEditorType.editorTypeVisualBuilder:
      'EDITOR_TYPE_VISUAL_BUILDER',
  ContentServiceV1PostEditorType.$unknown: r'$unknown',
};
