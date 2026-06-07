import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

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
import 'package:flutter_app/src/core/widgets/app_back_button.dart';
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/transport/http/status.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/comment_tree_utils.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/post_header.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/post_content.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/interaction_bar.dart';
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
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text(result.getMessage)));
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
      setState(() {
        _comments = result.items ?? [];
        _commentTree = buildCommentTree(_comments);
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

  Widget _buildView(BuildContext context, {required bool isMobile}) {
    final theme = Theme.of(context);
    final post = _post!;
    final commentCount = post.commentCount ?? _comments.length;

    final appBar = AppBar(
      backgroundColor: theme.colorScheme.surface,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
      leading: const AppBackButton(),
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
    );

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

    // Web 端由 WebShellLayout 提供 Scaffold，不再嵌套 Scaffold
    if (!isMobile) return body;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: appBar,
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
        ),
        SliverToBoxAdapter(child: SizedBox(height: 16.h)),
      ],
    );
  }

  // =================== Web 端 ===================

  Widget _buildWebBody(Post post, int commentCount) {
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
                    InteractionBar(post: post, isMobile: false),
                    CommentSection(
                      commentCount: commentCount,
                      commentTree: _commentTree,
                      allComments: _comments,
                      isMobile: false,
                      onReply: (c) => setState(() => _replyTo = c),
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
}
