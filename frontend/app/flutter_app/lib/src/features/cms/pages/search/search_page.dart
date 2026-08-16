import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Tag, ContentServiceV1SearchPostHit,
        InteractionServiceV1TargetType, InteractionServiceV1CounterMetric;
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';

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

  // suggestions 数据（热门搜索 + 推荐阅读）仍用 list 接口加载
  List<ContentServiceV1Post> _posts = [];
  List<ContentServiceV1Tag> _tags = [];

  // 点赞计数缓存：key 为 post.id。计数来自 interaction_counter 表
  // （post.likes 列已移除），由 InteractionService.GetCounts 批量查询。
  Map<int, int> _likeCounts = {};

  final _interactionService = InteractionService();

  // 真实搜索结果（来自 OpenSearch）
  List<ContentServiceV1SearchPostHit> _searchHits = [];
  bool _isSearching = false;

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
      _tags = (results[1] as ListTagResponse?)?.items ?? [];
      _isLoading = false;
    });
  }

  /// 执行真实全文搜索（调用 OpenSearch）
  Future<void> _runSearch(String query) async {
    final q = query.trim();
    if (q.isEmpty) return;
    setState(() {
      _query = q;
      _hasSearched = true;
      _isSearching = true;
    });
    final resp = await _postService.search(q);
    if (!mounted) return;
    setState(() {
      _searchHits = resp?.items ?? [];
      _isSearching = false;
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
        onSubmitted: (value) => _runSearch(value.trim()),
      ),
      actions: [
        if (_searchController.text.isNotEmpty)
          IconButton(icon: const Icon(Icons.clear, size: 20), onPressed: () {
            _searchController.clear();
            setState(() { _query = ''; _hasSearched = false; _searchHits = []; });
          }),
        TextButton(
          onPressed: () => _runSearch(_searchController.text.trim()),
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
                      _runSearch(name);
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
                      child: _SimplePostCard(
                          post: post,
                          likeCount: _likeCounts[post.id] ?? 0),
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

    if (_isSearching) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_searchHits.isEmpty) {
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
                    Text(
                      S.of(context).relatedPostsCount(_searchHits.length),
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w600,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                    const SizedBox(height: 8),
                    ..._searchHits.map(
                      (hit) => Padding(
                        padding: const EdgeInsets.only(bottom: 12),
                        child: _SearchHitCard(hit: hit),
                      ),
                    ),
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
  // （真实搜索在 _runSearch 中调用 PostService.search，命中结果存入 _searchHits）
}

/// 搜索结果卡片（基于 OpenSearch 命中，仅含 title + postId）
class _SearchHitCard extends StatefulWidget {
  final ContentServiceV1SearchPostHit hit;

  const _SearchHitCard({required this.hit});

  @override
  State<_SearchHitCard> createState() => _SearchHitCardState();
}

class _SearchHitCardState extends State<_SearchHitCard> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final hit = widget.hit;
    final title = hit.title ?? '';
    final isMobile = ResponsiveUtils.isMobile(context);

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: GestureDetector(
        onTap: () {
          final id = hit.postId;
          if (id != null) context.push('/post/$id');
        },
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
          child: Row(
            children: [
              Expanded(
                child: Text(
                  title,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w600,
                    color: theme.colorScheme.onSurface,
                  ),
                ),
              ),
              const SizedBox(width: 8),
              Icon(
                Icons.chevron_right,
                size: 20,
                color: theme.colorScheme.onSurface.withAlpha(100),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

/// 简化版文章卡片（推荐阅读用，基于完整 Post 数据）
class _SimplePostCard extends StatefulWidget {
  final ContentServiceV1Post post;
  final int likeCount;

  const _SimplePostCard({required this.post, required this.likeCount});

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
                  '${widget.likeCount}',
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
