import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:url_launcher/url_launcher.dart';

/// 导航栏链接
///
/// 用于 Web 端顶部导航栏的链接项，支持内部路由跳转和外部链接。
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
      if (uri != null) launchUrl(uri);
      return;
    }

    try {
      context.go(route);
    } catch (_) {
      context.go('/');
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = widget.isActive
        ? theme.colorScheme.primary
        : theme.colorScheme.onSurface;

    return GestureDetector(
      onTap: _handleTap,
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        onEnter: (_) => setState(() => _isHovered = true),
        onExit: (_) => setState(() => _isHovered = false),
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            border: widget.isActive
                ? Border(
                    bottom: BorderSide(
                      color: theme.colorScheme.primary,
                      width: 2,
                    ),
                  )
                : null,
          ),
          child: Text(
            widget.label,
            style: TextStyle(
              fontSize: 14,
              fontWeight: widget.isActive ? FontWeight.w600 : FontWeight.w400,
              color: _isHovered && !widget.isActive
                  ? theme.colorScheme.primary
                  : color,
            ),
          ),
        ),
      ),
    );
  }
}
