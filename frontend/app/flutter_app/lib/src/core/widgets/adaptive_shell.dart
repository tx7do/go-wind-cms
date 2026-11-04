import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart' show GetIt;
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/constants/router_paths.dart';
import 'package:flutter_app/src/core/layout/width_class.dart';
import 'package:flutter_app/src/core/repositories/user_auth_cache.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/features/cms/pages/home/widgets/nav_bar_link.dart';
import 'package:flutter_app/src/features/cms/services/nav_destinations.dart';

/// 归一路径：尾部斜杠视为同一路径（根路径除外）
String _normalized(String path) {
  if (path.length > 1 && path.endsWith('/')) return path.substring(0, path.length - 1);
  return path;
}

/// 底栏一个 tab 的规格。
///
/// [branchIndex] 绑定 `StatefulShellRoute` 的分支序号——tab 顺序可以由服务端
/// 导航重排，但分支序号不能跟着动，所以两者必须分开存。
@immutable
class ShellTab {
  const ShellTab({
    required this.branchIndex,
    required this.route,
    required this.label,
    required this.icon,
    required this.activeIcon,
  });

  final int branchIndex;
  final String route;
  final String label;
  final IconData icon;
  final IconData activeIcon;
}

/// 自适应壳层：窄屏底栏 + 宽屏顶栏，同一个 Scaffold。
///
/// 为什么值得单独一层：旧的实现把移动端 tab 状态放在 `HomeMobileView` 里，
/// `/explore`、`/profile` 在窄屏被改写成 `HomePage(initialRoute: ...)`，
/// 于是「切 tab」不是路由跳转，深链、返回键、滚动位置全都对不上；而宽屏
/// 换窄屏时 `ResponsiveLayout` 直接换掉整棵子树，State 全丢。
///
/// 这里分支由 go_router 的 `StatefulShellRoute.indexedStack` 持有，壳层只决定
/// 外框。Scaffold 的 body 槽自带 GlobalKey（见 scaffold.dart `_bodyKey`），
/// 所以跨断点收放 appBar / bottomNavigationBar 不会重挂分支子树。
class AdaptiveShell extends StatefulWidget {
  const AdaptiveShell({
    super.key,
    required this.navigationShell,
    required this.branchRoutes,
    required this.currentPath,
  });

  final StatefulNavigationShell navigationShell;

  /// 与 branches 同序的分支根路径；只有停在根路径上才显示底栏
  final List<String> branchRoutes;

  final String currentPath;

  /// 服务端导航与内置分支的合并规则（纯函数，便于钉测试）。
  ///
  /// 内置四支为底，服务端 MOBILE 导航只改「标签 / 图标 / 顺序」；
  /// 指向未知路由的项**不进底栏**（分支不存在，点了无处可去），
  /// 它们仍可从顶栏导航与搜索到达。
  static List<ShellTab> mergeTabs(
    List<ShellTab> defaults,
    List<NavSlot> serverSlots,
  ) {
    if (serverSlots.isEmpty) return defaults;

    final merged = <ShellTab>[];
    for (final slot in serverSlots) {
      final route = _normalized(slot.route);
      final match = defaults.where((d) => _normalized(d.route) == route);
      // 服务端指向未知路由的项跳过：分支不存在，点了无处可去
      if (match.isEmpty) continue;
      final target = match.first;
      if (merged.any((m) => m.branchIndex == target.branchIndex)) continue;
      merged.add(
        ShellTab(
          branchIndex: target.branchIndex,
          route: target.route,
          label: slot.label.isEmpty ? target.label : slot.label,
          icon: slot.icon,
          activeIcon: slot.activeIcon,
        ),
      );
    }
    // 服务端没提到的内置 tab 补在尾部，四支恒在（底栏缺支会让用户找不到入口）
    for (final d in defaults) {
      if (!merged.any((m) => m.branchIndex == d.branchIndex)) merged.add(d);
    }
    return merged;
  }

  @override
  State<AdaptiveShell> createState() => _AdaptiveShellState();
}

class _AdaptiveShellState extends State<AdaptiveShell> {
  NavDestinationStore get _store => NavDestinationStore.instance;

  @override
  void initState() {
    super.initState();
    _store.addListener(_onNavChanged);
    _store.ensureLoaded();
  }

  @override
  void dispose() {
    _store.removeListener(_onNavChanged);
    super.dispose();
  }

  void _onNavChanged() {
    if (mounted) setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    final compact = context.widthClass.useBottomNav;
    final tabs = _tabs(context);
    final activeBranch = widget.navigationShell.currentIndex;
    final selectedIndex = tabs.indexWhere((t) => t.branchIndex == activeBranch);
    final atBranchRoot = widget.branchRoutes
        .map(_normalized)
        .contains(_normalized(widget.currentPath));
    // 详情页等二级位置不收进底栏：移动端二级页自带返回栏，全屏更干净
    final showTabBar = compact && atBranchRoot && tabs.length > 1;

    return Scaffold(
      appBar: compact
          ? null
          : _TopBar(currentPath: _normalized(widget.currentPath)),
      // 宽屏正文 + 页脚；Scaffold 的 body 槽自带 GlobalKey，加页脚不会重挂分支子树
      body: compact
          ? widget.navigationShell
          : Column(
              children: [
                Expanded(child: widget.navigationShell),
                const _ShellFooter(),
              ],
            ),
      bottomNavigationBar: showTabBar
          ? _BottomTabBar(
              tabs: tabs,
              selectedIndex: selectedIndex < 0 ? 0 : selectedIndex,
              onSelected: (index) => widget.navigationShell
                  .goBranch(tabs[index].branchIndex),
            )
          : null,
    );
  }

  /// 内置四支为底，服务端 MOBILE 导航只改「标签 / 图标 / 顺序」。
  ///
  /// 服务端指向未知路由的项**不进底栏**（分支不存在，点了无处可去），
  /// 它们仍可从顶栏导航与搜索到达。
  List<ShellTab> _tabs(BuildContext context) {
    final s = S.of(context);
    final defaults = <ShellTab>[
      ShellTab(
        branchIndex: 0,
        route: AppRoutePath.initial,
        label: s.home,
        icon: Icons.home_outlined,
        activeIcon: Icons.home,
      ),
      ShellTab(
        branchIndex: 1,
        route: AppRoutePath.explore,
        label: s.discover,
        icon: Icons.explore_outlined,
        activeIcon: Icons.explore,
      ),
      ShellTab(
        branchIndex: 2,
        route: AppRoutePath.bookmarks,
        label: s.bookmarks,
        icon: Icons.bookmark_border,
        activeIcon: Icons.bookmark,
      ),
      ShellTab(
        branchIndex: 3,
        route: AppRoutePath.profile,
        label: s.me,
        icon: Icons.person_outline,
        activeIcon: Icons.person,
      ),
    ];
    return AdaptiveShell.mergeTabs(
      defaults,
      _store.slotsOf(NavigationLocation.mobile),
    );
  }
}

/// 宽屏顶栏：品牌 + 服务端 HEADER 导航 + 搜索/设置/登录
class _TopBar extends StatelessWidget implements PreferredSizeWidget {
  const _TopBar({required this.currentPath});

  final String currentPath;

  @override
  Size get preferredSize => const Size.fromHeight(56);

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final s = S.of(context);
    final space = theme.space;
    final links = NavDestinationStore.instance
        .itemsOf(NavigationLocation.header)
        .map((item) {
          final slot = NavSlot.from(item);
          if (slot == null) return const SizedBox.shrink();
          return NavBarLink(
            label: slot.label,
            route: slot.route,
            isExternal: slot.external,
            isOpenNewTab: slot.openNewTab,
            isActive: _normalized(slot.route) == currentPath,
          );
        })
        .toList();

    return AppBar(
      toolbarHeight: preferredSize.height,
      titleSpacing: space.lg,
      title: Row(
        children: [
          MouseRegion(
            cursor: SystemMouseCursors.click,
            child: GestureDetector(
              behavior: HitTestBehavior.opaque,
              onTap: () => context.go(AppRoutePath.initial),
              child: Text(s.appName, style: theme.textTheme.headlineSmall),
            ),
          ),
          SizedBox(width: space.xxl),
          if (links.isNotEmpty)
            Expanded(
              child: SingleChildScrollView(
                scrollDirection: Axis.horizontal,
                child: Row(children: links),
              ),
            ),
        ],
      ),
      actions: [
        IconButton(
          icon: const Icon(Icons.search),
          tooltip: s.search,
          onPressed: () => context.go(AppRoutePath.search),
        ),
        IconButton(
          icon: const Icon(Icons.settings_outlined),
          tooltip: s.settings,
          onPressed: () => context.go(AppRoutePath.settings),
        ),
        SizedBox(width: space.sm),
        const _AuthEntry(),
        SizedBox(width: space.lg),
      ],
    );
  }
}

/// 登录 / 我的入口
class _AuthEntry extends StatelessWidget {
  const _AuthEntry();

  @override
  Widget build(BuildContext context) {
    final s = S.of(context);
    return ValueListenableBuilder<bool>(
      valueListenable: GetIt.instance<UserAuthCache>().loginStateNotifier,
      builder: (context, hasLogin, _) {
        return MouseRegion(
          cursor: SystemMouseCursors.click,
          child: OutlinedButton(
            onPressed: () => context.go(
              hasLogin ? AppRoutePath.profile : AppRoutePath.login,
            ),
            child: Text(hasLogin ? s.me : s.login),
          ),
        );
      },
    );
  }
}

/// 宽屏页脚：服务端 FOOTER 导航 + 版权行。
///
/// 旧的 WebFooter 是从未被挂载的死文件；这里是壳层的正式页脚位。
/// 窄屏不放页脚（移动端不需要版权位，内容到底即结束）。
class _ShellFooter extends StatelessWidget {
  const _ShellFooter();

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final p = theme.palette;
    final s = S.of(context);
    final links = NavDestinationStore.instance
        .itemsOf(NavigationLocation.footer)
        .map((item) {
          final slot = NavSlot.from(item);
          if (slot == null) return const SizedBox.shrink();
          return NavBarLink(
            label: slot.label,
            route: slot.route,
            isExternal: slot.external,
            isOpenNewTab: slot.openNewTab,
          );
        })
        .toList();

    return Container(
      width: double.infinity,
      decoration: BoxDecoration(
        color: p.surface,
        border: Border(
          top: BorderSide(
            color: p.hairline,
            width: theme.borderWidth.hairline,
          ),
        ),
      ),
      padding: EdgeInsets.symmetric(
        vertical: theme.space.xxl + theme.space.sm,
        horizontal: theme.space.xxl,
      ),
      child: Column(
        children: [
          if (links.isNotEmpty) ...[
            Wrap(
              alignment: WrapAlignment.center,
              spacing: theme.space.lg,
              runSpacing: theme.space.sm,
              children: links,
            ),
            SizedBox(height: theme.space.md),
          ],
          Text(
            s.footerText,
            textAlign: TextAlign.center,
            style: theme.textTheme.labelSmall?.copyWith(height: 1.6),
          ),
        ],
      ),
    );
  }
}

/// 窄屏底栏。
///
/// 刻意不用 MD3 的胶囊指示器：选中态只换字色 + 一条 18×2 的主色短线，
/// 密度更高，也更接近国内 App 的底栏观感（doc §1 铁律：主色只用于行动与选中）。
class _BottomTabBar extends StatelessWidget implements PreferredSizeWidget {
  const _BottomTabBar({
    required this.tabs,
    required this.selectedIndex,
    required this.onSelected,
  });

  final List<ShellTab> tabs;
  final int selectedIndex;
  final ValueChanged<int> onSelected;

  static const double _barHeight = 54;

  @override
  Size get preferredSize => const Size.fromHeight(_barHeight);

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final p = theme.palette;
    final radius = theme.radius;

    return Material(
      color: theme.colorScheme.surface,
      child: SafeArea(
        top: false,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(height: theme.borderWidth.hairline, color: p.hairline),
            SizedBox(
              height: _barHeight,
              child: Row(
                children: [
                  for (int i = 0; i < tabs.length; i++)
                    Expanded(
                      child: _TabButton(
                        tab: tabs[i],
                        active: i == selectedIndex,
                        textStyle: theme.textTheme.labelSmall,
                        indicatorRadius: radius.fullAll,
                        onTap: () => onSelected(i),
                      ),
                    ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _TabButton extends StatelessWidget {
  const _TabButton({
    required this.tab,
    required this.active,
    required this.textStyle,
    required this.indicatorRadius,
    required this.onTap,
  });

  final ShellTab tab;
  final bool active;
  final TextStyle? textStyle;
  final BorderRadius indicatorRadius;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final p = Theme.of(context).palette;
    final color = active ? p.primary : p.inkTertiary;

    return InkResponse(
      onTap: onTap,
      radius: 32,
      overlayColor: WidgetStatePropertyAll(p.primary.withAlpha(0x14)),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(active ? tab.activeIcon : tab.icon, size: 22, color: color),
          const SizedBox(height: 3),
          Text(
            tab.label,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: textStyle?.copyWith(
              color: color,
              fontWeight: active ? FontWeight.w600 : FontWeight.w400,
            ),
          ),
          const SizedBox(height: 2),
          AnimatedContainer(
            duration: const Duration(milliseconds: 180),
            curve: Curves.easeOutCubic,
            width: active ? 18 : 0,
            height: 2,
            decoration: BoxDecoration(color: p.primary, borderRadius: indicatorRadius),
          ),
        ],
      ),
    );
  }
}
