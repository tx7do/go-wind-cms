import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_comment.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/widgets/comment_tree_utils.dart';

class CommentItem extends StatelessWidget {
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
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final indent = (depth * (isMobile ? 36.w : 36)).toDouble();
    final realChildren = findChildren(comment, allComments);

    return Padding(
      padding: EdgeInsets.only(left: indent, bottom: isMobile ? 14.h : 14),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          CircleAvatar(
            radius: isMobile ? 16.r : 16,
            backgroundColor: theme.colorScheme.primaryContainer,
            child: Text(
              (comment.authorName ?? '').isNotEmpty
                  ? comment.authorName![0]
                  : '?',
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
                    if (realChildren.isNotEmpty) ...[
                      SizedBox(width: isMobile ? 16.w : 16),
                      Icon(
                        Icons.chat_bubble_outline,
                        size: 12,
                        color: theme.colorScheme.onSurface.withAlpha(80),
                      ),
                      SizedBox(width: isMobile ? 3.w : 3),
                      Text(
                        '${comment.replyCount ?? realChildren.length}',
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
