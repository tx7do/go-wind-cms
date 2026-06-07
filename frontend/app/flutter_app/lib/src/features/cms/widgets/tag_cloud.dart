import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

typedef Tag = ContentServiceV1Tag;

/// 热门标签云组件
///
/// 右侧侧边栏极简版：统一柔和色调，限制数量，Wrap 错落排列
class TagCloud extends StatelessWidget {
  final List<Tag> tags;

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
                Icon(Icons.local_fire_department,
                    size: isMobile ? 18.sp : 18, color: Colors.orange),
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
            children: displayTags.map((tag) => _TagChip(tag: tag, isMobile: isMobile)).toList(),
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

  String get _name => getTagName(tag);

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: () {
        // TODO: Navigate to tag feed page
      },
      borderRadius: BorderRadius.circular(isMobile ? 16.r : 16),
      child: Container(
        padding: EdgeInsets.symmetric(
          horizontal: isMobile ? 12.w : 12,
          vertical: isMobile ? 6.h : 6,
        ),
        decoration: BoxDecoration(
          // 统一柔和淡灰色底，不使用五颜六色
          color: Colors.grey.shade50,
          borderRadius: BorderRadius.circular(isMobile ? 16.r : 16),
          border: Border.all(color: Colors.grey.shade200),
        ),
        child: Text(
          '# $_name',
          style: TextStyle(
            fontSize: isMobile ? 13.sp : 13,
            color: Colors.grey.shade700,
            fontWeight: FontWeight.w500,
          ),
        ),
      ),
    );
  }
}
