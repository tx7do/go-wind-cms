import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';

typedef Tag = ContentServiceV1Tag;

/// 热门标签云组件
class TagCloud extends StatelessWidget {
  final List<Tag> tags;

  const TagCloud({super.key, required this.tags});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isMobile = ResponsiveUtils.isMobile(context);

    return Padding(
      padding: isMobile
          ? EdgeInsets.fromLTRB(16.w, 8.h, 16.w, 4.h)
          : const EdgeInsets.fromLTRB(16, 8, 16, 4),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: EdgeInsets.only(bottom: isMobile ? 10.h : 10),
            child: Row(
              children: [
                Icon(Icons.local_fire_department,
                    size: isMobile ? 18.sp : 18, color: Colors.orange),
                SizedBox(width: isMobile ? 6.w : 6),
                Text(
                  '热门标签',
                  style: TextStyle(
                    fontSize: isMobile ? 14.sp : 14,
                    fontWeight: FontWeight.w600,
                    color: theme.colorScheme.onSurface,
                  ),
                ),
              ],
            ),
          ),
          Wrap(
            spacing: isMobile ? 8.w : 8,
            runSpacing: isMobile ? 8.h : 8,
            children: tags.map((tag) => _TagChip(tag: tag, isMobile: isMobile)).toList(),
          ),
        ],
      ),
    );
  }
}

class _TagChip extends StatelessWidget {
  final Tag tag;
  final bool isMobile;

  const _TagChip({required this.tag, required this.isMobile});

  Color _parseColor(String? hexColor) {
    if (hexColor == null) return Colors.grey;
    try {
      final hex = hexColor.replaceAll('#', '');
      return Color(int.parse('FF$hex', radix: 16));
    } catch (_) {
      return Colors.grey;
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final tagColor = _parseColor(tag.color);

    return ActionChip(
      onPressed: () {
        // TODO: Navigate to tag feed page
      },
      label: Text(
        '# ${(tag.translations ?? []).isNotEmpty ? tag.translations!.first.name ?? '' : ''}',
        style: TextStyle(
          fontSize: isMobile ? 13.sp : 13,
          color: tagColor,
          fontWeight: FontWeight.w500,
        ),
      ),
      avatar: Icon(Icons.tag, size: isMobile ? 14.sp : 14, color: tagColor),
      side: BorderSide(color: tagColor.withAlpha((0.3 * 255).round())),
      backgroundColor: tagColor.withAlpha((0.06 * 255).round()),
      padding: EdgeInsets.symmetric(horizontal: isMobile ? 4.w : 4, vertical: 0),
    );
  }
}
