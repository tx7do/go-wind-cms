import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Tag = ContentServiceV1Tag;

/// 文章卡片组件
///
/// 支持响应式布局和 Web 端 Hover 效果
class PostCard extends StatefulWidget {
  final Post post;
  final List<Category> categories;
  final List<Tag> tags;

  const PostCard({
    super.key,
    required this.post,
    this.categories = const [],
    this.tags = const [],
  });

  @override
  State<PostCard> createState() => _PostCardState();
}

class _PostCardState extends State<PostCard> {
  bool _isHovered = false;

  String get _title => (widget.post.translations ?? []).isNotEmpty
      ? widget.post.translations!.first.title ?? ''
      : '';

  String get _summary => (widget.post.translations ?? []).isNotEmpty
      ? widget.post.translations!.first.summary ?? ''
      : '';

  String? get _coverImage =>
      (widget.post.translations ?? []).isNotEmpty &&
          widget.post.translations!.first.thumbnail != null &&
          widget.post.translations!.first.thumbnail!.isNotEmpty
      ? widget.post.translations!.first.thumbnail
      : null;

  Category? _getCategory() {
    if ((widget.post.categoryIds ?? []).isEmpty) return null;
    final catId = widget.post.categoryIds!.first;
    try {
      return widget.categories.firstWhere((c) => c.id == catId);
    } catch (_) {
      return null;
    }
  }

  List<Tag> _getTags() {
    return widget.tags
        .where((t) =>
            t.id != null && (widget.post.tagIds ?? []).contains(t.id!))
        .toList();
  }

  DateTime? get _publishedAt {
    return widget.post.publishTime;
  }

  bool get _isMobile => ResponsiveUtils.isMobile(context);

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        curve: Curves.easeOut,
        transform: _isHovered && !_isMobile
            ? (Matrix4.identity()..translate(0.0, -2.0))
            : Matrix4.identity(),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(_isMobile ? 14.r : 14),
          boxShadow: _isHovered && !_isMobile
              ? [
                  BoxShadow(
                    color: Colors.black.withAlpha((0.08 * 255).round()),
                    blurRadius: 12,
                    offset: const Offset(0, 4),
                  ),
                ]
              : [],
        ),
        child: Card(
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(_isMobile ? 14.r : 14),
            side: BorderSide(
              color: _isHovered && !_isMobile
                  ? theme.colorScheme.primary.withAlpha((0.3 * 255).round())
                  : theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
              width: 1,
            ),
          ),
          color: _isHovered && !_isMobile
              ? theme.colorScheme.surfaceContainerLow
              : theme.colorScheme.surface,
          margin: EdgeInsets.zero,
          child: InkWell(
            onTap: () {
              // TODO: Navigate to post detail page
            },
            borderRadius: BorderRadius.circular(_isMobile ? 14.r : 14),
            child: Padding(
              padding: EdgeInsets.all(_isMobile ? 16.w : 16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _buildHeader(context, theme),
                  SizedBox(height: _isMobile ? 10.h : 10),
                  _buildContent(context, theme),
                  SizedBox(height: _isMobile ? 12.h : 12),
                  _buildFooter(context, theme),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildHeader(BuildContext context, ThemeData theme) {
    final category = _getCategory();

    return Row(
      children: [
        CircleAvatar(
          radius: _isMobile ? 16.r : 16,
          backgroundColor: theme.colorScheme.primaryContainer,
          child: Text(
            (widget.post.authorName ?? '').isNotEmpty ? (widget.post.authorName ?? '')[0] : '?',
            style: TextStyle(
              fontSize: _isMobile ? 13.sp : 13,
              fontWeight: FontWeight.w600,
              color: theme.colorScheme.onSurface,
            ),
          ),
        ),
        SizedBox(width: _isMobile ? 10.w : 10),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                widget.post.authorName ?? '',
                style: TextStyle(
                  fontSize: _isMobile ? 13.sp : 13,
                  fontWeight: FontWeight.w500,
                  color: theme.colorScheme.onSurface,
                ),
              ),
              SizedBox(height: _isMobile ? 2.h : 2),
              Text(
                _publishedAt != null ? _formatDate(_publishedAt!) : '',
                style: TextStyle(
                  fontSize: _isMobile ? 11.sp : 11,
                  color: theme.colorScheme.onSurface.withAlpha(120),
                ),
              ),
            ],
          ),
        ),
        if (category != null)
          Container(
            padding: EdgeInsets.symmetric(
              horizontal: _isMobile ? 10.w : 10,
              vertical: _isMobile ? 4.h : 4,
            ),
            decoration: BoxDecoration(
              color: theme.colorScheme.primaryContainer.withAlpha(
                (0.5 * 255).round(),
              ),
              borderRadius: BorderRadius.circular(_isMobile ? 8.r : 8),
            ),
            child: Text(
              (category.translations ?? []).isNotEmpty
                  ? category.translations!.first.name ?? ''
                  : '',
              style: TextStyle(
                fontSize: _isMobile ? 11.sp : 11,
                fontWeight: FontWeight.w500,
                color: theme.colorScheme.onSurface.withAlpha(180),
              ),
            ),
          ),
      ],
    );
  }

  Widget _buildContent(BuildContext context, ThemeData theme) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          _title,
          maxLines: 2,
          overflow: TextOverflow.ellipsis,
          style: TextStyle(
            fontSize: _isMobile ? 16.sp : 16,
            fontWeight: FontWeight.bold,
            color: theme.colorScheme.onSurface,
            height: 1.4,
          ),
        ),
        SizedBox(height: _isMobile ? 6.h : 6),
        Text(
          _summary,
          maxLines: 2,
          overflow: TextOverflow.ellipsis,
          style: TextStyle(
            fontSize: _isMobile ? 13.sp : 13,
            color: theme.colorScheme.onSurface.withAlpha(160),
            height: 1.5,
          ),
        ),
      ],
    );
  }

  Widget _buildFooter(BuildContext context, ThemeData theme) {
    final tags = _getTags();

    final tagText = tags
        .take(3)
        .map((tag) =>
            '#${(tag.translations ?? []).isNotEmpty ? (tag.translations ?? []).first.name ?? '' : ''}')
        .join('  ');

    return Row(
      children: [
        Expanded(
          child: Text(
            tagText,
            style: TextStyle(
              fontSize: _isMobile ? 12.sp : 12,
              color: theme.colorScheme.primary,
              fontWeight: FontWeight.w400,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ),
        _buildStat(
          Icons.comment_outlined,
          '${widget.post.commentCount ?? 0}',
          theme,
        ),
        SizedBox(width: _isMobile ? 14.w : 14),
        _buildStat(Icons.favorite_outline, '${widget.post.likes ?? 0}', theme),
      ],
    );
  }

  Widget _buildStat(IconData icon, String value, ThemeData theme) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(
          icon,
          size: _isMobile ? 16.sp : 16,
          color: theme.colorScheme.onSurface.withAlpha(120),
        ),
        SizedBox(width: _isMobile ? 3.w : 3),
        Text(
          value,
          style: TextStyle(
            fontSize: _isMobile ? 12.sp : 12,
            color: theme.colorScheme.onSurface.withAlpha(120),
          ),
        ),
      ],
    );
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);
    if (diff.inDays == 0) return '今天';
    if (diff.inDays == 1) return '昨天';
    if (diff.inDays < 7) return '${diff.inDays}天前';
    if (diff.inDays < 30) return '${(diff.inDays / 7).floor()}周前';
    return '${date.month}月${date.day}日';
  }
}
