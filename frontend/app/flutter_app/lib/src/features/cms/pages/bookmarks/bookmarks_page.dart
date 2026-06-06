import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_post_response.dart';
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
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
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
              '我的收藏',
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

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            pinned: true,
            elevation: 0,
            backgroundColor: theme.colorScheme.surface,
            surfaceTintColor: Colors.transparent,
            title: const Text(
              '我的收藏',
              style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
            ),
          ),
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
                          '已收藏 ${bookmarked.length} 篇文章',
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
                    child: Wrap(
                      spacing: 16,
                      runSpacing: 16,
                      children: bookmarked.map((post) {
                        return SizedBox(
                          width:
                              (Breakpoints.webContentMaxWidth -
                                  48 -
                                  (crossCount - 1) * 16) /
                              crossCount,
                          child: PostCard(post: post),
                        );
                      }).toList(),
                    ),
                  ),
                ),
              ),
            ),
          ],
          SliverToBoxAdapter(child: const SizedBox(height: 32)),
        ],
      ),
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
                '已收藏 ${bookmarked.length} 篇文章',
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
            '还没有收藏的文章',
            style: TextStyle(
              fontSize: 15,
              color: theme.colorScheme.onSurface.withAlpha(120),
            ),
          ),
          const SizedBox(height: 8),
          Text(
            '浏览文章时点击收藏按钮即可保存',
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
