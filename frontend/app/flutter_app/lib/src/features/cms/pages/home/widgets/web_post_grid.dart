import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';

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
    return Wrap(
      spacing: 16,
      runSpacing: 16,
      children: posts.map((post) {
        return SizedBox(
          width:
              (Breakpoints.webContentMaxWidth * 0.75 -
                  Breakpoints.webSidebarWidth -
                  Breakpoints.webContentPadding * 2 -
                  40 -
                  16) /
              2,
          child: PostCard(post: post, categories: categories, tags: tags),
        );
      }).toList(),
    );
  }
}
