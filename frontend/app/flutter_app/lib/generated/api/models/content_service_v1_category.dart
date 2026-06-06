// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_category.dart';
import 'content_service_v1_category_status.dart';
import 'content_service_v1_category_translation.dart';

part 'content_service_v1_category.g.dart';

/// 类别
@JsonSerializable()
class ContentServiceV1Category {
  const ContentServiceV1Category({
    this.id,
    this.status,
    this.sortOrder,
    this.isNav,
    this.icon,
    this.code,
    this.postCount,
    this.directPostCount,
    this.translations,
    this.availableLanguages,
    this.customFields,
    this.parentId,
    this.children,
    this.depth,
    this.path,
    this.createdBy,
    this.updatedBy,
    this.deletedBy,
    this.createdAt,
    this.updatedAt,
    this.deletedAt,
  });
  
  factory ContentServiceV1Category.fromJson(Map<String, Object?> json) => _$ContentServiceV1CategoryFromJson(json);
  
  /// 类别ID
  final int? id;

  /// 分类状态
  final ContentServiceV1CategoryStatus? status;

  /// 排序优先级（数值越小越靠前，同级分类间排序）
  final int? sortOrder;

  /// 是否显示在导航菜单
  final bool? isNav;

  /// 分类图标
  final String? icon;

  /// 唯一代码（如 slug、编码等，便于唯一标识分类）
  final String? code;

  /// 该分类下的文章总数（含子分类，可选）
  final int? postCount;

  /// 该分类下的直接文章数（不含子分类）
  final int? directPostCount;

  /// 多语言翻译列表
  final List<ContentServiceV1CategoryTranslation>? translations;

  /// 可用的语言代码列表（快速查询，避免遍历 translations）
  final List<String>? availableLanguages;

  /// 自定义字段，键值对形式，便于扩展
  final Map<String, String>? customFields;

  /// 父节点ID
  final int? parentId;

  /// 子节点树
  final List<ContentServiceV1Category>? children;

  /// 分类层级深度（0=顶级，1=二级，以此类推）
  final int? depth;

  /// 物化路径（Materialized Path），如 '1/5/23'，便于层级查询
  final String? path;

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

  Map<String, Object?> toJson() => _$ContentServiceV1CategoryToJson(this);
}
