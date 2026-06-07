import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/generated/l10n.dart';

/// 404 页面未找到
///
/// 当 GoRouter 路由匹配失败时展示，支持国际化。
class NotFoundPage extends StatelessWidget {
  const NotFoundPage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final loc = S.of(context);

    return Scaffold(
      body: SafeArea(
        child: Center(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 32),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // 404 大字
                Text(
                  '404',
                  style: TextStyle(
                    fontSize: 96,
                    fontWeight: FontWeight.bold,
                    color: theme.colorScheme.primary.withAlpha(80),
                    height: 1,
                  ),
                ),
                const SizedBox(height: 16),
                // 标题
                Text(
                  loc.pageNotFound,
                  style: theme.textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
                ),
                const SizedBox(height: 12),
                // 描述
                Text(
                  loc.pageNotFoundDesc,
                  textAlign: TextAlign.center,
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: theme.colorScheme.onSurface.withAlpha(160),
                  ),
                ),
                const SizedBox(height: 32),
                // 返回首页按钮
                FilledButton.tonal(
                  onPressed: () => context.go('/'),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 16,
                      vertical: 8,
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const Icon(Icons.home_outlined, size: 20),
                        const SizedBox(width: 8),
                        Text(loc.backToHome),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
