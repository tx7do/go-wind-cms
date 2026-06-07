import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'const.dart';
import 'fonts.dart';
import 'snackbar_color.dart' show SnackBarColor;
import 'app_theme_extension.dart';

class DarkColor {
  static const Color bgColor = Color.fromRGBO(18, 18, 18, 1.0);
  static const Color surfaceColor = Color.fromRGBO(28, 28, 30, 1.0);
  static const Color cardColor = Color.fromRGBO(35, 35, 38, 1.0);
  static const Color onSurface = Color.fromRGBO(232, 232, 234, 1.0);
  static const Color cursorColor = Color(0xFF3A7CA5);
}

ColorScheme _buildDarkColorScheme(Color seedColor) => ColorScheme.fromSeed(
  brightness: Brightness.dark,
  seedColor: seedColor,

  /// 主色调，用于突出显示和主要操作
  primary: seedColor,

  /// 主色调上的文字或图标颜色
  onPrimary: Colors.white,

  /// 主色调的容器背景色
  primaryContainer: Color.alphaBlend(
    seedColor.withAlpha((0.3 * 255).round()),
    DarkColor.surfaceColor,
  ),

  /// 主色调容器上的文字或图标颜色
  onPrimaryContainer: Color.alphaBlend(
    seedColor.withAlpha((0.15 * 255).round()),
    Colors.white,
  ),

  /// 背景颜色
  surface: DarkColor.surfaceColor,

  /// 背景颜色上的文字或图标颜色
  onSurface: DarkColor.onSurface,

  /// 次级容器
  secondaryContainer: Color.alphaBlend(
    seedColor.withAlpha((0.15 * 255).round()),
    DarkColor.surfaceColor,
  ),
  onSecondaryContainer: Color.alphaBlend(
    seedColor.withAlpha((0.1 * 255).round()),
    Colors.white,
  ),

  /// 错误状态的颜色
  error: const Color.fromRGBO(242, 80, 80, 1.0),

  /// 错误状态上的文字或图标颜色
  onError: Colors.white,
);

/// 标题栏主题
AppBarTheme _buildAppBarTheme(ColorScheme colorScheme) => AppBarTheme(
  backgroundColor: DarkColor.surfaceColor,
  titleSpacing: kDefaultPadding,
  elevation: 0,
  scrolledUnderElevation: 1,
  titleTextStyle: TextStyle(
    color: DarkColor.onSurface,
    fontSize: 20.sp,
    fontWeight: FontWeight.w700,
  ),
  iconTheme: IconThemeData(color: DarkColor.onSurface),
);

/// 凸出按钮主题
ElevatedButtonThemeData _buildElevatedButtonTheme(ColorScheme colorScheme) =>
    ElevatedButtonThemeData(
      style: ButtonStyle(
        elevation: WidgetStateProperty.all(0),
        backgroundColor: WidgetStateProperty.all(colorScheme.primary),
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

/// 文本选择
TextSelectionThemeData _buildTextSelectionTheme(Color seedColor) =>
    TextSelectionThemeData(
      selectionColor: seedColor.withAlpha((0.4 * 255).round()),
      selectionHandleColor: seedColor,
      cursorColor: DarkColor.cursorColor,
    );

/// 图标
const _iconTheme = IconThemeData(
  color: DarkColor.onSurface,
);

///
const _snackBarColor = SnackBarColor(
  error: Color.fromRGBO(242, 80, 80, 1),
  info: Color.fromRGBO(26, 181, 148, 1),
  warning: Color.fromRGBO(26, 181, 148, 1),
  success: Color.fromRGBO(26, 181, 148, 1),
);

/// 顶部搜索栏
final _searchBarTheme = SearchBarThemeData(
  backgroundColor: WidgetStateProperty.all(DarkColor.cardColor),
  elevation: WidgetStateProperty.all(0),
  shape: WidgetStateProperty.all(
    RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
  ),
);

/// 弹出菜单
PopupMenuThemeData _buildPopupMenuTheme(Color seedColor) => PopupMenuThemeData(
  color: DarkColor.cardColor,
  iconColor: seedColor,
  elevation: 4,
);

/// 分割线
const _dividerTheme = DividerThemeData(
  color: Color.fromRGBO(55, 55, 58, 1.0),
);

/// 构建 Scrollbar 主题
ScrollbarThemeData _buildScrollbarTheme(Color seedColor) => ScrollbarThemeData(
  thumbColor: WidgetStateProperty.resolveWith((states) {
    if (states.contains(WidgetState.hovered)) {
      return seedColor.withAlpha((0.5 * 255).round());
    }
    return seedColor.withAlpha((0.3 * 255).round());
  }),
  radius: const Radius.circular(4),
  thickness: WidgetStateProperty.all(6),
  thumbVisibility: WidgetStateProperty.all(true),
);

///
InputDecorationTheme _buildInputDecorationTheme(Color seedColor) =>
    const InputDecorationTheme(
      labelStyle: TextStyle(
        color: DarkColor.onSurface,
      ),
    );

const _textTheme = TextTheme();

///
ThemeData _buildDarkTheme(Color seedColor) {
  final colorScheme = _buildDarkColorScheme(seedColor);
  return ThemeData(
    platform: TargetPlatform.iOS,
    useMaterial3: true,
    fontFamily: kDefaultFontFamily,
    fontFamilyFallback: AppFont.fontFamilyFallback,
    textTheme: _textTheme,
    colorScheme: colorScheme,
    scaffoldBackgroundColor: DarkColor.bgColor,
    primaryColor: colorScheme.primary,
    appBarTheme: _buildAppBarTheme(colorScheme),
    textSelectionTheme: _buildTextSelectionTheme(seedColor),
    elevatedButtonTheme: _buildElevatedButtonTheme(colorScheme),
    iconTheme: _iconTheme,
    searchBarTheme: _searchBarTheme,
    popupMenuTheme: _buildPopupMenuTheme(seedColor),
    dividerTheme: _dividerTheme,
    scrollbarTheme: _buildScrollbarTheme(seedColor),
    inputDecorationTheme: _buildInputDecorationTheme(seedColor),
    cardTheme: CardThemeData(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.all(Radius.circular(14)),
      ),
      color: DarkColor.cardColor,
    ),
    chipTheme: const ChipThemeData(
      side: BorderSide.none,
    ),
    navigationBarTheme: NavigationBarThemeData(
      backgroundColor: DarkColor.surfaceColor,
      indicatorColor: colorScheme.primaryContainer,
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
          color: DarkColor.onSurface.withAlpha((0.6 * 255).round()),
        );
      }),
      iconTheme: WidgetStateProperty.resolveWith((states) {
        if (states.contains(WidgetState.selected)) {
          return IconThemeData(
            size: 24.sp,
            color: colorScheme.primary,
          );
        }
        return IconThemeData(
          size: 24.sp,
          color: DarkColor.onSurface.withAlpha((0.6 * 255).round()),
        );
      }),
    ),
    extensions: const <ThemeExtension<dynamic>>[
      AppThemeExtension(snackBarColor: _snackBarColor),
    ],
  );
}

/// 获取暗黑系主题
/// [seedColor] 种子颜色，用于动态生成配色方案
ThemeData getDarkTheme({Color? seedColor}) {
  final color = seedColor ?? const Color(0xFF3A7CA5);
  return _buildDarkTheme(color);
}
