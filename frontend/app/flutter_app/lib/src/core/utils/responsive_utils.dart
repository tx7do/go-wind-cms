import 'package:flutter/widgets.dart';

import 'package:flutter_app/src/core/layout/width_class.dart';

/// 响应式工具类（旧 API 兼容层）。
///
/// 分档口径全部转发到 [WidthClass]；新代码请直接用量宽反推的
/// [AdaptiveGrid.columns]，不要再按视口猜列数。
class ResponsiveUtils {
  ResponsiveUtils._();

  static WidthClass _classOf(BuildContext context) => context.widthClass;

  /// 获取自适应内边距
  static EdgeInsets padding(
    BuildContext context, {
    double mobile = 16,
    double tablet = 24,
    double web = 32,
  }) {
    final value = switch (_classOf(context)) {
      WidthClass.compact => mobile,
      WidthClass.medium => tablet,
      WidthClass.expanded || WidthClass.wide => web,
    };
    return EdgeInsets.all(value);
  }

  /// 获取自适应水平内边距
  static EdgeInsets horizontalPadding(
    BuildContext context, {
    double mobile = 16,
    double tablet = 24,
    double web = 32,
  }) {
    final value = switch (_classOf(context)) {
      WidthClass.compact => mobile,
      WidthClass.medium => tablet,
      WidthClass.expanded || WidthClass.wide => web,
    };
    return EdgeInsets.symmetric(horizontal: value);
  }

  /// 获取自适应字体大小
  static double fontSize(
    BuildContext context,
    double mobileSize, {
    double? tabletSize,
    double? webSize,
  }) {
    return switch (_classOf(context)) {
      WidthClass.compact => mobileSize,
      WidthClass.medium => tabletSize ?? mobileSize,
      WidthClass.expanded || WidthClass.wide => webSize ?? tabletSize ?? mobileSize,
    };
  }

  /// 获取自适应间距
  static double spacing(
    BuildContext context,
    double mobileSize, {
    double? tabletSize,
    double? webSize,
  }) => fontSize(context, mobileSize, tabletSize: tabletSize, webSize: webSize);

  /// 获取网格列数
  static int gridColumns(BuildContext context) =>
      _classOf(context).maxGridColumns;

  /// 获取文章列表网格列数（按视口反推，卡片最小宽 300dp）
  static int postGridColumns(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final wc = _classOf(context);
    return AdaptiveGrid.columns(
      availableWidth: _mainColumnWidth(wc, width),
      minItemWidth: AdaptiveGrid.minPostCardWidth,
      maxColumns: wc.maxGridColumns,
    );
  }

  /// 获取分类网格列数（信息密度低，块最小宽 160dp）
  static int categoryGridColumns(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    final wc = _classOf(context);
    return AdaptiveGrid.columns(
      availableWidth: _mainColumnWidth(wc, width),
      minItemWidth: AdaptiveGrid.minTileWidth,
      maxColumns: 4,
    );
  }

  /// 视口宽度 → 主列可用宽度（扣 gutter 与侧边栏）
  static double _mainColumnWidth(WidthClass wc, double viewportWidth) {
    final content = viewportWidth > wc.contentMaxWidth
        ? wc.contentMaxWidth
        : viewportWidth;
    final usable = content - wc.gutter * 2;
    return wc.showSidebar ? usable - wc.sidebarWidth - wc.sidebarGap : usable;
  }

  /// 判断当前是否为宽屏
  static bool isWideScreen(BuildContext context) =>
      !_classOf(context).isCompact;

  /// 判断当前是否为手机端
  static bool isMobile(BuildContext context) => _classOf(context).isCompact;

  /// 判断当前是否为平板端
  static bool isTablet(BuildContext context) =>
      _classOf(context) == WidthClass.medium;

  /// 获取内容区域最大宽度
  static double contentMaxWidth(BuildContext context) {
    final wc = _classOf(context);
    final width = MediaQuery.sizeOf(context).width;
    return width > wc.contentMaxWidth ? wc.contentMaxWidth : width;
  }
}
