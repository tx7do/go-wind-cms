import 'dart:async';

import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1Category;
import 'package:flutter_app/src/core/layout/width_class.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;

/// 焦点推荐轮播
///
/// 大面积氛围件：允许用 hero 渐变与图上深色压暗层（doc §2 例外条款），
/// 其余仍守 token——圆角 lg、主色只出现在指示器选中位。
class FeaturedCarousel extends StatefulWidget {
  final List<Post> posts;
  final List<Category> categories;

  const FeaturedCarousel({
    super.key,
    required this.posts,
    this.categories = const [],
  });

  @override
  State<FeaturedCarousel> createState() => _FeaturedCarouselState();
}

class _FeaturedCarouselState extends State<FeaturedCarousel> {
  int _currentPage = 0;
  Timer? _autoPlayTimer;

  @override
  void initState() {
    super.initState();
    _startAutoPlay();
  }

  @override
  void dispose() {
    _autoPlayTimer?.cancel();
    super.dispose();
  }

  void _startAutoPlay() {
    _autoPlayTimer = Timer.periodic(const Duration(seconds: 4), (_) {
      if (!mounted || widget.posts.isEmpty) return;
      setState(() {
        _currentPage = (_currentPage + 1) % widget.posts.length;
      });
    });
  }

  /// 用户交互后重启自动播放
  void _resetAutoPlay() {
    _autoPlayTimer?.cancel();
    _startAutoPlay();
  }

  void _goPage(int delta) {
    setState(() {
      _currentPage =
          (_currentPage + delta + widget.posts.length) % widget.posts.length;
    });
    _resetAutoPlay();
  }

  @override
  Widget build(BuildContext context) {
    if (widget.posts.isEmpty) return const SizedBox.shrink();

    final theme = Theme.of(context);
    final compact = context.widthClass.isCompact;
    final post = widget.posts[_currentPage];

    return Padding(
      padding: EdgeInsets.fromLTRB(
        compact ? 16 : 0,
        12,
        compact ? 16 : 0,
        4,
      ),
      child: Column(
        children: [
          // AnimatedSwitcher 而非 PageView：避免 CustomScrollView 里嵌 Viewport
          SizedBox(
            height: compact ? 160 : 200,
            child: GestureDetector(
              behavior: HitTestBehavior.opaque,
              onHorizontalDragEnd: (details) {
                final v = details.primaryVelocity;
                if (v == null) return;
                if (v < 0) {
                  _goPage(1); // 左滑 -> 下一页
                } else if (v > 0) {
                  _goPage(-1);
                }
              },
              child: AnimatedSwitcher(
                duration: theme.motion.slow,
                switchInCurve: Curves.easeInOut,
                switchOutCurve: Curves.easeInOut,
                transitionBuilder: (child, animation) =>
                    FadeTransition(opacity: animation, child: child),
                child: KeyedSubtree(
                  key: ValueKey(_currentPage),
                  child: _FeaturedCard(
                    post: post,
                    categories: widget.categories,
                  ),
                ),
              ),
            ),
          ),
          const SizedBox(height: 10),
          // 指示器
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: List.generate(widget.posts.length, (index) {
              final active = _currentPage == index;
              return AnimatedContainer(
                duration: theme.motion.base,
                curve: theme.motion.curve,
                margin: const EdgeInsets.symmetric(horizontal: 3),
                width: active ? 18 : 6,
                height: 6,
                decoration: BoxDecoration(
                  color: active
                      ? theme.palette.primary
                      : theme.palette.borderStrong,
                  borderRadius: theme.radius.fullAll,
                ),
              );
            }),
          ),
        ],
      ),
    );
  }
}

class _FeaturedCard extends StatelessWidget {
  const _FeaturedCard({required this.post, this.categories = const []});

  final Post post;
  final List<Category> categories;

  /// 图上玻璃签的不透明度：照片是主题无关的，这层白纱不随主题变
  static const int _glassAlpha = 0x3D;

  String get _title => getPostTitle(post);

  String? get _coverImage => getPostThumbnail(post);

  String get _categoryName {
    if ((post.categoryIds ?? []).isEmpty) return '';
    final catId = post.categoryIds!.first;
    try {
      final cat = categories.firstWhere((c) => c.id != null && c.id == catId);
      return getCategoryName(cat);
    } catch (e) {
      assert(() {
        debugPrint('[FeaturedCarousel] category lookup failed: $e');
        return true;
      }());
      return '';
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final p = theme.palette;
    final hasCover = _coverImage != null;

    return GestureDetector(
      behavior: HitTestBehavior.opaque,
      onTap: () => context.go('/post/${post.id}'),
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 4),
        decoration: BoxDecoration(
          borderRadius: theme.radius.lgAll,
          gradient: hasCover ? null : p.heroGradient,
        ),
        child: ClipRRect(
          borderRadius: theme.radius.lgAll,
          child: Stack(
            fit: StackFit.expand,
            children: [
              if (hasCover)
                CachedNetworkImage(
                  imageUrl: _coverImage!,
                  fit: BoxFit.cover,
                  placeholder: (context, url) => Container(color: p.surfaceLow),
                  errorWidget: (context, url, error) => Container(
                    color: p.surfaceLow,
                    alignment: Alignment.center,
                    child: Icon(
                      Icons.broken_image_outlined,
                      color: p.inkTertiary,
                    ),
                  ),
                )
              else ...[
                // 无图时的氛围圆：品牌淡底，给首屏一点呼吸感
                Positioned(
                  right: -20,
                  top: -20,
                  child: _circle(120, p.primaryTint),
                ),
                Positioned(
                  left: -30,
                  bottom: -30,
                  child: _circle(100, p.primarySoft),
                ),
              ],
              // 图上压暗层：保证白字可读（doc §2 允许的深色 scrim）
              if (hasCover)
                DecoratedBox(
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      colors: [
                        p.scrim.withAlpha(0x00),
                        p.scrim,
                      ],
                    ),
                  ),
                ),
              Padding(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    if (_categoryName.isNotEmpty) ...[
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 4,
                        ),
                        decoration: BoxDecoration(
                          color: hasCover
                              ? p.onImage.withAlpha(_glassAlpha)
                              : p.primarySoft,
                          borderRadius: theme.radius.smAll,
                        ),
                        child: Text(
                          _categoryName,
                          style: theme.textTheme.labelSmall?.copyWith(
                            fontWeight: FontWeight.w600,
                            color: hasCover ? p.onImage : p.primaryOnSoft,
                          ),
                        ),
                      ),
                      const SizedBox(height: 8),
                    ],
                    Text(
                      _title,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: theme.textTheme.headlineSmall?.copyWith(
                        color: hasCover ? p.onImage : p.inkStrong,
                        height: 1.35,
                      ),
                    ),
                    const SizedBox(height: 6),
                    Text(
                      post.authorName ?? '',
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: theme.textTheme.labelMedium?.copyWith(
                        color: hasCover
                            ? p.onImage.withAlpha(0xCC)
                            : p.inkSecondary,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _circle(double size, Color color) => Container(
        width: size,
        height: size,
        decoration: BoxDecoration(shape: BoxShape.circle, color: color),
      );
}
