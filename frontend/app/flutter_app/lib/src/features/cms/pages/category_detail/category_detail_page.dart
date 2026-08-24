import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Category,
        InteractionServiceV1TargetType, InteractionServiceV1CounterMetric;
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart' hide Category;
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart' hide Category;

typedef Post = ContentServiceV1Post;

/// 分类详情页：分类信息 + 该分类下的文章列表。
///
/// 文章按 category_id 服务端过滤（PostListPage 同款参数）；
/// 分类本身从 list 结果里取（list 响应必带 translations，
/// 而 get 未暴露 fieldMask，可能不回传翻译）。
class CategoryDetailPage extends StatefulWidget {
  final int categoryId;

  const CategoryDetailPage({super.key, required this.categoryId});

  @override
  State<CategoryDetailPage> createState() => _CategoryDetailPageState();
}

class _CategoryDetailPageState extends State<CategoryDetailPage> {
  final _postService = PostService();
  final _categoryService = CategoryService();
  final _interactionService = InteractionService();

  List<Post> _posts = [];
  List<Category> _categories = [];
  bool _isLoading = true;

  // 点赞计数缓存：key 为 post.id。计数来自 interaction_counter 表
  // （post.likes 列已移除），由 InteractionService.GetCounts 批量查询。
  Map<int, int> _likeCounts = {};

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final results = await Future.wait([
      _postService.list(PaginationQuery(formValues: {'category_id': widget.categoryId})),
      _categoryService.list(),
    ]);

    if (!mounted) return;

    final posts = (results[0] as ListPostResponse?)?.items ?? [];
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
      _categories = (results[1] as ListCategoryResponse?)?.items ?? [];
      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    return ResponsiveLayout(
      mobileBody: _buildView(context, isMobile: true),
      webBody: _buildView(context, isMobile: false),
    );
  }

  Widget _buildView(BuildContext context, {required bool isMobile}) {
    final theme = Theme.of(context);
    final category = _categories.firstWhere(
      (c) => c.id != null && c.id == widget.categoryId,
      orElse: () => ContentServiceV1Category(),
    );
    final categoryName = getCategoryName(category);

    final appBar = AppBar(
      backgroundColor: theme.colorScheme.surface,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
      leading: const AppBackButton(),
      title: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.category_outlined, size: isMobile ? 18.sp : 18, color: theme.colorScheme.primary),
          SizedBox(width: isMobile ? 4.w : 4),
          Flexible(
            child: Text(
              categoryName,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(fontSize: isMobile ? 18.sp : 18, fontWeight: FontWeight.bold),
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
            S.of(context).postsCountFull(category.postCount ?? 0),
            style: TextStyle(fontSize: isMobile ? 13.sp : 13, color: theme.colorScheme.onSurface.withAlpha(120)),
          ),
        ),
      ),
    );

    final body = _posts.isEmpty
        ? Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.article_outlined, size: 56, color: theme.colorScheme.onSurface.withAlpha(60)),
                const SizedBox(height: 14),
                Text(S.of(context).noRelatedPosts, style: TextStyle(fontSize: 15, color: theme.colorScheme.onSurface.withAlpha(120))),
              ],
            ),
          )
        : _buildBody(context, isMobile);

    // Web 端由 WebShellLayout 提供 Scaffold
    if (!isMobile) return body;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: appBar,
      body: body,
    );
  }

  Widget _buildBody(BuildContext context, bool isMobile) {
    if (!isMobile) {
      // Web 端：CustomScrollView 避免无界高度问题
      return CustomScrollView(
        slivers: [
          SliverPadding(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            sliver: SliverToBoxAdapter(
              child: Center(
                child: ConstrainedBox(
                  constraints: const BoxConstraints(maxWidth: Breakpoints.webContentMaxWidth),
                  child: Column(
                    children: _posts.map((post) => Padding(
                      padding: const EdgeInsets.only(bottom: 12),
                      child: PostCard(
                          post: post, likeCount: _likeCounts[post.id] ?? 0),
                    )).toList(),
                  ),
                ),
              ),
            ),
          ),
          const SliverToBoxAdapter(child: SizedBox(height: 32)),
        ],
      );
    }

    return ListView.builder(
      padding: EdgeInsets.symmetric(horizontal: 16.w, vertical: 8.h),
      itemCount: _posts.length,
      itemBuilder: (context, index) => Padding(
        padding: EdgeInsets.only(bottom: 12.h),
        child: PostCard(
            post: _posts[index],
            likeCount: _likeCounts[_posts[index].id] ?? 0),
      ),
    );
  }
}
