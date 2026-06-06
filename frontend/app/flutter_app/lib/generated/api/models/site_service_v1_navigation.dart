// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'site_service_v1_navigation_item.dart';
import 'site_service_v1_navigation_location.dart';

part 'site_service_v1_navigation.g.dart';

/// 导航
@JsonSerializable()
class SiteServiceV1Navigation {
  const SiteServiceV1Navigation({
    this.id,
    this.name,
    this.location,
    this.locale,
    this.isActive,
    this.items,
    this.createdBy,
    this.updatedBy,
    this.deletedBy,
    this.createdAt,
    this.updatedAt,
    this.deletedAt,
  });
  
  factory SiteServiceV1Navigation.fromJson(Map<String, Object?> json) => _$SiteServiceV1NavigationFromJson(json);
  
  /// 导航ID
  final int? id;

  /// 导航名称（如'主导航'、'页脚'）
  final String? name;

  /// 渲染位置
  final SiteServiceV1NavigationLocation? location;

  /// 关联的语言区域（如 'zh-CN', 'en-US'），空值表示默认/通用导航；仅用于后台管理筛选，不影响前端渲染逻辑
  final String? locale;

  /// 是否启用（禁用后前端不渲染）
  final bool? isActive;

  /// 导航项列表（按 sort_order 排序）
  final List<SiteServiceV1NavigationItem>? items;

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

  Map<String, Object?> toJson() => _$SiteServiceV1NavigationToJson(this);
}
