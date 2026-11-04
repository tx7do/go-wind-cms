import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:url_launcher/url_launcher.dart';

import 'package:flutter_app/src/core/themes/tokens.dart';

/// 顶栏导航链接
///
/// 支持内部路由跳转和外部链接；选中态是「主色字 + 主色下划线」，
/// 不做背景填充——顶栏的选中语言与底栏保持一致。
class NavBarLink extends StatefulWidget {
  final String label;
  final String? route;
  final bool isExternal;
  final bool isOpenNewTab;
  final bool isActive;

  const NavBarLink({
    super.key,
    required this.label,
    this.route,
    this.isExternal = false,
    this.isOpenNewTab = false,
    this.isActive = false,
  });

  @override
  State<NavBarLink> createState() => _NavBarLinkState();
}

class _NavBarLinkState extends State<NavBarLink> {
  bool _isHovered = false;

  void _handleTap() {
    final route = widget.route;
    if (route == null || route.isEmpty) return;

    if (widget.isExternal) {
      final uri = Uri.tryParse(route);
      if (uri == null) return;
      launchUrl(
        uri,
        mode: LaunchMode.platformDefault,
        webOnlyWindowName: widget.isOpenNewTab ? '_blank' : null,
      );
      return;
    }
    // 未匹配路由交给 errorBuilder（NotFoundPage），不吞错
    context.go(route);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final p = theme.palette;
    final active = widget.isActive;
    final color = active || _isHovered ? p.primary : p.ink;

    return GestureDetector(
      onTap: _handleTap,
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        onEnter: (_) => setState(() => _isHovered = true),
        onExit: (_) => setState(() => _isHovered = false),
        child: Container(
          padding: EdgeInsets.symmetric(
            horizontal: theme.space.md,
            vertical: theme.space.sm,
          ),
          decoration: BoxDecoration(
            border: active
                ? Border(
                    bottom: BorderSide(
                      color: p.primary,
                      width: theme.borderWidth.emphasis,
                    ),
                  )
                : null,
          ),
          child: Text(
            widget.label,
            style: theme.textTheme.labelLarge?.copyWith(color: color),
          ),
        ),
      ),
    );
  }
}
