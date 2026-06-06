import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/constants/breakpoints.dart';

/// 响应式布局组件
///
/// 根据当前屏幕宽度自动决定渲染哪套 UI：
/// - < 600dp  → mobileBody（手机单栏）
/// - 600~1024dp → tabletBody 或降级为 webBody（平板双栏）
/// - \> 1024dp → webBody（Web/桌面宽屏）
class ResponsiveLayout extends StatelessWidget {
  final Widget mobileBody;
  final Widget? tabletBody;
  final Widget webBody;

  const ResponsiveLayout({
    super.key,
    required this.mobileBody,
    required this.webBody,
    this.tabletBody,
  });

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final width = constraints.maxWidth;
        if (Breakpoints.isMobile(width)) {
          return mobileBody;
        } else if (Breakpoints.isTablet(width)) {
          return tabletBody ?? webBody;
        } else {
          return webBody;
        }
      },
    );
  }
}

/// Web 端居中内容容器
///
/// 限制最大宽度并居中，符合现代网页审美（类似 Medium、掘金）
class WebContentCenter extends StatelessWidget {
  final Widget child;
  final double maxWidth;
  final EdgeInsetsGeometry padding;

  const WebContentCenter({
    super.key,
    required this.child,
    this.maxWidth = Breakpoints.webContentMaxWidth,
    this.padding = EdgeInsets.zero,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: ConstrainedBox(
        constraints: BoxConstraints(maxWidth: maxWidth),
        child: Padding(
          padding: padding,
          child: child,
        ),
      ),
    );
  }
}
