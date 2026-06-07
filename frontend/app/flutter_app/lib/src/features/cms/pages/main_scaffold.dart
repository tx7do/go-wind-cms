import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_page.dart';
import 'package:flutter_app/src/features/cms/pages/explore/explore_page.dart';
import 'package:flutter_app/src/features/cms/pages/bookmarks/bookmarks_page.dart';
import 'package:flutter_app/src/features/cms/pages/profile/profile_page.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';

/// 底部导航项定义
List<_BottomNavItem> _bottomNavItems = [
  _BottomNavItem(
    icon: Icons.home_outlined,
    activeIcon: Icons.home,
    label: 'home',
    localizedName: (loc) => loc.home,
  ),
  _BottomNavItem(
    icon: Icons.explore_outlined,
    activeIcon: Icons.explore,
    label: 'discover',
    localizedName: (loc) => loc.discover,
  ),
  _BottomNavItem(
    icon: Icons.bookmark_border,
    activeIcon: Icons.bookmark,
    label: 'bookmarks',
    localizedName: (loc) => loc.bookmarks,
  ),
  _BottomNavItem(
    icon: Icons.person_outline,
    activeIcon: Icons.person,
    label: 'me',
    localizedName: (loc) => loc.me,
  ),
];

class _BottomNavItem {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  final String Function(S) localizedName;

  const _BottomNavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
    required this.localizedName,
  });
}

/// CMS 主布局 - 底部导航 + 页面切换
class CmsMainScaffold extends StatefulWidget {
  const CmsMainScaffold({super.key});

  @override
  State<CmsMainScaffold> createState() => _CmsMainScaffoldState();
}

class _CmsMainScaffoldState extends State<CmsMainScaffold> {
  int _currentIndex = 0;

  /// 底部导航对应的页面
  static const List<Widget> _pages = [
    HomePage(),
    ExplorePage(),
    BookmarksPage(),
    ProfilePage(),
  ];

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isWide = ResponsiveUtils.isWideScreen(context);

    return Scaffold(
      body: IndexedStack(index: _currentIndex, children: _pages),
      // Web/桌面端已有顶部导航栏，隐藏底部导航栏
      bottomNavigationBar: isWide
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
                  padding: EdgeInsets.symmetric(horizontal: 8.w, vertical: 6.h),
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
                    destinations: _bottomNavItems.map((item) {
                      final loc = S.of(context);
                      return NavigationDestination(
                        icon: Icon(item.icon, size: 24.sp),
                        selectedIcon: Icon(item.activeIcon, size: 24.sp),
                        label: item.localizedName(loc),
                      );
                    }).toList(),
                  ),
                ),
              ),
            ),
    );
  }
}
