import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_comment.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/features/cms/services/comment_service.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/transport/http/status.dart';
import 'package:flutter_app/src/features/cms/widgets/content_viewer.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;
typedef CommentType = CommentServiceV1Comment;

/// 将扁平评论列表构建为树形结构
///
/// 如果后端已返回 [children] 字段，直接返回顶层评论；
/// 否则根据 [parentId] 归组。
List<CommentType> buildCommentTree(List<CommentType> flatList) {
  // 检查是否已经有树形数据（第一条评论有 children 字段非空即视为树形）
  if (flatList.isNotEmpty && flatList.any((c) => c.children != null && c.children!.isNotEmpty)) {
    return flatList.where((c) => c.parentId == null || c.parentId == 0).toList();
  }

  // 扁平列表 → 按 parentId 归组
  final Map<int, List<CommentType>> childrenMap = {};
  final List<CommentType> roots = [];

  for (final comment in flatList) {
    final pid = comment.parentId;
    if (pid == null || pid == 0) {
      roots.add(comment);
    } else {
      childrenMap.putIfAbsent(pid, () => []).add(comment);
    }
  }

  // 递归构建（使用可变副本携带 children）
  List<CommentType> attachChildren(List<CommentType> parents) {
    return parents.map((parent) {
      final kids = childrenMap[parent.id ?? -1] ?? [];
      if (kids.isEmpty) return parent;
      // CommentType 是 const 构造函数，无法修改 children，
      // 直接返回原对象，渲染时通过 childrenMap 查找子节点
      return parent;
    }).toList();
  }

  attachChildren(roots);
  return roots;
}

/// 递归获取某评论的所有子评论（从扁平列表中查找）
List<CommentType> findChildren(
  CommentType comment,
  List<CommentType> flatList,
) {
  // 优先使用后端返回的 children
  if (comment.children != null && comment.children!.isNotEmpty) {
    return comment.children!;
  }
  // 否则从扁平列表中按 parentId 查找
  return flatList.where((c) => c.parentId == comment.id).toList();
}

/// 文章详情页
class PostDetailPage extends StatefulWidget {
  final int postId;

  const PostDetailPage({super.key, required this.postId});

  @override
  State<PostDetailPage> createState() => _PostDetailPageState();
}

class _PostDetailPageState extends State<PostDetailPage> {
  final _postService = PostService();
  final _categoryService = CategoryService();
  final _tagService = TagService();
  final _commentService = CommentService();

  Post? _post;
  List<Category> _categories = [];
  List<ContentServiceV1Tag> _tags = [];
  List<CommentServiceV1Comment> _comments = [];
  List<CommentServiceV1Comment> _commentTree = [];
  CommentServiceV1Comment? _replyTo;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final results = await Future.wait([
      _postService.get(widget.postId),
      _categoryService.list(),
      _tagService.list(),
      _commentService.list(PaginationQuery(
        skipLocale: true,
        formValues: {'object_id': widget.postId},
      )),
    ]);

    if (!mounted) return;

    final post = results[0] as Post?;
    final categories = (results[1] as ListCategoryResponse?)?.items ?? [];
    final tags = (results[2] as ListTagResponse?)?.items ?? [];
    final comments = (results[3] as ListCommentResponse?)?.items ?? [];

    setState(() {
      _post = post;
      _categories = categories;
      _tags = tags;
      _comments = comments;
      _commentTree = buildCommentTree(comments);
      _isLoading = false;
    });
  }

  /// 发送评论
  Future<void> _sendComment(String content, {int? parentId, int? replyToId}) async {
    final newComment = CommentServiceV1Comment(
      objectId: widget.postId,
      content: content,
      parentId: parentId,
      replyToId: replyToId,
    );
    final result = await _commentService.create(newComment);
    if (result is Status) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(result.getMessage)),
        );
      }
      return;
    }
    // 发送成功，刷新评论列表
    await _refreshComments();
  }

  /// 仅刷新评论
  Future<void> _refreshComments() async {
    final result = await _commentService.list(PaginationQuery(
      skipLocale: true,
      formValues: {'object_id': widget.postId},
    ));
    if (!mounted) return;
    if (result is ListCommentResponse) {
      setState(() {
        _comments = result.items ?? [];
        _commentTree = buildCommentTree(_comments);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    return ResponsiveLayout(
      mobileBody: _buildView(context, isMobile: true),
      webBody: _buildView(context, isMobile: false),
    );
  }

  Widget _buildView(BuildContext context, {required bool isMobile}) {
    final theme = Theme.of(context);
    final post = _post!;

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
                ? _buildMobileBody(context, post, _commentTree, _comments)
                : _buildWebBody(context, post, _commentTree, _comments),
          ),
          _CommentInputBar(
            isMobile: isMobile,
            replyTo: _replyTo,
            onReplyChanged: (c) => setState(() => _replyTo = c),
            onSend: _sendComment,
          ),
        ],
      ),
    );
  }

  // =================== 手机端 ===================

  Widget _buildMobileBody(
    BuildContext context,
    Post post,
    List<CommentServiceV1Comment> commentTree,
    List<CommentServiceV1Comment> allComments,
  ) {
    return CustomScrollView(
      slivers: [
        SliverToBoxAdapter(child: _PostHeader(post: post, isMobile: true, categories: _categories)),
        SliverToBoxAdapter(child: _PostContent(post: post, isMobile: true)),
        _buildTagsSliver(post, isMobile: true),
        SliverToBoxAdapter(child: _InteractionBar(post: post, isMobile: true)),
        _buildCommentsSliver(context, post, commentTree, allComments, isMobile: true),
        SliverToBoxAdapter(child: SizedBox(height: 16.h)),
      ],
    );
  }

  // =================== Web 端 ===================

  Widget _buildWebBody(
    BuildContext context,
    Post post,
    List<CommentServiceV1Comment> commentTree,
    List<CommentServiceV1Comment> allComments,
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
                    _PostHeader(post: post, isMobile: false, categories: _categories),
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
                            S.of(context).commentsCount(post.commentCount ?? allComments.length),
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w600,
                              color: theme.colorScheme.onSurface,
                            ),
                          ),
                        ],
                      ),
                    ),
                    // 评论列表（树形）
                    ...commentTree.expand(
                      (c) => _buildCommentWithChildren(c, allComments, isMobile: false),
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

  Widget _buildTagsSliver(Post post, {required bool isMobile}) {
    final tags = _tags.where((t) => post.tagIds != null && t.id != null && post.tagIds!.contains(t.id!)).toList();
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

  Widget _buildTagsWidget(Post post, {required bool isMobile}) {
    final tags = _tags.where((t) => post.tagIds != null && t.id != null && post.tagIds!.contains(t.id!)).toList();
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
    Post post,
    List<CommentType> commentTree,
    List<CommentType> allComments, {
    required bool isMobile,
  }) {
    final theme = Theme.of(context);
    final count = post.commentCount ?? allComments.length;
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
                  S.of(context).commentsCount(count),
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
              (context, index) {
                final widgets = _buildCommentWithChildren(
                  commentTree[index],
                  allComments,
                  isMobile: isMobile,
                ).toList();
                return Column(children: widgets);
              },
              childCount: commentTree.length,
            ),
          ),
        ),
      ],
    );
  }

  // =================== 树形评论递归构建 ===================

  /// 递归展开一条评论及其所有子评论为 Widget 列表
  Iterable<Widget> _buildCommentWithChildren(
    CommentType comment,
    List<CommentType> allComments, {
    required bool isMobile,
    int depth = 0,
  }) sync* {
    yield _CommentItem(
      comment: comment,
      isMobile: isMobile,
      depth: depth,
      allComments: allComments,
      onReply: (c) => setState(() => _replyTo = c),
    );
    final children = findChildren(comment, allComments);
    for (final child in children) {
      yield* _buildCommentWithChildren(
        child,
        allComments,
        isMobile: isMobile,
        depth: depth + 1,
      );
    }
  }
}

// =================== 子组件 ===================

class _PostHeader extends StatelessWidget {
  final Post post;
  final bool isMobile;
  final List<Category> categories;

  const _PostHeader({required this.post, required this.isMobile, required this.categories});

  String get _title => getPostTitle(post);

  String get _categoryName {
    if ((post.categoryIds ?? []).isEmpty) return '';
    final catId = post.categoryIds!.first;
    try {
      final cat = categories.firstWhere((c) => c.id != null && c.id == catId);
      return getCategoryName(cat);
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
                          ? _formatDate(context, post.publishTime!)
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

  String _formatDate(BuildContext context, DateTime date) {
    final loc = S.of(context);
    return loc.yearMonthDay(date.year, date.month, date.day);
  }
}

class _PostContent extends StatelessWidget {
  final Post post;
  final bool isMobile;

  const _PostContent({required this.post, required this.isMobile});

  String get _summary => getPostSummary(post);

  String get _content => getPostContent(post);

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
          // 正文 — 使用 ContentViewer 根据 editorType 自适应渲染
          ContentViewer(
            content: _content,
            editorType: post.editorType,
            isMobile: isMobile,
          ),
        ],
      ),
    );
  }
}

class _InteractionBar extends StatelessWidget {
  final Post post;
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
              '${post.visits ?? 0}',
              S.of(context).views,
            ),
            _InteractionItem(Icons.favorite_outline, '${post.likes ?? 0}', S.of(context).likes),
            _InteractionItem(
              Icons.comment_outlined,
              '${post.commentCount ?? 0}',
              S.of(context).comments,
            ),
            _InteractionItem(Icons.share_outlined, S.of(context).share, ''),
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
  final int depth;
  final List<CommentServiceV1Comment> allComments;
  final void Function(CommentServiceV1Comment comment)? onReply;

  const _CommentItem({
    required this.comment,
    required this.isMobile,
    this.depth = 0,
    required this.allComments,
    this.onReply,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final indent = (depth * (isMobile ? 36.w : 36)).toDouble();
    final hasChildren = findChildren(comment, allComments).isNotEmpty;

    return Padding(
      padding: EdgeInsets.only(
        left: indent,
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
                    Flexible(
                      child: Text(
                        comment.authorName ?? '',
                        style: TextStyle(
                          fontSize: isMobile ? 13.sp : 13,
                          fontWeight: FontWeight.w500,
                          color: theme.colorScheme.onSurface,
                        ),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      comment.createdAt != null
                          ? _formatDate(context, comment.createdAt!)
                          : '',
                      style: TextStyle(
                        fontSize: isMobile ? 11.sp : 11,
                        color: theme.colorScheme.onSurface.withAlpha(100),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: isMobile ? 4.h : 4),
                // 如果是回复（depth > 0），显示回复的 "@某人"
                if (depth > 0 && comment.replyToId != null)
                  Padding(
                    padding: const EdgeInsets.only(bottom: 2),
                    child: Text(
                      _replyToText,
                      style: TextStyle(
                        fontSize: isMobile ? 12.sp : 12,
                        color: theme.colorScheme.primary.withAlpha(180),
                      ),
                    ),
                  ),
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
                      '${comment.likeCount ?? 0}',
                      style: TextStyle(
                        fontSize: isMobile ? 11.sp : 11,
                        color: theme.colorScheme.onSurface.withAlpha(100),
                      ),
                    ),
                    SizedBox(width: isMobile ? 16.w : 16),
                    GestureDetector(
                      onTap: () => onReply?.call(comment),
                      child: Text(
                        S.of(context).reply,
                        style: TextStyle(
                          fontSize: isMobile ? 11.sp : 11,
                          color: theme.colorScheme.primary,
                        ),
                      ),
                    ),
                    // 如果有子评论，显示回复数
                    if (hasChildren) ...[
                      SizedBox(width: isMobile ? 16.w : 16),
                      Icon(
                        Icons.chat_bubble_outline,
                        size: 12,
                        color: theme.colorScheme.onSurface.withAlpha(80),
                      ),
                      SizedBox(width: isMobile ? 3.w : 3),
                      Text(
                        '${comment.replyCount ?? findChildren(comment, allComments).length}',
                        style: TextStyle(
                          fontSize: isMobile ? 11.sp : 11,
                          color: theme.colorScheme.onSurface.withAlpha(80),
                        ),
                      ),
                    ],
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  /// 获取回复目标的作者名
  String get _replyToText {
    if (comment.replyToId == null) return '';
    final target = allComments.where((c) => c.id == comment.replyToId);
    if (target.isEmpty) return '';
    final name = target.first.authorName ?? '';
    if (name.isEmpty) return '';
    return '@$name';
  }

  String _formatDate(BuildContext context, DateTime date) {
    final loc = S.of(context);
    final now = DateTime.now();
    final diff = now.difference(date);
    if (diff.inDays == 0) return loc.today;
    if (diff.inDays == 1) return loc.yesterday;
    if (diff.inDays < 7) return loc.daysAgo(diff.inDays);
    return loc.monthDay(date.month, date.day);
  }
}

class _CommentInputBar extends StatefulWidget {
  final bool isMobile;
  final CommentServiceV1Comment? replyTo;
  final void Function(CommentServiceV1Comment? comment)? onReplyChanged;
  final Future<void> Function(String content, {int? parentId, int? replyToId}) onSend;

  const _CommentInputBar({
    required this.isMobile,
    this.replyTo,
    this.onReplyChanged,
    required this.onSend,
  });

  @override
  State<_CommentInputBar> createState() => _CommentInputBarState();
}

class _CommentInputBarState extends State<_CommentInputBar> {
  final _controller = TextEditingController();
  final _focusNode = FocusNode();
  bool _isSending = false;

  @override
  void initState() {
    super.initState();
    // 当回复目标变化时，自动聚焦输入框
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (widget.replyTo != null) {
        _focusNode.requestFocus();
      }
    });
  }

  @override
  void didUpdateWidget(covariant _CommentInputBar oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.replyTo != oldWidget.replyTo && widget.replyTo != null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _focusNode.requestFocus();
      });
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  Future<void> _handleSend() async {
    final text = _controller.text.trim();
    if (text.isEmpty || _isSending) return;

    setState(() => _isSending = true);
    try {
      await widget.onSend(
        text,
        parentId: widget.replyTo?.id,
        replyToId: widget.replyTo?.id,
      );
      _controller.clear();
      widget.onReplyChanged?.call(null);
    } finally {
      if (mounted) setState(() => _isSending = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final replyTo = widget.replyTo;

    return Container(
      padding: EdgeInsets.fromLTRB(
        widget.isMobile ? 16.w : 16,
        10,
        widget.isMobile ? 16.w : 16,
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
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // 回复提示条
                if (replyTo != null)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                    margin: const EdgeInsets.only(bottom: 8),
                    decoration: BoxDecoration(
                      color: theme.colorScheme.primaryContainer.withAlpha(80),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Row(
                      children: [
                        Icon(Icons.reply, size: 16, color: theme.colorScheme.primary),
                        const SizedBox(width: 4),
                        Expanded(
                          child: Text(
                            '${S.of(context).reply} ${replyTo.authorName ?? ''}',
                            style: TextStyle(
                              fontSize: 12,
                              color: theme.colorScheme.primary,
                            ),
                          ),
                        ),
                        GestureDetector(
                          onTap: () => widget.onReplyChanged?.call(null),
                          child: Icon(
                            Icons.close,
                            size: 16,
                            color: theme.colorScheme.onSurface.withAlpha(120),
                          ),
                        ),
                      ],
                    ),
                  ),
                Row(
                  children: [
                    Expanded(
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 14),
                        decoration: BoxDecoration(
                          color: theme.colorScheme.surfaceContainerHighest
                              .withAlpha((0.5 * 255).round()),
                          borderRadius: BorderRadius.circular(20),
                        ),
                        child: TextField(
                          controller: _controller,
                          focusNode: _focusNode,
                          enabled: !_isSending,
                          style: TextStyle(
                            fontSize: widget.isMobile ? 14.sp : 14,
                            color: theme.colorScheme.onSurface,
                          ),
                          decoration: InputDecoration(
                            border: InputBorder.none,
                            hintText: S.of(context).writeComment,
                            hintStyle: TextStyle(
                              fontSize: widget.isMobile ? 14.sp : 14,
                              color: theme.colorScheme.onSurface.withAlpha(100),
                            ),
                            isDense: true,
                            contentPadding: const EdgeInsets.symmetric(vertical: 8),
                          ),
                          textInputAction: TextInputAction.send,
                          onSubmitted: (_) => _handleSend(),
                        ),
                      ),
                    ),
                    const SizedBox(width: 10),
                    IconButton(
                      icon: _isSending
                          ? SizedBox(
                              width: 18,
                              height: 18,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                color: theme.colorScheme.primary,
                              ),
                            )
                          : Icon(
                              Icons.send,
                              size: 22,
                              color: theme.colorScheme.primary,
                            ),
                      onPressed: _isSending ? null : _handleSend,
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
