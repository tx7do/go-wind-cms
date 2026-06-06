// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_tag_status.dart';
import 'content_service_v1_tag_translation.dart';

part 'content_service_v1_tag.g.dart';

/// 标签
@JsonSerializable()
class ContentServiceV1Tag {
  const ContentServiceV1Tag({
    this.id,
    this.status,
    this.color,
    this.icon,
    this.group,
    this.sortOrder,
    this.isFeatured,
    this.code,
    this.postCount,
    this.translations,
    this.availableLanguages,
    this.createdBy,
    this.updatedBy,
    this.deletedBy,
    this.createdAt,
    this.updatedAt,
    this.deletedAt,
  });
  
  factory ContentServiceV1Tag.fromJson(Map<String, Object?> json) => _$ContentServiceV1TagFromJson(json);
  
  /// 标签ID
  final int? id;

  /// 标签状态
  final ContentServiceV1TagStatus? status;

  /// 标签颜色（HEX 格式，如 '#3498db'，用于标签云/前端展示）
  final String? color;

  /// 标签图标（可选，支持图标名称如 'fas fa-code' 或 SVG URL）
  final String? icon;

  /// 标签分组（可选，如 '技术栈'、'行业'，用于后台分类管理，不影响前台展示）
  final String? group;

  /// 排序优先级（数值越小越靠前，同组内排序）
  final int? sortOrder;

  /// 是否推荐
  final bool? isFeatured;

  /// 唯一代码（如 slug、编码等，便于唯一标识标签）
  final String? code;

  /// 使用该标签的文章总数（跨所有语言，仅统计 status=ACTIVE 的文章）
  final int? postCount;

  /// 多语言翻译列表
  final List<ContentServiceV1TagTranslation>? translations;

  /// 可用的语言代码列表（快速查询，避免遍历 translations）
  final List<String>? availableLanguages;

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

  Map<String, Object?> toJson() => _$ContentServiceV1TagToJson(this);
}
