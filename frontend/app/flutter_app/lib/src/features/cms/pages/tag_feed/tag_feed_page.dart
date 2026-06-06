import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/src/features/cms/data/mock_data.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';

typedef Post = ContentServiceV1Post;

/// 标签文章列表页
class TagFeedPage extends StatelessWidget {
  final int tagId;

  const TagFeedPage({super.key, required this.tagId});

  @override
  Widget build(BuildContext context) {
    return ResponsiveLayout(
      mobileBody: _buildView(context, isMobile: true),
      webBody: _buildView(context, isMobile: false),
    );
  }

  Widget _buildView(BuildContext context, {required bool isMobile}) {
    final theme = Theme.of(context);
    final tag = mockTags.firstWhere(
      (t) => t.id != null && t.id == tagId,
      orElse: () => mockTags.first,
    );
    final tagName = (tag.translations ?? []).isNotEmpty
        ? tag.translations!.first.name ?? ''
        : '';
    final posts = mockPosts.where((p) => (p.tagIds ?? []).contains(tag.id!)).toList();

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, size: 22),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.tag,
              size: isMobile ? 18.sp : 18,
              color: theme.colorScheme.primary,
            ),
            SizedBox(width: isMobile ? 4.w : 4),
            Text(
              tagName,
              style: TextStyle(
                fontSize: isMobile ? 18.sp : 18,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
        centerTitle: true,
        bottom: PreferredSize(
          preferredSize: Size.fromHeight(isMobile ? 32.h : 32),
          child: Padding(
            padding: EdgeInsets.only(bottom: isMobile ? 12.h : 12),
            child: Text(
              '${tag.postCount} 篇相关文章',
              style: TextStyle(
                fontSize: isMobile ? 13.sp : 13,
                color: theme.colorScheme.onSurface.withAlpha(120),
              ),
            ),
          ),
        ),
      ),
      body: posts.isEmpty
          ? Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    Icons.article_outlined,
                    size: 56,
                    color: theme.colorScheme.onSurface.withAlpha(60),
                  ),
                  const SizedBox(height: 14),
                  Text(
                    '暂无相关文章',
                    style: TextStyle(
                      fontSize: 15,
                      color: theme.colorScheme.onSurface.withAlpha(120),
                    ),
                  ),
                ],
              ),
            )
          : _buildBody(context, posts, isMobile),
    );
  }

  Widget _buildBody(BuildContext context, List<Post> posts, bool isMobile) {
    if (!isMobile) {
      // Web 端：居中 + 网格
      return Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(
            maxWidth: Breakpoints.webContentMaxWidth,
          ),
          child: ListView.builder(
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 8),
            itemCount: posts.length,
            itemBuilder: (context, index) => Padding(
              padding: const EdgeInsets.only(bottom: 12),
              child: PostCard(post: posts[index]),
            ),
          ),
        ),
      );
    }

    return ListView.builder(
      padding: EdgeInsets.symmetric(horizontal: 16.w, vertical: 8.h),
      itemCount: posts.length,
      itemBuilder: (context, index) => Padding(
        padding: EdgeInsets.only(bottom: 12.h),
        child: PostCard(post: posts[index]),
      ),
    );
  }
}
