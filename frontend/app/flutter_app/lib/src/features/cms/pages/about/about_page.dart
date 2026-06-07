import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/widgets/responsive_layout.dart';

/// 关于我们页面
///
/// 如果后端 Navigation 配置了 linkType=CUSTOM / url=/about，
/// 则此页面会被路由匹配到。
/// 如果后端有 slug='about' 的 page，可以在导航中用 linkType=PAGE + objectId 指向。
class AboutPage extends StatelessWidget {
  const AboutPage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ResponsiveLayout(
      mobileBody: _buildView(context, theme, isMobile: true),
      webBody: _buildView(context, theme, isMobile: false),
    );
  }

  Widget _buildView(BuildContext context, ThemeData theme, {required bool isMobile}) {
    final s = S.of(context);
    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, size: 22),
          onPressed: () => GoRouter.of(context).go('/'),
        ),
        title: Text(s.about),
      ),
      body: Center(
        child: ConstrainedBox(
          constraints: BoxConstraints(
            maxWidth: isMobile ? double.infinity : 720,
          ),
          child: SingleChildScrollView(
            padding: EdgeInsets.symmetric(
              horizontal: isMobile ? 20 : 48,
              vertical: 40,
            ),
            child: _buildStaticContent(theme, s, isMobile: isMobile),
          ),
        ),
      ),
    );
  }

  Widget _buildStaticContent(ThemeData theme, S s, {required bool isMobile}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        // Hero 区域
        Container(
          width: isMobile ? 72 : 96,
          height: isMobile ? 72 : 96,
          decoration: BoxDecoration(
            color: theme.colorScheme.primaryContainer,
            borderRadius: BorderRadius.circular(20),
          ),
          child: Icon(
            Icons.info_outline_rounded,
            size: isMobile ? 36 : 48,
            color: theme.colorScheme.primary,
          ),
        ),
        const SizedBox(height: 20),
        Text(
          s.appName,
          style: TextStyle(
            fontSize: isMobile ? 24 : 32,
            fontWeight: FontWeight.bold,
            color: theme.colorScheme.onSurface,
          ),
        ),
        const SizedBox(height: 8),
        Text(
          s.aboutSubtitle,
          textAlign: TextAlign.center,
          style: TextStyle(
            fontSize: isMobile ? 14 : 16,
            color: theme.colorScheme.onSurface.withAlpha(160),
            height: 1.5,
          ),
        ),
        SizedBox(height: isMobile ? 32 : 48),

        // 特性卡片
        _FeatureCard(
          icon: Icons.article_outlined,
          title: s.aboutFeature1Title,
          description: s.aboutFeature1Desc,
          theme: theme,
        ),
        const SizedBox(height: 16),
        _FeatureCard(
          icon: Icons.language_outlined,
          title: s.aboutFeature2Title,
          description: s.aboutFeature2Desc,
          theme: theme,
        ),
        const SizedBox(height: 16),
        _FeatureCard(
          icon: Icons.devices_outlined,
          title: s.aboutFeature3Title,
          description: s.aboutFeature3Desc,
          theme: theme,
        ),
        SizedBox(height: isMobile ? 32 : 48),

        // 技术栈
        Text(
          s.aboutTechStack,
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w600,
            color: theme.colorScheme.onSurface.withAlpha(180),
          ),
        ),
        const SizedBox(height: 12),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          alignment: WrapAlignment.center,
          children: [
            _TechChip(label: 'Go', theme: theme),
            _TechChip(label: 'Flutter', theme: theme),
            _TechChip(label: 'gRPC', theme: theme),
            _TechChip(label: 'PostgreSQL', theme: theme),
          ],
        ),
      ],
    );
  }
}

/// 特性卡片
class _FeatureCard extends StatelessWidget {
  final IconData icon;
  final String title;
  final String description;
  final ThemeData theme;

  const _FeatureCard({
    required this.icon,
    required this.title,
    required this.description,
    required this.theme,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(14),
        side: BorderSide(
          color: theme.colorScheme.onSurface.withAlpha((0.06 * 255).round()),
        ),
      ),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Row(
          children: [
            Container(
              width: 44,
              height: 44,
              decoration: BoxDecoration(
                color: theme.colorScheme.primaryContainer.withAlpha(120),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(icon, size: 22, color: theme.colorScheme.primary),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w600,
                      color: theme.colorScheme.onSurface,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    description,
                    style: TextStyle(
                      fontSize: 13,
                      color: theme.colorScheme.onSurface.withAlpha(160),
                      height: 1.4,
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

/// 技术标签
class _TechChip extends StatelessWidget {
  final String label;
  final ThemeData theme;

  const _TechChip({required this.label, required this.theme});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerHighest.withAlpha(120),
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 13,
          color: theme.colorScheme.onSurface.withAlpha(180),
        ),
      ),
    );
  }
}
