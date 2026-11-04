import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_cloud.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Tag;

typedef Tag = ContentServiceV1Tag;

/// 右侧边栏（极简版：热门标签云）
class WebSidebar extends StatelessWidget {
  final List<Tag> tags;

  const WebSidebar({super.key, required this.tags});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final p = theme.palette;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Card(
          elevation: 0,
          margin: EdgeInsets.zero,
          color: p.surface,
          shape: RoundedRectangleBorder(
            borderRadius: theme.radius.mdAll,
            side: BorderSide(
              color: p.border,
              width: theme.borderWidth.thin,
            ),
          ),
          child: Padding(
            padding: EdgeInsets.all(theme.space.cardPadding),
            child: TagCloud(tags: tags),
          ),
        ),
      ],
    );
  }
}
