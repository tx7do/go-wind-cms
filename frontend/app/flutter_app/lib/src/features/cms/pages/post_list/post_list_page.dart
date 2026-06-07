import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_post_response.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_category_response.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_tag_response.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';
import 'package:flutter_app/src/core/services/pagination_query.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;

/// 文章列表页
///
/// 支持按分类、标签过滤，从 API 服务端获取数据。
class PostListPage extends StatefulWidget {
  /// 可选：按分类 ID 过滤
  final int? categoryId;

  /// 可选：按标签 ID 过滤
  final int? tagId;

  const PostListPage({super.key, this.categoryId, this.tagId});

  @override
  State<PostListPage> createState() => _PostListPageState();
}

class _PostListPageState extends State<PostListPage> {
  final _postService = PostService();
  final _categoryService = CategoryService();
  final _tagService = TagService();

  List<Post> _posts = [];
  List<Category> _categories = [];
  List<ContentServiceV1Tag> _tags = [];
  bool _isLoading = true;

  /// 当前选中的分类（可为 null 表示全部）
  int? _selectedCategoryId;

  @override
  void initState() {
    super.initState();
    _selectedCategoryId = widget.categoryId;
    _loadData();
  }

  Future<void> _loadData() async {
    final filterValues = <String, dynamic>{};
    if (_selectedCategoryId != null) {
      filterValues['category_id'] = _selectedCategoryId;
    }
    if (widget.tagId != null) {
      filterValues['tag_id'] = widget.tagId;
    }

    final results = await Future.wait([
      _postService.list(PaginationQuery(formValues: filterValues)),
      _categoryService.list(),
      _tagService.list(),
    ]);

    if (!mounted) return;

    setState(() {
      _posts = (results[0] as ContentServiceV1ListPostResponse?)?.items ?? [];
      _categories = (results[1] as ContentServiceV1ListCategoryResponse?)?.items ?? [];
      _tags = (results[2] as ContentServiceV1ListTagResponse?)?.items ?? [];
      _isLoading = false;
    });
  }

  String get _pageTitle {
    if (widget.tagId != null) {
      final tag = _tags.where((t) => t.id == widget.tagId);
      if (tag.isNotEmpty) return getTagName(tag.first);
    }
    if (_selectedCategoryId != null) {
      final cat = _categories.where((c) => c.id == _selectedCategoryId);
      if (cat.isNotEmpty) return getCategoryName(cat.first);
    }
    return S.of(context).allPosts;
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
              Icons.article_outlined,
              size: isMobile ? 20.sp : 20,
              color: theme.colorScheme.primary,
            ),
            SizedBox(width: isMobile ? 6.w : 6),
            Flexible(
              child: Text(
                _pageTitle,
                style: TextStyle(
                  fontSize: isMobile ? 18.sp : 18,
                  fontWeight: FontWeight.bold,
                ),
                overflow: TextOverflow.ellipsis,
              ),
            ),
          ],
        ),
        centerTitle: true,
        bottom: PreferredSize(
          preferredSize: Size.fromHeight(isMobile ? 28.h : 28),
          child: Padding(
            padding: EdgeInsets.only(bottom: isMobile ? 10.h : 10),
            child: Text(
              loc.postsCountFull(_posts.length),
              style: TextStyle(
                fontSize: isMobile ? 13.sp : 13,
                color: theme.colorScheme.onSurface.withAlpha(120),
              ),
            ),
          ),
        ),
      ),
      body: _posts.isEmpty
          ? Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    Icons.article_outlined,
                    size: 56,
                    color: theme.colorScheme.onSurface.withAlpha(60),
                  ),
                  const SizedBox(height: 14),
                  Text(
                    loc.noRelatedPosts,
                    style: TextStyle(
                      fontSize: 15,
                      color: theme.colorScheme.onSurface.withAlpha(120),
                    ),
                  ),
                ],
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
          child: ListView.builder(
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 8),
            itemCount: _posts.length,
            itemBuilder: (context, index) => Padding(
              padding: const EdgeInsets.only(bottom: 12),
              child: PostCard(
                post: _posts[index],
                categories: _categories,
                tags: _tags,
              ),
            ),
          ),
        ),
      );
    }

    return ListView.builder(
      padding: EdgeInsets.symmetric(horizontal: 16.w, vertical: 8.h),
      itemCount: _posts.length,
      itemBuilder: (context, index) => Padding(
        padding: EdgeInsets.only(bottom: 12.h),
        child: PostCard(
          post: _posts[index],
          categories: _categories,
          tags: _tags,
        ),
      ),
    );
  }
}
