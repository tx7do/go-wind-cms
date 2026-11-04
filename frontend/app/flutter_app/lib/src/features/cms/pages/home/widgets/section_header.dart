import 'package:flutter/material.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';

/// 分区标题
///
/// 全站统一的区块语言：3px 主色左条 + titleLarge + 可选计数（等宽数字）。
class SectionHeader extends StatelessWidget {
  final String title;

  /// 右侧计数；null 则不显示
  final int? count;

  const SectionHeader({super.key, required this.title, this.count});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Row(
      children: [
        Container(
          width: 3,
          height: 16,
          decoration: BoxDecoration(
            color: theme.palette.primary,
            borderRadius: theme.radius.fullAll,
          ),
        ),
        SizedBox(width: theme.space.sm),
        Expanded(
          child: Text(title, style: theme.textTheme.titleLarge, maxLines: 1, overflow: TextOverflow.ellipsis),
        ),
        if (count != null) Text(S.of(context).postsCount(count!), style: theme.textTheme.labelSmall),
      ],
    );
  }
}
