import 'package:flutter/material.dart';

import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Category,
        ContentServiceV1Tag;

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Tag = ContentServiceV1Tag;

/// Web 端文章网格
///
/// 使用 Wrap 实现自适应网格布局，每行两列。
class WebPostGrid extends StatelessWidget {
  final List<Post> posts;
  final List<Category> categories;
  final List<Tag> tags;

  const WebPostGrid({
    super.key,
    required this.posts,
    required this.categories,
    required this.tags,
  });

  @override
  Widget build(BuildContext context) {
    // 使用 LayoutBuilder + Row/Column 布局替代 GridView，
    // 避免 Web 端 viewport hitTestChildren null 错误
    return LayoutBuilder(
      builder: (context, constraints) {
        final crossAxisSpacing = 16.0;
        final mainAxisSpacing = 16.0;
        final childAspectRatio = 1.1;
        final crossAxisCount = 2;
        final availableWidth = constraints.maxWidth - crossAxisSpacing * (crossAxisCount - 1);
        final childWidth = availableWidth / crossAxisCount;
        final childHeight = childWidth / childAspectRatio;

        final rows = <Widget>[];
        for (var i = 0; i < posts.length; i += crossAxisCount) {
          final rowChildren = <Widget>[];
          for (var j = 0; j < crossAxisCount && i + j < posts.length; j++) {
            rowChildren.add(
              SizedBox(
                width: childWidth,
                height: childHeight,
                child: PostCard(
                  post: posts[i + j],
                  categories: categories,
                  tags: tags,
                ),
              ),
            );
            if (j < crossAxisCount - 1 && i + j + 1 < posts.length) {
              rowChildren.add(SizedBox(width: crossAxisSpacing));
            }
          }
          rows.add(
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: rowChildren,
            ),
          );
          if (i + crossAxisCount < posts.length) {
            rows.add(SizedBox(height: mainAxisSpacing));
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
