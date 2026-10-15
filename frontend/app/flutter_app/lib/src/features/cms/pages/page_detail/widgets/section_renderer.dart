import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:url_launcher/url_launcher.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show
        ContentServiceV1EditorType,
        ContentServiceV1Section,
        ContentServiceV1SectionTranslation,
        ContentServiceV1SectionType;
import 'package:flutter_app/src/core/utils/translation_helpers.dart';
import 'package:flutter_app/src/features/cms/widgets/content_viewer.dart';

/// Section 内容键契约（三端渲染器共用，与 react SectionRenderer.tsx / vue
/// SectionRenderer.vue 保持同步）：语言相关内容存 translations[].content，
/// 主内容键一律为 `value`；语言无关样式存 config。缺译回退：匹配当前语言 →
/// 第一条翻译。各类型键明细见 react 端文件头表格。

/// 取区块当前语言的翻译：匹配当前语言 → 回退第一条
ContentServiceV1SectionTranslation? pickSectionTranslation(
  ContentServiceV1Section section,
) {
  final translations = section.translations;
  if (translations == null || translations.isEmpty) return null;
  final locale = getCurrentLocale();
  if (locale != null) {
    for (final tr in translations) {
      if (tr.languageCode == locale) return tr;
    }
  }
  return translations.first;
}

String sectionValue(ContentServiceV1Section section) =>
    pickSectionTranslation(section)?.content?['value'] ?? '';

/// GALLERY/CAROUSEL 图片列表：value 为 JSON 数组（URL 或 {src,caption}）或换行分隔 URL
List<({String src, String? caption})> sectionImageList(
  ContentServiceV1Section section,
) {
  final raw = sectionValue(section).trim();
  if (raw.isEmpty) return const [];
  if (raw.startsWith('[')) {
    try {
      final decoded = jsonDecode(raw);
      if (decoded is List) {
        final result = <({String src, String? caption})>[];
        for (final item in decoded) {
          if (item is String) {
            result.add((src: item, caption: null));
          } else if (item is Map) {
            final src = (item['src'] ?? item['url']) as String? ?? '';
            if (src.isNotEmpty) {
              result.add((src: src, caption: item['caption'] as String?));
            }
          }
        }
        return result;
      }
    } catch (_) {
      // JSON 解析失败时按纯 URL 文本处理，不中断整页渲染
    }
  }
  return raw
      .split(RegExp(r'\r?\n'))
      .map((line) => line.trim())
      .where((line) => line.isNotEmpty)
      .map((src) => (src: src, caption: null))
      .toList();
}

double _configDouble(
  ContentServiceV1Section section,
  String key, {
  required double fallback,
  required double min,
  required double max,
}) {
  final raw = double.tryParse(section.config?[key] ?? '');
  if (raw == null) return fallback;
  return raw.clamp(min, max).toDouble();
}

/// CMS Page 区块渲染器：按 sortOrder 排序渲染全部区块，
/// 未知类型按 Markdown 兜底渲染 value。
class SectionRenderer extends StatelessWidget {
  final List<ContentServiceV1Section>? sections;
  final bool isMobile;

  const SectionRenderer({
    super.key,
    this.sections,
    this.isMobile = false,
  });

  List<ContentServiceV1Section> get _sorted {
    final list = [...(sections ?? const <ContentServiceV1Section>[])];
    list.sort((a, b) => (a.sortOrder ?? 0).compareTo(b.sortOrder ?? 0));
    return list;
  }

  int _titleLevel(ContentServiceV1Section section) => _configDouble(
        section,
        'level',
        fallback: 2,
        min: 1,
        max: 6,
      ).round();

  @override
  Widget build(BuildContext context) {
    final sorted = _sorted;
    if (sorted.isEmpty) return const SizedBox.shrink();

    final gap = isMobile ? 16.h : 24.0;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        for (var i = 0; i < sorted.length; i++) ...[
          if (i > 0) SizedBox(height: gap),
          _buildSection(context, sorted[i]),
        ],
      ],
    );
  }

  Widget _buildSection(BuildContext context, ContentServiceV1Section section) {
    final theme = Theme.of(context);
    final type = section.type;

    // 富文本 / HTML
    if (type == ContentServiceV1SectionType.sectionTypeRichText ||
        type == ContentServiceV1SectionType.sectionTypeHtml) {
      return ContentViewer(
        content: sectionValue(section),
        editorType: ContentServiceV1EditorType.editorTypeRichText,
        isMobile: isMobile,
      );
    }

    // Markdown / 自定义/未知类型兜底
    if (type == ContentServiceV1SectionType.sectionTypeMarkdown ||
        type == ContentServiceV1SectionType.sectionTypeCustom) {
      return ContentViewer(
        content: sectionValue(section),
        editorType: ContentServiceV1EditorType.editorTypeMarkdown,
        isMobile: isMobile,
      );
    }

    if (type == ContentServiceV1SectionType.sectionTypeTitle) {
      final level = _titleLevel(section);
      final sizes = isMobile
          ? [22.0, 20.0, 18.0, 16.0, 15.0, 14.0]
          : [32.0, 28.0, 24.0, 20.0, 18.0, 16.0];
      return Text(
        sectionValue(section),
        style: theme.textTheme.titleLarge?.copyWith(
          fontSize: sizes[level - 1],
          fontWeight: FontWeight.bold,
        ),
      );
    }

    if (type == ContentServiceV1SectionType.sectionTypeImage) {
      final src = sectionValue(section);
      if (src.isEmpty) return const SizedBox.shrink();
      return ClipRRect(
        borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
        child: Image.network(
          src,
          fit: BoxFit.cover,
          width: double.infinity,
          errorBuilder: (_, _, _) => Container(
            height: 180,
            color: theme.colorScheme.surfaceContainerHighest,
            alignment: Alignment.center,
            child: const Icon(Icons.broken_image_outlined),
          ),
        ),
      );
    }

    if (type == ContentServiceV1SectionType.sectionTypeGallery ||
        type == ContentServiceV1SectionType.sectionTypeCarousel) {
      final images = sectionImageList(section);
      if (images.isEmpty) return const SizedBox.shrink();
      final height = isMobile ? 180.0 : 240.0;
      return SizedBox(
        height: height,
        child: ListView.separated(
          scrollDirection: Axis.horizontal,
          itemCount: images.length,
          separatorBuilder: (_, _) => SizedBox(width: isMobile ? 12.w : 12),
          itemBuilder: (context, index) {
            return ClipRRect(
              borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
              child: Image.network(
                images[index].src,
                fit: BoxFit.cover,
                width: height * 4 / 3,
                height: height,
                errorBuilder: (_, _, _) => Container(
                  color: theme.colorScheme.surfaceContainerHighest,
                  child: const Icon(Icons.broken_image_outlined),
                ),
              ),
            );
          },
        ),
      );
    }

    if (type == ContentServiceV1SectionType.sectionTypeVideo) {
      final src = sectionValue(section);
      if (src.isEmpty) return const SizedBox.shrink();
      // 无内嵌播放器依赖：渲染预览卡，点击经 url_launcher 打开
      return _VideoLaunchTile(url: src, isMobile: isMobile);
    }

    if (type == ContentServiceV1SectionType.sectionTypeButton) {
      final label = sectionValue(section);
      final href = section.config?['href'] ??
          section.config?['url'] ??
          section.config?['link'];
      if (label.isEmpty || href == null || href.isEmpty) {
        return const SizedBox.shrink();
      }
      final outline = section.config?['variant'] == 'outline';
      return FilledButton.tonal(
        onPressed: () async {
          final uri = Uri.tryParse(href);
          if (uri != null) {
            try {
              await launchUrl(uri, mode: LaunchMode.externalApplication);
            } catch (e) {
              // 打开失败仅提示，不中断页面
              debugPrint('launch button url failed: $e');
            }
          }
        },
        style: outline
            ? FilledButton.styleFrom(
                backgroundColor: theme.colorScheme.surface,
                foregroundColor: theme.colorScheme.primary,
              )
            : null,
        child: Text(label),
      );
    }

    if (type == ContentServiceV1SectionType.sectionTypeDivider) {
      return Divider(
        height: 1,
        color: theme.colorScheme.onSurface.withAlpha(30),
      );
    }

    if (type == ContentServiceV1SectionType.sectionTypeSpacer) {
      return SizedBox(
        height: _configDouble(
          section,
          'height',
          fallback: 48,
          min: 8,
          max: 400,
        ),
      );
    }

    if (type == ContentServiceV1SectionType.sectionTypeCode) {
      return Container(
        width: double.infinity,
        padding: EdgeInsets.all(isMobile ? 14.w : 14),
        decoration: BoxDecoration(
          color: theme.colorScheme.surfaceContainerHighest,
          borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
        ),
        child: SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: Text(
            sectionValue(section),
            style: TextStyle(
              fontFamily: 'monospace',
              fontSize: isMobile ? 12.sp : 13,
              height: 1.5,
            ),
          ),
        ),
      );
    }

    if (type == ContentServiceV1SectionType.sectionTypeForm) {
      // 表单提交后端（form_submission）落地前的占位：可见但不误导访客输入
      final label = pickSectionTranslation(section)?.content?['title'] ??
          section.name ??
          '';
      if (label.isEmpty) return const SizedBox.shrink();
      return Container(
        width: double.infinity,
        padding: EdgeInsets.all(isMobile ? 14.w : 16),
        decoration: BoxDecoration(
          border: Border.all(color: theme.colorScheme.onSurface.withAlpha(40)),
          borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
        ),
        child: Text(
          label,
          textAlign: TextAlign.center,
          style: TextStyle(
            fontSize: isMobile ? 13.sp : 13,
            color: theme.colorScheme.onSurface.withAlpha(150),
          ),
        ),
      );
    }

    // 未知类型：有内容则按 Markdown 兜底
    return ContentViewer(
      content: sectionValue(section),
      editorType: ContentServiceV1EditorType.editorTypeMarkdown,
      isMobile: isMobile,
    );
  }
}

/// 视频区块占位卡：点击经 url_launcher 打开（无内嵌播放依赖）
class _VideoLaunchTile extends StatelessWidget {
  final String url;
  final bool isMobile;

  const _VideoLaunchTile({required this.url, required this.isMobile});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ClipRRect(
      borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
      child: AspectRatio(
        aspectRatio: 16 / 9,
        child: Material(
          color: theme.colorScheme.surfaceContainerHighest,
          child: InkWell(
            onTap: () async {
              final uri = Uri.tryParse(url);
              if (uri != null) {
                try {
                  await launchUrl(uri, mode: LaunchMode.externalApplication);
                } catch (e) {
                  debugPrint('launch video url failed: $e');
                }
              }
            },
            child: Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.play_circle_outline,
                      size: isMobile ? 48.r : 48, color: theme.colorScheme.primary),
                  SizedBox(height: isMobile ? 8.h : 8),
                  Text(
                    url,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: TextStyle(
                      fontSize: isMobile ? 12.sp : 12,
                      color: theme.colorScheme.onSurface.withAlpha(150),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
