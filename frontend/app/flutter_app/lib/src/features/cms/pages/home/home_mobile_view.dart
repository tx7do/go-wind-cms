import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:get_it/get_it.dart';

import 'package:flutter_app/src/core/widgets/app_bottom_nav_bar.dart';
import 'package:flutter_app/src/core/preference/user_preference_cache.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_content_view.dart';
import 'package:flutter_app/src/features/cms/pages/explore/explore_page.dart';
import 'package:flutter_app/src/features/cms/pages/bookmarks/bookmarks_page.dart';
import 'package:flutter_app/src/features/cms/pages/profile/profile_page.dart';
import 'package:flutter_app/src/features/cms/services/navigation_service.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show SiteServiceV1Navigation, SiteServiceV1Navigation$Location;

/// 首页 - 手机端视图（带底部导航，动态加载）
class HomeMobileView extends StatefulWidget {
  /// 初始显示的路由（如从 /profile 切回时自动定位到对应 tab）
  final String? initialRoute;

  const HomeMobileView({super.key, this.initialRoute});

  @override
  State<HomeMobileView> createState() => _HomeMobileViewState();
}

class _HomeMobileViewState extends State<HomeMobileView> {
  final _navService = NavigationService();

  List<BottomNavItem> _navItems = [];
  int _currentIndex = 0;
  bool _isLoading = true;

  /// 获取当前用户语言偏好
  String? get _currentLocale {
    try {
      final lang = GetIt.instance<UserPreferenceCache>().language;
      if (lang.isEmpty) return null;
      return lang.replaceAll('_', '-');
    } catch (_) {
      return null;
    }
  }

  @override
  void initState() {
    super.initState();
    _loadNavItems();
  }

  Future<void> _loadNavItems() async {
    final result = await _navService.list();
    if (!mounted) return;

    final navigations =
        (result as ListNavigationResponse?)?.items ?? <SiteServiceV1Navigation>[];
    final items = getFlatNavItems(
      navigations,
      SiteServiceV1Navigation$Location.mobile,
      locale: _currentLocale,
    );

    setState(() {
      _navItems = items
          .map((item) => BottomNavItem.fromNavigationItem(item))
          .toList();
      _isLoading = false;
    });

    // 初始路由同步
    if (widget.initialRoute != null) {
      syncTabFromRoute(widget.initialRoute!);
    }
  }

  @override
  Widget build(BuildContext context) {
    // 无导航数据时显示首页内容（无底部导航）
    if (_isLoading || _navItems.isEmpty) {
      return const Scaffold(body: HomeContentView());
    }

    return Scaffold(
      body: AnimatedSwitcher(
        duration: const Duration(milliseconds: 250),
        switchInCurve: Curves.easeInOut,
        switchOutCurve: Curves.easeInOut,
        child: KeyedSubtree(
          key: ValueKey(_currentIndex),
          child: _buildPage(_currentIndex),
        ),
      ),
      bottomNavigationBar: AppBottomNavBar(
        selectedIndex: _currentIndex,
        onDestinationSelected: _onNavTap,
        items: _navItems,
      ),
    );
  }

  void _onNavTap(int index) {
    if (index == _currentIndex) return;
    setState(() => _currentIndex = index);
  }

  /// 根据当前路由路径自动切换到对应的 tab
  ///
  /// 从路由 builder 传入 [initialRoute]，安全地匹配 tab
  void syncTabFromRoute(String route) {
    for (int i = 0; i < _navItems.length; i++) {
      if (_navItems[i].route == route && i != _currentIndex) {
        setState(() => _currentIndex = i);
        return;
      }
    }
  }

  Widget _buildPage(int index) {
    final route = index < _navItems.length ? _navItems[index].route : null;
    if (route == '/' || route == null || route.isEmpty) {
      return const HomeContentView();
    }
    return _RoutePage(route: route);
  }
}

/// 根据路由路径显示对应页面
class _RoutePage extends StatelessWidget {
  final String route;
  const _RoutePage({required this.route});

  @override
  Widget build(BuildContext context) {
    // 已知页面路由映射
    switch (route) {
      case '/explore':
        return const ExplorePage();
      case '/bookmarks':
        return const BookmarksPage();
      case '/profile':
        return const ProfilePage();
      default:
        // 兜底：尝试通过 GoRouter 跳转
        return Center(
          child: TextButton(
            onPressed: () => context.go(route),
            child: Text('Go to $route'),
          ),
        );
    }
  }
}
