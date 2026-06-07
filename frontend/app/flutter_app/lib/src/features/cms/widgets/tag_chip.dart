import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

/// 统一标签芯片组件
///
/// 适用于所有 tag 展示场景，自动处理：
/// - 响应式尺寸（[isMobile] 控制字号、间距）
/// - Web 端 hover 效果
/// - 点击导航到 `/tag/:id`（可通过 [onTap] 自定义）
/// - 可选显示文章数量（[showPostCount]）
class TagChip extends StatefulWidget {
  final ContentServiceV1Tag tag;
  final bool isMobile;

  /// 是否显示文章数量后缀，如 `# Flutter (12)`
  final bool showPostCount;

  /// 自定义点击回调；为 null 时默认导航到 `/tag/:id`
  final VoidCallback? onTap;

  /// 紧凑模式（减小 padding），适用于文章详情页内嵌标签
  final bool compact;

  const TagChip({
    super.key,
    required this.tag,
    required this.isMobile,
    this.showPostCount = false,
    this.onTap,
    this.compact = false,
  });

  @override
  State<TagChip> createState() => _TagChipState();
}

class _TagChipState extends State<TagChip> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = getTagName(widget.tag);
    final fontSize = widget.compact
        ? (widget.isMobile ? 12.sp : 12.0)
        : (widget.isMobile ? 14.sp : 13.0);

    final label = Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Text(
          '# $name',
          style: TextStyle(fontSize: fontSize, fontWeight: FontWeight.w500),
        ),
        if (widget.showPostCount &&
            widget.tag.postCount != null &&
            widget.tag.postCount! > 0) ...[
          SizedBox(width: widget.isMobile ? 4.w.toDouble() : 4.0),
          Text(
            '(${widget.tag.postCount})',
            style: TextStyle(
              fontSize: widget.compact
                  ? fontSize
                  : (widget.isMobile ? 12.sp.toDouble() : 12.0),
              color: theme.colorScheme.onSurface.withAlpha(140),
            ),
          ),
        ],
      ],
    );

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: ActionChip(
        onPressed: widget.onTap ?? _defaultNavigate,
        backgroundColor: _isHovered
            ? theme.colorScheme.primaryContainer.withAlpha(120)
            : theme.colorScheme.surfaceContainerLow,
        side: BorderSide(
          color: theme.colorScheme.onSurface.withAlpha((0.08 * 255).round()),
        ),
        labelPadding: widget.compact
            ? EdgeInsets.symmetric(
                horizontal: widget.isMobile ? 4.w : 4,
                vertical: widget.isMobile ? 1.h : 1,
              )
            : EdgeInsets.symmetric(
                horizontal: widget.isMobile ? 4.w : 6,
                vertical: widget.isMobile ? 2.h : 2,
              ),
        visualDensity: widget.compact ? VisualDensity.compact : null,
        label: label,
      ),
    );
  }

  void _defaultNavigate() {
    if (widget.tag.id != null) {
      context.go('/tag/${widget.tag.id}');
    }
  }
}
