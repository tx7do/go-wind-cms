import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/src/features/cms/pages/home/home_page.dart';
import 'package:flutter_app/src/features/cms/pages/explore/explore_page.dart';
import 'package:flutter_app/src/features/cms/pages/bookmarks/bookmarks_page.dart';
import 'package:flutter_app/src/features/cms/pages/profile/profile_page.dart';
import 'package:flutter_app/src/features/cms/services/navigation_service.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_list_navigation_response.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';

/// CMS 主布局 - 底部导航 + 页面切换
class CmsMainScaffold extends StatefulWidget {
  const CmsMainScaffold({super.key});

  @override
  State<CmsMainScaffold> createState() => _CmsMainScaffoldState();
}

class _CmsMainScaffoldState extends State<CmsMainScaffold> {
  int _currentIndex = 0;
  final _navService = NavigationService();

  List<SiteServiceV1Navigation> _navigations = [];

  /// 底部导航对应的页面
  static const List<Widget> _pages = [
    HomePage(),
    ExplorePage(),
    BookmarksPage(),
    ProfilePage(),
  ];

  @override
  void initState() {
    super.initState();
    _loadNavigations();
  }

  Future<void> _loadNavigations() async {
    final result = await _navService.list();
    if (!mounted) return;
    setState(() {
      _navigations = (result as ListNavigationResponse?)?.items ?? [];
    });
  }

  /// 从 API 数据获取移动端底部导航项
  List get _navItems =>
      getFlatNavItems(_navigations, NavigationLocation.mobile);

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isWide = ResponsiveUtils.isWideScreen(context);

    return Scaffold(
      body: IndexedStack(
        index: _currentIndex,
        children: _pages,
      ),
      // Web/桌面端已有顶部导航栏，隐藏底部导航栏
      bottomNavigationBar: isWide
          ? null
          : _navItems.length < 2
              ? null
              : Container(
              decoration: BoxDecoration(
                color: theme.colorScheme.surface,
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withAlpha((0.05 * 255).round()),
                    blurRadius: 10,
                    offset: const Offset(0, -2),
                  ),
                ],
              ),
              child: SafeArea(
                child: Padding(
                  padding:
                      EdgeInsets.symmetric(horizontal: 8.w, vertical: 6.h),
                  child: NavigationBar(
                    selectedIndex: _currentIndex,
                    onDestinationSelected: (index) {
                      setState(() {
                        _currentIndex = index;
                      });
                    },
                    height: 60.h,
                    backgroundColor: Colors.transparent,
                    elevation: 0,
                    indicatorColor: theme.colorScheme.primaryContainer,
                    labelBehavior:
                        NavigationDestinationLabelBehavior.alwaysShow,
                    destinations: _navItems.asMap().entries.map((entry) {
                      final index = entry.key;
                      final item = entry.value;
                      return NavigationDestination(
                        icon: Icon(resolveNavIcon(item.icon), size: 24.sp),
                        selectedIcon: Icon(resolveNavIcon((item.icon ?? '').replaceAll('_outlined', '')), size: 24.sp),
                        label: item.title ?? '',
                      );
                    }).toList(),
                  ),
                ),
              ),
            ),
    );
  }
}
