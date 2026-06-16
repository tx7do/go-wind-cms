import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Tag;
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/features/cms/widgets/post_card.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

typedef Post = ContentServiceV1Post;

/// 标签文章列表页
class TagFeedPage extends StatefulWidget {
  final int tagId;

  const TagFeedPage({super.key, required this.tagId});

  @override
  State<TagFeedPage> createState() => _TagFeedPageState();
}

class _TagFeedPageState extends State<TagFeedPage> {
  final _postService = PostService();
  final _tagService = TagService();

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
      _postService.list(),
      _tagService.list(),
    ]);

    if (!mounted) return;

    setState(() {
      _posts = (results[0] as ListPostResponse?)?.items ?? [];
      _tags = (results[1] as ListTagResponse?)?.items ?? [];
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
    final tag = _tags.firstWhere(
      (t) => t.id != null && t.id == widget.tagId,
      orElse: () => ContentServiceV1Tag(),
    );
    final tagName = getTagName(tag);
    final posts = _posts
        .where((p) => (p.tagIds ?? []).contains(tag.id))
        .toList();

    final appBar = AppBar(
      backgroundColor: theme.colorScheme.surface,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
      leading: const AppBackButton(),
      title: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.tag, size: isMobile ? 18.sp : 18, color: theme.colorScheme.primary),
          SizedBox(width: isMobile ? 4.w : 4),
          Text(tagName, style: TextStyle(fontSize: isMobile ? 18.sp : 18, fontWeight: FontWeight.bold)),
        ],
      ),
      centerTitle: true,
      bottom: PreferredSize(
        preferredSize: Size.fromHeight(isMobile ? 32.h : 32),
        child: Padding(
          padding: EdgeInsets.only(bottom: isMobile ? 12.h : 12),
          child: Text(
            S.of(context).relatedPostsCountFull(tag.postCount ?? 0),
            style: TextStyle(fontSize: isMobile ? 13.sp : 13, color: theme.colorScheme.onSurface.withAlpha(120)),
          ),
        ),
      ),
    );

    final body = posts.isEmpty
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
        : _buildBody(context, posts, isMobile);

    // Web 端由 WebShellLayout 提供 Scaffold
    if (!isMobile) return body;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: appBar,
      body: body,
    );
  }

  Widget _buildBody(BuildContext context, List<Post> posts, bool isMobile) {
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
                    children: posts.map((post) => Padding(
                      padding: const EdgeInsets.only(bottom: 12),
                      child: PostCard(post: post),
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
      itemCount: posts.length,
      itemBuilder: (context, index) => Padding(
        padding: EdgeInsets.only(bottom: 12.h),
        child: PostCard(post: posts[index]),
      ),
    );
  }
}
