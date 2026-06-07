import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/widgets/app_bottom_nav_bar.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_content_view.dart';
import 'package:flutter_app/src/features/cms/pages/explore/explore_page.dart';
import 'package:flutter_app/src/features/cms/pages/bookmarks/bookmarks_page.dart';
import 'package:flutter_app/src/features/cms/pages/profile/profile_page.dart';

final _navItems = [
  BottomNavItem(
    icon: Icons.home_outlined,
    activeIcon: Icons.home,
    localizedName: (loc) => loc.home,
  ),
  BottomNavItem(
    icon: Icons.explore_outlined,
    activeIcon: Icons.explore,
    localizedName: (loc) => loc.discover,
  ),
  BottomNavItem(
    icon: Icons.bookmark_border,
    activeIcon: Icons.bookmark,
    localizedName: (loc) => loc.bookmarks,
  ),
  BottomNavItem(
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
    return Scaffold(
      body: _buildPage(_currentIndex),
      bottomNavigationBar: AppBottomNavBar(
        selectedIndex: _currentIndex,
        onDestinationSelected: (index) {
          setState(() {
            _currentIndex = index;
          });
        },
        items: _navItems,
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
