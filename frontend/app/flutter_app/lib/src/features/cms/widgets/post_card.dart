import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Category, ContentServiceV1Tag;
import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef Tag = ContentServiceV1Tag;

/// 文章卡片
///
/// 一套观感吃所有宽度：卡片本身不感知断点，列数交给 AdaptiveGrid。
/// 视觉纪律（doc §1/§5）：描边分层不投影、主色只出现在分类签与悬停描边、
/// 计数与时间用等宽数字。
///
/// likeCount 由父列表组件批量查询 InteractionService.GetCounts 后注入
/// （post.likes 列已移除，计数来自 interaction_counter 表）。
class PostCard extends StatefulWidget {
  final Post post;
  final List<Category> categories;
  final List<Tag> tags;
  final int likeCount;

  const PostCard({
    super.key,
    required this.post,
    this.categories = const [],
    this.tags = const [],
    this.likeCount = 0,
  });

  @override
  State<PostCard> createState() => _PostCardState();
}

class _PostCardState extends State<PostCard> {
  bool _isHovered = false;

  String get _title => getPostTitle(widget.post);

  String get _summary => getPostSummary(widget.post);

  String? get _coverImage => getPostThumbnail(widget.post);

  Category? _getCategory() {
    if ((widget.post.categoryIds ?? []).isEmpty) return null;
    final catId = widget.post.categoryIds!.first;
    try {
      return widget.categories.firstWhere((c) => c.id == catId);
    } catch (e) {
      assert(() {
        debugPrint('[PostCard] category lookup failed: $e');
        return true;
      }());
      return null;
    }
  }

  List<Tag> _getTags() {
    return widget.tags
        .where(
          (t) => t.id != null && (widget.post.tagIds ?? []).contains(t.id!),
        )
        .toList();
  }

  DateTime? get _publishedAt {
    final t = widget.post.publishTime;
    return t != null ? DateTime.tryParse(t) : null;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final p = theme.palette;
    final radius = theme.radius;
    final space = theme.space;
    final hovered = _isHovered;

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: AnimatedContainer(
        duration: theme.motion.fast,
        curve: theme.motion.curve,
        transform: hovered
            ? Matrix4.translationValues(0, -theme.motion.lift, 0)
            : Matrix4.identity(),
        decoration: BoxDecoration(
          borderRadius: radius.mdAll,
          boxShadow: hovered ? theme.shadow.hover : theme.shadow.card,
        ),
        child: Card(
          elevation: 0,
          margin: EdgeInsets.zero,
          color: p.surface,
          shape: RoundedRectangleBorder(
            borderRadius: radius.mdAll,
            side: BorderSide(
              color: hovered ? p.primaryBorder : p.border,
              width: theme.borderWidth.thin,
            ),
          ),
          child: InkWell(
            onTap: () {
              if (widget.post.id != null) {
                context.go('/post/${widget.post.id}');
              }
            },
            borderRadius: radius.mdAll,
            child: Padding(
              padding: EdgeInsets.all(space.cardPadding - space.xs),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 无封面也渲染占位块：网格里所有卡视觉对齐
                  _Cover(imageUrl: _coverImage),
                  SizedBox(height: space.md),
                  _buildHeader(context, theme),
                  SizedBox(height: space.sm + 2),
                  _buildContent(context, theme),
                  SizedBox(height: space.md),
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
    final p = theme.palette;
    final category = _getCategory();
    final author = widget.post.authorName ?? '';

    return Row(
      children: [
        CircleAvatar(
          radius: 11,
          backgroundColor: p.primarySoft,
          child: Text(
            author.isNotEmpty ? author.characters.first : '?',
            style: theme.textTheme.labelSmall?.copyWith(
              fontSize: 11,
              fontWeight: FontWeight.w600,
              color: p.primaryOnSoft,
            ),
          ),
        ),
        const SizedBox(width: 6),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                author,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
                style: theme.textTheme.labelMedium,
              ),
              Text(
                _publishedAt != null ? _formatDate(_publishedAt!) : '',
                style: theme.textTheme.labelSmall,
              ),
            ],
          ),
        ),
        if (category != null) ...[
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
            decoration: BoxDecoration(
              color: p.primarySoft,
              borderRadius: theme.radius.xsAll,
            ),
            child: Text(
              getCategoryName(category),
              style: theme.textTheme.labelSmall?.copyWith(
                color: p.primaryOnSoft,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ],
    );
  }

  Widget _buildContent(BuildContext context, ThemeData theme) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: [
        Text(
          _title,
          maxLines: 2,
          overflow: TextOverflow.ellipsis,
          style: theme.textTheme.titleMedium?.copyWith(height: 1.4),
        ),
        const SizedBox(height: 6),
        Text(
          _summary,
          maxLines: 2,
          overflow: TextOverflow.ellipsis,
          style: theme.textTheme.bodySmall?.copyWith(height: 1.6),
        ),
      ],
    );
  }

  Widget _buildFooter(BuildContext context, ThemeData theme) {
    final tags = _getTags();
    final tagText = tags.take(3).map((tag) => '#${getTagName(tag)}').join('  ');

    return Row(
      children: [
        Expanded(
          child: Text(
            tagText,
            style: theme.textTheme.labelSmall,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ),
        _buildStat(Icons.favorite_outline, '${widget.likeCount}', theme),
      ],
    );
  }

  Widget _buildStat(IconData icon, String value, ThemeData theme) {
    final style = theme.textTheme.labelSmall;
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 14, color: style?.color),
        const SizedBox(width: 3),
        Text(value, style: style),
      ],
    );
  }

  String _formatDate(DateTime date) {
    final loc = S.of(context);
    final now = DateTime.now();
    final diff = now.difference(date);
    if (diff.inDays == 0) return loc.today;
    if (diff.inDays == 1) return loc.yesterday;
    if (diff.inDays < 7) return loc.daysAgo(diff.inDays);
    if (diff.inDays < 30) return loc.weeksAgo((diff.inDays / 7).floor());
    return loc.monthDay(date.month, date.day);
  }
}

/// 卡片题图：16:9 稳定宽高比（网格里所有卡对齐）。
/// 无图/加载失败显示占位：主色淡底 + 线性图标，不用灰块
class _Cover extends StatelessWidget {
  const _Cover({this.imageUrl});

  final String? imageUrl;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final p = theme.palette;

    Widget fallback({IconData? icon}) => AspectRatio(
          aspectRatio: 16 / 9,
          child: Container(
            color: p.primaryTint,
            alignment: Alignment.center,
            child: Icon(
              icon ?? Icons.article_outlined,
              size: 28,
              color: p.inkTertiary,
            ),
          ),
        );

    final url = imageUrl;
    if (url == null || url.isEmpty) return fallback();

    return ClipRRect(
      borderRadius: theme.radius.mediaTop(theme.radius.md - 1),
      child: CachedNetworkImage(
        imageUrl: url,
        width: double.infinity,
        fit: BoxFit.cover,
        placeholder: (context, url) => fallback(),
        errorWidget: (context, url, error) =>
            fallback(icon: Icons.broken_image_outlined),
      ),
    );
  }
}
