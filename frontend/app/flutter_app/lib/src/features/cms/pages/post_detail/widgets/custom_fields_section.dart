import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_widget_from_html/flutter_widget_from_html.dart';
import 'package:get_it/get_it.dart' show GetIt;
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show
        ApiClient,
        ContentServiceV1ContentModel,
        ContentServiceV1ContentModelComponent,
        ContentServiceV1ComponentConfig,
        ContentServiceV1ContentReferenceSummary,
        ContentServiceV1FieldDefinition,
        ContentServiceV1FieldDefinition$Type,
        ContentServiceV1ResolveContentReferencesRequest,
        MediaServiceV1GetMediaAssetRequest,
        MediaServiceV1MediaAsset;
import 'package:flutter_app/src/core/utilities/logger.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';

/// 自定义字段渲染（内容建模前台读侧，flutter 端）：
/// 按文章分类链绑定的内容模型字段定义，渲染 post.customFields。
/// 值约定为字符串（admin 编辑器 `Record<string,string>` 口径）：
///   - RICHTEXT → HtmlWidget（flutter_widget_from_html）；
///   - IMAGE/FILE → `media_asset:<id>` 引用经公开媒体端点换 URL（直链 URL 直用）；
///   - RELATION → `post:<id>`/`page:<id>` 引用解析为公开摘要，可点击跳目标详情；
///     草稿/软删目标解析缺失，回退原样字符串；
///   - 其余（TEXT/NUMBER 等）原样展示。
/// 模型未绑定或加载失败静默不渲染（自定义字段是增强信息，不阻断正文）。
class CustomFieldsSection extends StatefulWidget {
  final int postId;
  final Map<String, String>? customFields;
  final bool isMobile;

  const CustomFieldsSection({
    super.key,
    required this.postId,
    required this.customFields,
    required this.isMobile,
  });

  @override
  State<CustomFieldsSection> createState() => _CustomFieldsSectionState();
}

class _CustomFieldsSectionState extends State<CustomFieldsSection> {
  final _mediaLogger = Logger('CustomFieldsSection.media');
  late final Future<ContentServiceV1ContentModel?> _modelFuture;

  @override
  void initState() {
    super.initState();
    _modelFuture = PostService().getModelByPost(postId: widget.postId);
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<ContentServiceV1ContentModel?>(
      future: _modelFuture,
      builder: (context, snapshot) {
        final model = snapshot.data;
        final fields = model?.fields ?? const <ContentServiceV1FieldDefinition>[];
        final custom = widget.customFields ?? const {};
        if (fields.isEmpty || custom.isEmpty) return const SizedBox.shrink();

        // 仅渲染模型定义过且文章实际携带的字段，按模型字段顺序
        final rows = <_Row>[];
        for (final f in fields) {
          final v = custom[f.name ?? ''];
          if (v == null || v.isEmpty) continue;
          rows.add(_Row(
            name: f.name ?? '',
            label: (f.label?.isNotEmpty ?? false) ? f.label! : (f.name ?? ''),
            type: f.type ?? ContentServiceV1FieldDefinition$Type.fieldTypeText,
            value: v,
          ));
        }
        if (rows.isEmpty) return const SizedBox.shrink();

        final theme = Theme.of(context);
        return Container(
          margin: EdgeInsets.only(bottom: widget.isMobile ? 24.h : 24),
          padding: EdgeInsets.all(widget.isMobile ? 14.w : 14),
          decoration: BoxDecoration(
            color: theme.colorScheme.surfaceContainerHighest.withAlpha(80),
            borderRadius: BorderRadius.circular(widget.isMobile ? 12.r : 12),
            border: Border.all(color: theme.colorScheme.outlineVariant),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                model?.name ?? '',
                style: TextStyle(
                  fontSize: widget.isMobile ? 15.sp : 15,
                  fontWeight: FontWeight.w600,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              SizedBox(height: widget.isMobile ? 10.h : 10),
              for (final row in rows)
                Padding(
                  padding: EdgeInsets.only(bottom: widget.isMobile ? 10.h : 10),
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      SizedBox(
                        width: widget.isMobile ? 96.w : 120,
                        child: Text(
                          row.label,
                          style: TextStyle(
                            fontSize: widget.isMobile ? 12.sp : 12,
                            color: theme.colorScheme.onSurface.withAlpha(150),
                          ),
                        ),
                      ),
                      Expanded(child: _buildValue(context, model, row)),
                    ],
                  ),
                ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildValue(
      BuildContext context, ContentServiceV1ContentModel? model, _Row row) {
    final theme = Theme.of(context);
    switch (row.type) {
      case ContentServiceV1FieldDefinition$Type.fieldTypeRichtext:
        return HtmlWidget(
          row.value,
          textStyle: TextStyle(
            fontSize: widget.isMobile ? 13.sp : 13,
            color: theme.colorScheme.onSurface,
            height: 1.6,
          ),
        );
      case ContentServiceV1FieldDefinition$Type.fieldTypeImage:
      case ContentServiceV1FieldDefinition$Type.fieldTypeFile:
        final ref = _parseMediaRef(row.value);
        if (ref != null) {
          return FutureBuilder<MediaServiceV1MediaAsset?>(
            future: _fetchMediaAsset(ref),
            builder: (context, snap) {
              final url = snap.data?.url ?? '';
              if (row.type ==
                      ContentServiceV1FieldDefinition$Type.fieldTypeImage &&
                  url.isNotEmpty) {
                return ClipRRect(
                  borderRadius: BorderRadius.circular(widget.isMobile ? 12.r : 12),
                  child: Image.network(
                    url,
                    width: double.infinity,
                    fit: BoxFit.cover,
                    errorBuilder: (_, _, _) =>
                        Text('#${row.value}', style: _mutedStyle(theme)),
                  ),
                );
              }
              if (url.isNotEmpty) {
                return Text(url, style: _linkStyle(theme));
              }
              return Text('media_asset:$ref', style: _mutedStyle(theme));
            },
          );
        }
        if (row.value.startsWith('http')) {
          return row.type == ContentServiceV1FieldDefinition$Type.fieldTypeImage
              ? Image.network(
                  row.value,
                  width: double.infinity,
                  fit: BoxFit.cover,
                  errorBuilder: (_, _, _) =>
                      Text(row.value, style: _mutedStyle(theme)),
                )
              : Text(row.value, style: _linkStyle(theme));
        }
        return Text(row.value, style: _mutedStyle(theme));
      case ContentServiceV1FieldDefinition$Type.fieldTypeRelation:
        final m = RegExp(r'^(post|page):(\d+)$').firstMatch(row.value.trim());
        if (m == null) return Text(row.value, style: _regularStyle(theme));
        return FutureBuilder<ContentServiceV1ContentReferenceSummary?>(
          future: _resolveReference(row.value.trim()),
          builder: (context, snap) {
            final target = snap.data;
            final title = target?.title ?? '';
            if (target == null || title.isEmpty) {
              // 解析缺失（草稿/软删/不可见）：回退原样字符串
              return Text(row.value, style: _mutedStyle(theme));
            }
            final id = target.targetId ?? 0;
            // flutter 端 post/page 详情均按 id 路由
            final route =
                target.targetEntityType == 'page' ? '/page/$id' : '/post/$id';
            return GestureDetector(
              onTap: id > 0 ? () => context.go(route) : null,
              child: Text(title, style: _linkStyle(theme)),
            );
          },
        );
      case ContentServiceV1FieldDefinition$Type.fieldTypeComponent:
      case ContentServiceV1FieldDefinition$Type.fieldTypeDynamicZone:
        return _buildComponentValue(theme, model, row);
      default:
        return Text(row.value, style: _regularStyle(theme));
    }
  }

  /// COMPONENT/DYNAMIC_ZONE 值渲染：解析实例 JSON，按模型携带的组件 schema
  /// 渲染键值卡（动态区按 _component 分发）；解析失败回退原样字符串。
  Widget _buildComponentValue(
      ThemeData theme, ContentServiceV1ContentModel? model, _Row row) {
    if (model == null) return Text(row.value, style: _mutedStyle(theme));
    final cfg = _componentConfigOf(model, row.name);
    if (cfg == null) return Text(row.value, style: _regularStyle(theme));
    final expectArray = cfg.repeatable == true ||
        row.type == ContentServiceV1FieldDefinition$Type.fieldTypeDynamicZone;
    final raw = row.value.trim();
    if (raw.isEmpty) return const SizedBox.shrink();
    dynamic decoded;
    try {
      decoded = jsonDecode(raw);
    } catch (_) {
      return Text(row.value, style: _mutedStyle(theme));
    }
    List<dynamic> instances;
    if (expectArray) {
      if (decoded is! List) return Text(row.value, style: _mutedStyle(theme));
      instances = decoded;
    } else {
      if (decoded is! Map<String, dynamic>) {
        return Text(row.value, style: _mutedStyle(theme));
      }
      instances = [decoded];
    }
    if (instances.isEmpty) return const SizedBox.shrink();

    ContentServiceV1ContentModelComponent? schemaOf(String code) {
      for (final c in model.components ?? const []) {
        if (c.definition?.code == code) return c;
      }
      return null;
    }

    final cards = <Widget>[];
    for (var i = 0; i < instances.length; i++) {
      final inst = instances[i];
      if (inst is! Map) continue;
      final code = row.type == ContentServiceV1FieldDefinition$Type.fieldTypeDynamicZone
          ? (inst['_component'] as String? ?? '')
          : (cfg.componentCode ?? '');
      final schema = schemaOf(code);
      final children = <Widget>[
        Text(
          schema?.definition?.name ?? code,
          style: TextStyle(
            fontSize: widget.isMobile ? 12.sp : 12,
            fontWeight: FontWeight.w600,
            color: theme.colorScheme.onSurface.withAlpha(180),
          ),
        ),
      ];
      for (final f in schema?.fields ?? const []) {
        final raw2 = inst[f.name ?? ''];
        if (raw2 == null || raw2 !is String || raw2.isEmpty) continue;
        final isImg = f.type == ContentServiceV1FieldDefinition$Type.fieldTypeImage &&
            raw2.startsWith('http');
        children.add(Padding(
          padding: EdgeInsets.only(top: widget.isMobile ? 4.h : 4),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(f.label ?? f.name ?? '',
                  style: TextStyle(
                    fontSize: widget.isMobile ? 11.sp : 11,
                    color: theme.colorScheme.onSurface.withAlpha(130),
                  )),
              if (isImg)
                ClipRRect(
                  borderRadius:
                      BorderRadius.circular(widget.isMobile ? 8.r : 8),
                  child: Image.network(raw2,
                      width: double.infinity, fit: BoxFit.cover),
                )
              else
                Text(raw2,
                    style: TextStyle(
                      fontSize: widget.isMobile ? 13.sp : 13,
                      color: theme.colorScheme.onSurface,
                    )),
            ],
          ),
        ));
      }
      cards.add(Container(
        margin: EdgeInsets.only(bottom: widget.isMobile ? 6.h : 6),
        padding: EdgeInsets.all(widget.isMobile ? 8.w : 8),
        decoration: BoxDecoration(
          border: Border.all(color: theme.colorScheme.outlineVariant),
          borderRadius: BorderRadius.circular(widget.isMobile ? 8.r : 8),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: children,
        ),
      ));
    }
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: cards);
  }

  /// 从模型 fields 里找字段对应的 componentConfig
  ContentServiceV1ComponentConfig? _componentConfigOf(
      ContentServiceV1ContentModel model, String? fieldName) {
    for (final f in model.fields ?? const []) {
      if (f.name == fieldName) return f.componentConfig;
    }
    return null;
  }

  /// 解析 `media_asset:<id>` 引用；非引用形态返回 null
  int? _parseMediaRef(String value) {
    final m = RegExp(r'^media_asset:(\d+)$').firstMatch(value.trim());
    return m != null ? int.parse(m.group(1)!) : null;
  }

  Future<MediaServiceV1MediaAsset?> _fetchMediaAsset(int id) async {
    try {
      return await GetIt.instance<ApiClient>().mediaAssetService.get(
            MediaServiceV1GetMediaAssetRequest(id: id),
          );
    } catch (e) {
      _mediaLogger.w('fetch media asset $id failed: $e');
      return null;
    }
  }

  /// 解析 RELATION 引用（"post:12"/"page:7"）为目标内容摘要；
  /// 不可见/不存在的引用返回 null（调用方回退原样字符串）。
  Future<ContentServiceV1ContentReferenceSummary?> _resolveReference(
      String ref) async {
    try {
      final resp = await GetIt.instance<ApiClient>().postService
          .resolveContentReferences(
        ContentServiceV1ResolveContentReferencesRequest(refs: [ref]),
      );
      final items = resp.items;
      return (items == null || items.isEmpty) ? null : items.first;
    } catch (e) {
      _mediaLogger.w('resolve reference $ref failed: $e');
      return null;
    }
  }

  TextStyle _mutedStyle(ThemeData theme) => TextStyle(
      fontSize: widget.isMobile ? 12.sp : 12,
      color: theme.colorScheme.onSurface.withAlpha(120));
  TextStyle _linkStyle(ThemeData theme) => TextStyle(
      fontSize: widget.isMobile ? 12.sp : 12, color: theme.colorScheme.primary);
  TextStyle _regularStyle(ThemeData theme) => TextStyle(
      fontSize: widget.isMobile ? 13.sp : 13,
      color: theme.colorScheme.onSurface);
}

class _Row {
  final String name;
  final String label;
  final ContentServiceV1FieldDefinition$Type type;
  final String value;

  const _Row({
    required this.name,
    required this.label,
    required this.type,
    required this.value,
  });
}
