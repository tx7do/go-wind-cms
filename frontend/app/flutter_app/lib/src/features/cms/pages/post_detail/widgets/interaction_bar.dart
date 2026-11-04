import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/repositories/user_auth_cache.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post;
import 'package:flutter_app/src/features/cms/services/interaction_service.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show InteractionServiceV1TargetType, InteractionServiceV1LikeResponse,
        InteractionServiceV1WatchResponse, InteractionServiceV1CounterMetric,
        InteractionServiceV1GetInteractionStatusResponse,
        InteractionServiceV1GetCountsResponse;

typedef Post = ContentServiceV1Post;

/// 文章交互条（点赞 / 收藏）
///
/// 横排（正文底部，移动端与平板）与纵排（桌面右栏悬浮位）两种形态共用同一份
/// 交互态。likes 与 bookmark 为可交互按钮，状态由 InteractionService 提供：
///   - 进入时调 getInteractionStatus 取初始 {liked, watched} 态
///   - 计数来自 interaction_counter 表（post.likes 列已移除），getCounts 批量取
///   - 点击后由响应体乐观更新；未登录时后端返回 401，按钮保持未激活态
/// viewer 身份由后端鉴权上下文确定，前端无需传 user_id。
class InteractionBar extends StatefulWidget {
  final Post post;
  final bool isMobile;

  /// 桌面右栏纵排形态
  final bool vertical;

  const InteractionBar({
    super.key,
    required this.post,
    required this.isMobile,
    this.vertical = false,
  });

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
    // 匿名访客必然 401（status 需 viewer 身份），直接保持未激活态，
    // 不发这条注定失败的请求（计数走匿名可访问的 counts:list，不受影响）
    if (!GetIt.instance<UserAuthCache>().loginStateNotifier.value) return;
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
    final st = statusResp?.statuses?[pid.toString()];
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
    // 在 setState 闭包外解包，规避闭包捕获导致空安全提升失效
    final liked = resp.liked;
    final likeCount = resp.likeCount;
    setState(() {
      _liked = liked ?? _liked;
      _likeCount = likeCount ?? _likeCount;
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
    final watched = resp.watched;
    final watchCount = resp.watchCount;
    setState(() {
      _watched = watched ?? _watched;
      _watchCount = watchCount ?? _watchCount;
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final p = theme.palette;
    final compact = widget.vertical;

    return Padding(
      padding: EdgeInsets.symmetric(
        horizontal: widget.isMobile && !compact ? 20 : 0,
        vertical: compact ? 0 : 12,
      ),
      child: Container(
        width: compact ? 64 : null,
        padding: EdgeInsets.symmetric(
          vertical: compact ? theme.space.sm : theme.space.md,
          horizontal: compact ? theme.space.xxs : theme.space.lg,
        ),
        decoration: BoxDecoration(
          color: p.surface,
          borderRadius: theme.radius.mdAll,
          border: Border.all(color: p.border, width: theme.borderWidth.thin),
        ),
        child: compact
            ? Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  _likeItem(),
                  SizedBox(
                    width: 32,
                    child: Divider(height: theme.space.md, color: p.hairline),
                  ),
                  _watchItem(),
                ],
              )
            : Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [_likeItem(), _watchItem()],
              ),
      ),
    );
  }

  Widget _likeItem() => InteractionItem(
        icon: _liked ? Icons.favorite : Icons.favorite_outline,
        value: '$_likeCount',
        label: widget.vertical ? '' : S.of(context).likes,
        isPressed: _liked,
        onTap: _toggleLike,
      );

  Widget _watchItem() => InteractionItem(
        icon: _watched ? Icons.bookmark : Icons.bookmark_border,
        value: '$_watchCount',
        label: widget.vertical ? '' : S.of(context).myBookmarks,
        isPressed: _watched,
        onTap: _toggleWatch,
      );
}

class InteractionItem extends StatefulWidget {
  final IconData icon;
  final String value;
  final String label;

  /// 是否为可交互按钮态（点赞/收藏）。false 表示纯展示 item。
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
    final p = theme.palette;

    final color = widget.isPressed || _isHovered ? p.primary : p.inkSecondary;

    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => _isHovered = true),
      onExit: (_) => setState(() => _isHovered = false),
      child: InkWell(
        onTap: widget.onTap == null ? null : () => widget.onTap!(),
        borderRadius: theme.radius.smAll,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          child: Column(
            children: [
              Icon(widget.icon, size: 20, color: color),
              const SizedBox(height: 4),
              Text(
                widget.value,
                style: theme.textTheme.labelMedium?.copyWith(
                  fontWeight: FontWeight.w600,
                  color: widget.isPressed ? p.primary : p.ink,
                ),
              ),
              if (widget.label.isNotEmpty) ...[
                const SizedBox(height: 2),
                Text(widget.label, style: theme.textTheme.labelSmall),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
