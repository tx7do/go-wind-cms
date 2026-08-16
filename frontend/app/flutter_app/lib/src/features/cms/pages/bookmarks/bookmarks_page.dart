import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1ListPostResponse,
        InteractionServiceV1TargetType, InteractionServiceV1CounterMetric;
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';

import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';

typedef ListPostResponse = ContentServiceV1ListPostResponse;

/// 收藏页
class BookmarksPage extends StatefulWidget {
  const BookmarksPage({super.key});

  @override
  State<BookmarksPage> createState() => _BookmarksPageState();
}

class _BookmarksPageState extends State<BookmarksPage> {
  final _interactionService = InteractionService();

  List<ContentServiceV1Post> _posts = [];
  bool _isLoading = true;

  // 点赞计数缓存：key 为 post.id。计数来自 interaction_counter 表
  // （post.likes 列已移除），由 InteractionService.GetCounts 批量查询。
  Map<int, int> _likeCounts = {};

  // 收藏列表由 InteractionService.ListWatchedPosts 返回当前 viewer 收藏的 post。
  // viewer 身份由后端鉴权上下文确定，需登录态；未登录时后端返回 401，列表为空。
  List<ContentServiceV1Post> get _bookmarkedPosts => _posts;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final result = await _interactionService.listWatchedPosts();

    if (!mounted) return;

    final posts = (result as ListPostResponse?)?.items ?? [];
    final ids = posts
        .where((p) => p.id != null)
        .map((p) => p.id!)
        .toList(growable: false);
    Map<int, int> likeCounts = {};
    if (ids.isNotEmpty) {
      final countsResp = await _interactionService.getCounts(
          InteractionServiceV1TargetType.targetTypePost,
          ids,
          [InteractionServiceV1CounterMetric.counterMetricLike]);
      likeCounts = InteractionService.extractCountMap(
          countsResp,
          ids,
          InteractionServiceV1CounterMetric.counterMetricLike);
    }
    if (!mounted) return;

    setState(() {
      _posts = posts;
      _likeCounts = likeCounts;
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
                    likeCounts: _likeCounts,
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
              child: PostCard(
                  post: bookmarked[index],
                  likeCount:
                      _likeCounts[bookmarked[index].id] ?? 0),
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
  final Map<int, int> likeCounts;

  const _WebPostGrid({
    required this.posts,
    required this.crossAxisCount,
    required this.likeCounts,
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
                child: PostCard(
                    post: posts[i + j],
                    likeCount: likeCounts[posts[i + j].id] ?? 0),
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
