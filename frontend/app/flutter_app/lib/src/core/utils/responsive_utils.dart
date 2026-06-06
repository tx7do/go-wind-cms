import 'package:flutter/widgets.dart';

import 'package:flutter_app/src/core/constants/breakpoints.dart';

/// 响应式工具类
///
/// 提供根据屏幕宽度自适应的间距、字体等方法。
/// - 手机端：使用 base 值
/// - 平板/Web端：使用 base 值（不随窗口等比例放大）
///
/// 取代全量使用 flutter_screenutil 的 .w .h .sp 模式
class ResponsiveUtils {
  ResponsiveUtils._();

  /// 获取自适应内边距
  static EdgeInsets padding(
    BuildContext context, {
    double mobile = 16,
    double tablet = 24,
    double web = 32,
  }) {
    final width = MediaQuery.sizeOf(context).width;
    if (Breakpoints.isMobile(width)) {
      return EdgeInsets.all(mobile);
    } else if (Breakpoints.isTablet(width)) {
      return EdgeInsets.all(tablet);
    } else {
      return EdgeInsets.all(web);
    }
  }

  /// 获取自适应水平内边距
  static EdgeInsets horizontalPadding(
    BuildContext context, {
    double mobile = 16,
    double tablet = 24,
    double web = 32,
  }) {
    final width = MediaQuery.sizeOf(context).width;
    if (Breakpoints.isMobile(width)) {
      return EdgeInsets.symmetric(horizontal: mobile);
    } else if (Breakpoints.isTablet(width)) {
      return EdgeInsets.symmetric(horizontal: tablet);
    } else {
      return EdgeInsets.symmetric(horizontal: web);
    }
  }

  /// 获取自适应字体大小
  ///
  /// 手机端使用 screenutil 的 .sp（适配不同手机分辨率）
  /// 平板/Web端直接使用固定值
  static double fontSize(
    BuildContext context,
    double mobileSize, {
    double? tabletSize,
    double? webSize,
  }) {
    final width = MediaQuery.sizeOf(context).width;
    if (Breakpoints.isMobile(width)) {
      return mobileSize;
    } else if (Breakpoints.isTablet(width)) {
      return tabletSize ?? mobileSize;
    } else {
      return webSize ?? tabletSize ?? mobileSize;
    }
  }

  /// 获取自适应间距
  static double spacing(
    BuildContext context,
    double mobileSize, {
    double? tabletSize,
    double? webSize,
  }) {
    final width = MediaQuery.sizeOf(context).width;
    if (Breakpoints.isMobile(width)) {
      return mobileSize;
    } else if (Breakpoints.isTablet(width)) {
      return tabletSize ?? mobileSize;
    } else {
      return webSize ?? tabletSize ?? mobileSize;
    }
  }

  /// 获取网格列数
  static int gridColumns(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    if (Breakpoints.isMobile(width)) return 1;
    if (Breakpoints.isTablet(width)) return 2;
    return 3;
  }

  /// 获取文章列表网格列数
  static int postGridColumns(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    if (Breakpoints.isMobile(width)) return 1;
    if (width < 900) return 2;
    return 3;
  }

  /// 获取分类网格列数
  static int categoryGridColumns(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    if (Breakpoints.isMobile(width)) return 2;
    if (width < 900) return 3;
    return 4;
  }

  /// 判断当前是否为宽屏
  static bool isWideScreen(BuildContext context) {
    return Breakpoints.isWeb(MediaQuery.sizeOf(context).width);
  }

  /// 判断当前是否为手机端
  static bool isMobile(BuildContext context) {
    return Breakpoints.isMobile(MediaQuery.sizeOf(context).width);
  }

  /// 判断当前是否为平板端
  static bool isTablet(BuildContext context) {
    return Breakpoints.isTablet(MediaQuery.sizeOf(context).width);
  }

  /// 获取内容区域最大宽度
  static double contentMaxWidth(BuildContext context) {
    final width = MediaQuery.sizeOf(context).width;
    if (Breakpoints.isWeb(width)) return Breakpoints.webContentMaxWidth;
    if (Breakpoints.isTablet(width)) return width - 32;
    return width;
  }
}
