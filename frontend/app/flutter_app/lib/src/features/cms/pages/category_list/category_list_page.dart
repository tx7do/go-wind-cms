import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_category_response.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

typedef Category = ContentServiceV1Category;

/// 分类列表页
///
/// 展示所有分类，点击可跳转到该分类下的文章列表。
class CategoryListPage extends StatefulWidget {
  const CategoryListPage({super.key});

  @override
  State<CategoryListPage> createState() => _CategoryListPageState();
}

class _CategoryListPageState extends State<CategoryListPage> {
  final _categoryService = CategoryService();

  List<Category> _categories = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final result = await _categoryService.list();
    if (!mounted) return;

    setState(() {
      _categories = (result as ContentServiceV1ListCategoryResponse?)?.items ?? [];
      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    return ResponsiveLayout(
      mobileBody: _buildView(isMobile: true),
      webBody: _buildView(isMobile: false),
    );
  }

  Widget _buildView({required bool isMobile}) {
    final theme = Theme.of(context);
    final loc = S.of(context);

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, size: 22),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.folder_outlined,
              size: isMobile ? 20.sp : 20,
              color: theme.colorScheme.primary,
            ),
            SizedBox(width: isMobile ? 6.w : 6),
            Text(
              loc.browseCategories,
              style: TextStyle(
                fontSize: isMobile ? 18.sp : 18,
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
        centerTitle: true,
      ),
      body: _categories.isEmpty
          ? Center(
              child: Text(
                loc.noRelatedPosts,
                style: TextStyle(
                  fontSize: 15,
                  color: theme.colorScheme.onSurface.withAlpha(120),
                ),
              ),
            )
          : _buildBody(isMobile),
    );
  }

  Widget _buildBody(bool isMobile) {
    if (!isMobile) {
      return Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(
            maxWidth: Breakpoints.webContentMaxWidth,
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
            child: Wrap(
              spacing: 12,
              runSpacing: 12,
              children: _categories.map((cat) {
                return _CategoryCard(
                  category: cat,
                  isMobile: false,
                  onTap: () => _navigateToPosts(cat),
                );
              }).toList(),
            ),
          ),
        ),
      );
    }

    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 16.w, vertical: 12.h),
      child: GridView.builder(
        gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 2,
          mainAxisSpacing: 10.h,
          crossAxisSpacing: 10.w,
          childAspectRatio: 2.2,
        ),
        itemCount: _categories.length,
        itemBuilder: (context, index) {
          return _CategoryCard(
            category: _categories[index],
            isMobile: true,
            onTap: () => _navigateToPosts(_categories[index]),
          );
        },
      ),
    );
  }

  void _navigateToPosts(Category category) {
    if (category.id != null) {
      context.go('/posts?categoryId=${category.id}');
    }
  }
}

class _CategoryCard extends StatefulWidget {
  final Category category;
  final bool isMobile;
  final VoidCallback onTap;

  const _CategoryCard({
    required this.category,
    required this.isMobile,
    required this.onTap,
  });

  @override
  State<_CategoryCard> createState() => _CategoryCardState();
}

class _CategoryCardState extends State<_CategoryCard> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final name = getCategoryName(widget.category);
    final desc = getCategoryDescription(widget.category);

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: Material(
        color: _isHovered
            ? theme.colorScheme.surfaceContainerHighest
            : theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(widget.isMobile ? 14.r : 14),
        child: InkWell(
          onTap: widget.onTap,
          borderRadius: BorderRadius.circular(widget.isMobile ? 14.r : 14),
          child: Container(
            width: widget.isMobile ? null : 200,
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(widget.isMobile ? 14.r : 14),
              border: Border.all(
                color: theme.colorScheme.onSurface.withAlpha((0.08 * 255).round()),
                width: 1,
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
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  if (desc.isNotEmpty) ...[
                    SizedBox(height: widget.isMobile ? 4.h : 4),
                    Text(
                      desc,
                      style: TextStyle(
                        fontSize: widget.isMobile ? 12.sp : 12,
                        color: theme.colorScheme.onSurface.withAlpha(140),
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                  SizedBox(height: widget.isMobile ? 4.h : 4),
                  Text(
                    S.of(context).postsCountFull(widget.category.postCount ?? 0),
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
