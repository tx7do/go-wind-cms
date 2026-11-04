import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Tag;
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

/// 统一标签芯片
///
/// 圆角小方块而非 MD3 药丸——标签是「元信息链接」，选中/悬停才沾主色
/// （doc §1：主色只给行动与选中）。
///
/// 适用于所有 tag 展示场景：
/// - 点击导航到 `/tag/:id`（可通过 [onTap] 自定义）
/// - 可选显示文章数量（[showPostCount]）
/// - 紧凑模式（[compact]，减小 padding），适用于文章详情页内嵌标签
///
/// [isMobile] 仅为兼容既有调用方保留，观感不再随断点缩放。
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
    this.isMobile = false,
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
    final p = theme.palette;
    final radius = theme.radius.smAll;
    final hovered = _isHovered;

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: Material(
        color: hovered ? p.primarySoft : p.surfaceLow,
        shape: RoundedRectangleBorder(
          borderRadius: radius,
          side: BorderSide(
            color: hovered ? p.primaryBorder : p.hairline,
            width: theme.borderWidth.thin,
          ),
        ),
        child: InkWell(
          borderRadius: radius,
          onTap: widget.onTap ?? _defaultNavigate,
          child: Padding(
            padding: EdgeInsets.symmetric(
              horizontal: widget.compact ? 8 : 10,
              vertical: widget.compact ? 3 : 5,
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  '# ${getTagName(widget.tag)}',
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: (widget.compact
                          ? theme.textTheme.labelSmall
                          : theme.textTheme.labelMedium)
                      ?.copyWith(
                        color: hovered ? p.primaryOnSoft : p.inkSecondary,
                        fontWeight: FontWeight.w500,
                      ),
                ),
                if (widget.showPostCount &&
                    widget.tag.postCount != null &&
                    widget.tag.postCount! > 0) ...[
                  const SizedBox(width: 4),
                  Text(
                    '(${widget.tag.postCount})',
                    style: theme.textTheme.labelSmall,
                  ),
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }

  void _defaultNavigate() {
    if (widget.tag.id != null) {
      context.go('/tag/${widget.tag.id}');
    }
  }
}
