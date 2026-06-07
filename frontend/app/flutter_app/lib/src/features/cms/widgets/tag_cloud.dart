import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_chip.dart';

/// 热门标签云组件
///
/// 右侧侧边栏极简版：统一柔和色调，限制数量，Wrap 错落排列
class TagCloud extends StatelessWidget {
  final List<ContentServiceV1Tag> tags;

  /// 最大显示标签数量，超出部分隐藏
  final int maxTags;

  const TagCloud({super.key, required this.tags, this.maxTags = 10});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isMobile = ResponsiveUtils.isMobile(context);
    final displayTags = tags.take(maxTags).toList();

    return Padding(
      padding: isMobile
          ? EdgeInsets.fromLTRB(16.w, 8.h, 16.w, 4.h)
          : const EdgeInsets.fromLTRB(0, 0, 0, 0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: EdgeInsets.only(bottom: isMobile ? 10.h : 14),
            child: Row(
              children: [
                Icon(
                  Icons.local_fire_department,
                  size: isMobile ? 18.sp : 18,
                  color: Colors.orange,
                ),
                SizedBox(width: isMobile ? 6.w : 6),
                Text(
                  S.of(context).hotTags,
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
            runSpacing: isMobile ? 8.h : 10,
            children: displayTags
                .map((tag) => TagChip(tag: tag, isMobile: isMobile))
                .toList(),
          ),
        ],
      ),
    );
  }
}
