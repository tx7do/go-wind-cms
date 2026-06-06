// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'comment_service_v1_comment.dart';
import 'comment_service_v1_comment_author_type.dart';
import 'comment_service_v1_comment_content_type.dart';
import 'comment_service_v1_comment_status.dart';

part 'comment_service_v1_comment.g.dart';

/// 评论
@JsonSerializable()
class CommentServiceV1Comment {
  const CommentServiceV1Comment({
    this.id,
    this.contentType,
    this.objectId,
    this.content,
    this.authorId,
    this.authorName,
    this.authorEmail,
    this.authorUrl,
    this.authorType,
    this.status,
    this.likeCount,
    this.dislikeCount,
    this.replyCount,
    this.ipAddress,
    this.location,
    this.userAgent,
    this.detectedLanguage,
    this.isSpam,
    this.isSticky,
    this.parentId,
    this.children,
    this.replyToId,
    this.createdBy,
    this.updatedBy,
    this.deletedBy,
    this.createdAt,
    this.updatedAt,
    this.deletedAt,
  });
  
  factory CommentServiceV1Comment.fromJson(Map<String, Object?> json) => _$CommentServiceV1CommentFromJson(json);
  
  /// 评论ID
  final int? id;

  /// 内容类型
  final CommentServiceV1CommentContentType? contentType;

  /// 对象ID
  final int? objectId;

  /// 评论内容
  final String? content;

  /// 评论作者ID，0表示游客
  final int? authorId;

  /// 评论作者名称（游客填写）
  final String? authorName;

  /// 评论作者邮箱（游客填写）
  final String? authorEmail;

  /// 评论作者网址（游客填写）
  final String? authorUrl;

  /// 作者类型（用于前端高亮管理员回复）
  final CommentServiceV1CommentAuthorType? authorType;

  /// 评论状态
  final CommentServiceV1CommentStatus? status;

  /// 点赞数
  final int? likeCount;

  /// 点踩数
  final int? dislikeCount;

  /// 回复数
  final int? replyCount;

  /// 评论者 IP（脱敏存储，如 192.168.1.*）
  final String? ipAddress;

  /// 评论者地理位置（根据 IP 反查，可能不准确）
  final String? location;

  /// User-Agent
  final String? userAgent;

  /// 自动检测的语言代码（ISO 639-1，如 'zh', 'en'）
  final String? detectedLanguage;

  /// 是否标记为垃圾评论
  final bool? isSpam;

  /// 是否置顶评论
  final bool? isSticky;

  /// 父节点ID
  final int? parentId;

  /// 子节点树
  final List<CommentServiceV1Comment>? children;

  /// 回复的评论ID（记录 '@某人' 关系，可能与 parent_id 不同）
  final int? replyToId;

  /// 创建者用户ID
  final int? createdBy;

  /// 更新者用户ID
  final int? updatedBy;

  /// 删除者用户ID
  final int? deletedBy;

  /// 创建时间
  final DateTime? createdAt;

  /// 更新时间
  final DateTime? updatedAt;

  /// 删除时间
  final DateTime? deletedAt;

  Map<String, Object?> toJson() => _$CommentServiceV1CommentToJson(this);
}
