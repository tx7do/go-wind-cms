import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'const.dart';
import 'fonts.dart';
import 'snackbar_color.dart' show SnackBarColor;
import 'app_theme_extension.dart';

/// 默认主题色
const Color kDefaultSeedColor = Color(0xFF3A7CA5);

class LightColor {
  static const Color bgColor = Color.fromRGBO(245, 247, 250, 1.0);
  static const Color primaryColor = Color.fromRGBO(58, 124, 165, 1.0);
  static const Color onPrimaryColor = Color.fromRGBO(255, 255, 255, 1.0);
  static const Color inputTextColor = Color.fromRGBO(10, 25, 12, 1.0);
  static const Color inputFillColor = Color.fromRGBO(255, 255, 255, 1.0);
  static const Color cursorColor = Color(0xFF3A7CA5);
  static const Color accentColor = Color.fromRGBO(255, 152, 67, 1.0);
}

ColorScheme _buildLightColorScheme(Color seedColor) => ColorScheme.fromSeed(
  brightness: Brightness.light,
  seedColor: seedColor,

  /// 主色调，用于突出显示和主要操作
  primary: LightColor.primaryColor,

  /// 主色调上的文字或图标颜色
  onPrimary: LightColor.onPrimaryColor,

  /// 主色调的容器背景色，更淡一些以提供对比
  primaryContainer: const Color.fromRGBO(210, 232, 245, 1.0),

  /// 主色调容器上的文字或图标颜色
  onPrimaryContainer: const Color.fromRGBO(20, 60, 85, 1.0),

  /// 背景颜色
  surface: const Color.fromARGB(255, 255, 255, 255),

  // 背景颜色上的文字或图标颜色
  onSurface: const Color.fromRGBO(30, 30, 30, 1.0),

  /// 次级容器
  secondaryContainer: const Color.fromRGBO(255, 237, 220, 1.0),
  onSecondaryContainer: const Color.fromRGBO(120, 60, 10, 1.0),

  /// 错误状态的颜色
  error: Colors.red,

  /// 错误状态上的文字或图标颜色
  onError: Colors.white,
);

/// 标题栏主题
final _appBarTheme = AppBarTheme(
  backgroundColor: const Color.fromARGB(255, 255, 255, 255),
  titleSpacing: kDefaultPadding,
  elevation: 0,
  scrolledUnderElevation: 1,
  titleTextStyle: TextStyle(
    color: const Color.fromRGBO(30, 30, 30, 1.0),
    fontSize: 20.sp,
    fontWeight: FontWeight.w700,
  ),
  iconTheme: const IconThemeData(color: Color.fromRGBO(30, 30, 30, 1.0)),
);

/// 凸出按钮主题
final _elevatedButtonTheme = ElevatedButtonThemeData(
  style: ButtonStyle(
    elevation: WidgetStateProperty.all(0),
    backgroundColor: WidgetStateProperty.all(LightColor.primaryColor),
    foregroundColor: WidgetStateProperty.all(Colors.white),
    textStyle: WidgetStateProperty.all(
      TextStyle(
        fontSize: 16.sp,
        fontWeight: FontWeight.w600,
      ),
    ),
    shape: WidgetStateProperty.all(
      RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
    ),
  ),
);

/// 文本
const _textSelectionThemeData = TextSelectionThemeData(
  selectionColor: Color.fromARGB(255, 111, 181, 210),
  selectionHandleColor: Color.fromARGB(255, 111, 181, 210),
  cursorColor: LightColor.cursorColor,
);

/// 图标
const _iconTheme = IconThemeData(
  color: Color.fromRGBO(30, 30, 30, 1.0),
);

///
const _snackBarColor = SnackBarColor(
  error: Color.fromRGBO(211, 45, 81, 1),
  info: Color.fromRGBO(26, 181, 148, 1),
  warning: Color.fromRGBO(26, 181, 148, 1),
  success: Color.fromRGBO(26, 181, 148, 1),
);

/// 顶部搜索栏
final _searchBarTheme = SearchBarThemeData(
  backgroundColor: WidgetStateProperty.all(Colors.white),
  elevation: WidgetStateProperty.all(0),
  shape: WidgetStateProperty.all(
    RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
  ),
);

/// 弹出菜单
const _popupMenuTheme = PopupMenuThemeData(
  color: Colors.white,
  iconColor: LightColor.primaryColor,
  elevation: 4,
);

/// 分割线
const _dividerTheme = DividerThemeData(
  color: Color.fromRGBO(230, 230, 235, 1.0),
);

/// 构建 Scrollbar 主题
ScrollbarThemeData _buildScrollbarTheme(Color primaryColor) => ScrollbarThemeData(
  thumbColor: WidgetStateProperty.resolveWith((states) {
    if (states.contains(WidgetState.hovered)) {
      return primaryColor.withAlpha((0.5 * 255).round());
    }
    return primaryColor.withAlpha((0.25 * 255).round());
  }),
  radius: const Radius.circular(4),
  thickness: WidgetStateProperty.all(6),
  thumbVisibility: WidgetStateProperty.all(true),
);

///
const _inputDecorationTheme = InputDecorationTheme(
    labelStyle: TextStyle(
  color: LightColor.primaryColor,
));

const _textTheme = TextTheme();

///
ThemeData _buildLightTheme(Color seedColor) {
  final colorScheme = _buildLightColorScheme(seedColor);
  return ThemeData(
    platform: TargetPlatform.iOS,
    useMaterial3: true,
    fontFamily: kDefaultFontFamily,
    fontFamilyFallback: AppFont.fontFamilyFallback,
    textTheme: _textTheme,
    colorScheme: colorScheme,
    scaffoldBackgroundColor: LightColor.bgColor,
    primaryColor: colorScheme.primary,
    appBarTheme: _appBarTheme,
    textSelectionTheme: _textSelectionThemeData,
    elevatedButtonTheme: _elevatedButtonTheme,
    iconTheme: _iconTheme,
    searchBarTheme: _searchBarTheme,
    popupMenuTheme: _popupMenuTheme,
    dividerTheme: _dividerTheme,
    scrollbarTheme: _buildScrollbarTheme(colorScheme.primary),
    inputDecorationTheme: _inputDecorationTheme,
    cardTheme: CardThemeData(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.all(Radius.circular(14)),
      ),
      color: Colors.white,
    ),
    chipTheme: const ChipThemeData(
      side: BorderSide.none,
    ),
    navigationBarTheme: NavigationBarThemeData(
      labelTextStyle: WidgetStateProperty.resolveWith((states) {
        if (states.contains(WidgetState.selected)) {
          return TextStyle(
            fontSize: 12.sp,
            fontWeight: FontWeight.w600,
            color: colorScheme.primary,
          );
        }
        return TextStyle(
          fontSize: 12.sp,
          fontWeight: FontWeight.w400,
        );
      }),
      iconTheme: WidgetStateProperty.resolveWith((states) {
        if (states.contains(WidgetState.selected)) {
          return IconThemeData(
            size: 24.sp,
            color: colorScheme.primary,
          );
        }
        return IconThemeData(size: 24.sp);
      }),
    ),
    extensions: const <ThemeExtension<dynamic>>[
      AppThemeExtension(snackBarColor: _snackBarColor),
    ]);
}

/// 获取光亮系主题
/// [seedColor] 种子颜色，用于动态生成配色方案
ThemeData getLightTheme({Color? seedColor}) {
  final color = seedColor ?? kDefaultSeedColor;
  return _buildLightTheme(color);
}
