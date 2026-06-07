import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/constants/index.dart';

import 'dark_theme.dart' show getDarkTheme;
import 'light_theme.dart' show getLightTheme;

export 'dark_theme.dart' show getDarkTheme, DarkColor;
export 'light_theme.dart' show getLightTheme, LightColor, kDefaultSeedColor;

/// 获取当前的主题数据
ThemeData getCurrentTheme(ThemeMode themeMode, {Color? seedColor}) {
  switch (themeMode) {
    case ThemeMode.light:
      return getLightTheme(seedColor: seedColor);
    case ThemeMode.dark:
      return getDarkTheme(seedColor: seedColor);
    case ThemeMode.system:
      // 注意：此函数已废弃，实际由 MaterialApp 的 theme/darkTheme/themeMode 自动处理
      return getLightTheme(seedColor: seedColor);
  }
}

class AppStyle {
  static TextStyle navAppBarTitleStyle = TextStyle(
    color:
        Theme.of(Global.applicationKey.currentContext!).colorScheme.onPrimary,
    fontSize: 16.0,
    fontWeight: FontWeight.w600,
  );
}
