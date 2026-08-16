import 'package:flutter/material.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/features/cms/widgets/featured_carousel.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Category,
        ContentServiceV1Tag, InteractionServiceV1TargetType,
        InteractionServiceV1CounterMetric;
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';

import 'widgets/section_header.dart';
import 'widgets/web_post_grid.dart';
import 'widgets/web_sidebar.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Tag = ContentServiceV1Tag;

/// 首页 - Web/平板端视图
class HomeWebView extends StatefulWidget {
  const HomeWebView({super.key});

  @override
  State<HomeWebView> createState() => _HomeWebViewState();
}

class _HomeWebViewState extends State<HomeWebView> {
  final _postService = PostService();
  final _categoryService = CategoryService();
  final _tagService = TagService();
  final _scrollController = ScrollController();

  List<Post> _posts = [];
  List<Category> _categories = [];
  List<Tag> _tags = [];
  bool _isLoading = true;
  bool _isLoadingMore = false;
  bool _hasMorePosts = true;
  int _currentPage = 1;
  static const int _pageSize = 10;

  // 点赞计数缓存：key 为 post.id。计数来自 interaction_counter 表
  // （post.likes 列已移除），由 InteractionService.GetCounts 批量查询。
  // 分页追加新 post 后整体重查，保证 _likeCounts 覆盖当前 _posts 全集。
  Map<int, int> _likeCounts = {};

  final _interactionService = InteractionService();

  /// 刷新当前 _posts 全集的点赞计数到 _likeCounts。
  Future<void> _refreshLikeCounts() async {
    final ids = _posts
        .where((p) => p.id != null)
        .map((p) => p.id!)
        .toList(growable: false);
    if (ids.isEmpty) {
      _likeCounts = {};
      return;
    }
    final countsResp = await _interactionService.getCounts(
        InteractionServiceV1TargetType.targetTypePost,
        ids,
        [InteractionServiceV1CounterMetric.counterMetricLike]);
    _likeCounts = InteractionService.extractCountMap(
        countsResp,
        ids,
        InteractionServiceV1CounterMetric.counterMetricLike);
  }

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    _loadData();
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScroll);
    _scrollController.dispose();
    super.dispose();
  }

  /// 滚动监听：距离底部 200px 时预加载下一页
  void _onScroll() {
    if (_scrollController.hasClients &&
        _scrollController.position.maxScrollExtent > 0 &&
        _scrollController.position.pixels >=
            _scrollController.position.maxScrollExtent - 200 &&
        !_isLoadingMore &&
        _hasMorePosts &&
        !_isLoading) {
      _loadMorePosts();
    }
  }

  Future<void> _loadData() async {
    final results = await Future.wait([
      _postService.listPaged(page: _currentPage, pageSize: _pageSize),
      _categoryService.list(),
      _tagService.list(),
    ]);

    if (!mounted) return;

    final postResponse = results[0] as ListPostResponse?;
    final items = postResponse?.items ?? [];
    final total = postResponse?.total ?? 0;

    _posts = items;
    await _refreshLikeCounts();
    if (!mounted) return;

    setState(() {
      _categories = (results[1] as ListCategoryResponse?)?.items ?? [];
      _tags = (results[2] as ListTagResponse?)?.items ?? [];
      _hasMorePosts = items.length >= _pageSize && _posts.length < total;
      _isLoading = false;
    });
  }

  Future<void> _loadMorePosts() async {
    if (_isLoadingMore || !_hasMorePosts) return;

    setState(() => _isLoadingMore = true);

    final response = await _postService.listPaged(
      page: _currentPage + 1,
      pageSize: _pageSize,
    );

    if (!mounted) return;

    final items = response?.items ?? [];
    final total = response?.total ?? 0;

    _posts.addAll(items);
    await _refreshLikeCounts();
    if (!mounted) return;

    setState(() {
      _currentPage++;
      _hasMorePosts = items.length >= _pageSize && _posts.length < total;
      _isLoadingMore = false;
    });
  }

  @override
  Widget build(BuildContext context) {

    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    return CustomScrollView(
        controller: _scrollController,
        slivers: [
          // 主内容区域
          SliverToBoxAdapter(
            child: WebContentCenter(
              maxWidth: Breakpoints.webContentMaxWidth,
              padding: const EdgeInsets.symmetric(
                horizontal: Breakpoints.webContentPadding,
                vertical: 28,
              ),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 左侧主内容
                  Expanded(
                    flex: 3,
                    child: Column(
                      children: [
                        FeaturedCarousel(
                          posts: _posts
                              .where((p) => p.isFeatured == true)
                              .toList(),
                          categories: _categories,
                        ),
                        const SizedBox(height: 32),
                        SectionHeader(
                          title: S.of(context).latestPosts,
                          count: _posts.length,
                        ),
                        const SizedBox(height: 16),
                        WebPostGrid(
                          posts: _posts,
                          categories: _categories,
                          tags: _tags,
                          likeCounts: _likeCounts,
                        ),
                        // 加载更多指示器
                        if (_isLoadingMore)
                          const Padding(
                            padding: EdgeInsets.symmetric(vertical: 24),
                            child: Center(
                              child: SizedBox(
                                width: 20,
                                height: 20,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                ),
                              ),
                            ),
                          ),
                        if (!_hasMorePosts && _posts.isNotEmpty)
                          Padding(
                            padding: const EdgeInsets.symmetric(vertical: 24),
                            child: Center(
                              child: Text(
                                S.of(context).allLoaded,
                                style: TextStyle(
                                  fontSize: 13,
                                  color: Theme.of(
                                    context,
                                  ).colorScheme.onSurface.withAlpha(80),
                                ),
                              ),
                            ),
                          ),
                      ],
                    ),
                  ),
                  const SizedBox(width: 40),
                  // 右侧固定侧边栏
                  SizedBox(
                    width: Breakpoints.webSidebarWidth,
                    child: WebSidebar(tags: _tags),
                  ),
                ],
              ),
            ),
          ),
        ],
    );
  }
}
