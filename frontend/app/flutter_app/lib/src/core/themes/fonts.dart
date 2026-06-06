import 'package:flutter_app/src/core/utilities/platform.dart';

import 'dynamic_font.dart';

class AppFont {
  const AppFont._();

  static const systemFont = "system-font";
  static bool systemFontLoaded = false;

  /// Emoji
  static final emojiFontFamily = [
    // iOS and MacOs.
    'Apple Color Emoji',
    // Android, ChromeOS, Ubuntu and some other Linux distros.
    'Noto Color Emoji',
    // Windows.
    'Segoe UI Emoji',
  ];

  /// iOS & macOS
  static const List<String> appleFontFamily = [
    // 英文字体
    '.SF UI Text',
    '.SF UI Display',
    '.AppleSystemUIFont',
    // 中文字体
    'PingFang SC', 'PingFang TC', 'PingFang HK',
  ];

  /// Android
  static const List<String> androidFontFamily = [
    // 英文字体
    'Roboto',
    // 中文字体
    'Source Han Sans',
    'Noto',
  ];

  /// windows
  static const List<String> windowsFontFamily = [
    'Microsoft YaHei',
  ];

  /// 小米
  static const List<String> xiaomiFontFamily = [
    'miui',
    'mipro',
  ];

  /// VIVO
  static final vivoSystemFont = DynamicFont.file(
    fontFamily: systemFont,
    filepath: '/system/fonts/DroidSansFallbackMonster.ttf',
  );

  /// 荣耀
  static final honorSystemFont = DynamicFont.file(
    fontFamily: systemFont,
    filepath: '/system/fonts/DroidSansChinese.ttf',
  );

  static List<String> get fontFamilyFallback {
    return [
      systemFont,
      if (PlatformUtils.isApple) ...appleFontFamily,
      if (PlatformUtils.isAndroid) ...xiaomiFontFamily,
      if (PlatformUtils.isWindows) ...windowsFontFamily,
      ...emojiFontFamily,
    ];
  }
}
