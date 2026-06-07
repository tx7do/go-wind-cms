import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';

/// 收藏页
class BookmarksPage extends StatefulWidget {
  const BookmarksPage({super.key});

  @override
  State<BookmarksPage> createState() => _BookmarksPageState();
}

class _BookmarksPageState extends State<BookmarksPage> {
  final _postService = PostService();

  List<ContentServiceV1Post> _posts = [];
  bool _isLoading = true;

  // TODO: 等待收藏 API 实现后，替换为真实的收藏列表
  // 当前暂从 API 获取文章列表，取前3篇作为占位
  List<ContentServiceV1Post> get _bookmarkedPosts => _posts.take(3).toList();

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final result = await _postService.list();

    if (!mounted) return;

    setState(() {
      _posts = (result as ListPostResponse?)?.items ?? [];
      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    return ResponsiveLayout(
      mobileBody: _buildMobileView(),
      webBody: _buildWebView(),
    );
  }

  // =================== 手机端 ===================

  Widget _buildMobileView() {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            pinned: true,
            floating: true,
            elevation: 0,
            backgroundColor: theme.colorScheme.surface,
            surfaceTintColor: Colors.transparent,
            title: Text(
              S.of(context).myBookmarks,
              style: TextStyle(
                fontSize: 22.sp,
                fontWeight: FontWeight.bold,
                color: theme.colorScheme.onSurface,
              ),
            ),
          ),
          ..._buildContentSlivers(isMobile: true),
        ],
      ),
    );
  }

  // =================== Web 端 ===================

  Widget _buildWebView() {
    final theme = Theme.of(context);
    final crossCount = ResponsiveUtils.postGridColumns(context);
    final bookmarked = _bookmarkedPosts;

    return CustomScrollView(
      slivers: [
        if (bookmarked.isEmpty)
          SliverFillRemaining(child: _buildEmptyState(theme))
        else ...[
          SliverToBoxAdapter(
            child: Center(
              child: ConstrainedBox(
                constraints: const BoxConstraints(
                  maxWidth: Breakpoints.webContentMaxWidth,
                ),
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(24, 8, 24, 12),
                  child: Row(
                    children: [
                      Icon(
                        Icons.bookmark,
                        size: 18,
                        color: theme.colorScheme.primary,
                      ),
                      const SizedBox(width: 6),
                      Text(
                        S.of(context).bookmarkedCount(bookmarked.length),
                        style: TextStyle(
                          fontSize: 13,
                          color: theme.colorScheme.onSurface.withAlpha(160),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
          SliverToBoxAdapter(
            child: Center(
              child: ConstrainedBox(
                constraints: const BoxConstraints(
                  maxWidth: Breakpoints.webContentMaxWidth,
                ),
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 24),
                  child: _WebPostGrid(
                    posts: bookmarked,
                    crossAxisCount: crossCount,
                  ),
                ),
              ),
            ),
          ),
        ],
        SliverToBoxAdapter(child: const SizedBox(height: 32)),
      ],
    );
  }

  // =================== 共享组件 ===================

  List<Widget> _buildContentSlivers({required bool isMobile}) {
    final theme = Theme.of(context);
    final bookmarked = _bookmarkedPosts;

    if (bookmarked.isEmpty) {
      return [SliverFillRemaining(child: _buildEmptyState(theme))];
    }

    return [
      SliverToBoxAdapter(
        child: Padding(
          padding: EdgeInsets.fromLTRB(
            isMobile ? 16.w : 24,
            8,
            isMobile ? 16.w : 24,
            12,
          ),
          child: Row(
            children: [
              Icon(
                Icons.bookmark,
                size: isMobile ? 18.sp : 18,
                color: theme.colorScheme.primary,
              ),
              SizedBox(width: isMobile ? 6.w : 6),
              Text(
                S.of(context).bookmarkedCount(bookmarked.length),
                style: TextStyle(
                  fontSize: isMobile ? 13.sp : 13,
                  color: theme.colorScheme.onSurface.withAlpha(160),
                ),
              ),
            ],
          ),
        ),
      ),
      SliverPadding(
        padding: EdgeInsets.symmetric(horizontal: isMobile ? 16.w : 24),
        sliver: SliverList(
          delegate: SliverChildBuilderDelegate(
            (context, index) => Padding(
              padding: EdgeInsets.only(bottom: isMobile ? 12.h : 12),
              child: PostCard(post: bookmarked[index]),
            ),
            childCount: bookmarked.length,
          ),
        ),
      ),
    ];
  }

  Widget _buildEmptyState(ThemeData theme) {
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            Icons.bookmark_border,
            size: 64,
            color: theme.colorScheme.onSurface.withAlpha(60),
          ),
          const SizedBox(height: 16),
          Text(
            S.of(context).noBookmarks,
            style: TextStyle(
              fontSize: 15,
              color: theme.colorScheme.onSurface.withAlpha(120),
            ),
          ),
          const SizedBox(height: 8),
          Text(
            S.of(context).bookmarkHint,
            style: TextStyle(
              fontSize: 13,
              color: theme.colorScheme.onSurface.withAlpha(80),
            ),
          ),
        ],
      ),
    );
  }
}

/// Web 端文章网格（避免 GridView viewport hitTest 错误）
class _WebPostGrid extends StatelessWidget {
  final List<ContentServiceV1Post> posts;
  final int crossAxisCount;

  const _WebPostGrid({
    required this.posts,
    required this.crossAxisCount,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        const crossAxisSpacing = 16.0;
        const mainAxisSpacing = 16.0;
        const childAspectRatio = 1.1;
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
                child: PostCard(post: posts[i + j]),
              ),
            );
            if (j < crossAxisCount - 1 && i + j + 1 < posts.length) {
              rowChildren.add(const SizedBox(width: crossAxisSpacing));
            }
          }
          rows.add(Row(crossAxisAlignment: CrossAxisAlignment.start, children: rowChildren));
          if (i + crossAxisCount < posts.length) {
            rows.add(const SizedBox(height: mainAxisSpacing));
          }
        }

        return Column(crossAxisAlignment: CrossAxisAlignment.start, children: rows);
      },
    );
  }
}
