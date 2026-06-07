import 'package:flutter/material.dart';
import 'package:get_it/get_it.dart';
import 'package:go_router/go_router.dart';
import 'package:url_launcher/url_launcher.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/preference/user_preference_cache.dart';
import 'package:flutter_app/src/features/cms/services/navigation_service.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation_location.dart';

typedef NavigationLocation = SiteServiceV1NavigationLocation;

/// 页脚
///
/// 包含底部导航链接（从后端 FOOTER 类型导航加载）和版权信息。
class WebFooter extends StatelessWidget {
  final List<SiteServiceV1Navigation> navigations;

  const WebFooter({super.key, this.navigations = const []});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final s = S.of(context);
    final footerItems = getFlatNavItems(
      navigations,
      NavigationLocation.footer,
      locale: _currentLocale(context),
    );

    return Container(
      padding: const EdgeInsets.symmetric(vertical: 32),
      child: Column(
        children: [
          // 底部导航链接
          if (footerItems.isNotEmpty)
            Padding(
              padding: const EdgeInsets.only(bottom: 20),
              child: Wrap(
                spacing: 24,
                runSpacing: 8,
                alignment: WrapAlignment.center,
                children: footerItems.map((item) {
                  final route = resolveNavRoute(item);
                  final isExternal = isExternalLink(item);
                  return _FooterLink(
                    label: item.title ?? '',
                    route: route,
                    isExternal: isExternal,
                    isOpenNewTab: item.isOpenNewTab == true,
                  );
                }).toList(),
              ),
            ),
          // 版权信息
          Text(
            s.footerText,
            style: TextStyle(
              fontSize: 12,
              color: theme.colorScheme.onSurface.withAlpha(80),
              height: 1.6,
            ),
          ),
        ],
      ),
    );
  }

  String? _currentLocale(BuildContext context) {
    try {
      final lang = GetIt.instance<UserPreferenceCache>().language;
      if (lang.isEmpty) return null;
      return lang.replaceAll('_', '-');
    } catch (_) {
      return null;
    }
  }
}

/// 页脚导航链接
class _FooterLink extends StatefulWidget {
  final String label;
  final String? route;
  final bool isExternal;
  final bool isOpenNewTab;

  const _FooterLink({
    required this.label,
    this.route,
    this.isExternal = false,
    this.isOpenNewTab = false,
  });

  @override
  State<_FooterLink> createState() => _FooterLinkState();
}

class _FooterLinkState extends State<_FooterLink> {
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
    return GestureDetector(
      onTap: _handleTap,
      child: MouseRegion(
        cursor: SystemMouseCursors.click,
        onEnter: (_) => setState(() => _isHovered = true),
        onExit: (_) => setState(() => _isHovered = false),
        child: Text(
          widget.label,
          style: TextStyle(
            fontSize: 13,
            color: _isHovered
                ? theme.colorScheme.primary
                : theme.colorScheme.onSurface.withAlpha(140),
          ),
        ),
      ),
    );
  }
}
