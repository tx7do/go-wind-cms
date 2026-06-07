import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/constants/router_paths.dart';
import 'package:flutter_app/src/core/preference/user_preference_cache.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/services/navigation_service.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/features/cms/widgets/featured_carousel.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation_location.dart';

import 'widgets/nav_bar_link.dart';
import 'widgets/section_header.dart';
import 'widgets/web_post_grid.dart';
import 'widgets/web_sidebar.dart';
import 'widgets/web_footer.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Tag = ContentServiceV1Tag;
typedef NavigationLocation = SiteServiceV1NavigationLocation;

/// 首页 - Web/平板端视图
class HomeWebView extends StatefulWidget {
  const HomeWebView({super.key});

  @override
  State<HomeWebView> createState() => _HomeWebViewState();
}

class _HomeWebViewState extends State<HomeWebView> {
  final _navService = NavigationService();
  final _postService = PostService();
  final _categoryService = CategoryService();
  final _tagService = TagService();
  final _scrollController = ScrollController();

  List<SiteServiceV1Navigation> _navigations = [];
  List<Post> _posts = [];
  List<Category> _categories = [];
  List<Tag> _tags = [];
  bool _isLoading = true;
  bool _isLoadingMore = false;
  bool _hasMorePosts = true;
  int _currentPage = 1;
  static const int _pageSize = 10;

  /// 获取当前用户语言偏好，转换为 API 格式 (zh_CN → zh-CN)
  String? get _currentLocale {
    try {
      final lang = GetIt.instance<UserPreferenceCache>().language;
      if (lang.isEmpty) return null;
      return lang.replaceAll('_', '-');
    } catch (_) {
      return null;
    }
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
      _navService.list(),
      _postService.listPaged(page: _currentPage, pageSize: _pageSize),
      _categoryService.list(),
      _tagService.list(),
    ]);

    if (!mounted) return;

    final postResponse = results[1] as ListPostResponse?;
    final items = postResponse?.items ?? [];
    final total = int.tryParse(postResponse?.total ?? '') ?? 0;

    setState(() {
      _navigations = (results[0] as ListNavigationResponse?)?.items ?? [];
      _posts = items;
      _categories = (results[2] as ListCategoryResponse?)?.items ?? [];
      _tags = (results[3] as ListTagResponse?)?.items ?? [];
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
    final total = int.tryParse(response?.total ?? '') ?? 0;

    setState(() {
      _currentPage++;
      _posts.addAll(items);
      _hasMorePosts = items.length >= _pageSize && _posts.length < total;
      _isLoadingMore = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (_isLoading) {
      return Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        body: const Center(child: CircularProgressIndicator()),
      );
    }

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: CustomScrollView(
        controller: _scrollController,
        slivers: [
          // Web 端 AppBar
          SliverAppBar(
            pinned: true,
            elevation: 0,
            backgroundColor: theme.colorScheme.surface,
            surfaceTintColor: Colors.transparent,
            title: MouseRegion(
              cursor: SystemMouseCursors.click,
              child: Text(
                S.of(context).appName,
                style: TextStyle(
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
                  color: theme.colorScheme.onSurface,
                ),
              ),
            ),
            centerTitle: false,
            actions: [
              // 可滚动的导航链接区域（Flexible 不会抢占 title 空间）
              Flexible(
                child: SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Row(
                    children:
                        getFlatNavItems(
                              _navigations,
                              NavigationLocation.header,
                              locale: _currentLocale,
                            )
                            .map(
                              (item) => NavBarLink(
                                label: item.title ?? '',
                                route: resolveNavRoute(item),
                                isExternal: isExternalLink(item),
                                isOpenNewTab: item.isOpenNewTab == true,
                                isActive: false,
                              ),
                            )
                            .toList(),
                  ),
                ),
              ),
              const SizedBox(width: 8),
              IconButton(
                icon: const Icon(Icons.search, size: 22),
                onPressed: () => context.go('/search'),
                tooltip: S.of(context).search,
              ),
              IconButton(
                icon: const Icon(Icons.settings_outlined, size: 22),
                onPressed: () => context.go(AppRoutePath.login),
                tooltip: S.of(context).settings,
              ),
              const SizedBox(width: 4),
              MouseRegion(
                cursor: SystemMouseCursors.click,
                child: OutlinedButton(
                  onPressed: () => context.go(AppRoutePath.login),
                  style: OutlinedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 8,
                    ),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: Text(S.of(context).login),
                ),
              ),
              const SizedBox(width: 8),
            ],
          ),

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

          // 页脚
          SliverToBoxAdapter(child: WebFooter(navigations: _navigations)),
        ],
      ),
    );
  }
}
