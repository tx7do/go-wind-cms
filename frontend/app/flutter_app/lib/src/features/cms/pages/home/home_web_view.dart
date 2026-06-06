import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/services/navigation_service.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/features/cms/services/comment_service.dart';
import 'package:flutter_app/src/features/cms/widgets/featured_carousel.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_cloud.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_comment.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation_location.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_list_navigation_response.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_post_response.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_category_response.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_tag_response.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_list_comment_response.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Tag = ContentServiceV1Tag;
typedef Comment = CommentServiceV1Comment;
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
  final _commentService = CommentService();

  List<SiteServiceV1Navigation> _navigations = [];
  List<Post> _posts = [];
  List<Category> _categories = [];
  List<Tag> _tags = [];
  List<Comment> _comments = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final results = await Future.wait([
      _navService.list(),
      _postService.list(),
      _categoryService.list(),
      _tagService.list(),
      _commentService.list(),
    ]);

    if (!mounted) return;

    setState(() {
      _navigations = (results[0] as ListNavigationResponse?)?.items ?? [];
      _posts = (results[1] as ListPostResponse?)?.items ?? [];
      _categories = (results[2] as ListCategoryResponse?)?.items ?? [];
      _tags = (results[3] as ListTagResponse?)?.items ?? [];
      _comments = (results[4] as ListCommentResponse?)?.items ?? [];
      _isLoading = false;
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
                'Wind CMS',
                style: TextStyle(
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
                  color: theme.colorScheme.onSurface,
                ),
              ),
            ),
            centerTitle: false,
            actions: [
              ...getFlatNavItems(
                _navigations,
                NavigationLocation.header,
              ).map(
                (item) => _NavBarLink(item.title ?? '', resolveNavRoute(item) == '/'),
              ),
              const SizedBox(width: 16),
              IconButton(
                icon: const Icon(Icons.search, size: 22),
                onPressed: () {},
                tooltip: '搜索',
              ),
              const SizedBox(width: 16),
              IconButton(
                icon: const Icon(Icons.settings_outlined, size: 22),
                onPressed: () {},
                tooltip: '设置',
              ),
              const SizedBox(width: 8),
              MouseRegion(
                cursor: SystemMouseCursors.click,
                child: OutlinedButton(
                  onPressed: () {},
                  style: OutlinedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 8,
                    ),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text('登录'),
                ),
              ),
              const SizedBox(width: 16),
            ],
          ),

          // 主内容区域
          SliverToBoxAdapter(
            child: WebContentCenter(
              maxWidth: Breakpoints.webContentMaxWidth,
              padding: const EdgeInsets.symmetric(
                horizontal: Breakpoints.webContentPadding,
                vertical: 20,
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
                          posts: _posts.where((p) => p.isFeatured == true).toList(),
                          categories: _categories,
                        ),
                        const SizedBox(height: 20),
                        _SectionHeader(title: '最新文章', count: _posts.length),
                        const SizedBox(height: 12),
                        _WebPostGrid(
                          posts: _posts,
                          categories: _categories,
                          tags: _tags,
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(width: 24),
                  // 右侧固定侧边栏
                  SizedBox(
                    width: Breakpoints.webSidebarWidth,
                    child: _Sidebar(
                      categories: _categories,
                      tags: _tags,
                      comments: _comments,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

/// 导航栏链接
class _NavBarLink extends StatefulWidget {
  final String label;
  final bool isActive;

  const _NavBarLink(this.label, this.isActive);

  @override
  State<_NavBarLink> createState() => _NavBarLinkState();
}

class _NavBarLinkState extends State<_NavBarLink> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = widget.isActive
        ? theme.colorScheme.primary
        : theme.colorScheme.onSurface;

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          border: widget.isActive
              ? Border(
                  bottom: BorderSide(
                    color: theme.colorScheme.primary,
                    width: 2,
                  ),
                )
              : null,
        ),
        child: Text(
          widget.label,
          style: TextStyle(
            fontSize: 14,
            fontWeight: widget.isActive ? FontWeight.w600 : FontWeight.w400,
            color: _isHovered && !widget.isActive
                ? theme.colorScheme.primary
                : color,
          ),
        ),
      ),
    );
  }
}

/// 分区标题
class _SectionHeader extends StatelessWidget {
  final String title;
  final int count;

  const _SectionHeader({required this.title, required this.count});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Row(
      children: [
        Container(
          width: 4,
          height: 18,
          decoration: BoxDecoration(
            color: theme.colorScheme.primary,
            borderRadius: BorderRadius.circular(2),
          ),
        ),
        const SizedBox(width: 8),
        Text(
          title,
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
            color: theme.colorScheme.onSurface,
          ),
        ),
        const Spacer(),
        Text(
          '$count 篇',
          style: TextStyle(
            fontSize: 13,
            color: theme.colorScheme.onSurface.withAlpha(128),
          ),
        ),
      ],
    );
  }
}

/// Web 端文章网格
class _WebPostGrid extends StatelessWidget {
  final List<Post> posts;
  final List<Category> categories;
  final List<Tag> tags;

  const _WebPostGrid({
    required this.posts,
    required this.categories,
    required this.tags,
  });

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 16,
      runSpacing: 16,
      children: posts.map((post) {
        return SizedBox(
          width:
              (Breakpoints.webContentMaxWidth * 0.75 -
                  Breakpoints.webSidebarWidth -
                  Breakpoints.webContentPadding * 2 -
                  24 -
                  16) /
              2,
          child: PostCard(
            post: post,
            categories: categories,
            tags: tags,
          ),
        );
      }).toList(),
    );
  }
}

/// 右侧边栏
class _Sidebar extends StatelessWidget {
  final List<Category> categories;
  final List<Tag> tags;
  final List<Comment> comments;

  const _Sidebar({
    required this.categories,
    required this.tags,
    required this.comments,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // 分类卡片
        Card(
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
            side: BorderSide(
              color: theme.colorScheme.onSurface.withAlpha(
                (0.06 * 255).round(),
              ),
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(
                      Icons.folder_outlined,
                      size: 18,
                      color: theme.colorScheme.primary,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      '文章分类',
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: categories.map((cat) {
                    final name = (cat.translations ?? []).isNotEmpty
                        ? (cat.translations ?? []).first.name ?? ''
                        : '';
                    return ActionChip(
                      onPressed: () {},
                      label: Text(name, style: const TextStyle(fontSize: 13)),
                      padding: const EdgeInsets.symmetric(horizontal: 4),
                    );
                  }).toList(),
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),

        // 标签云
        Card(
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
            side: BorderSide(
              color: theme.colorScheme.onSurface.withAlpha(
                (0.06 * 255).round(),
              ),
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: TagCloud(tags: tags),
          ),
        ),
        const SizedBox(height: 16),

        // 热门评论
        Card(
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(14),
            side: BorderSide(
              color: theme.colorScheme.onSurface.withAlpha(
                (0.06 * 255).round(),
              ),
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(
                      Icons.comment_outlined,
                      size: 18,
                      color: theme.colorScheme.primary,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      '热门评论',
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                ...comments
                    .take(3)
                    .map(
                      (comment) => Padding(
                        padding: const EdgeInsets.only(bottom: 12),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              comment.authorName ?? '',
                              style: TextStyle(
                                fontSize: 12,
                                fontWeight: FontWeight.w500,
                                color: theme.colorScheme.onSurface,
                              ),
                            ),
                            const SizedBox(height: 4),
                            Text(
                              comment.content ?? '',
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis,
                              style: TextStyle(
                                fontSize: 13,
                                color: theme.colorScheme.onSurface.withAlpha(
                                  160,
                                ),
                                height: 1.4,
                              ),
                            ),
                            Divider(
                              height: 16,
                              color: theme.colorScheme.onSurface.withAlpha(
                                (0.06 * 255).round(),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
