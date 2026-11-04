
import 'package:flutter_app/src/core/layout/width_class.dart';

/// 设备断点常量（旧 API，逐步由 [WidthClass] 取代）。
///
/// 分档口径的唯一事实源在 width_class.dart，这里只做兼容转发：
///
/// | 设备 | 宽度 | 档 | 布局 |
/// |---|---|---|---|
/// | 手机 | < 600dp | compact | 单栏 + 底部标签栏 |
/// | 小平板/折叠屏 | 600–900dp | medium | 双栏卡片，无侧边栏 |
/// | 桌面窄窗 | 900–1200dp | expanded | 主列 + 侧边栏 |
/// | 桌面宽窗 | ≥ 1200dp | wide | 主列 + 加宽侧边栏 |
class Breakpoints {
  Breakpoints._();

  /// 手机端最大宽度（= compact/medium 边界）
  static const double mobile = 600;

  /// 平板端最大宽度（历史口径，现由 expanded/wide 细分）
  static const double tablet = 1024;

  /// 宽屏内容最大宽度（= WidthClass.wide.contentMaxWidth）
  static const double webContentMaxWidth = 1280;

  /// 侧边栏宽度（= WidthClass.expanded.sidebarWidth）
  static const double webSidebarWidth = 280;

  /// 内容区水平留白（= WidthClass.expanded.gutter）
  static const double webContentPadding = 32;

  static bool isMobile(double width) =>
      WidthClassOf.fromWidth(width).isCompact;

  static bool isTablet(double width) {
    final wc = WidthClassOf.fromWidth(width);
    return wc == WidthClass.medium;
  }

  static bool isWeb(double width) =>
      !WidthClassOf.fromWidth(width).isCompact;
}
