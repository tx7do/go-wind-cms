import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/src/features/cms/data/mock_data.dart';
import 'package:flutter_app/src/features/cms/widgets/featured_carousel.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_cloud.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/src/features/cms/widgets/category_tabs.dart';

typedef Post = ContentServiceV1Post;

/// 首页 - 手机端视图（原有单栏瀑布流布局）
class HomeMobileView extends StatefulWidget {
  const HomeMobileView({super.key});

  @override
  State<HomeMobileView> createState() => _HomeMobileViewState();
}

class _HomeMobileViewState extends State<HomeMobileView>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  int _currentCategoryIndex = 0;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: mockCategories.length, vsync: this);
    _tabController.addListener(() {
      if (!_tabController.indexIsChanging) {
        setState(() {
          _currentCategoryIndex = _tabController.index;
        });
      }
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  List<Post> get _filteredPosts {
    if (_currentCategoryIndex == 0) return mockPosts;
    final category = mockCategories[_currentCategoryIndex];
    return mockPosts.where((p) => (p.categoryIds ?? []).contains(category.id!)).toList();
  }

  List<Post> get _featuredPosts =>
      mockPosts.where((p) => p.isFeatured == true).toList();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

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
                'Wind CMS',
                style: TextStyle(
                  fontSize: 22.sp,
                  fontWeight: FontWeight.bold,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              actions: [
                IconButton(
                  icon: Icon(Icons.search, size: 24.sp),
                  onPressed: () {},
                  tooltip: '搜索',
                ),
                SizedBox(width: 8.w),
              ],
              bottom: PreferredSize(
                preferredSize: Size.fromHeight(46.h),
                child: Container(
                  color: theme.colorScheme.surface,
                  child: CategoryTabs(
                    categories: mockCategories,
                    tabController: _tabController,
                  ),
                ),
              ),
            ),
          ];
        },
        body: TabBarView(
          controller: _tabController,
          children: List.generate(mockCategories.length, (index) {
            final posts = index == 0 ? mockPosts : _filteredPosts;
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
          SliverToBoxAdapter(child: FeaturedCarousel(posts: _featuredPosts)),
          SliverToBoxAdapter(child: TagCloud(tags: mockTags)),
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
                  isRecommend ? '最新文章' : '相关文章',
                  style: TextStyle(
                    fontSize: 16.sp,
                    fontWeight: FontWeight.w600,
                    color: Theme.of(context).colorScheme.onSurface,
                  ),
                ),
                const Spacer(),
                Text(
                  '${posts.length} 篇',
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
                child: PostCard(post: posts[index]),
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
