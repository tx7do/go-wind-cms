import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/src/features/cms/services/navigation_service.dart';

/// 底部导航项定义
class BottomNavItem {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  final String? route;

  const BottomNavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
    this.route,
  });

  /// 从服务端 NavigationItem 构造
  factory BottomNavItem.fromNavigationItem(NavigationItem item) {
    final iconName = item.icon ?? '';
    final activeIconName = iconName.replaceAll('_outlined', '');
    return BottomNavItem(
      icon: resolveNavIcon(iconName, fallback: Icons.article_outlined),
      activeIcon: resolveNavIcon(activeIconName, fallback: Icons.article),
      label: item.title ?? '',
      route: resolveNavRoute(item),
    );
  }
}

/// 应用底部导航栏
///
/// 移动端专用的底部导航组件，封装了样式、阴影、安全区域等。
class AppBottomNavBar extends StatelessWidget {
  final int selectedIndex;
  final ValueChanged<int> onDestinationSelected;
  final List<BottomNavItem> items;

  const AppBottomNavBar({
    super.key,
    required this.selectedIndex,
    required this.onDestinationSelected,
    required this.items,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Container(
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
          // 用 Theme 覆盖去掉水波纹
          child: Theme(
            data: theme.copyWith(
              splashFactory: NoSplash.splashFactory,
              highlightColor: Colors.transparent,
              splashColor: Colors.transparent,
            ),
            child: NavigationBar(
              selectedIndex: selectedIndex,
              onDestinationSelected: onDestinationSelected,
              height: 60.h,
              backgroundColor: Colors.transparent,
              elevation: 0,
              indicatorColor: theme.colorScheme.primaryContainer,
              surfaceTintColor: Colors.transparent,
              labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
              destinations: items.map((item) {
                return NavigationDestination(
                  icon: Icon(item.icon, size: 24.sp),
                  selectedIcon: Icon(item.activeIcon, size: 24.sp),
                  label: item.label,
                );
              }).toList(),
            ),
          ),
        ),
      ),
    );
  }
}
