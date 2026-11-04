import 'package:flutter/material.dart';

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Category;
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/core/utils/translation_helpers.dart';

typedef Category = ContentServiceV1Category;

/// 分类横滑 Tab
///
/// 选中态 = 主色字 + 2.5px 下划线，未选中沉到 inkSecondary——
/// 分类是「筛选行动」，是主色允许出现的少数位置（doc §1）。
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
      padding: EdgeInsets.symmetric(horizontal: theme.space.md),
      labelPadding: EdgeInsets.symmetric(horizontal: theme.space.lg),
      indicatorSize: TabBarIndicatorSize.label,
      indicatorColor: theme.palette.primary,
      indicatorWeight: 2.5,
      indicatorPadding: const EdgeInsets.only(bottom: 4),
      labelColor: theme.palette.inkStrong,
      unselectedLabelColor: theme.palette.inkSecondary,
      labelStyle: theme.textTheme.titleMedium,
      unselectedLabelStyle: theme.textTheme.bodyMedium,
      dividerColor: theme.palette.hairline,
      tabs: categories.map((cat) => Tab(text: getCategoryName(cat))).toList(),
    );
  }
}
