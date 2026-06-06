import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/data/mock_data.dart';

/// 搜索页
class SearchPage extends StatefulWidget {
  const SearchPage({super.key});

  @override
  State<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends State<SearchPage> {
  final _searchController = TextEditingController();
  String _query = '';
  bool _hasSearched = false;

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
        titleSpacing: 0,
        title: TextField(
          controller: _searchController,
          autofocus: true,
          style: const TextStyle(fontSize: 15),
          decoration: InputDecoration(
            hintText: '搜索文章、标签...',
            hintStyle: TextStyle(
              fontSize: 15,
              color: theme.colorScheme.onSurface.withAlpha(100),
            ),
            border: InputBorder.none,
            isDense: true,
            contentPadding: const EdgeInsets.symmetric(vertical: 10),
          ),
          onSubmitted: (value) {
            setState(() {
              _query = value.trim();
              _hasSearched = true;
            });
          },
        ),
        actions: [
          if (_searchController.text.isNotEmpty)
            IconButton(
              icon: const Icon(Icons.clear, size: 20),
              onPressed: () {
                _searchController.clear();
                setState(() {
                  _query = '';
                  _hasSearched = false;
                });
              },
            ),
          TextButton(
            onPressed: () {
              setState(() {
                _query = _searchController.text.trim();
                _hasSearched = true;
              });
            },
            child: const Text('搜索', style: TextStyle(fontSize: 14)),
          ),
        ],
      ),
      body: _hasSearched
          ? _buildResults(context, isMobile)
          : _buildSuggestions(context, isMobile),
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
                '热门搜索',
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
                children: mockTags.take(6).map((tag) {
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
                '推荐阅读',
                style: TextStyle(
                  fontSize: 15,
                  fontWeight: FontWeight.w600,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              const SizedBox(height: 12),
              ...mockPosts
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
              '没有找到 "$_query" 相关内容',
              style: TextStyle(
                fontSize: 15,
                color: theme.colorScheme.onSurface.withAlpha(120),
              ),
            ),
          ],
        ),
      );
    }

    return ListView(
      padding: EdgeInsets.all(hPad),
      children: [
        Center(
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
                    '相关标签',
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
                      final name = (tag.translations ?? []).isNotEmpty
                          ? (tag.translations ?? []).first.name ?? ''
                          : '';
                      return ActionChip(
                        onPressed: () {},
                        label: Text(
                          '# $name (${tag.postCount})',
                          style: const TextStyle(fontSize: 13),
                        ),
                      );
                    }).toList(),
                  ),
                  const SizedBox(height: 16),
                ],
                if (filteredPosts.isNotEmpty) ...[
                  Text(
                    '相关文章 (${filteredPosts.length})',
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
      ],
    );
  }

  // =================== 搜索逻辑 ===================

  List get _filteredPosts {
    if (_query.isEmpty) return [];
    final q = _query.toLowerCase();
    return mockPosts.where((post) {
      final title = (post.translations ?? []).isNotEmpty
          ? (post.translations ?? []).first.title ?? ''
          : '';
      final summary = (post.translations ?? []).isNotEmpty
          ? (post.translations ?? []).first.summary ?? ''
          : '';
      final tagNames = mockTags
          .where((t) => post.tagIds != null && t.id != null && (post.tagIds as List).contains(t.id!))
          .map(
            (t) => (t.translations ?? []).isNotEmpty ? (t.translations ?? []).first.name ?? '' : '',
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
    return mockTags.where((tag) {
      final name = (tag.translations ?? []).isNotEmpty
          ? (tag.translations ?? []).first.name ?? ''
          : '';
      return name.toLowerCase().contains(q);
    }).toList();
  }
}

/// 简化版文章卡片（搜索结果用）
class _SimplePostCard extends StatefulWidget {
  final dynamic post;

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
                  post.authorName,
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
                  '${post.likes}',
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
