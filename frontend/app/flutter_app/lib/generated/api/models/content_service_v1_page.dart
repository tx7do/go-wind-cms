// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_page.dart';
import 'content_service_v1_page_editor_type.dart';
import 'content_service_v1_page_status.dart';
import 'content_service_v1_page_translation.dart';
import 'content_service_v1_page_type.dart';

part 'content_service_v1_page.g.dart';

/// 页面
@JsonSerializable()
class ContentServiceV1Page {
  const ContentServiceV1Page({
    this.id,
    this.status,
    this.type,
    this.editorType,
    this.slug,
    this.authorId,
    this.authorName,
    this.disallowComment,
    this.redirectUrl,
    this.showInNavigation,
    this.sortOrder,
    this.template,
    this.isCustomTemplate,
    this.visits,
    this.customFields,
    this.translations,
    this.availableLanguages,
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
  
  factory ContentServiceV1Page.fromJson(Map<String, Object?> json) => _$ContentServiceV1PageFromJson(json);
  
  /// 页面ID
  final int? id;

  /// 页面状态
  final ContentServiceV1PageStatus? status;

  /// 页面类型（首页/404/普通页面等，影响路由和渲染逻辑）
  final ContentServiceV1PageType? type;

  /// 编辑器类型
  final ContentServiceV1PageEditorType? editorType;

  /// 页面唯一标识（默认语言，URL 友好，决定主路由）
  final String? slug;

  /// 页面作者ID
  final int? authorId;

  /// 页面作者名称
  final String? authorName;

  /// 是否禁止评论（默认 true，页面通常关闭评论）
  final bool? disallowComment;

  /// 重定向 URL
  final String? redirectUrl;

  /// 是否在主导航中显示（用于自动生成菜单）
  final bool? showInNavigation;

  /// 导航排序（数值越小越靠前）
  final int? sortOrder;

  /// 页面模板名称（如 'about-us', 'contact'，覆盖全局默认模板）
  final String? template;

  /// 是否使用自定义模板代码（而非预设模板）
  final bool? isCustomTemplate;

  /// 页面访问次数
  final int? visits;

  /// 自定义字段，键值对形式，便于扩展
  final Map<String, String>? customFields;

  /// 多语言翻译列表
  final List<ContentServiceV1PageTranslation>? translations;

  /// 可用的语言代码列表（快速查询，避免遍历 translations）
  final List<String>? availableLanguages;

  /// 父页面ID（0 表示顶级页面，用于构建站点树形结构）
  final int? parentId;

  /// 子节点树
  final List<ContentServiceV1Page>? children;

  /// 页面层级深度（0=顶级，1=二级，用于面包屑导航）
  final int? depth;

  /// 物化路径（如 '1/5/23'，便于层级查询）
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

  Map<String, Object?> toJson() => _$ContentServiceV1PageToJson(this);
}
