import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post;
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show InteractionServiceV1TargetType, InteractionServiceV1LikeResponse,
        InteractionServiceV1WatchResponse, InteractionServiceV1CounterMetric,
        InteractionServiceV1GetInteractionStatusResponse,
        InteractionServiceV1GetCountsResponse;

typedef Post = ContentServiceV1Post;

/// 文章底部交互条
///
/// likes 与 bookmark 为可交互按钮，状态由 InteractionService（点赞/收藏计数内聚子系统）提供：
///   - 进入时调 getInteractionStatus 取初始 {liked, watched} 态
///   - 点赞计数：由 interaction_counter 表提供（post.likes 列已移除），进入时调
///     getCounts 取初始计数；点击 likes 按钮后由 LikeResponse.likeCount 乐观更新
///   - 收藏计数：由 interaction_counter 表提供（metric=COUNTER_METRIC_WATCH），进入时调
///     getCounts 取初始计数；点击 bookmark 按钮后由 WatchResponse.watchCount 乐观更新
///   - 点击 bookmark 按钮：按当前态调 watch/unwatch
/// viewer 身份由后端鉴权上下文确定，前端无需传 user_id。未登录时后端返回 401，
/// 按钮保持未激活态。
class InteractionBar extends StatefulWidget {
  final Post post;
  final bool isMobile;

  const InteractionBar({super.key, required this.post, required this.isMobile});

  @override
  State<InteractionBar> createState() => _InteractionBarState();
}

class _InteractionBarState extends State<InteractionBar> {
  final _interactionService = InteractionService();
  bool _liked = false;
  bool _watched = false;
  int _likeCount = 0;
  int _watchCount = 0;

  @override
  void initState() {
    super.initState();
    _loadStatus();
  }

  Future<void> _loadStatus() async {
    final pid = widget.post.id;
    if (pid == null) return;
    // 并行取交互状态（{liked, watched}）与点赞/收藏计数（来自 interaction_counter 表，
    // post 表的散落计数列已移除）。
    final results = await Future.wait([
      _interactionService.getStatus(
          InteractionServiceV1TargetType.targetTypePost, [pid]),
      _interactionService.getCounts(
          InteractionServiceV1TargetType.targetTypePost,
          [pid],
          [
            InteractionServiceV1CounterMetric.counterMetricLike,
            InteractionServiceV1CounterMetric.counterMetricWatch
          ]),
    ]);
    if (!mounted) return;
    final statusResp =
        results[0] as InteractionServiceV1GetInteractionStatusResponse?;
    final countsResp = results[1] as InteractionServiceV1GetCountsResponse?;
    final st = statusResp?.status?[pid.toString()];
    final likeCount = InteractionService.extractCount(
        countsResp, pid, InteractionServiceV1CounterMetric.counterMetricLike);
    final watchCount = InteractionService.extractCount(countsResp, pid,
        InteractionServiceV1CounterMetric.counterMetricWatch);
    setState(() {
      if (st != null) {
        _liked = st.liked ?? _liked;
        _watched = st.watched ?? _watched;
      }
      _likeCount = likeCount;
      _watchCount = watchCount;
    });
  }

  Future<void> _toggleLike() async {
    final pid = widget.post.id;
    if (pid == null) return;
    final InteractionServiceV1LikeResponse? resp;
    if (_liked) {
      resp = await _interactionService.unlike(
          InteractionServiceV1TargetType.targetTypePost, pid);
    } else {
      resp = await _interactionService.like(
          InteractionServiceV1TargetType.targetTypePost, pid);
    }
    if (!mounted || resp == null) return;
    setState(() {
      _liked = resp.liked ?? _liked;
      _likeCount = resp.likeCount ?? _likeCount;
    });
  }

  Future<void> _toggleWatch() async {
    final pid = widget.post.id;
    if (pid == null) return;
    final InteractionServiceV1WatchResponse? resp;
    if (_watched) {
      resp = await _interactionService.unwatch(pid);
    } else {
      resp = await _interactionService.watch(pid);
    }
    if (!mounted || resp == null) return;
    setState(() {
      _watched = resp.watched ?? _watched;
      _watchCount = resp.watchCount ?? _watchCount;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: EdgeInsets.symmetric(
        horizontal: widget.isMobile ? 20.w : 0,
        vertical: 12,
      ),
      child: Container(
        padding: EdgeInsets.symmetric(
          vertical: widget.isMobile ? 12.h : 12,
          horizontal: widget.isMobile ? 16.w : 16,
        ),
        decoration: BoxDecoration(
          color: theme.colorScheme.surface,
          borderRadius: BorderRadius.circular(widget.isMobile ? 14.r : 14),
          border: Border.all(
            color: theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
          ),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            InteractionItem(
              icon: _liked ? Icons.favorite : Icons.favorite_outline,
              value: '$_likeCount',
              label: S.of(context).likes,
              isPressed: _liked,
              onTap: _toggleLike,
            ),
            InteractionItem(
              icon: _watched ? Icons.bookmark : Icons.bookmark_border,
              value: '$_watchCount',
              label: S.of(context).myBookmarks,
              isPressed: _watched,
              onTap: _toggleWatch,
            ),
          ],
        ),
      ),
    );
  }
}

class InteractionItem extends StatefulWidget {
  final IconData icon;
  final String value;
  final String label;

  /// 是否为可交互按钮态（点赞/收藏）。false 表示纯展示 item（visits/comments/share）。
  final bool isPressed;

  /// 点击回调，仅 isPressed=true 时有意义。
  final Future<void> Function()? onTap;

  const InteractionItem({
    super.key,
    required this.icon,
    required this.value,
    required this.label,
    this.isPressed = false,
    this.onTap,
  });

  @override
  State<InteractionItem> createState() => _InteractionItemState();
}

class _InteractionItemState extends State<InteractionItem> {
  bool _isHovered = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    final color = widget.isPressed
        ? theme.colorScheme.primary
        : (_isHovered
            ? theme.colorScheme.primary
            : theme.colorScheme.onSurface.withAlpha(160));

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: InkWell(
        onTap: widget.onTap == null ? null : () => widget.onTap!(),
        borderRadius: BorderRadius.circular(8),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          child: Column(
            children: [
              Icon(
                widget.icon,
                size: 20,
                color: color,
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
