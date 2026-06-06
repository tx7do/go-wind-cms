import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/comment_service_v1_comment.dart';
import 'package:flutter_app/src/features/cms/data/mock_data.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';

/// 文章详情页
class PostDetailPage extends StatelessWidget {
  final int postId;

  const PostDetailPage({super.key, required this.postId});

  @override
  Widget build(BuildContext context) {
    return ResponsiveLayout(
      mobileBody: _buildView(context, isMobile: true),
      webBody: _buildView(context, isMobile: false),
    );
  }

  Widget _buildView(BuildContext context, {required bool isMobile}) {
    final theme = Theme.of(context);
    final post = mockPosts.firstWhere(
      (p) => p.id != null && p.id == postId,
      orElse: () => mockPosts.first,
    );
    final comments = mockComments.where((c) => c.objectId != null && c.objectId == post.id).toList();

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
        actions: [
          IconButton(
            icon: const Icon(Icons.share_outlined, size: 22),
            onPressed: () {},
          ),
          IconButton(
            icon: const Icon(Icons.bookmark_border, size: 22),
            onPressed: () {},
          ),
          const SizedBox(width: 4),
        ],
      ),
      body: Column(
        children: [
          Expanded(
            child: isMobile
                ? _buildMobileBody(context, post, comments)
                : _buildWebBody(context, post, comments),
          ),
          _CommentInputBar(isMobile: isMobile),
        ],
      ),
    );
  }

  // =================== 手机端 ===================

  Widget _buildMobileBody(
    BuildContext context,
    dynamic post,
    List<CommentServiceV1Comment> comments,
  ) {
    return CustomScrollView(
      slivers: [
        SliverToBoxAdapter(child: _PostHeader(post: post, isMobile: true)),
        SliverToBoxAdapter(child: _PostContent(post: post, isMobile: true)),
        _buildTagsSliver(post, isMobile: true),
        SliverToBoxAdapter(child: _InteractionBar(post: post, isMobile: true)),
        _buildCommentsSliver(context, comments, isMobile: true),
        SliverToBoxAdapter(child: SizedBox(height: 16.h)),
      ],
    );
  }

  // =================== Web 端 ===================

  Widget _buildWebBody(
    BuildContext context,
    dynamic post,
    List<CommentServiceV1Comment> comments,
  ) {
    final theme = Theme.of(context);

    return CustomScrollView(
      slivers: [
        SliverToBoxAdapter(
          child: Center(
            child: ConstrainedBox(
              constraints: const BoxConstraints(
                maxWidth: Breakpoints.webContentMaxWidth,
              ),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _PostHeader(post: post, isMobile: false),
                    _PostContent(post: post, isMobile: false),
                    _buildTagsWidget(post, isMobile: false),
                    _InteractionBar(post: post, isMobile: false),
                    // 评论区标题
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 16, 0, 8),
                      child: Row(
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
                            '评论 (${comments.length})',
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                              color: theme.colorScheme.onSurface,
                            ),
                          ),
                        ],
                      ),
                    ),
                    // 评论列表
                    ...comments.map(
                      (c) => _CommentItem(comment: c, isMobile: false),
                    ),
                    const SizedBox(height: 16),
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  // =================== 标签 ===================

  Widget _buildTagsSliver(dynamic post, {required bool isMobile}) {
    final tags = mockTags.where((t) => post.tagIds != null && t.id != null && post.tagIds!.contains(t.id!)).toList();
    if (tags.isEmpty) return const SliverToBoxAdapter(child: SizedBox.shrink());

    return SliverToBoxAdapter(
      child: Padding(
        padding: EdgeInsets.fromLTRB(
          isMobile ? 20.w : 0,
          8,
          isMobile ? 20.w : 0,
          0,
        ),
        child: Wrap(
          spacing: isMobile ? 8.w : 8,
          runSpacing: isMobile ? 6.h : 6,
          children: tags.map((tag) {
            return ActionChip(
              onPressed: () {},
              label: Text(
                '# ${(tag.translations ?? []).isNotEmpty ? (tag.translations ?? []).first.name ?? "" : ""}',
                style: TextStyle(fontSize: isMobile ? 12.sp : 12),
              ),
              padding: EdgeInsets.symmetric(horizontal: isMobile ? 4.w : 4),
              visualDensity: VisualDensity.compact,
            );
          }).toList(),
        ),
      ),
    );
  }

  Widget _buildTagsWidget(dynamic post, {required bool isMobile}) {
    final tags = mockTags.where((t) => post.tagIds != null && t.id != null && post.tagIds!.contains(t.id!)).toList();
    if (tags.isEmpty) return const SizedBox.shrink();

    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 8, 0, 0),
      child: Wrap(
        spacing: 8,
        runSpacing: 6,
        children: tags.map((tag) {
          return ActionChip(
            onPressed: () {},
            label: Text(
              '# ${(tag.translations ?? []).isNotEmpty ? (tag.translations ?? []).first.name ?? "" : ""}',
              style: const TextStyle(fontSize: 12),
            ),
            padding: const EdgeInsets.symmetric(horizontal: 4),
            visualDensity: VisualDensity.compact,
          );
        }).toList(),
      ),
    );
  }

  // =================== 评论区（手机端 Sliver） ===================

  Widget _buildCommentsSliver(
    BuildContext context,
    List<Comment> comments, {
    required bool isMobile,
  }) {
    final theme = Theme.of(context);
    return SliverMainAxisGroup(
      slivers: [
        SliverToBoxAdapter(
          child: Padding(
            padding: EdgeInsets.fromLTRB(
              isMobile ? 20.w : 0,
              16,
              isMobile ? 20.w : 0,
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
                  '评论 (${comments.length})',
                  style: TextStyle(
                    fontSize: isMobile ? 16.sp : 16,
                    fontWeight: FontWeight.w600,
                    color: theme.colorScheme.onSurface,
                  ),
                ),
              ],
            ),
          ),
        ),
        SliverPadding(
          padding: EdgeInsets.symmetric(horizontal: isMobile ? 20.w : 0),
          sliver: SliverList(
            delegate: SliverChildBuilderDelegate(
              (context, index) =>
                  _CommentItem(comment: comments[index], isMobile: isMobile),
              childCount: comments.length,
            ),
          ),
        ),
      ],
    );
  }
}

// =================== 子组件 ===================

class _PostHeader extends StatelessWidget {
  final dynamic post;
  final bool isMobile;

  const _PostHeader({required this.post, required this.isMobile});

  String get _title =>
      post.translations?.isNotEmpty == true ? post.translations!.first.title ?? '' : '';

  String get _categoryName {
    if ((post.categoryIds ?? []).isEmpty) return '';
    final catId = post.categoryIds!.first;
    try {
      final cat = mockCategories.firstWhere((c) => c.id != null && c.id == catId);
      return (cat.translations ?? []).isNotEmpty ? cat.translations!.first.name ?? '' : '';
    } catch (_) {
      return '';
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: EdgeInsets.fromLTRB(
        isMobile ? 20.w : 0,
        8,
        isMobile ? 20.w : 0,
        0,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            _title,
            style: TextStyle(
              fontSize: isMobile ? 22.sp : 24,
              fontWeight: FontWeight.bold,
              color: theme.colorScheme.onSurface,
              height: 1.4,
            ),
          ),
          SizedBox(height: isMobile ? 14.h : 14),
          Row(
            children: [
              CircleAvatar(
                radius: isMobile ? 18.r : 18,
                backgroundColor: theme.colorScheme.primaryContainer,
                child: Text(
                  post.authorName?.isNotEmpty == true ? post.authorName![0] : '?',
                  style: TextStyle(
                    fontSize: isMobile ? 14.sp : 14,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
              SizedBox(width: isMobile ? 10.w : 10),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      post.authorName ?? '',
                      style: TextStyle(
                        fontSize: isMobile ? 14.sp : 14,
                        fontWeight: FontWeight.w500,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                    SizedBox(height: isMobile ? 2.h : 2),
                    Text(
                      post.publishTime != null
                          ? _formatDate(post.publishTime!)
                          : '',
                      style: TextStyle(
                        fontSize: isMobile ? 12.sp : 12,
                        color: theme.colorScheme.onSurface.withAlpha(120),
                      ),
                    ),
                  ],
                ),
              ),
              if (_categoryName.isNotEmpty)
                Container(
                  padding: EdgeInsets.symmetric(
                    horizontal: isMobile ? 12.w : 12,
                    vertical: isMobile ? 5.h : 5,
                  ),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.primaryContainer.withAlpha(
                      (0.5 * 255).round(),
                    ),
                    borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
                  ),
                  child: Text(
                    _categoryName,
                    style: TextStyle(fontSize: isMobile ? 12.sp : 12),
                  ),
                ),
            ],
          ),
          Divider(
            height: isMobile ? 24.h : 24,
            thickness: 1,
            color: theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) =>
      '${date.year}年${date.month}月${date.day}日';
}

class _PostContent extends StatelessWidget {
  final dynamic post;
  final bool isMobile;

  const _PostContent({required this.post, required this.isMobile});

  String get _summary =>
      post.translations?.isNotEmpty == true ? post.translations!.first.summary ?? '' : '';

  String get _content =>
      post.translations?.isNotEmpty == true ? post.translations!.first.content ?? '' : '';

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: EdgeInsets.fromLTRB(
        isMobile ? 20.w : 0,
        0,
        isMobile ? 20.w : 0,
        8,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 摘要
          Container(
            padding: EdgeInsets.all(isMobile ? 14.w : 14),
            decoration: BoxDecoration(
              color: theme.colorScheme.primaryContainer.withAlpha(
                (0.2 * 255).round(),
              ),
              borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
              border: Border(
                left: BorderSide(color: theme.colorScheme.primary, width: 3),
              ),
            ),
            child: Text(
              _summary,
              style: TextStyle(
                fontSize: isMobile ? 14.sp : 14,
                color: theme.colorScheme.onSurface.withAlpha(180),
                height: 1.6,
                fontStyle: FontStyle.italic,
              ),
            ),
          ),
          SizedBox(height: isMobile ? 16.h : 16),
          // 正文
          Text(
            _content,
            style: TextStyle(
              fontSize: isMobile ? 15.sp : 15,
              color: theme.colorScheme.onSurface,
              height: 1.8,
            ),
          ),
        ],
      ),
    );
  }
}

class _InteractionBar extends StatelessWidget {
  final dynamic post;
  final bool isMobile;

  const _InteractionBar({required this.post, required this.isMobile});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: EdgeInsets.symmetric(
        horizontal: isMobile ? 20.w : 0,
        vertical: 12,
      ),
      child: Container(
        padding: EdgeInsets.symmetric(
          vertical: isMobile ? 12.h : 12,
          horizontal: isMobile ? 16.w : 16,
        ),
        decoration: BoxDecoration(
          color: theme.colorScheme.surface,
          borderRadius: BorderRadius.circular(isMobile ? 14.r : 14),
          border: Border.all(
            color: theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
          ),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            _InteractionItem(
              Icons.remove_red_eye_outlined,
              '${post.visits}',
              '浏览',
            ),
            _InteractionItem(Icons.favorite_outline, '${post.likes}', '点赞'),
            _InteractionItem(
              Icons.comment_outlined,
              '${post.commentCount}',
              '评论',
            ),
            _InteractionItem(Icons.share_outlined, '分享', ''),
          ],
        ),
      ),
    );
  }
}

class _InteractionItem extends StatefulWidget {
  final IconData icon;
  final String value;
  final String label;

  const _InteractionItem(this.icon, this.value, this.label);

  @override
  State<_InteractionItem> createState() => _InteractionItemState();
}

class _InteractionItemState extends State<_InteractionItem> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: InkWell(
        onTap: () {},
        borderRadius: BorderRadius.circular(8),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          child: Column(
            children: [
              Icon(
                widget.icon,
                size: 20,
                color: _isHovered
                    ? theme.colorScheme.primary
                    : theme.colorScheme.onSurface.withAlpha(160),
              ),
              const SizedBox(height: 4),
              Text(
                widget.value,
                style: TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w500,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              if (widget.label.isNotEmpty) ...[
                const SizedBox(height: 2),
                Text(
                  widget.label,
                  style: TextStyle(
                    fontSize: 10,
                    color: theme.colorScheme.onSurface.withAlpha(100),
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

class _CommentItem extends StatelessWidget {
  final CommentServiceV1Comment comment;
  final bool isMobile;

  const _CommentItem({required this.comment, required this.isMobile});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isReply = comment.parentId != null && comment.parentId != 0;

    return Padding(
      padding: EdgeInsets.only(
        left: isReply ? (isMobile ? 40.w : 40) : 0,
        bottom: isMobile ? 14.h : 14,
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CircleAvatar(
            radius: isMobile ? 16.r : 16,
            backgroundColor: theme.colorScheme.primaryContainer,
            child: Text(
              (comment.authorName ?? '').isNotEmpty ? comment.authorName![0] : '?',
              style: TextStyle(fontSize: isMobile ? 12.sp : 12),
            ),
          ),
          SizedBox(width: isMobile ? 10.w : 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Text(
                      comment.authorName ?? '',
                      style: TextStyle(
                        fontSize: isMobile ? 13.sp : 13,
                        fontWeight: FontWeight.w500,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                    const Spacer(),
                    Text(
                      comment.createdAt != null
                          ? _formatDate(comment.createdAt!)
                          : '',
                      style: TextStyle(
                        fontSize: isMobile ? 11.sp : 11,
                        color: theme.colorScheme.onSurface.withAlpha(100),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: isMobile ? 4.h : 4),
                Text(
                  comment.content ?? '',
                  style: TextStyle(
                    fontSize: isMobile ? 14.sp : 14,
                    color: theme.colorScheme.onSurface.withAlpha(200),
                    height: 1.5,
                  ),
                ),
                SizedBox(height: isMobile ? 6.h : 6),
                Row(
                  children: [
                    Icon(
                      Icons.favorite_outline,
                      size: 14,
                      color: theme.colorScheme.onSurface.withAlpha(100),
                    ),
                    SizedBox(width: isMobile ? 3.w : 3),
                    Text(
                      '${comment.likeCount}',
                      style: TextStyle(
                        fontSize: isMobile ? 11.sp : 11,
                        color: theme.colorScheme.onSurface.withAlpha(100),
                      ),
                    ),
                    SizedBox(width: isMobile ? 16.w : 16),
                    Text(
                      '回复',
                      style: TextStyle(
                        fontSize: isMobile ? 11.sp : 11,
                        color: theme.colorScheme.primary,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);
    if (diff.inDays == 0) return '今天';
    if (diff.inDays == 1) return '昨天';
    if (diff.inDays < 7) return '${diff.inDays}天前';
    return '${date.month}月${date.day}日';
  }
}

class _CommentInputBar extends StatelessWidget {
  final bool isMobile;

  const _CommentInputBar({required this.isMobile});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
      padding: EdgeInsets.fromLTRB(
        isMobile ? 16.w : 16,
        10,
        isMobile ? 16.w : 16,
        10,
      ),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(
          top: BorderSide(
            color: theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
          ),
        ),
      ),
      child: SafeArea(
        top: false,
        child: Center(
          child: ConstrainedBox(
            constraints: const BoxConstraints(
              maxWidth: Breakpoints.webContentMaxWidth,
            ),
            child: Row(
              children: [
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 14,
                      vertical: 8,
                    ),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.surfaceContainerHighest
                          .withAlpha((0.5 * 255).round()),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Text(
                      '写下你的评论...',
                      style: TextStyle(
                        fontSize: isMobile ? 14.sp : 14,
                        color: theme.colorScheme.onSurface.withAlpha(100),
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 10),
                IconButton(
                  icon: Icon(
                    Icons.send,
                    size: 22,
                    color: theme.colorScheme.primary,
                  ),
                  onPressed: () {},
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
