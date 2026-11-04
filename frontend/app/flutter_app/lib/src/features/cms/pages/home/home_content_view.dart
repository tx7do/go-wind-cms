import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Category,
        ContentServiceV1Tag, InteractionServiceV1TargetType,
        InteractionServiceV1CounterMetric;
import 'package:flutter_app/src/features/cms/widgets/featured_carousel.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_cloud.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/src/features/cms/widgets/category_tabs.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/features/cms/pages/home/widgets/section_header.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Tag = ContentServiceV1Tag;

/// 首页内容视图（分类 Tab + 文章列表 + 轮播图）
///
/// 首页正文视图（窄屏模式的内容主体，宽屏用 HomeWebView）。
class HomeContentView extends StatefulWidget {
  const HomeContentView({super.key});

  @override
  State<HomeContentView> createState() => _HomeContentViewState();
}

class _HomeContentViewState extends State<HomeContentView>
    with SingleTickerProviderStateMixin {
  final _categoryService = CategoryService();
  final _postService = PostService();
  final _tagService = TagService();

  TabController? _tabController;
  int _currentCategoryIndex = 0;

  List<Category> _categories = [];
  List<Post> _posts = [];
  List<Tag> _tags = [];

  // 点赞计数缓存：key 为 post.id。计数来自 interaction_counter 表
  // （post.likes 列已移除），由 InteractionService.GetCounts 批量查询。
  Map<int, int> _likeCounts = {};

  final _interactionService = InteractionService();

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  @override
  void dispose() {
    _tabController?.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    final results = await Future.wait([
      _categoryService.list(),
      _postService.list(),
      _tagService.list(),
    ]);

    if (!mounted) return;

    final posts = (results[1] as ListPostResponse?)?.items ?? [];
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
      _categories = (results[0] as ListCategoryResponse?)?.items ?? [];
      // 按 sortOrder 排序
      _categories.sort((a, b) => (a.sortOrder ?? 0).compareTo(b.sortOrder ?? 0));

      _posts = posts;
      _likeCounts = likeCounts;
      _tags = (results[2] as ListTagResponse?)?.items ?? [];

      // 插入"推荐"虚拟分类
      if (_categories.isNotEmpty) {
        _categories.insert(
          0,
          Category(
            id: 0,
            sortOrder: -1,
            translations: [
              CategoryTranslation(
                languageCode: 'zh',
                name: S.current.recommend,
                slug: 'recommend',
              ),
            ],
          ),
        );
      }

      _tabController?.dispose();
      _tabController = TabController(length: _categories.length, vsync: this);
      _tabController!.addListener(() {
        if (!_tabController!.indexIsChanging) {
          setState(() => _currentCategoryIndex = _tabController!.index);
        }
      });
    });
  }

  List<Post> get _filteredPosts {
    if (_currentCategoryIndex == 0) return _posts;
    if (_currentCategoryIndex >= _categories.length) return [];
    final category = _categories[_currentCategoryIndex];
    return _posts
        .where((p) => (p.categoryIds ?? []).contains(category.id))
        .toList();
  }

  List<Post> get _featuredPosts =>
      _posts.where((p) => p.isFeatured == true).toList();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    // 壳层（AdaptiveShell）在 tab 根位置提供 Scaffold，这里只管内容
    if (_categories.isEmpty || _tabController == null) {
      return const Center(child: CircularProgressIndicator());
    }

    final space = theme.space;
    final p = theme.palette;

    return NestedScrollView(
      headerSliverBuilder: (context, innerBoxIsScrolled) {
        return [
          SliverAppBar(
            floating: true,
            snap: true,
            pinned: false,
            title: Text(S.of(context).appName),
            actions: [
              IconButton(
                icon: const Icon(Icons.search),
                onPressed: () => context.go('/search'),
                tooltip: S.of(context).search,
              ),
              SizedBox(width: space.sm),
            ],
            bottom: PreferredSize(
              preferredSize: const Size.fromHeight(46),
              child: ColoredBox(
                color: p.surface,
                child: CategoryTabs(
                  categories: _categories,
                  tabController: _tabController!,
                ),
              ),
            ),
          ),
        ];
      },
      body: TabBarView(
        controller: _tabController,
        children: List.generate(_categories.length, (index) {
          final posts = index == 0 ? _posts : _filteredPosts;
          return _buildTabContent(context, posts, index == 0);
        }),
      ),
    );
  }

  Widget _buildTabContent(
    BuildContext context,
    List<Post> posts,
    bool isRecommend,
  ) {
    final space = Theme.of(context).space;
    return CustomScrollView(
      slivers: [
        if (isRecommend) ...[
          SliverToBoxAdapter(
            child: FeaturedCarousel(
              posts: _featuredPosts,
              categories: _categories,
            ),
          ),
          SliverToBoxAdapter(child: TagCloud(tags: _tags)),
        ],

        SliverToBoxAdapter(
          child: Padding(
            padding: EdgeInsets.fromLTRB(
              space.cardPadding,
              space.lg,
              space.cardPadding,
              space.sm,
            ),
            child: SectionHeader(
              title: isRecommend
                  ? S.of(context).latestPosts
                  : S.of(context).relatedArticles,
              count: posts.length,
            ),
          ),
        ),

        SliverPadding(
          padding: EdgeInsets.symmetric(horizontal: space.lg),
          sliver: SliverList(
            delegate: SliverChildBuilderDelegate(
              (context, index) => Padding(
                padding: EdgeInsets.only(bottom: space.listRowGap),
                child: PostCard(
                  post: posts[index],
                  likeCount: _likeCounts[posts[index].id] ?? 0,
                  categories: _categories,
                  tags: _tags,
                ),
              ),
              childCount: posts.length,
            ),
          ),
        ),

        SliverToBoxAdapter(child: SizedBox(height: space.xxl)),
      ],
    );
  }
}
