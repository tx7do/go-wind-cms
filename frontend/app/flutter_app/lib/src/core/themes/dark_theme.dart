import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'const.dart';
import 'fonts.dart';
import 'utils.dart';
import 'snackbar_color.dart' show SnackBarColor;
import 'app_theme_extension.dart';

class DarkColor {
  static const Color bgColor = Color.fromRGBO(40, 40, 40, 1.0);
  static const Color primaryColor = Color.fromRGBO(26, 26, 26, 1);
  static const Color onPrimaryColor = Color.fromRGBO(208, 208, 208, 1.0);
  static const Color inputTextColor = Color.fromRGBO(212, 212, 212, 1.0);
  static const Color inputFillColor = Color.fromRGBO(44, 44, 44, 1.0);
  static const Color cursorColor = Colors.green;
}

final ColorScheme _darkColorScheme = ColorScheme.fromSeed(
  brightness: Brightness.dark,
  seedColor: Colors.black,
  primary: DarkColor.primaryColor,
  // 主色调保持不变，以确保品牌一致性
  onPrimary: DarkColor.onPrimaryColor,
  // 主色调上的文字颜色改为黑色或深色
  primaryContainer: DarkColor.primaryColor.withAlpha((0.8 * 255).round()),
  // 主色调的容器背景色，更淡一些以提供对比
  onPrimaryContainer: Colors.black54,
  // 主色调容器上的文字或图标颜色
  surface: DarkColor.bgColor,
  // 背景颜色改为更深的黑色调
  onSurface: Colors.white70,
  // 背景颜色上的文字颜色保持一定的透明度
  error: Colors.red,
  // 错误颜色
  onError: Colors.white, // 错误颜色上的文字颜色改为黑色
);

final _appBarTheme = AppBarTheme(
  backgroundColor: Colors.black,
  titleSpacing: kDefaultPadding,
  elevation: 2,
  titleTextStyle: TextStyle(
    color: Color.fromARGB(255, 255, 255, 255),
    fontSize: 20.sp,
  ),
  iconTheme: IconThemeData(color: Color.fromARGB(255, 255, 255, 255)),
);

final _elevatedButtonTheme = ElevatedButtonThemeData(
  style: ButtonStyle(
    elevation: WidgetStateProperty.all(0),
    backgroundColor: WidgetStateProperty.all(Colors.black),
    textStyle: WidgetStateProperty.all(
      TextStyle(
        fontSize: 18.sp,
        color: Color.fromARGB(255, 255, 255, 255),
      ),
    ),
  ),
);

/// 文本
const _textSelectionThemeData = TextSelectionThemeData(
  selectionColor: Color.fromARGB(255, 111, 181, 210),
  selectionHandleColor: Color.fromARGB(255, 111, 181, 210),
  cursorColor: DarkColor.cursorColor,
);

/// 图标
const _iconTheme = IconThemeData(
  color: DarkColor.primaryColor,
);

const _snackBarColor = SnackBarColor(
  error: Color.fromRGBO(211, 45, 81, 1),
  info: Color.fromRGBO(5, 165, 133, 1),
  warning: Color.fromRGBO(5, 165, 133, 1),
  success: Color.fromRGBO(5, 165, 133, 1),
);

/// 搜索栏
final _searchBarTheme = SearchBarThemeData(
  backgroundColor: WidgetStateProperty.all(DarkColor.primaryColor),
);

/// 弹出菜单
const _popupMenuTheme = PopupMenuThemeData(
  color: Color.fromRGBO(40, 40, 40, 0.8),
  iconColor: DarkColor.primaryColor,
);

/// 分割线
const _dividerTheme = DividerThemeData(
  color: Color.fromRGBO(45, 45, 45, 0.9),
);

/// 滚动条主题（暗色模式）
final _scrollbarTheme = ScrollbarThemeData(
  thumbColor: WidgetStateProperty.resolveWith((states) {
    if (states.contains(WidgetState.hovered)) {
      return Colors.white.withAlpha((0.4 * 255).round());
    }
    return Colors.white.withAlpha((0.2 * 255).round());
  }),
  radius: const Radius.circular(4),
  thickness: WidgetStateProperty.all(6),
  thumbVisibility: WidgetStateProperty.all(true),
);

const _textTheme = TextTheme();

const _inputDecorationTheme = InputDecorationTheme(
    labelStyle: TextStyle(
  color: Colors.blue,
));

final _darkTheme = ThemeData(
  platform: TargetPlatform.iOS,
  useMaterial3: true,
  colorScheme: _darkColorScheme,
  fontFamily: kDefaultFontFamily,
  fontFamilyFallback: AppFont.fontFamilyFallback,
  textTheme: _textTheme,
  scaffoldBackgroundColor: DarkColor.bgColor,
  primaryColor: Colors.black,
  primarySwatch: createMaterialColor(const Color(0xFF223344)),
  appBarTheme: _appBarTheme,
  textSelectionTheme: _textSelectionThemeData,
  elevatedButtonTheme: _elevatedButtonTheme,
  iconTheme: _iconTheme,
  searchBarTheme: _searchBarTheme,
  popupMenuTheme: _popupMenuTheme,
  dividerTheme: _dividerTheme,
  scrollbarTheme: _scrollbarTheme,
  inputDecorationTheme: _inputDecorationTheme,
  extensions: const <ThemeExtension<dynamic>>[
    AppThemeExtension(snackBarColor: _snackBarColor),
  ],
);

/// 获取暗黑系主题
/// [seedColor] 种子颜色，用于动态生成配色方案
ThemeData getDarkTheme({Color? seedColor}) {
  return _darkTheme;
}
