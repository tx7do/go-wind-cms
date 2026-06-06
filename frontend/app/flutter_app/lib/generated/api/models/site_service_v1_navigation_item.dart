// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'site_service_v1_navigation_item.dart';
import 'site_service_v1_navigation_item_link_type.dart';

part 'site_service_v1_navigation_item.g.dart';

/// 导航项目项
@JsonSerializable()
class SiteServiceV1NavigationItem {
  const SiteServiceV1NavigationItem({
    this.id,
    this.navigationId,
    this.title,
    this.description,
    this.icon,
    this.url,
    this.linkType,
    this.objectId,
    this.sortOrder,
    this.isOpenNewTab,
    this.isInvalid,
    this.requiredPermission,
    this.parentId,
    this.children,
    this.createdBy,
    this.updatedBy,
    this.deletedBy,
    this.createdAt,
    this.updatedAt,
    this.deletedAt,
  });
  
  factory SiteServiceV1NavigationItem.fromJson(Map<String, Object?> json) => _$SiteServiceV1NavigationItemFromJson(json);
  
  /// 导航项目项ID
  final int? id;

  /// 所属导航项目菜单ID
  final int? navigationId;

  /// 显示文本（单语言，多语言通过多套导航项目实现）
  final String? title;

  /// 描述文本（用于悬停提示或移动端）
  final String? description;

  /// 图标（支持 Font Awesome 格式如 'fas fa-home'，或 SVG URL）
  final String? icon;

  /// 目标 URL（link_type=CUSTOM/EXTERNAL 时必填；关联类型可自动生成）
  final String? url;

  /// 链接类型
  final SiteServiceV1NavigationItemLinkType? linkType;

  /// 关联对象ID（link_type=POST/PAGE/CATEGORY 时使用）
  final int? objectId;

  /// 排序优先级（数值越小越靠前）
  final int? sortOrder;

  /// 是否在新标签页打开
  final bool? isOpenNewTab;

  /// 是否无效（如关联的文章/页面被删除或设置为不可见，前端可根据此字段决定是否隐藏该导航项目项）
  final bool? isInvalid;

  /// 访问权限标识（如 view_admin_panel，前端可根据用户权限决定是否显示该导航项目项）
  final String? requiredPermission;

  /// 父菜单项ID（0 表示顶级项）
  final int? parentId;

  /// 子节点树
  final List<SiteServiceV1NavigationItem>? children;

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

  Map<String, Object?> toJson() => _$SiteServiceV1NavigationItemToJson(this);
}
