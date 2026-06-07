import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_content_view.dart';
import 'package:flutter_app/src/features/cms/pages/explore/explore_page.dart';
import 'package:flutter_app/src/features/cms/pages/bookmarks/bookmarks_page.dart';
import 'package:flutter_app/src/features/cms/pages/profile/profile_page.dart';

/// 底部导航项定义
class _BottomNavItem {
  final IconData icon;
  final IconData activeIcon;
  final String Function(S) localizedName;

  const _BottomNavItem({
    required this.icon,
    required this.activeIcon,
    required this.localizedName,
  });
}

final _bottomNavItems = [
  _BottomNavItem(
    icon: Icons.home_outlined,
    activeIcon: Icons.home,
    localizedName: (loc) => loc.home,
  ),
  _BottomNavItem(
    icon: Icons.explore_outlined,
    activeIcon: Icons.explore,
    localizedName: (loc) => loc.discover,
  ),
  _BottomNavItem(
    icon: Icons.bookmark_border,
    activeIcon: Icons.bookmark,
    localizedName: (loc) => loc.bookmarks,
  ),
  _BottomNavItem(
    icon: Icons.person_outline,
    activeIcon: Icons.person,
    localizedName: (loc) => loc.me,
  ),
];

/// 首页 - 手机端视图（带底部导航）
class HomeMobileView extends StatefulWidget {
  const HomeMobileView({super.key});

  @override
  State<HomeMobileView> createState() => _HomeMobileViewState();
}

class _HomeMobileViewState extends State<HomeMobileView> {
  int _currentIndex = 0;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final loc = S.of(context);

    return Scaffold(
      body: _buildPage(_currentIndex),
      bottomNavigationBar: Container(
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
              labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
              destinations: _bottomNavItems.map((item) {
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

  Widget _buildPage(int index) {
    switch (index) {
      case 0:
        return const HomeContentView();
      case 1:
        return const ExplorePage();
      case 2:
        return const BookmarksPage();
      case 3:
        return const ProfilePage();
      default:
        return const HomeContentView();
    }
  }
}
