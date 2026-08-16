import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show CommentServiceV1Comment, InteractionServiceV1TargetType,
        InteractionServiceV1CounterMetric;
import 'package:flutter_app/src/features/cms/services/comment_service.dart';
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';

typedef Comment = CommentServiceV1Comment;

/// 我的评论页
///
/// 展示当前用户发表的所有评论，按时间倒序排列。
class MyCommentsPage extends StatefulWidget {
  const MyCommentsPage({super.key});

  @override
  State<MyCommentsPage> createState() => _MyCommentsPageState();
}

class _MyCommentsPageState extends State<MyCommentsPage> {
  final _commentService = CommentService();

  List<Comment> _comments = [];
  bool _isLoading = true;

  // 评论点赞计数缓存：key 为 comment.id。计数来自 interaction_counter 表
  // （comment.like_count 列已移除），由 InteractionService.GetCounts 批量查询。
  Map<int, int> _likeCounts = {};

  final _interactionService = InteractionService();

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final result = await _commentService.list();
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
        _likeCounts = ccounts;
        _isLoading = false;
      });
    } else {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isMobile = ResponsiveUtils.isMobile(context);
    final loc = S.of(context);

    final body = _isLoading
        ? const Center(child: CircularProgressIndicator())
        : _comments.isEmpty
            ? Center(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      Icons.comment_outlined,
                      size: 56,
                      color: theme.colorScheme.onSurface.withAlpha(60),
                    ),
                    const SizedBox(height: 14),
                    Text(
                      loc.noCommentsYet,
                      style: TextStyle(
                        fontSize: isMobile ? 15.sp : 15,
                        color: theme.colorScheme.onSurface.withAlpha(120),
                      ),
                    ),
                  ],
                ),
              )
            : isMobile
                ? RefreshIndicator(
                    onRefresh: _loadData,
                    child: ListView.separated(
                      padding: EdgeInsets.all(16.w),
                      itemCount: _comments.length,
                      separatorBuilder: (_, __) => Divider(
                        height: 24.h,
                        color: theme.colorScheme.onSurface.withAlpha(30),
                      ),
                      itemBuilder: (context, index) {
                        return _CommentCard(
                          comment: _comments[index],
                          isMobile: isMobile,
                          likeCount:
                              _likeCounts[_comments[index].id] ?? 0,
                          onTap: () {
                            final objectId = _comments[index].objectId;
                            if (objectId != null) {
                              context.push('/post/$objectId');
                            }
                          },
                        );
                      },
                    ),
                  )
                : CustomScrollView(
                    slivers: [
                      SliverPadding(
                        padding: const EdgeInsets.symmetric(horizontal: 24),
                        sliver: SliverToBoxAdapter(
                          child: Center(
                            child: ConstrainedBox(
                              constraints: const BoxConstraints(maxWidth: 600),
                              child: Column(
                                children: List.generate(_comments.length, (index) {
                                  return Column(
                                    children: [
                                      _CommentCard(
                                        comment: _comments[index],
                                        isMobile: false,
                                        likeCount:
                                            _likeCounts[_comments[index].id] ?? 0,
                                        onTap: () {
                                          final objectId = _comments[index].objectId;
                                          if (objectId != null) {
                                            context.push('/post/$objectId');
                                          }
                                        },
                                      ),
                                      if (index < _comments.length - 1)
                                        Divider(
                                          height: 24,
                                          color: theme.colorScheme.onSurface.withAlpha(30),
                                        ),
                                    ],
                                  );
                                }),
                              ),
                            ),
                          ),
                        ),
                      ),
                      const SliverToBoxAdapter(child: SizedBox(height: 32)),
                    ],
                  );

    // Web 端由 WebShellLayout 提供 Scaffold
    if (!isMobile) return body;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        leading: const AppBackButton(),
        title: Text(
          loc.myComments,
          style: TextStyle(
            fontSize: 18.sp,
            fontWeight: FontWeight.bold,
          ),
        ),
        centerTitle: true,
      ),
      body: body,
    );
  }
}

/// 评论卡片
class _CommentCard extends StatelessWidget {
  final Comment comment;
  final bool isMobile;
  final int likeCount;
  final VoidCallback? onTap;

  const _CommentCard({
    required this.comment,
    required this.isMobile,
    required this.likeCount,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return GestureDetector(
      onTap: onTap,
      behavior: HitTestBehavior.opaque,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 头部：头像 + 作者名 + 时间
          Row(
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
                    Text(
                      comment.authorName ?? '',
                      style: TextStyle(
                        fontSize: isMobile ? 13.sp : 13,
                        fontWeight: FontWeight.w500,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                    if (comment.createdAt != null)
                      Text(
                        _formatDate(context, DateTime.parse(comment.createdAt!)),
                        style: TextStyle(
                          fontSize: isMobile ? 11.sp : 11,
                          color: theme.colorScheme.onSurface.withAlpha(100),
                        ),
                      ),
                  ],
                ),
              ),
              // 跳转图标
              Icon(
                Icons.open_in_new,
                size: isMobile ? 16.sp : 16,
                color: theme.colorScheme.onSurface.withAlpha(60),
              ),
            ],
          ),
          SizedBox(height: isMobile ? 10.h : 10),

          // 评论内容
          Container(
            width: double.infinity,
            padding: EdgeInsets.all(isMobile ? 12.w : 12),
            decoration: BoxDecoration(
              color: theme.colorScheme.surfaceContainerHighest
                  .withAlpha((0.3 * 255).round()),
              borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
            ),
            child: Text(
              comment.content ?? '',
              style: TextStyle(
                fontSize: isMobile ? 14.sp : 14,
                color: theme.colorScheme.onSurface.withAlpha(200),
                height: 1.6,
              ),
            ),
          ),

          // 底部：互动信息（点赞计数由 interaction_counter 表提供，comment.like_count 列已移除）
          Padding(
            padding: EdgeInsets.only(top: isMobile ? 8.h : 8),
            child: Row(
              children: [
                Icon(
                  Icons.favorite_outline,
                  size: 14,
                  color: theme.colorScheme.onSurface.withAlpha(100),
                ),
                SizedBox(width: isMobile ? 3.w : 3),
                Text(
                  '${likeCount}',
                  style: TextStyle(
                    fontSize: isMobile ? 11.sp : 11,
                    color: theme.colorScheme.onSurface.withAlpha(100),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
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
