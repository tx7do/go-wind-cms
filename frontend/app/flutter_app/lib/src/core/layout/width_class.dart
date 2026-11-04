import 'package:flutter/widgets.dart';

/// 宽度档——本端所有响应式决策的唯一分档口径。
///
/// 旧的三档（600 / 1024）缺一档「窄桌面」：600dp 一刀切下去，7 寸平板
/// 横屏直接拿到 1140dp 桌面三栏，卡片被压成一条，正文行长冲到 1000px+。
/// 这里拆四档，每档只负责一件事：
///
/// | 档 | 宽度 | 布局 | 导航 |
/// |---|---|---|---|
/// | compact | < 600 | 单栏 | 底部标签栏 |
/// | medium | 600–900 | 双栏卡片，无侧边栏 | 顶部导航（紧凑） |
/// | expanded | 900–1200 | 主列 + 侧边栏 | 顶部导航 |
/// | wide | ≥ 1200 | 三栏（侧边栏加宽） | 顶部导航 |
enum WidthClass { compact, medium, expanded, wide }

/// 正文行长上限（doc §3：阅读栏 ≤ 760dp，任何档都不例外）。
const double kArticleLineWidth = 760;

extension WidthClassOf on WidthClass {
  static const double _mediumStart = 600;
  static const double _expandedStart = 900;
  static const double _wideStart = 1200;

  /// 按视口宽度取档。边界值归入更高一档（与 Material 断点口径一致）。
  static WidthClass fromWidth(double width) {
    if (width >= _wideStart) return WidthClass.wide;
    if (width >= _expandedStart) return WidthClass.expanded;
    if (width >= _mediumStart) return WidthClass.medium;
    return WidthClass.compact;
  }

  bool get isCompact => this == WidthClass.compact;
  bool get isWide => this == WidthClass.wide;

  /// 侧边栏只在 expanded/wide 出现；medium 档宁可让卡片多一列。
  bool get showSidebar => index >= WidthClass.expanded.index;

  /// 移动端形态：底部标签栏 + 无顶栏面包屑。
  bool get useBottomNav => isCompact;

  /// 内容区最大宽度（不含 gutter）。compact 不限宽（本身已满宽）。
  double get contentMaxWidth => switch (this) {
    WidthClass.compact => double.infinity,
    WidthClass.medium => 720,
    WidthClass.expanded => 1040,
    WidthClass.wide => 1280,
  };

  /// 内容区左右留白。
  double get gutter => switch (this) {
    WidthClass.compact => 16,
    WidthClass.medium => 24,
    WidthClass.expanded => 32,
    WidthClass.wide => 40,
  };

  double get sidebarWidth => switch (this) {
    WidthClass.compact || WidthClass.medium => 0,
    WidthClass.expanded => 280,
    WidthClass.wide => 300,
  };

  /// 主列与侧边栏之间的沟槽。
  double get sidebarGap => showSidebar ? 32 : 0;

  /// 栅格列数上限（实际列数由可用宽度反推，见 [AdaptiveGrid.columns]）。
  int get maxGridColumns => switch (this) {
    WidthClass.compact => 1,
    WidthClass.medium => 2,
    WidthClass.expanded => 3,
    WidthClass.wide => 3,
  };
}

/// 栅格反推：给「最小卡片宽」，算当前可用宽度能放几列。
///
/// 写死 `crossAxisCount: 2` 的问题在于卡片宽度随视口连续变化，而列数
/// 是离散的——必须在真正测量到的宽度上算，否则要么留大片空白要么挤爆。
abstract final class AdaptiveGrid {
  /// 文章卡片的最小可读宽度（题图 + 两行摘要 + 元信息）。
  static const double minPostCardWidth = 300;

  /// 分类/标签这类信息密度低的块的最小宽度。
  static const double minTileWidth = 160;

  static int columns({
    required double availableWidth,
    required double minItemWidth,
    double gap = 16,
    int maxColumns = 4,
  }) {
    if (availableWidth.isInfinite || availableWidth <= 0) return 1;
    if (maxColumns < 1) return 1;
    // n 列需要 n*min + (n-1)*gap；解出最大的 n
    final fit = ((availableWidth + gap) / (minItemWidth + gap)).floor();
    return fit.clamp(1, maxColumns);
  }
}

/// 从视口取当前档（页面级用这个；组件级要用测量到的宽度，见 AdaptiveGrid）。
extension WidthClassContext on BuildContext {
  WidthClass get widthClass =>
      WidthClassOf.fromWidth(MediaQuery.sizeOf(this).width);
}
