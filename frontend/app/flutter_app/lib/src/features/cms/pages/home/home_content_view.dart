import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/src/features/cms/widgets/featured_carousel.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_cloud.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/src/features/cms/widgets/category_tabs.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Tag = ContentServiceV1Tag;

/// 首页内容视图（分类 Tab + 文章列表 + 轮播图）
///
/// 嵌入在 HomeMobileView 的第一个 tab 中。
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

    setState(() {
      _categories = (results[0] as ListCategoryResponse?)?.items ?? [];
      // 按 sortOrder 排序
      _categories.sort((a, b) => (a.sortOrder ?? 0).compareTo(b.sortOrder ?? 0));
      
      _posts = (results[1] as ListPostResponse?)?.items ?? [];
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

    if (_categories.isEmpty || _tabController == null) {
      return Scaffold(
        backgroundColor: theme.scaffoldBackgroundColor,
        body: const Center(child: CircularProgressIndicator()),
      );
    }

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: NestedScrollView(
        headerSliverBuilder: (context, innerBoxIsScrolled) {
          return [
            SliverAppBar(
              floating: true,
              snap: true,
              pinned: false,
              elevation: 0,
              backgroundColor: theme.colorScheme.surface,
              surfaceTintColor: Colors.transparent,
              title: Text(
                S.of(context).appName,
                style: TextStyle(
                  fontSize: 22.sp,
                  fontWeight: FontWeight.bold,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              actions: [
                IconButton(
                  icon: Icon(Icons.search, size: 24.sp),
                  onPressed: () => context.go('/search'),
                  tooltip: S.of(context).search,
                ),
                SizedBox(width: 8.w),
              ],
              bottom: PreferredSize(
                preferredSize: Size.fromHeight(46.h),
                child: Container(
                  color: theme.colorScheme.surface,
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
      ),
    );
  }

  Widget _buildTabContent(
    BuildContext context,
    List<Post> posts,
    bool isRecommend,
  ) {
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
            padding: EdgeInsets.fromLTRB(20.w, 16.h, 20.w, 8.h),
            child: Row(
              children: [
                Container(
                  width: 4.w,
                  height: 18.h,
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.primary,
                    borderRadius: BorderRadius.circular(2.r),
                  ),
                ),
                SizedBox(width: 8.w),
                Text(
                  isRecommend
                      ? S.of(context).latestPosts
                      : S.of(context).relatedArticles,
                  style: TextStyle(
                    fontSize: 16.sp,
                    fontWeight: FontWeight.w600,
                    color: Theme.of(context).colorScheme.onSurface,
                  ),
                ),
                const Spacer(),
                Text(
                  S.of(context).postsCount(posts.length),
                  style: TextStyle(
                    fontSize: 13.sp,
                    color: Theme.of(
                      context,
                    ).colorScheme.onSurface.withAlpha(128),
                  ),
                ),
              ],
            ),
          ),
        ),

        SliverPadding(
          padding: EdgeInsets.symmetric(horizontal: 16.w),
          sliver: SliverList(
            delegate: SliverChildBuilderDelegate(
              (context, index) => Padding(
                padding: EdgeInsets.only(bottom: 12.h),
                child: PostCard(
                  post: posts[index],
                  categories: _categories,
                  tags: _tags,
                ),
              ),
              childCount: posts.length,
            ),
          ),
        ),

        SliverToBoxAdapter(child: SizedBox(height: 24.h)),
      ],
    );
  }
}
