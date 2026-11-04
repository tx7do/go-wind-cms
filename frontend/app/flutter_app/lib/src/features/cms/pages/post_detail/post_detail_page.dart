import 'package:flutter/material.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show CommentServiceV1Comment, ContentServiceV1Post,
        ContentServiceV1Category, ContentServiceV1Tag,
        InteractionServiceV1TargetType, InteractionServiceV1CounterMetric;
import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';
import 'package:flutter_app/src/features/cms/services/post_service.dart';
import 'package:flutter_app/src/features/cms/services/category_service.dart';
import 'package:flutter_app/src/features/cms/services/tag_service.dart';
import 'package:flutter_app/src/features/cms/services/comment_service.dart';
import 'package:flutter_app/src/core/layout/width_class.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/interaction_bar.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/transport/http/api_error_messages.dart';
import 'package:flutter_app/src/core/transport/http/status.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/comment_tree_utils.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/post_header.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/post_content.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/post_tags.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/comment_section.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/comment_input_bar.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;

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

  // 评论点赞计数缓存：key 为 comment.id。计数来自 interaction_counter 表
  // （comment.like_count 列已移除），由 InteractionService.GetCounts 批量查询。
  Map<int, int> _commentLikeCounts = {};

  final _interactionService = InteractionService();

  @override
  void initState() {
    super.initState();
    _loadData();
    // 浏览量上报（PV +1，best-effort 静默）：每次进入详情记一次
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _interactionService.recordView(
          InteractionServiceV1TargetType.targetTypePost, widget.postId);
    });
  }

  Future<void> _loadData() async {
    final results = await Future.wait([
      _postService.get(widget.postId),
      _categoryService.list(),
      _tagService.list(),
      _commentService.list(
        PaginationQuery(
          skipLocale: true,
          formValues: {'object_id': widget.postId},
        ),
      ),
    ]);

    if (!mounted) return;

    final post = results[0] as Post?;
    final categories = (results[1] as ListCategoryResponse?)?.items ?? [];
    final tags = (results[2] as ListTagResponse?)?.items ?? [];
    final comments = (results[3] as ListCommentResponse?)?.items ?? [];
    final cids = comments
        .where((c) => c.id != null)
        .map((c) => c.id!)
        .toList(growable: false);
    Map<int, int> ccounts = {};
    if (cids.isNotEmpty) {
      final cresp = await _interactionService.getCounts(
          InteractionServiceV1TargetType.targetTypeComment,
          cids,
          [InteractionServiceV1CounterMetric.counterMetricLike]);
      ccounts = InteractionService.extractCountMap(
          cresp,
          cids,
          InteractionServiceV1CounterMetric.counterMetricLike);
    }
    if (!mounted) return;

    setState(() {
      _post = post;
      _categories = categories;
      _tags = tags;
      _comments = comments;
      _commentLikeCounts = ccounts;
      _commentTree = buildCommentTree(comments);
      _isLoading = false;
    });
  }

  /// 发送评论
  Future<void> _sendComment(
    String content, {
    int? parentId,
    int? replyToId,
  }) async {
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
          SnackBar(
            content: Text(localizedApiErrorMessage(S.of(context), result)),
          ),
        );
      }
      return;
    }
    await _refreshComments();
  }

  /// 仅刷新评论
  Future<void> _refreshComments() async {
    final result = await _commentService.list(
      PaginationQuery(
        skipLocale: true,
        formValues: {'object_id': widget.postId},
      ),
    );
    if (!mounted) return;
    if (result is ListCommentResponse) {
      final comments = result.items ?? [];
      final cids = comments
          .where((c) => c.id != null)
          .map((c) => c.id!)
          .toList(growable: false);
      Map<int, int> ccounts = {};
      if (cids.isNotEmpty) {
        final cresp = await _interactionService.getCounts(
            InteractionServiceV1TargetType.targetTypeComment,
            cids,
            [InteractionServiceV1CounterMetric.counterMetricLike]);
        ccounts = InteractionService.extractCountMap(
            cresp,
            cids,
            InteractionServiceV1CounterMetric.counterMetricLike);
      }
      if (!mounted) return;
      setState(() {
        _comments = comments;
        _commentLikeCounts = ccounts;
        _commentTree = buildCommentTree(comments);
      });
    }
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

  /// 评论输入栏在双视图之外渲染：拖拽窗口跨断点时草稿不丢
  Widget _buildView(BuildContext context, {required bool isMobile}) {
    final theme = Theme.of(context);
    final post = _post!;
    final commentCount = _comments.length;

    final body = Column(
      children: [
        Expanded(
          child: isMobile
              ? _buildMobileBody(post, commentCount)
              : _buildWebBody(post, commentCount),
        ),
        CommentInputBar(
          isMobile: isMobile,
          replyTo: _replyTo,
          onReplyChanged: (c) => setState(() => _replyTo = c),
          onSend: _sendComment,
        ),
      ],
    );

    // Web 端：顶部轻量操作栏（返回 + 分享）。
    // 细线铺满全宽（内容块本身限宽，边线不该跟着悬在半空）
    if (!isMobile) {
      final p = theme.palette;
      return Column(
        children: [
          Container(
            width: double.infinity,
            decoration: BoxDecoration(
              color: p.surface,
              border: Border(
                bottom: BorderSide(
                  color: p.hairline,
                  width: theme.borderWidth.hairline,
                ),
              ),
            ),
            child: Center(
              child: ConstrainedBox(
                constraints: const BoxConstraints(
                  maxWidth: kArticleLineWidth + 96,
                ),
                child: Padding(
                  padding: EdgeInsets.symmetric(horizontal: theme.space.xxl),
                  child: Row(
                    children: [
                      AppBackButton(color: p.ink),
                      const Spacer(),
                      IconButton(
                        icon: Icon(
                          Icons.share_outlined,
                          size: 20,
                          color: p.inkSecondary,
                        ),
                        onPressed: () {},
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ),
          Expanded(child: body),
        ],
      );
    }

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(leading: const AppBackButton()),
      body: body,
    );
  }

  // =================== 手机端 ===================

  Widget _buildMobileBody(Post post, int commentCount) {
    return CustomScrollView(
      slivers: [
        SliverToBoxAdapter(
          child: PostHeader(
            post: post,
            isMobile: true,
            categories: _categories,
          ),
        ),
        SliverToBoxAdapter(child: PostContent(post: post, isMobile: true)),
        PostTags(
          allTags: _tags,
          tagIds: post.tagIds,
          isMobile: true,
          asSliver: true,
        ),
        SliverToBoxAdapter(child: InteractionBar(post: post, isMobile: true)),
        CommentSection(
          commentCount: commentCount,
          commentTree: _commentTree,
          allComments: _comments,
          isMobile: true,
          asSliver: true,
          onReply: (c) => setState(() => _replyTo = c),
          commentLikeCounts: _commentLikeCounts,
        ),
        SliverToBoxAdapter(
          child: SizedBox(height: Theme.of(context).space.lg),
        ),
      ],
    );
  }

  // =================== Web 端 ===================

  /// 正文列 ~760dp（中文一行 32–38 字的舒适行长），
  /// 点赞/收藏悬浮右栏——桌面端把行动收口到视线端点，不插在正文和评论中间。
  Widget _buildWebBody(Post post, int commentCount) {
    final theme = Theme.of(context);
    final space = theme.space;

    return CustomScrollView(
      slivers: [
        SliverToBoxAdapter(
          child: Center(
            child: ConstrainedBox(
              constraints: const BoxConstraints(
                maxWidth: kArticleLineWidth + 96,
              ),
              child: Padding(
                padding: EdgeInsets.symmetric(horizontal: space.xxl),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          PostHeader(
                            post: post,
                            isMobile: false,
                            categories: _categories,
                          ),
                          PostContent(post: post, isMobile: false),
                          PostTags(
                            allTags: _tags,
                            tagIds: post.tagIds,
                            isMobile: false,
                          ),
                          CommentSection(
                            commentCount: commentCount,
                            commentTree: _commentTree,
                            allComments: _comments,
                            isMobile: false,
                            onReply: (c) => setState(() => _replyTo = c),
                            commentLikeCounts: _commentLikeCounts,
                          ),
                          SizedBox(height: space.lg),
                        ],
                      ),
                    ),
                    SizedBox(width: space.xxl),
                    _ActionRail(post: post),
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }
}

/// 桌面端右栏：点赞 / 收藏纵排悬浮位（InteractionBar vertical 形态，态逻辑复用）
class _ActionRail extends StatelessWidget {
  const _ActionRail({required this.post});

  final Post post;

  @override
  Widget build(BuildContext context) {
    // 顶端留出呼吸位，不贴着操作条细线
    return Padding(
      padding: EdgeInsets.only(top: Theme.of(context).space.sm),
      child: InteractionBar(post: post, isMobile: false, vertical: true),
    );
  }
}
