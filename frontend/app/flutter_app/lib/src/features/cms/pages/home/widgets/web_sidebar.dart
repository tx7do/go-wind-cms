import 'package:flutter/material.dart';

import 'package:flutter_app/src/features/cms/widgets/tag_cloud.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';

typedef Tag = ContentServiceV1Tag;

/// 右侧边栏（极简版：只保留热门标签云）
class WebSidebar extends StatelessWidget {
  final List<Tag> tags;

  const WebSidebar({super.key, required this.tags});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Card(
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
            side: BorderSide(
              color: theme.colorScheme.onSurface.withAlpha(
                (0.06 * 255).round(),
              ),
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: TagCloud(tags: tags),
          ),
        ),
      ],
    );
  }
}
