import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show CommentServiceV1Comment;
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/comment_tree_utils.dart';
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show InteractionServiceV1TargetType, InteractionServiceV1LikeResponse;

/// 评论项。点赞由 interaction 子系统提供：initState 取初始态，
/// 点击调 like/unlike，后端在同事务内更新 ledger 与 comment.like_count 缓存并回读，
/// 前端用返回值更新本地显示。
class CommentItem extends StatefulWidget {
  final CommentServiceV1Comment comment;
  final bool isMobile;
  final int depth;
  final List<CommentServiceV1Comment> allComments;
  final void Function(CommentServiceV1Comment comment)? onReply;

  const CommentItem({
    super.key,
    required this.comment,
    required this.isMobile,
    this.depth = 0,
    required this.allComments,
    this.onReply,
  });

  @override
  State<CommentItem> createState() => _CommentItemState();
}

class _CommentItemState extends State<CommentItem> {
  final _interactionService = InteractionService();
  bool _isLiked = false;
  int _likeCount = 0;

  @override
  void initState() {
    super.initState();
    _likeCount = widget.comment.likeCount ?? 0;
    _loadStatus();
  }

  Future<void> _loadStatus() async {
    final cid = widget.comment.id;
    if (cid == null) return;
    final resp = await _interactionService.getStatus(
        InteractionServiceV1TargetType.targetTypeComment, [cid]);
    if (!mounted || resp == null) return;
    final st = resp.status?[cid.toString()];
    if (st != null) {
      setState(() {
        _isLiked = st.liked ?? false;
      });
    }
  }

  Future<void> _toggleLike() async {
    final cid = widget.comment.id;
    if (cid == null) return;
    final InteractionServiceV1LikeResponse? resp;
    if (_isLiked) {
      resp = await _interactionService.unlike(
          InteractionServiceV1TargetType.targetTypeComment, cid);
    } else {
      resp = await _interactionService.like(
          InteractionServiceV1TargetType.targetTypeComment, cid);
    }
    if (!mounted || resp == null) return;
    setState(() {
      _isLiked = resp.liked ?? _isLiked;
      _likeCount = resp.likeCount ?? _likeCount;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final indent = (widget.depth * (widget.isMobile ? 36.w : 36)).toDouble();
    final realChildren = findChildren(widget.comment, widget.allComments);

    return Padding(
      padding: EdgeInsets.only(left: indent, bottom: widget.isMobile ? 14.h : 14),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CircleAvatar(
            radius: widget.isMobile ? 16.r : 16,
            backgroundColor: theme.colorScheme.primaryContainer,
            child: Text(
              (widget.comment.authorName ?? '').isNotEmpty
                  ? widget.comment.authorName![0]
                  : '?',
              style: TextStyle(fontSize: widget.isMobile ? 12.sp : 12),
            ),
          ),
          SizedBox(width: widget.isMobile ? 10.w : 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Flexible(
                      child: Text(
                        widget.comment.authorName ?? '',
                        style: TextStyle(
                          fontSize: widget.isMobile ? 13.sp : 13,
                          fontWeight: FontWeight.w500,
                          color: theme.colorScheme.onSurface,
                        ),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      widget.comment.createdAt != null
                          ? _formatDate(context, DateTime.parse(widget.comment.createdAt!))
                          : '',
                      style: TextStyle(
                        fontSize: widget.isMobile ? 11.sp : 11,
                        color: theme.colorScheme.onSurface.withAlpha(100),
                      ),
                    ),
                  ],
                ),
                SizedBox(height: widget.isMobile ? 4.h : 4),
                // 如果是回复（depth > 0），显示回复的 "@某人"
                if (widget.depth > 0 && widget.comment.replyToId != null)
                  Padding(
                    padding: const EdgeInsets.only(bottom: 2),
                    child: Text(
                      _replyToText,
                      style: TextStyle(
                        fontSize: widget.isMobile ? 12.sp : 12,
                        color: theme.colorScheme.primary.withAlpha(180),
                      ),
                    ),
                  ),
                Text(
                  widget.comment.content ?? '',
                  style: TextStyle(
                    fontSize: widget.isMobile ? 14.sp : 14,
                    color: theme.colorScheme.onSurface.withAlpha(200),
                    height: 1.5,
                  ),
                ),
                SizedBox(height: widget.isMobile ? 6.h : 6),
                Row(
                  children: [
                    GestureDetector(
                      onTap: _toggleLike,
                      child: Icon(
                        _isLiked ? Icons.favorite : Icons.favorite_outline,
                        size: 14,
                        color: _isLiked
                            ? theme.colorScheme.primary
                            : theme.colorScheme.onSurface.withAlpha(100),
                      ),
                    ),
                    SizedBox(width: widget.isMobile ? 3.w : 3),
                    Text(
                      '$_likeCount',
                      style: TextStyle(
                        fontSize: widget.isMobile ? 11.sp : 11,
                        color: _isLiked
                            ? theme.colorScheme.primary
                            : theme.colorScheme.onSurface.withAlpha(100),
                      ),
                    ),
                    SizedBox(width: widget.isMobile ? 16.w : 16),
                    GestureDetector(
                      onTap: () => widget.onReply?.call(widget.comment),
                      child: Text(
                        S.of(context).reply,
                        style: TextStyle(
                          fontSize: widget.isMobile ? 11.sp : 11,
                          color: theme.colorScheme.primary,
                        ),
                      ),
                    ),
                    // 如果有子评论，显示回复数
                    if (realChildren.isNotEmpty) ...[
                      SizedBox(width: widget.isMobile ? 16.w : 16),
                      Icon(
                        Icons.chat_bubble_outline,
                        size: 12,
                        color: theme.colorScheme.onSurface.withAlpha(80),
                      ),
                      SizedBox(width: widget.isMobile ? 3.w : 3),
                      Text(
                        '${widget.comment.replyCount ?? realChildren.length}',
                        style: TextStyle(
                          fontSize: widget.isMobile ? 11.sp : 11,
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
    if (widget.comment.replyToId == null) return '';
    final target = widget.allComments.where((c) => c.id == widget.comment.replyToId);
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
