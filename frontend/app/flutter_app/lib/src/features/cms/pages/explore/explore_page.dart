import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart'
    show TagService, ListTagResponse;
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_chip.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;

/// 发现/分类页
class ExplorePage extends StatefulWidget {
  const ExplorePage({super.key});

  @override
  State<ExplorePage> createState() => _ExplorePageState();
}

class _ExplorePageState extends State<ExplorePage> {
  final _postService = PostService();
  final _categoryService = CategoryService();
  final _tagService = TagService();

  Category? _selectedCategory;
  List<Category> _categories = [];
  List<Post> _posts = [];
  List<ContentServiceV1Tag> _tags = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
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
      _isLoading = false;
    });
  }

  List<Category> get _displayCategories => _categories.skip(1).toList();

  List<Post> get _filteredPosts {
    if (_selectedCategory == null) return _posts;
    return _posts
        .where((p) => (p.categoryIds ?? []).contains(_selectedCategory!.id))
        .toList();
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
    final categories = _displayCategories;
    final filteredPosts = _filteredPosts;

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
              S.of(context).discover,
              style: TextStyle(
                fontSize: 22.sp,
                fontWeight: FontWeight.bold,
                color: theme.colorScheme.onSurface,
              ),
            ),
          ),
          _buildCategoryGridSliver(categories, isMobile: true),
          _buildTagSectionSliver(isMobile: true),
          _buildPostListSliver(filteredPosts, isMobile: true),
          SliverToBoxAdapter(child: SizedBox(height: 24.h)),
        ],
      ),
    );
  }

  // =================== Web 端 ===================

  Widget _buildWebView() {
    final theme = Theme.of(context);
    final categories = _displayCategories;
    final filteredPosts = _filteredPosts;
    final crossCount = ResponsiveUtils.postGridColumns(context);

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            pinned: true,
            elevation: 0,
            backgroundColor: theme.colorScheme.surface,
            surfaceTintColor: Colors.transparent,
            title: Text(
              S.of(context).discover,
              style: const TextStyle(fontSize: 22, fontWeight: FontWeight.bold),
            ),
          ),
          SliverToBoxAdapter(
            child: Center(
              child: ConstrainedBox(
                constraints: const BoxConstraints(
                  maxWidth: Breakpoints.webContentMaxWidth,
                ),
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(24, 16, 24, 0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        S.of(context).browseCategories,
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: theme.colorScheme.onSurface,
                        ),
                      ),
                      const SizedBox(height: 12),
                      _CategoryGrid(
                        categories: categories,
                        crossAxisCount: ResponsiveUtils.categoryGridColumns(
                          context,
                        ),
                        selectedCategory: _selectedCategory,
                        onSelect: (cat) {
                          setState(() {
                            _selectedCategory = _selectedCategory?.id == cat.id
                                ? null
                                : cat;
                          });
                        },
                      ),
                      const SizedBox(height: 24),
                      Text(
                        S.of(context).hotTags,
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                          color: theme.colorScheme.onSurface,
                        ),
                      ),
                      const SizedBox(height: 12),
                      _TagWrap(tags: _tags),
                      const SizedBox(height: 24),
                      _SectionRow(
                        title: _selectedCategory == null
                            ? S.of(context).allPosts
                            : getCategoryName(_selectedCategory),
                        count: filteredPosts.length,
                      ),
                      const SizedBox(height: 12),
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
                    posts: filteredPosts,
                    categories: _categories,
                    tags: _tags,
                    crossAxisCount: crossCount,
                  ),
                ),
              ),
            ),
          ),
          SliverToBoxAdapter(child: const SizedBox(height: 32)),
        ],
      ),
    );
  }

  // =================== 共享组件 ===================

  Widget _buildCategoryGridSliver(
    List<Category> categories, {
    required bool isMobile,
  }) {
    return SliverPadding(
      padding: EdgeInsets.symmetric(horizontal: isMobile ? 16.w : 24),
      sliver: SliverGrid(
        gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: isMobile
              ? 2
              : ResponsiveUtils.categoryGridColumns(context),
          mainAxisSpacing: isMobile ? 10.h : 12,
          crossAxisSpacing: isMobile ? 10.w : 12,
          childAspectRatio: 2.2,
        ),
        delegate: SliverChildBuilderDelegate((context, index) {
          final category = categories[index];
          final isSelected = _selectedCategory?.id == category.id;
          return _CategoryGridItem(
            category: category,
            isSelected: isSelected,
            onTap: () {
              setState(() {
                _selectedCategory = isSelected ? null : category;
              });
            },
            isMobile: isMobile,
          );
        }, childCount: categories.length),
      ),
    );
  }

  Widget _buildTagSectionSliver({required bool isMobile}) {
    final theme = Theme.of(context);
    return SliverToBoxAdapter(
      child: Padding(
        padding: EdgeInsets.fromLTRB(
          isMobile ? 16.w : 24,
          isMobile ? 20.h : 20,
          isMobile ? 16.w : 24,
          4,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              S.of(context).hotTags,
              style: TextStyle(
                fontSize: isMobile ? 16.sp : 16,
                fontWeight: FontWeight.w600,
                color: theme.colorScheme.onSurface,
              ),
            ),
            SizedBox(height: isMobile ? 12.h : 12),
            _TagWrap(tags: _tags),
          ],
        ),
      ),
    );
  }

  Widget _buildPostListSliver(List<Post> posts, {required bool isMobile}) {
    final theme = Theme.of(context);
    return MultiSliver([
      SliverToBoxAdapter(
        child: Padding(
          padding: EdgeInsets.fromLTRB(
            isMobile ? 16.w : 24,
            isMobile ? 20.h : 20,
            isMobile ? 16.w : 24,
            8,
          ),
          child: Row(
            children: [
              Container(
                width: isMobile ? 4.w : 4,
                height: isMobile ? 18.h : 18,
                decoration: BoxDecoration(
                  color: theme.colorScheme.primary,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              SizedBox(width: isMobile ? 8.w : 8),
              Text(
                _selectedCategory == null
                    ? S.of(context).allPosts
                    : getCategoryName(_selectedCategory),
                style: TextStyle(
                  fontSize: isMobile ? 16.sp : 16,
                  fontWeight: FontWeight.w600,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              const Spacer(),
              Text(
                S.of(context).postsCount(posts.length),
                style: TextStyle(
                  fontSize: isMobile ? 13.sp : 13,
                  color: theme.colorScheme.onSurface.withAlpha(128),
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
                post: posts[index],
                categories: _categories,
                tags: _tags,
              ),
            ),
            childCount: posts.length,
          ),
        ),
      ),
    ]);
  }
}

// =================== 辅助组件 ===================

class MultiSliver extends StatelessWidget {
  final List<Widget> children;

  const MultiSliver(this.children, {super.key});

  @override
  Widget build(BuildContext context) {
    return SliverMainAxisGroup(slivers: children);
  }
}

class _CategoryGrid extends StatelessWidget {
  final List<Category> categories;
  final int crossAxisCount;
  final Category? selectedCategory;
  final ValueChanged<Category> onSelect;

  const _CategoryGrid({
    required this.categories,
    required this.crossAxisCount,
    this.selectedCategory,
    required this.onSelect,
  });

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 10,
      runSpacing: 10,
      children: categories.map((cat) {
        final isSelected = selectedCategory?.id == cat.id;
        return _CategoryGridItem(
          category: cat,
          isSelected: isSelected,
          onTap: () => onSelect(cat),
          isMobile: false,
        );
      }).toList(),
    );
  }
}

class _CategoryGridItem extends StatefulWidget {
  final Category category;
  final bool isSelected;
  final VoidCallback onTap;
  final bool isMobile;

  const _CategoryGridItem({
    required this.category,
    required this.isSelected,
    required this.onTap,
    required this.isMobile,
  });

  @override
  State<_CategoryGridItem> createState() => _CategoryGridItemState();
}

class _CategoryGridItemState extends State<_CategoryGridItem> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = getCategoryName(widget.category);

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: Material(
        color: widget.isSelected
            ? theme.colorScheme.primaryContainer
            : _isHovered
            ? theme.colorScheme.surfaceContainerHighest
            : theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(widget.isMobile ? 14.r : 14),
        child: InkWell(
          onTap: widget.onTap,
          borderRadius: BorderRadius.circular(widget.isMobile ? 14.r : 14),
          child: Container(
            width: widget.isMobile ? null : 180,
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(widget.isMobile ? 14.r : 14),
              border: Border.all(
                color: widget.isSelected
                    ? theme.colorScheme.primary
                    : theme.colorScheme.onSurface.withAlpha(
                        (0.08 * 255).round(),
                      ),
                width: widget.isSelected ? 1.5 : 1,
              ),
            ),
            child: Padding(
              padding: EdgeInsets.all(widget.isMobile ? 14.w : 14),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    name,
                    style: TextStyle(
                      fontSize: widget.isMobile ? 15.sp : 15,
                      fontWeight: FontWeight.w600,
                      color: theme.colorScheme.onSurface,
                    ),
                  ),
                  SizedBox(height: widget.isMobile ? 4.h : 4),
                  Text(
                    S
                        .of(context)
                        .postsCountFull(widget.category.postCount ?? 0),
                    style: TextStyle(
                      fontSize: widget.isMobile ? 12.sp : 12,
                      color: theme.colorScheme.onSurface.withAlpha(140),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _TagWrap extends StatelessWidget {
  final List<ContentServiceV1Tag> tags;

  const _TagWrap({required this.tags});

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 8,
      runSpacing: 8,
      children: tags.map((tag) {
        return TagChip(tag: tag, isMobile: false);
      }).toList(),
    );
  }
}

class _SectionRow extends StatelessWidget {
  final String title;
  final int count;

  const _SectionRow({required this.title, required this.count});

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
          S.of(context).postsCount(count),
          style: TextStyle(
            fontSize: 13,
            color: theme.colorScheme.onSurface.withAlpha(128),
          ),
        ),
      ],
    );
  }
}

class _WebPostGrid extends StatelessWidget {
  final List<Post> posts;
  final List<Category> categories;
  final List<ContentServiceV1Tag> tags;
  final int crossAxisCount;

  const _WebPostGrid({
    required this.posts,
    required this.categories,
    required this.tags,
    required this.crossAxisCount,
  });

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 16,
      runSpacing: 16,
      children: posts.map((post) {
        return SizedBox(
          width:
              (Breakpoints.webContentMaxWidth -
                  48 -
                  (crossAxisCount - 1) * 16) /
              crossAxisCount,
          child: PostCard(post: post, categories: categories, tags: tags),
        );
      }).toList(),
    );
  }
}
