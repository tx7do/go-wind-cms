import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;

/// 焦点推荐轮播组件
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
  final PageController _pageController = PageController();
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
    _pageController.dispose();
    super.dispose();
  }

  void _startAutoPlay() {
    _autoPlayTimer = Timer.periodic(const Duration(seconds: 4), (_) {
      if (_pageController.hasClients) {
        final nextPage = (_currentPage + 1) % widget.posts.length;
        _pageController.animateToPage(
          nextPage,
          duration: const Duration(milliseconds: 500),
          curve: Curves.easeInOut,
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    if (widget.posts.isEmpty) return const SizedBox.shrink();

    final theme = Theme.of(context);
    final isMobile = ResponsiveUtils.isMobile(context);

    return Padding(
      padding: isMobile
          ? EdgeInsets.fromLTRB(16.w, 12.h, 16.w, 4.h)
          : const EdgeInsets.fromLTRB(0, 12, 0, 4),
      child: Column(
        children: [
          // 轮播图
          SizedBox(
            height: isMobile ? 160.h : 180,
            child: PageView.builder(
              controller: _pageController,
              onPageChanged: (index) {
                setState(() => _currentPage = index);
              },
              itemCount: widget.posts.length,
              itemBuilder: (context, index) {
                final post = widget.posts[index];
                return _FeaturedCard(post: post, categories: widget.categories);
              },
            ),
          ),
          SizedBox(height: isMobile ? 10.h : 10),
          // 指示器
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: List.generate(widget.posts.length, (index) {
              return AnimatedContainer(
                duration: const Duration(milliseconds: 300),
                margin: EdgeInsets.symmetric(horizontal: isMobile ? 3.w : 3),
                width: _currentPage == index
                    ? (isMobile ? 20.w : 20)
                    : (isMobile ? 6.w : 6),
                height: isMobile ? 6.h : 6,
                decoration: BoxDecoration(
                  color: _currentPage == index
                      ? theme.colorScheme.primary
                      : theme.colorScheme.onSurface.withAlpha(40),
                  borderRadius: BorderRadius.circular(isMobile ? 3.r : 3),
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
  final Post post;
  final List<Category> categories;

  const _FeaturedCard({required this.post, this.categories = const []});

  String get _title => (post.translations ?? []).isNotEmpty
      ? post.translations!.first.title ?? ''
      : '';

  /// 获取文章关联的分类名称
  String get _categoryName {
    if ((post.categoryIds ?? []).isEmpty) return '';
    final catId = post.categoryIds!.first;
    try {
      final cat = categories.firstWhere((c) => c.id != null && c.id == catId);
      return (cat.translations ?? []).isNotEmpty
          ? cat.translations!.first.name ?? ''
          : '';
    } catch (_) {
      return '';
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isMobile = ResponsiveUtils.isMobile(context);

    return Container(
      margin: EdgeInsets.symmetric(horizontal: isMobile ? 4.w : 4),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(isMobile ? 16.r : 16),
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            theme.colorScheme.primaryContainer,
            theme.colorScheme.primary.withAlpha((0.3 * 255).round()),
          ],
        ),
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(isMobile ? 16.r : 16),
        child: Stack(
          children: [
            // 背景装饰
            Positioned(
              right: -20,
              top: -20,
              child: Container(
                width: 120,
                height: 120,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: theme.colorScheme.primary.withAlpha(
                    (0.1 * 255).round(),
                  ),
                ),
              ),
            ),
            Positioned(
              left: -30,
              bottom: -30,
              child: Container(
                width: 100,
                height: 100,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: theme.colorScheme.primary.withAlpha(
                    (0.08 * 255).round(),
                  ),
                ),
              ),
            ),
            // 内容
            Padding(
              padding: EdgeInsets.all(isMobile ? 20.w : 20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  if (_categoryName.isNotEmpty)
                    Container(
                      padding: EdgeInsets.symmetric(
                        horizontal: isMobile ? 10.w : 10,
                        vertical: isMobile ? 4.h : 4,
                      ),
                      decoration: BoxDecoration(
                        color: theme.colorScheme.primary.withAlpha(
                          (0.2 * 255).round(),
                        ),
                        borderRadius: BorderRadius.circular(
                          isMobile ? 12.r : 12,
                        ),
                      ),
                      child: Text(
                        _categoryName,
                        style: TextStyle(
                          fontSize: isMobile ? 11.sp : 11,
                          fontWeight: FontWeight.w500,
                          color: theme.colorScheme.onSurface,
                        ),
                      ),
                    ),
                  SizedBox(height: isMobile ? 8.h : 8),
                  Text(
                    _title,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: TextStyle(
                      fontSize: isMobile ? 16.sp : 16,
                      fontWeight: FontWeight.bold,
                      color: theme.colorScheme.onSurface,
                      height: 1.3,
                    ),
                  ),
                  SizedBox(height: isMobile ? 6.h : 6),
                  Row(
                    children: [
                      Text(
                        post.authorName ?? '',
                        style: TextStyle(
                          fontSize: isMobile ? 12.sp : 12,
                          color: theme.colorScheme.onSurface.withAlpha(160),
                        ),
                      ),
                      SizedBox(width: isMobile ? 12.w : 12),
                      Icon(
                        Icons.remove_red_eye_outlined,
                        size: isMobile ? 14.sp : 14,
                        color: theme.colorScheme.onSurface.withAlpha(120),
                      ),
                      SizedBox(width: isMobile ? 3.w : 3),
                      Text(
                        '${post.visits}',
                        style: TextStyle(
                          fontSize: isMobile ? 12.sp : 12,
                          color: theme.colorScheme.onSurface.withAlpha(120),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
