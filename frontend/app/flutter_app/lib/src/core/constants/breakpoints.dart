/// 设备断点常量
///
/// 设备类型     | 屏幕宽度范围         | 布局策略
/// ----------- | ------------------- | ---------------------------
/// 手机 Mobile | < 600 dp            | 纵向单栏瀑布流
/// 平板 Tablet | 600 dp ~ 1024 dp    | 双栏布局
/// 网页 Web    | > 1024 dp           | 三栏/最大宽度居中
class Breakpoints {
  Breakpoints._();

  /// 手机端最大宽度
  static const double mobile = 600;

  /// 平板端最大宽度
  static const double tablet = 1024;

  /// Web 端内容最大宽度
  static const double webContentMaxWidth = 1200;

  /// Web 端侧边栏宽度
  static const double webSidebarWidth = 280;

  /// Web 端内容区间距
  static const double webContentPadding = 24;

  /// 判断当前是否为手机端
  static bool isMobile(double width) => width < mobile;

  /// 判断当前是否为平板端
  static bool isTablet(double width) => width >= mobile && width < tablet;

  /// 判断当前是否为 Web/桌面端
  static bool isWeb(double width) => width >= tablet;
}
