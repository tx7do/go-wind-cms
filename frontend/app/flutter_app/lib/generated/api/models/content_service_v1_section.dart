// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'content_service_v1_section_type.dart';

part 'content_service_v1_section.g.dart';

/// 页面区块（内容构建模块）
@JsonSerializable()
class ContentServiceV1Section {
  const ContentServiceV1Section({
    this.type,
    this.name,
    this.sortOrder,
    this.config,
    this.content,
  });
  
  factory ContentServiceV1Section.fromJson(Map<String, Object?> json) => _$ContentServiceV1SectionFromJson(json);
  
  /// 区块类型
  final ContentServiceV1SectionType? type;

  /// 区块名称（后台标识用）
  final String? name;

  /// 排序（越小越靠前）
  final int? sortOrder;

  /// 区块通用配置（边距、样式、类名等）
  final Map<String, String>? config;

  /// 区块内容（键值对，适配不同类型区块）
  final Map<String, String>? content;

  Map<String, Object?> toJson() => _$ContentServiceV1SectionToJson(this);
}
