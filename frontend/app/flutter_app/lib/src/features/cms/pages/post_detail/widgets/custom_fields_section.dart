import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_widget_from_html/flutter_widget_from_html.dart';
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show
        ApiClient,
        ContentServiceV1ContentModel,
        ContentServiceV1FieldDefinition,
        ContentServiceV1FieldDefinition$Type,
        MediaServiceV1GetMediaAssetRequest,
        MediaServiceV1MediaAsset;
import 'package:flutter_app/src/core/utilities/logger.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';

/// 自定义字段渲染（内容建模前台读侧，flutter 端）：
/// 按文章分类链绑定的内容模型字段定义，渲染 post.customFields。
/// 值约定为字符串（admin 编辑器 `Record<string,string>` 口径）：
///   - RICHTEXT → HtmlWidget（flutter_widget_from_html）；
///   - IMAGE/FILE → `media_asset:<id>` 引用经公开媒体端点换 URL（直链 URL 直用）；
///   - 其余（TEXT/NUMBER/RELATION 等）原样展示。
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
                      Expanded(child: _buildValue(context, row)),
                    ],
                  ),
                ),
            ],
          ),
        );
      },
    );
  }

  Widget _buildValue(BuildContext context, _Row row) {
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
      default:
        return Text(row.value, style: _regularStyle(theme));
    }
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
