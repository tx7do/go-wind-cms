import 'package:flutter/material.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Tag;
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/widgets/tag_chip.dart';

/// 搜索页
class SearchPage extends StatefulWidget {
  const SearchPage({super.key});

  @override
  State<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends State<SearchPage> {
  final _searchController = TextEditingController();
  final _postService = PostService();
  final _tagService = TagService();

  String _query = '';
  bool _hasSearched = false;
  bool _isLoading = true;

  List<ContentServiceV1Post> _posts = [];
  List<ContentServiceV1Tag> _tags = [];

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
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ResponsiveLayout(
      mobileBody: _buildView(isMobile: true),
      webBody: _buildView(isMobile: false),
    );
  }

  Widget _buildView({required bool isMobile}) {
    final theme = Theme.of(context);

    final appBar = AppBar(
      backgroundColor: theme.colorScheme.surface,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
      leading: const AppBackButton(),
      titleSpacing: 0,
      title: TextField(
        controller: _searchController,
        autofocus: true,
        style: const TextStyle(fontSize: 15),
        decoration: InputDecoration(
          hintText: S.of(context).searchHint,
          hintStyle: TextStyle(fontSize: 15, color: theme.colorScheme.onSurface.withAlpha(100)),
          border: InputBorder.none,
          isDense: true,
          contentPadding: const EdgeInsets.symmetric(vertical: 10),
        ),
        onSubmitted: (value) {
          setState(() { _query = value.trim(); _hasSearched = true; });
        },
      ),
      actions: [
        if (_searchController.text.isNotEmpty)
          IconButton(icon: const Icon(Icons.clear, size: 20), onPressed: () {
            _searchController.clear();
            setState(() { _query = ''; _hasSearched = false; });
          }),
        TextButton(
          onPressed: () { setState(() { _query = _searchController.text.trim(); _hasSearched = true; }); },
          child: Text(S.of(context).search, style: const TextStyle(fontSize: 14)),
        ),
      ],
    );

    final body = _hasSearched
        ? _buildResults(context, isMobile)
        : _buildSuggestions(context, isMobile);

    if (_isLoading) {
      final loadingBody = const Center(child: CircularProgressIndicator());
      // Web 端不嵌套 Scaffold
      if (!isMobile) return loadingBody;
      return Scaffold(backgroundColor: theme.scaffoldBackgroundColor, appBar: appBar, body: loadingBody);
    }

    // Web 端由 WebShellLayout 提供 Scaffold
    if (!isMobile) return body;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: appBar,
      body: body,
    );
  }

  Widget _buildSuggestions(BuildContext context, bool isMobile) {
    final theme = Theme.of(context);
    final hPad = isMobile ? 16.0 : 24.0;

    return SingleChildScrollView(
      padding: EdgeInsets.all(hPad),
      child: Center(
        child: ConstrainedBox(
          constraints: BoxConstraints(
            maxWidth: isMobile
                ? double.infinity
                : Breakpoints.webContentMaxWidth,
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                S.of(context).hotSearch,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w600,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              const SizedBox(height: 12),
              Wrap(
                spacing: 8,
                runSpacing: 8,
                children: _tags.take(6).map((tag) {
                  final name = (tag.translations ?? []).isNotEmpty
                      ? (tag.translations ?? []).first.name ?? ''
                      : '';
                  return ActionChip(
                    onPressed: () {
                      _searchController.text = name;
                      setState(() {
                        _query = name;
                        _hasSearched = true;
                      });
                    },
                    label: Text(name, style: const TextStyle(fontSize: 13)),
                  );
                }).toList(),
              ),
              const SizedBox(height: 24),
              Text(
                S.of(context).recommendedReading,
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w600,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              const SizedBox(height: 12),
              ..._posts
                  .take(3)
                  .map(
                    (post) => Padding(
                      padding: const EdgeInsets.only(bottom: 12),
                      child: _SimplePostCard(post: post),
                    ),
                  ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildResults(BuildContext context, bool isMobile) {
    final theme = Theme.of(context);
    final hPad = isMobile ? 16.0 : 24.0;
    final filteredPosts = _filteredPosts;
    final filteredTags = _filteredTags;

    if (filteredPosts.isEmpty && filteredTags.isEmpty) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.search_off,
              size: 56,
              color: theme.colorScheme.onSurface.withAlpha(60),
            ),
            const SizedBox(height: 14),
            Text(
              S.of(context).noSearchResults(_query),
              style: TextStyle(
                fontSize: 15,
                color: theme.colorScheme.onSurface.withAlpha(120),
              ),
            ),
          ],
        ),
      );
    }

    return CustomScrollView(
      slivers: [
        SliverPadding(
          padding: EdgeInsets.all(hPad),
          sliver: SliverToBoxAdapter(
            child: Center(
              child: ConstrainedBox(
                constraints: BoxConstraints(
                  maxWidth: isMobile
                      ? double.infinity
                      : Breakpoints.webContentMaxWidth,
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (filteredTags.isNotEmpty) ...[
                      Text(
                        S.of(context).relatedTags,
                        style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w600,
                          color: theme.colorScheme.onSurface,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Wrap(
                        spacing: 8,
                        runSpacing: 8,
                        children: filteredTags.map((tag) {
                          return TagChip(
                            tag: tag,
                            isMobile: isMobile,
                            showPostCount: true,
                          );
                        }).toList(),
                      ),
                      const SizedBox(height: 16),
                    ],
                    if (filteredPosts.isNotEmpty) ...[
                      Text(
                        S.of(context).relatedPostsCount(filteredPosts.length),
                        style: TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w600,
                          color: theme.colorScheme.onSurface,
                        ),
                      ),
                      const SizedBox(height: 8),
                      ...filteredPosts.map(
                        (post) => Padding(
                          padding: const EdgeInsets.only(bottom: 12),
                          child: _SimplePostCard(post: post),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  // =================== 搜索逻辑 ===================

  List get _filteredPosts {
    if (_query.isEmpty) return [];
    final q = _query.toLowerCase();
    return _posts.where((post) {
      final title = (post.translations ?? []).isNotEmpty
          ? (post.translations ?? []).first.title ?? ''
          : '';
      final summary = (post.translations ?? []).isNotEmpty
          ? (post.translations ?? []).first.summary ?? ''
          : '';
      final tagNames = _tags
          .where(
            (t) =>
                post.tagIds != null &&
                t.id != null &&
                (post.tagIds as List).contains(t.id!),
          )
          .map(
            (t) => (t.translations ?? []).isNotEmpty
                ? (t.translations ?? []).first.name ?? ''
                : '',
          )
          .toList();
      return title.toLowerCase().contains(q) ||
          summary.toLowerCase().contains(q) ||
          tagNames.any((n) => n.toLowerCase().contains(q));
    }).toList();
  }

  List get _filteredTags {
    if (_query.isEmpty) return [];
    final q = _query.toLowerCase();
    return _tags.where((tag) {
      final name = (tag.translations ?? []).isNotEmpty
          ? (tag.translations ?? []).first.name ?? ''
          : '';
      return name.toLowerCase().contains(q);
    }).toList();
  }
}

/// 简化版文章卡片（搜索结果用）
class _SimplePostCard extends StatefulWidget {
  final ContentServiceV1Post post;

  const _SimplePostCard({required this.post});

  @override
  State<_SimplePostCard> createState() => _SimplePostCardState();
}

class _SimplePostCardState extends State<_SimplePostCard> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final post = widget.post;
    final title = (post.translations ?? []).isNotEmpty
        ? (post.translations ?? []).first.title ?? ''
        : '';
    final summary = (post.translations ?? []).isNotEmpty
        ? (post.translations ?? []).first.summary ?? ''
        : '';
    final isMobile = ResponsiveUtils.isMobile(context);

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: _isHovered && !isMobile
              ? theme.colorScheme.surfaceContainerLow
              : theme.colorScheme.surface,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: _isHovered && !isMobile
                ? theme.colorScheme.primary.withAlpha((0.2 * 255).round())
                : theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(
                fontSize: 15,
                fontWeight: FontWeight.w600,
                color: theme.colorScheme.onSurface,
              ),
            ),
            const SizedBox(height: 6),
            Text(
              summary,
              maxLines: 2,
              overflow: TextOverflow.ellipsis,
              style: TextStyle(
                fontSize: 13,
                color: theme.colorScheme.onSurface.withAlpha(160),
                height: 1.4,
              ),
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Text(
                  post.authorName ?? '',
                  style: TextStyle(
                    fontSize: 12,
                    color: theme.colorScheme.onSurface.withAlpha(120),
                  ),
                ),
                const Spacer(),
                Icon(
                  Icons.favorite_outline,
                  size: 14,
                  color: theme.colorScheme.onSurface.withAlpha(100),
                ),
                const SizedBox(width: 3),
                Text(
                  '${post.likes ?? 0}',
                  style: TextStyle(
                    fontSize: 12,
                    color: theme.colorScheme.onSurface.withAlpha(100),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
