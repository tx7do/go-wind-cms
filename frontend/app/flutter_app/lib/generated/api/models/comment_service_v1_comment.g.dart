// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'comment_service_v1_comment.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

CommentServiceV1Comment _$CommentServiceV1CommentFromJson(
  Map<String, dynamic> json,
) => CommentServiceV1Comment(
  id: (json['id'] as num?)?.toInt(),
  contentType: json['contentType'] == null
      ? null
      : CommentServiceV1CommentContentType.fromJson(
          json['contentType'] as String,
        ),
  objectId: (json['objectId'] as num?)?.toInt(),
  content: json['content'] as String?,
  authorId: (json['authorId'] as num?)?.toInt(),
  authorName: json['authorName'] as String?,
  authorEmail: json['authorEmail'] as String?,
  authorUrl: json['authorUrl'] as String?,
  authorType: json['authorType'] == null
      ? null
      : CommentServiceV1CommentAuthorType.fromJson(
          json['authorType'] as String,
        ),
  status: json['status'] == null
      ? null
      : CommentServiceV1CommentStatus.fromJson(json['status'] as String),
  likeCount: (json['likeCount'] as num?)?.toInt(),
  dislikeCount: (json['dislikeCount'] as num?)?.toInt(),
  replyCount: (json['replyCount'] as num?)?.toInt(),
  ipAddress: json['ipAddress'] as String?,
  location: json['location'] as String?,
  userAgent: json['userAgent'] as String?,
  detectedLanguage: json['detectedLanguage'] as String?,
  isSpam: json['isSpam'] as bool?,
  isSticky: json['isSticky'] as bool?,
  parentId: (json['parentId'] as num?)?.toInt(),
  children: (json['children'] as List<dynamic>?)
      ?.map((e) => CommentServiceV1Comment.fromJson(e as Map<String, dynamic>))
      .toList(),
  replyToId: (json['replyToId'] as num?)?.toInt(),
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

Map<String, dynamic> _$CommentServiceV1CommentToJson(
  CommentServiceV1Comment instance,
) => <String, dynamic>{
  'id': instance.id,
  'contentType':
      _$CommentServiceV1CommentContentTypeEnumMap[instance.contentType],
  'objectId': instance.objectId,
  'content': instance.content,
  'authorId': instance.authorId,
  'authorName': instance.authorName,
  'authorEmail': instance.authorEmail,
  'authorUrl': instance.authorUrl,
  'authorType': _$CommentServiceV1CommentAuthorTypeEnumMap[instance.authorType],
  'status': _$CommentServiceV1CommentStatusEnumMap[instance.status],
  'likeCount': instance.likeCount,
  'dislikeCount': instance.dislikeCount,
  'replyCount': instance.replyCount,
  'ipAddress': instance.ipAddress,
  'location': instance.location,
  'userAgent': instance.userAgent,
  'detectedLanguage': instance.detectedLanguage,
  'isSpam': instance.isSpam,
  'isSticky': instance.isSticky,
  'parentId': instance.parentId,
  'children': instance.children,
  'replyToId': instance.replyToId,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
};

const _$CommentServiceV1CommentContentTypeEnumMap = {
  CommentServiceV1CommentContentType.contentTypeUnspecified:
      'CONTENT_TYPE_UNSPECIFIED',
  CommentServiceV1CommentContentType.contentTypePost: 'CONTENT_TYPE_POST',
  CommentServiceV1CommentContentType.contentTypePage: 'CONTENT_TYPE_PAGE',
  CommentServiceV1CommentContentType.contentTypeProduct: 'CONTENT_TYPE_PRODUCT',
  CommentServiceV1CommentContentType.$unknown: r'$unknown',
};

const _$CommentServiceV1CommentAuthorTypeEnumMap = {
  CommentServiceV1CommentAuthorType.authorTypeUnspecified:
      'AUTHOR_TYPE_UNSPECIFIED',
  CommentServiceV1CommentAuthorType.authorTypeGuest: 'AUTHOR_TYPE_GUEST',
  CommentServiceV1CommentAuthorType.authorTypeUser: 'AUTHOR_TYPE_USER',
  CommentServiceV1CommentAuthorType.authorTypeAdmin: 'AUTHOR_TYPE_ADMIN',
  CommentServiceV1CommentAuthorType.authorTypeModerator:
      'AUTHOR_TYPE_MODERATOR',
  CommentServiceV1CommentAuthorType.$unknown: r'$unknown',
};

const _$CommentServiceV1CommentStatusEnumMap = {
  CommentServiceV1CommentStatus.statusUnspecified: 'STATUS_UNSPECIFIED',
  CommentServiceV1CommentStatus.statusPending: 'STATUS_PENDING',
  CommentServiceV1CommentStatus.statusApproved: 'STATUS_APPROVED',
  CommentServiceV1CommentStatus.statusRejected: 'STATUS_REJECTED',
  CommentServiceV1CommentStatus.statusSpam: 'STATUS_SPAM',
  CommentServiceV1CommentStatus.$unknown: r'$unknown',
};
