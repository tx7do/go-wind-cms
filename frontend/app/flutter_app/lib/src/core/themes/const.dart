import 'package:flutter_app/src/core/utilities/platform.dart';

const double kDefaultPadding = 16.0;

const String kDefaultFontFamily = 'Noto Sans SC';

const double kMainSpace = 10.0;

double kMainLineWidth(bool isDarkMode) => isDarkMode ? 0.5 : 1.0;

/// 获取默认的字体
String getDefaultFontFamily() {
  if (PlatformUtils.isApple) {
    return 'PingFang SC';
  } else if (PlatformUtils.isAndroid) {
    return 'Roboto';
  } else if (PlatformUtils.isWindows) {
    return 'Microsoft YaHei';
  } else if (PlatformUtils.isLinux) {
    return 'Roboto';
  } else if (PlatformUtils.isWeb) {
    return 'Noto Sans SC';
  } else {
    return kDefaultFontFamily;
  }
}
