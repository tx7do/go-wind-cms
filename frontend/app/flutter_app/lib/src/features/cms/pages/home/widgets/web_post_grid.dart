import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/layout/width_class.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Category,
        ContentServiceV1Tag;

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Tag = ContentServiceV1Tag;

/// 文章网格
///
/// Row/Column 分行布局（规避 Web 端 Sliver 里嵌 GridView 的 hitTest 问题），
/// 列数由 [AdaptiveGrid.columns] 按实测宽度反推（300dp 起排，上限随档位）。
class WebPostGrid extends StatelessWidget {
  final List<Post> posts;
  final List<Category> categories;
  final List<Tag> tags;
  final Map<int, int> likeCounts;

  const WebPostGrid({
    super.key,
    required this.posts,
    required this.categories,
    required this.tags,
    required this.likeCounts,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final gap = theme.space.gridGap;

    // LayoutBuilder + Row/Column 布局替代 GridView，
    // 避免 Web 端 viewport hitTestChildren null 错误
    return LayoutBuilder(
      builder: (context, constraints) {
        final crossAxisCount = AdaptiveGrid.columns(
          availableWidth: constraints.maxWidth,
          minItemWidth: AdaptiveGrid.minPostCardWidth,
          gap: gap,
          maxColumns: context.widthClass.maxGridColumns,
        );

        final rows = <Widget>[];
        for (var i = 0; i < posts.length; i += crossAxisCount) {
          final rowChildren = <Widget>[];
          for (var j = 0; j < crossAxisCount && i + j < posts.length; j++) {
            rowChildren.add(
              Expanded(
                child: PostCard(
                  post: posts[i + j],
                  likeCount: likeCounts[posts[i + j].id] ?? 0,
                  categories: categories,
                  tags: tags,
                ),
              ),
            );
            if (j < crossAxisCount - 1 && i + j + 1 < posts.length) {
              rowChildren.add(SizedBox(width: gap));
            }
          }
          rows.add(
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: rowChildren,
            ),
          );
          if (i + crossAxisCount < posts.length) {
            rows.add(SizedBox(height: gap));
          }
        }

        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: rows,
        );
      },
    );
  }
}
