import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/constants/breakpoints.dart';
import 'package:flutter_app/src/core/layout/width_class.dart';

/// 响应式布局组件（旧 API，保留给尚未迁移的页面）。
///
/// 新页面请用 [AdaptiveContent]：单一组件树按档换布局，而不是两套 UI
/// 在断点处互换（互换会重建 State，拖窗口就丢滚动位置）。
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
        return switch (WidthClassOf.fromWidth(constraints.maxWidth)) {
          WidthClass.compact => mobileBody,
          WidthClass.medium => tabletBody ?? webBody,
          WidthClass.expanded || WidthClass.wide => webBody,
        };
      },
    );
  }
}

/// Web 端居中内容容器（旧 API，请优先用 [AdaptiveContent]）
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
