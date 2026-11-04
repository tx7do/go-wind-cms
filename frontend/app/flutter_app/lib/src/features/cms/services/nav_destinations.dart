import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart' show IconData, Icons;
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show SiteServiceV1Navigation, SiteServiceV1Navigation$Location;
import 'package:flutter_app/src/core/preference/user_preference_cache.dart';
import 'package:flutter_app/src/features/cms/services/navigation_service.dart';

typedef Navigation = SiteServiceV1Navigation;
typedef NavigationLocation = SiteServiceV1Navigation$Location;

/// 一个已解析、可直接渲染的导航位（顶栏链接 / 底栏 tab / 侧栏 / 页脚共用）。
@immutable
class NavSlot {
  const NavSlot({
    required this.label,
    required this.route,
    required this.icon,
    required this.activeIcon,
    this.external = false,
    this.openNewTab = false,
  });

  final String label;
  final String route;
  final IconData icon;
  final IconData activeIcon;
  final bool external;
  final bool openNewTab;

  static NavSlot? from(NavigationItem item) {
    final route = resolveNavRoute(item);
    if (route == null || route.isEmpty) return null;
    final iconName = item.icon ?? '';
    return NavSlot(
      label: item.title ?? '',
      route: route,
      icon: resolveNavIcon(iconName, fallback: Icons.article_outlined),
      // 服务端按 MD3 惯例成对给图标（xxx_outlined / xxx）
      activeIcon: resolveNavIcon(
        iconName.replaceAll('_outlined', ''),
        fallback: Icons.article,
      ),
      external: isExternalLink(item),
      openNewTab: item.isOpenNewTab == true,
    );
  }
}

/// 全站导航数据源：一次拉取，多处消费（顶栏、底栏、侧栏、页脚）。
///
/// 此前 WebShellLayout 与 HomeMobileView 各自 `_loadNavigations()`，同一
/// 接口在壳层里被反复请求，locale 兜底也各写一遍。
class NavDestinationStore extends ChangeNotifier {
  NavDestinationStore._();

  static final NavDestinationStore instance = NavDestinationStore._();

  final NavigationService _service = NavigationService();

  List<Navigation> _all = const [];
  bool _loading = false;
  bool _loaded = false;

  bool get isLoading => _loading;
  bool get isLoaded => _loaded;

  /// 当前语言偏好（zh_CN → zh-CN）；取不到时返回 null，走通用导航
  String? get _locale {
    try {
      final lang = GetIt.instance<UserPreferenceCache>().language;
      if (lang.isEmpty) return null;
      return lang.replaceAll('_', '-');
    } catch (e) {
      // 偏好未注册（测试/极早期启动）时按通用导航渲染，但要留痕
      debugPrint('[NavDestinationStore] preference unavailable: $e');
      return null;
    }
  }

  Future<void> ensureLoaded() async {
    if (_loading || _loaded) return;
    _loading = true;
    notifyListeners();
    try {
      final result = await _service.list();
      _all = (result as ListNavigationResponse?)?.items ?? const [];
      _loaded = true;
    } catch (e, s) {
      // 不吞错：失败要进日志，UI 侧退回内置默认导航
      debugPrint('[NavDestinationStore] navigation list failed: $e');
      debugPrintStack(stackTrace: s);
      _all = const [];
    } finally {
      _loading = false;
    }
    notifyListeners();
  }

  /// 指定位置的扁平导航项（已按 sortOrder 排好）
  List<NavSlot> slotsOf(NavigationLocation location) {
    return itemsOf(location)
        .map(NavSlot.from)
        .whereType<NavSlot>()
        .where((s) => !s.external)
        .toList();
  }

  /// 含外链的原始项（顶栏需要外链）
  List<NavigationItem> itemsOf(NavigationLocation location) =>
      getFlatNavItems(_all, location, locale: _locale);

  void invalidate() {
    _loaded = false;
    _all = const [];
    notifyListeners();
  }
}
