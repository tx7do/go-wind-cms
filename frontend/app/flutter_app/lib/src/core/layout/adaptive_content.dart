import 'package:flutter/material.dart';

import 'width_class.dart';

/// 自适应内容容器：限宽 + 留白 + 可选侧边栏，一档一套。
///
/// 取代此前「页面各自 Center + ConstrainedBox(1140)」的写法——那样
/// 600~900dp 之间没有任何形态差，平板直接吃到桌面三栏。
class AdaptiveContent extends StatelessWidget {
  const AdaptiveContent({
    super.key,
    required this.child,
    this.sidebar,
    this.maxWidth,
    this.alignTop = true,
  });

  /// 主列内容
  final Widget child;

  /// 侧边栏；仅在 expanded/wide 档渲染，其余档自动丢弃（不做 display:none 式隐藏）
  final Widget? sidebar;

  /// 覆写主列限宽（文章正文传 [kArticleLineWidth]）
  final double? maxWidth;

  /// 内容是否贴顶（默认 true；单块居中型页面置 false）
  final bool alignTop;

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final width = constraints.maxWidth;
        final wc = WidthClassOf.fromWidth(width);
        final gutter = wc.gutter;

        final main = ConstrainedBox(
          constraints: BoxConstraints(maxWidth: maxWidth ?? wc.contentMaxWidth),
          child: child,
        );

        final showSidebar = sidebar != null && wc.showSidebar;

        final row = showSidebar
            ? Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Expanded(child: main),
                  SizedBox(width: wc.sidebarGap),
                  SizedBox(width: wc.sidebarWidth, child: sidebar),
                ],
              )
            : Center(child: main);

        return Padding(
          padding: EdgeInsets.symmetric(horizontal: gutter),
          child: alignTop ? row : Align(alignment: Alignment.center, child: row),
        );
      },
    );
  }
}

/// 按宽度档切换子树（取代旧的 ResponsiveLayout 双套 UI 分发）。
///
/// 与旧实现的关键差别：分档只看**测量到的宽度**，且 compact 之外的三档
/// 共享同一棵组件树，只在 [AdaptiveContent] 内部按档换布局——窗口拖拽
/// 不会重建 State，收藏/滚动位置不丢。
class WidthBuilder extends StatelessWidget {
  const WidthBuilder({super.key, required this.builder});

  final Widget Function(BuildContext context, WidthClass width, double availableWidth) builder;

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) =>
          builder(context, WidthClassOf.fromWidth(constraints.maxWidth), constraints.maxWidth),
    );
  }
}
