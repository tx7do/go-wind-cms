import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post;

typedef Post = ContentServiceV1Post;

class InteractionBar extends StatelessWidget {
  final Post post;
  final bool isMobile;

  const InteractionBar({super.key, required this.post, required this.isMobile});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: EdgeInsets.symmetric(
        horizontal: isMobile ? 20.w : 0,
        vertical: 12,
      ),
      child: Container(
        padding: EdgeInsets.symmetric(
          vertical: isMobile ? 12.h : 12,
          horizontal: isMobile ? 16.w : 16,
        ),
        decoration: BoxDecoration(
          color: theme.colorScheme.surface,
          borderRadius: BorderRadius.circular(isMobile ? 14.r : 14),
          border: Border.all(
            color: theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
          ),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            InteractionItem(
              icon: Icons.remove_red_eye_outlined,
              value: '${post.visits ?? 0}',
              label: S.of(context).views,
            ),
            InteractionItem(
              icon: Icons.favorite_outline,
              value: '${post.likes ?? 0}',
              label: S.of(context).likes,
            ),
            InteractionItem(
              icon: Icons.comment_outlined,
              value: '${post.commentCount ?? 0}',
              label: S.of(context).comments,
            ),
            InteractionItem(
              icon: Icons.share_outlined,
              value: S.of(context).share,
              label: '',
            ),
          ],
        ),
      ),
    );
  }
}

class InteractionItem extends StatefulWidget {
  final IconData icon;
  final String value;
  final String label;

  const InteractionItem({
    super.key,
    required this.icon,
    required this.value,
    required this.label,
  });

  @override
  State<InteractionItem> createState() => _InteractionItemState();
}

class _InteractionItemState extends State<InteractionItem> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: InkWell(
        onTap: () {},
        borderRadius: BorderRadius.circular(8),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          child: Column(
            children: [
              Icon(
                widget.icon,
                size: 20,
                color: _isHovered
                    ? theme.colorScheme.primary
                    : theme.colorScheme.onSurface.withAlpha(160),
              ),
              const SizedBox(height: 4),
              Text(
                widget.value,
                style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w500,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              if (widget.label.isNotEmpty) ...[
                const SizedBox(height: 2),
                Text(
                  widget.label,
                  style: TextStyle(
                    fontSize: 10,
                    color: theme.colorScheme.onSurface.withAlpha(100),
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
