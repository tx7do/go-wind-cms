import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

typedef Post = ContentServiceV1Post;
typedef Category = ContentServiceV1Category;

class PostHeader extends StatelessWidget {
  final Post post;
  final bool isMobile;
  final List<Category> categories;

  const PostHeader({
    super.key,
    required this.post,
    required this.isMobile,
    required this.categories,
  });

  String get _title => getPostTitle(post);

  String get _categoryName {
    if ((post.categoryIds ?? []).isEmpty) return '';
    final catId = post.categoryIds!.first;
    try {
      final cat = categories.firstWhere((c) => c.id != null && c.id == catId);
      return getCategoryName(cat);
    } catch (_) {
      return '';
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: EdgeInsets.fromLTRB(
        isMobile ? 20.w : 0,
        8,
        isMobile ? 20.w : 0,
        0,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            _title,
            style: TextStyle(
              fontSize: isMobile ? 22.sp : 24,
              fontWeight: FontWeight.bold,
              color: theme.colorScheme.onSurface,
              height: 1.4,
            ),
          ),
          SizedBox(height: isMobile ? 14.h : 14),
          Row(
            children: [
              CircleAvatar(
                radius: isMobile ? 18.r : 18,
                backgroundColor: theme.colorScheme.primaryContainer,
                child: Text(
                  post.authorName?.isNotEmpty == true
                      ? post.authorName![0]
                      : '?',
                  style: TextStyle(
                    fontSize: isMobile ? 14.sp : 14,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
              SizedBox(width: isMobile ? 10.w : 10),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      post.authorName ?? '',
                      style: TextStyle(
                        fontSize: isMobile ? 14.sp : 14,
                        fontWeight: FontWeight.w500,
                        color: theme.colorScheme.onSurface,
                      ),
                    ),
                    SizedBox(height: isMobile ? 2.h : 2),
                    Text(
                      post.publishTime != null
                          ? _formatDate(context, post.publishTime!)
                          : '',
                      style: TextStyle(
                        fontSize: isMobile ? 12.sp : 12,
                        color: theme.colorScheme.onSurface.withAlpha(120),
                      ),
                    ),
                  ],
                ),
              ),
              if (_categoryName.isNotEmpty)
                Container(
                  padding: EdgeInsets.symmetric(
                    horizontal: isMobile ? 12.w : 12,
                    vertical: isMobile ? 5.h : 5,
                  ),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.primaryContainer.withAlpha(
                      (0.5 * 255).round(),
                    ),
                    borderRadius:
                        BorderRadius.circular(isMobile ? 12.r : 12),
                  ),
                  child: Text(
                    _categoryName,
                    style: TextStyle(fontSize: isMobile ? 12.sp : 12),
                  ),
                ),
            ],
          ),
          Divider(
            height: isMobile ? 24.h : 24,
            thickness: 1,
            color: theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
          ),
        ],
      ),
    );
  }

  String _formatDate(BuildContext context, DateTime date) {
    final loc = S.of(context);
    return loc.yearMonthDay(date.year, date.month, date.day);
  }
}
