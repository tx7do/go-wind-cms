import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:flutter_widget_from_html/flutter_widget_from_html.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_post_editor_type.dart';

/// 文章内容查看器
///
/// 根据 [editorType] 自适应渲染不同编辑器类型的内容：
/// - [ContentServiceV1PostEditorType.editorTypeMarkdown]: Markdown 内容（基础文本渲染，后续可接入 flutter_markdown）
/// - [ContentServiceV1PostEditorType.editorTypeRichText]: 富文本 HTML（基础文本渲染，后续可接入 flutter_html）
/// - [ContentServiceV1PostEditorType.editorTypePlainText]: 纯文本
/// - [ContentServiceV1PostEditorType.editorTypeCode]: 代码块（等宽字体 + 代码样式）
/// - [ContentServiceV1PostEditorType.editorTypeJsonBlock]: JSON 结构化内容
/// - [ContentServiceV1PostEditorType.editorTypeVisualBuilder]: 可视化构建器（暂按 JSON 展示）
/// - 未指定 / unknown: 按 Markdown 渲染
class ContentViewer extends StatelessWidget {
  /// 文章原始内容
  final String content;

  /// 编辑器类型
  final ContentServiceV1PostEditorType? editorType;

  /// 是否手机端布局
  final bool isMobile;

  const ContentViewer({
    super.key,
    required this.content,
    this.editorType,
    this.isMobile = false,
  });

  @override
  Widget build(BuildContext context) {
    if (content.isEmpty) return const SizedBox.shrink();

    return switch (editorType) {
      ContentServiceV1PostEditorType.editorTypeMarkdown => _buildMarkdown(context),
      ContentServiceV1PostEditorType.editorTypeRichText => _buildRichText(context),
      ContentServiceV1PostEditorType.editorTypePlainText => _buildPlainText(context),
      ContentServiceV1PostEditorType.editorTypeCode => _buildCodeBlock(context),
      ContentServiceV1PostEditorType.editorTypeJsonBlock => _buildJsonBlock(context),
      ContentServiceV1PostEditorType.editorTypeVisualBuilder => _buildJsonBlock(context),
      // unspecified / null / unknown → 默认按 Markdown 渲染
      _ => _buildMarkdown(context),
    };
  }

  // ─── Markdown 渲染 ─────────────────────────────────────

  Widget _buildMarkdown(BuildContext context) {
    final theme = Theme.of(context);
    return MarkdownBody(
      data: content,
      selectable: true,
      styleSheet: MarkdownStyleSheet(
        p: TextStyle(
          fontSize: isMobile ? 15.sp : 15,
          color: theme.colorScheme.onSurface,
          height: 1.8,
        ),
        h1: TextStyle(
          fontSize: isMobile ? 24.sp : 24,
          fontWeight: FontWeight.bold,
          color: theme.colorScheme.onSurface,
          height: 1.4,
        ),
        h2: TextStyle(
          fontSize: isMobile ? 22.sp : 22,
          fontWeight: FontWeight.bold,
          color: theme.colorScheme.onSurface,
          height: 1.4,
        ),
        h3: TextStyle(
          fontSize: isMobile ? 20.sp : 20,
          fontWeight: FontWeight.bold,
          color: theme.colorScheme.onSurface,
          height: 1.4,
        ),
        h4: TextStyle(
          fontSize: isMobile ? 18.sp : 18,
          fontWeight: FontWeight.w600,
          color: theme.colorScheme.onSurface,
          height: 1.4,
        ),
        h5: TextStyle(
          fontSize: isMobile ? 16.sp : 16,
          fontWeight: FontWeight.w600,
          color: theme.colorScheme.onSurface,
          height: 1.4,
        ),
        h6: TextStyle(
          fontSize: isMobile ? 15.sp : 15,
          fontWeight: FontWeight.w600,
          color: theme.colorScheme.onSurface,
          height: 1.4,
        ),
        code: TextStyle(
          fontSize: isMobile ? 13.sp : 13,
          fontFamily: 'RobotoMono',
          fontFamilyFallback: const ['Courier', 'monospace'],
          color: theme.colorScheme.primary,
          backgroundColor: theme.colorScheme.surfaceContainerHighest,
        ),
        codeblockDecoration: BoxDecoration(
          color: theme.brightness == Brightness.dark
              ? const Color(0xFF1E1E1E)
              : const Color(0xFFF5F5F5),
          borderRadius: BorderRadius.circular(isMobile ? 10.r : 10),
          border: Border.all(
            color: theme.colorScheme.onSurface.withAlpha((0.08 * 255).round()),
          ),
        ),
        codeblockAlign: WrapAlignment.start,
        blockquote: TextStyle(
          fontSize: isMobile ? 15.sp : 15,
          color: theme.colorScheme.onSurface.withAlpha(180),
          height: 1.8,
          fontStyle: FontStyle.italic,
        ),
        blockquoteDecoration: BoxDecoration(
          color: theme.colorScheme.primaryContainer.withAlpha((0.15 * 255).round()),
          borderRadius: BorderRadius.circular(isMobile ? 8.r : 8),
        ),
        listBullet: TextStyle(
          fontSize: isMobile ? 15.sp : 15,
          color: theme.colorScheme.primary,
        ),
        listIndent: isMobile ? 20.w : 20,
        tableBody: TextStyle(
          fontSize: isMobile ? 14.sp : 14,
          color: theme.colorScheme.onSurface,
        ),
        tableHead: TextStyle(
          fontSize: isMobile ? 14.sp : 14,
          fontWeight: FontWeight.w600,
          color: theme.colorScheme.onSurface,
        ),
        a: TextStyle(
          fontSize: isMobile ? 15.sp : 15,
          color: theme.colorScheme.primary,
          decoration: TextDecoration.underline,
        ),
        img: TextStyle(fontSize: 0), // 图片通过 MarkdownBody 内置渲染
      ),
    );
  }

  // ─── 富文本 HTML 渲染 ──────────────────────────────────

  Widget _buildRichText(BuildContext context) {
    final theme = Theme.of(context);
    return HtmlWidget(
      content,
      textStyle: TextStyle(
        fontSize: isMobile ? 15.sp : 15,
        color: theme.colorScheme.onSurface,
        height: 1.8,
      ),
      customWidgetBuilder: (element) => null,
      onTapUrl: (url) {
        // TODO: 处理链接点击，可跳转内部路由或外部浏览器
        return false;
      },
    );
  }

  // ─── 纯文本渲染 ────────────────────────────────────────

  Widget _buildPlainText(BuildContext context) {
    final theme = Theme.of(context);
    return SelectableText(
      content,
      style: TextStyle(
        fontSize: isMobile ? 15.sp : 15,
        color: theme.colorScheme.onSurface,
        height: 1.8,
      ),
    );
  }

  // ─── 代码块渲染 ────────────────────────────────────────

  Widget _buildCodeBlock(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Container(
      width: double.infinity,
      padding: EdgeInsets.all(isMobile ? 14.w : 16),
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF1E1E1E) : const Color(0xFFF5F5F5),
        borderRadius: BorderRadius.circular(isMobile ? 10.r : 10),
        border: Border.all(
          color: theme.colorScheme.onSurface.withAlpha((0.08 * 255).round()),
        ),
      ),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Text(
          content,
          style: TextStyle(
            fontSize: isMobile ? 13.sp : 13,
            fontFamily: 'RobotoMono',
            // Roboto Mono 不支持中文，fallback 到系统等宽字体
            fontFamilyFallback: const ['Courier', 'monospace'],
            color: isDark ? const Color(0xFFD4D4D4) : const Color(0xFF333333),
            height: 1.6,
          ),
        ),
      ),
    );
  }

  // ─── JSON 结构化渲染 ──────────────────────────────────

  Widget _buildJsonBlock(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    // 尝试格式化 JSON
    String displayContent;
    try {
      final decoded = jsonDecode(content);
      displayContent =
          const JsonEncoder.withIndent('  ').convert(decoded);
    } catch (_) {
      displayContent = content;
    }

    return Container(
      width: double.infinity,
      padding: EdgeInsets.all(isMobile ? 14.w : 16),
      decoration: BoxDecoration(
        color: isDark ? const Color(0xFF1A1A2E) : const Color(0xFFF8F9FA),
        borderRadius: BorderRadius.circular(isMobile ? 10.r : 10),
        border: Border.all(
          color: theme.colorScheme.primary.withAlpha((0.2 * 255).round()),
        ),
      ),
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Text(
          displayContent,
          style: TextStyle(
            fontSize: isMobile ? 12.sp : 12,
            fontFamily: 'RobotoMono',
            fontFamilyFallback: const ['Courier', 'monospace'],
            color: isDark ? const Color(0xFFCCCCCC) : const Color(0xFF444444),
            height: 1.5,
          ),
        ),
      ),
    );
  }
}
