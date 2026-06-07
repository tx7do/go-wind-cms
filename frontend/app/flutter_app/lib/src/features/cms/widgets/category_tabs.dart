import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

typedef Category = ContentServiceV1Category;

/// 分类横滑 Tab 组件
class CategoryTabs extends StatelessWidget {
  final List<Category> categories;
  final TabController tabController;

  const CategoryTabs({
    super.key,
    required this.categories,
    required this.tabController,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return TabBar(
      controller: tabController,
      isScrollable: true,
      tabAlignment: TabAlignment.start,
      padding: EdgeInsets.symmetric(horizontal: 12.w),
      labelPadding: EdgeInsets.symmetric(horizontal: 16.w),
      indicatorSize: TabBarIndicatorSize.label,
      indicatorColor: theme.colorScheme.primary,
      indicatorWeight: 3,
      indicatorPadding: EdgeInsets.only(bottom: 4.h),
      labelColor: theme.colorScheme.onSurface,
      unselectedLabelColor: theme.colorScheme.onSurface.withAlpha(128),
      labelStyle: TextStyle(fontSize: 15.sp, fontWeight: FontWeight.w600),
      unselectedLabelStyle: TextStyle(
        fontSize: 14.sp,
        fontWeight: FontWeight.w400,
      ),
      tabs: categories.map((cat) => Tab(text: getCategoryName(cat))).toList(),
    );
  }
}
