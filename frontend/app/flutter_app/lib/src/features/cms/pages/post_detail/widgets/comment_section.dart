import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show CommentServiceV1Comment;
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/comment_tree_utils.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/comment_item.dart';

/// 评论区组件
///
/// 支持两种渲染模式：
/// - [asSliver] = true → 返回 SliverMainAxisGroup（用于 CustomScrollView）
/// - [asSliver] = false → 返回普通 Column Widget（用于 Web 端 Column）
class CommentSection extends StatelessWidget {
  final int commentCount;
  final List<CommentType> commentTree;
  final List<CommentServiceV1Comment> allComments;
  final bool isMobile;
  final bool asSliver;
  final void Function(CommentServiceV1Comment comment)? onReply;

  /// 评论点赞计数映射：key 为 comment.id。计数来自 interaction_counter 表
  /// （comment.like_count 列已移除），由父页批量查询 GetCounts(COMMENT, LIKE)。
  final Map<int, int> commentLikeCounts;

  const CommentSection({
    super.key,
    required this.commentCount,
    required this.commentTree,
    required this.allComments,
    required this.isMobile,
    this.asSliver = false,
    this.onReply,
    required this.commentLikeCounts,
  });

  /// 评论区标题栏（蓝色竖条 + "评论 (N)"）
  Widget _buildHeader(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
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
            S.of(context).commentsCount(commentCount),
            style: TextStyle(
              fontSize: isMobile ? 16.sp : 16,
              fontWeight: FontWeight.w600,
              color: theme.colorScheme.onSurface,
            ),
          ),
        ],
      ),
    );
  }

  /// 递归展开评论树
  Iterable<Widget> _expandTree(CommentType comment, {int depth = 0}) sync* {
    yield CommentItem(
      comment: comment,
      isMobile: isMobile,
      depth: depth,
      allComments: allComments,
      onReply: onReply,
      likeCount: commentLikeCounts[comment.id] ?? 0,
    );
    final children = findChildren(comment, allComments);
    for (final child in children) {
      yield* _expandTree(child, depth: depth + 1);
    }
  }

  /// 所有评论展开后的 Widget 列表
  List<Widget> get _allCommentWidgets {
    return commentTree.expand((c) => _expandTree(c)).toList();
  }

  // =================== Sliver 模式 ===================

  Widget _buildSliver(BuildContext context) {
    return SliverMainAxisGroup(
      slivers: [
        SliverToBoxAdapter(child: _buildHeader(context)),
        SliverPadding(
          padding: EdgeInsets.symmetric(horizontal: isMobile ? 20.w : 0),
          sliver: SliverList(
            delegate: SliverChildBuilderDelegate((context, index) {
              final widgets = _expandTree(commentTree[index]).toList();
              return Column(children: widgets);
            }, childCount: commentTree.length),
          ),
        ),
      ],
    );
  }

  // =================== 非 Sliver 模式 ===================

  Widget _buildColumn(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      mainAxisSize: MainAxisSize.min,
      children: [_buildHeader(context), ..._allCommentWidgets],
    );
  }

  @override
  Widget build(BuildContext context) {
    return asSliver ? _buildSliver(context) : _buildColumn(context);
  }
}
