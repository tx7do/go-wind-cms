import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/constants/router_paths.dart';
import 'package:flutter_app/src/core/preference/user_preference_cache.dart';
import 'package:flutter_app/src/core/repositories/user_auth_cache.dart';
import 'package:flutter_app/src/features/cms/services/navigation_service.dart';
import 'package:flutter_app/src/features/cms/pages/home/widgets/nav_bar_link.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation_location.dart';

typedef NavigationLocation = SiteServiceV1NavigationLocation;

/// Web 端持久化顶部导航栏 Shell 布局
///
/// 包裹所有 Web 端页面，提供一致的顶部导航栏。
/// 移动端由 ShellRoute 直接返回 child，不经过此组件。
class WebShellLayout extends StatefulWidget {
  final Widget child;
  final String currentPath;

  const WebShellLayout({
    super.key,
    required this.child,
    required this.currentPath,
  });

  @override
  State<WebShellLayout> createState() => _WebShellLayoutState();
}

class _WebShellLayoutState extends State<WebShellLayout> {
  final _navService = NavigationService();
  List<SiteServiceV1Navigation> _navigations = [];

  /// 获取当前用户语言偏好，转换为 API 格式 (zh_CN → zh-CN)
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
    _loadNavigations();
  }

  Future<void> _loadNavigations() async {
    final result = await _navService.list();
    if (!mounted) return;
    setState(() {
      _navigations = (result as ListNavigationResponse?)?.items ?? [];
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        elevation: 0,
        backgroundColor: theme.colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        title: MouseRegion(
          cursor: SystemMouseCursors.click,
          child: GestureDetector(
            onTap: () => context.go('/'),
            child: Text(
              S.of(context).appName,
              style: TextStyle(
                fontSize: 22,
                fontWeight: FontWeight.bold,
                color: theme.colorScheme.onSurface,
              ),
            ),
          ),
        ),
        centerTitle: false,
        actions: [
          // 可滚动的导航链接区域
          Flexible(
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Row(
                children:
                    getFlatNavItems(
                      _navigations,
                      NavigationLocation.header,
                      locale: _currentLocale,
                    )
                    .map((item) {
                      final route = resolveNavRoute(item);
                      return NavBarLink(
                        label: item.title ?? '',
                        route: route,
                        isExternal: isExternalLink(item),
                        isOpenNewTab: item.isOpenNewTab == true,
                        isActive: route == widget.currentPath,
                      );
                    })
                    .toList(),
              ),
            ),
          ),
          const SizedBox(width: 8),
          // 搜索
          IconButton(
            icon: const Icon(Icons.search, size: 22),
            onPressed: () => context.go('/search'),
            tooltip: S.of(context).search,
          ),
          // 设置
          IconButton(
            icon: const Icon(Icons.settings_outlined, size: 22),
            onPressed: () => context.go(AppRoutePath.settings),
            tooltip: S.of(context).settings,
          ),
          const SizedBox(width: 4),
          // 登录/个人中心
          ValueListenableBuilder<bool>(
            valueListenable:
                GetIt.instance<UserAuthCache>().loginStateNotifier,
            builder: (context, hasLogin, _) {
              return MouseRegion(
                cursor: SystemMouseCursors.click,
                child: OutlinedButton(
                  onPressed: () => context.go(
                    hasLogin ? AppRoutePath.profile : AppRoutePath.login,
                  ),
                  style: OutlinedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 8,
                    ),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: Text(
                    hasLogin ? S.of(context).me : S.of(context).login,
                  ),
                ),
              );
            },
          ),
          const SizedBox(width: 8),
        ],
      ),
      body: widget.child,
    );
  }
}
