import 'package:flutter/material.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Tag;
import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/layout/width_class.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/features/cms/pages/home/widgets/section_header.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_chip.dart';

/// 热门标签云
///
/// 区块语言（SectionHeader）+ Wrap 排列，限制数量防侧栏失控。
class TagCloud extends StatelessWidget {
  final List<ContentServiceV1Tag> tags;

  /// 最大显示标签数量，超出部分隐藏
  final int maxTags;

  const TagCloud({super.key, required this.tags, this.maxTags = 10});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final compact = context.widthClass.isCompact;
    final displayTags = tags.take(maxTags).toList();

    return Padding(
      padding: compact
          ? EdgeInsets.fromLTRB(16, theme.space.sm, 16, theme.space.xs)
          : EdgeInsets.zero,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: EdgeInsets.only(bottom: theme.space.md + 2),
            child: SectionHeader(title: S.of(context).hotTags),
          ),
          Wrap(
            spacing: theme.space.sm,
            runSpacing: theme.space.sm + 2,
            children: displayTags
                .map((tag) => TagChip(tag: tag))
                .toList(),
          ),
        ],
      ),
    );
  }
}
