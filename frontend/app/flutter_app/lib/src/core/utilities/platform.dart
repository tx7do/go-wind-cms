import 'dart:io' show Platform;
import 'package:flutter/foundation.dart' show kIsWeb;

/// 系统平台判断工具集
class PlatformUtils {
  PlatformUtils._(); // Private constructor to prevent instantiation

  /// 是否网页
  static bool get isWeb => _isWeb();

  /// 是否手机系统
  static bool get isMobile => _isMobile();

  /// 是否桌面系统
  static bool get isDesktop => _isDesktop();

  /// 是否Android系统
  static bool get isAndroid => _isAndroid();

  /// 是否IOS系统
  static bool get isIOS => _isIOS();

  /// 是否MacOS系统
  static bool get isMacOS => _isMacOS();

  /// 是否Apple系统
  static bool get isApple => _isMacOS() || _isIOS();

  /// 是否Windows系统
  static bool get isWindows => _isWindows();

  /// 是否Fuchsia系统
  static bool get isFuchsia => _isFuchsia();

  /// 是否Linux系统
  static bool get isLinux => _isLinux();

  /// 是否鸿蒙系统
  static bool get isOHOS => _isOHOS();

  /// 获取平台名
  static String get getPlatform => _getPlatform();

  /// 获取操作系统名
  static String get getOperatingSystem => _getOperatingSystem();

  /// 是否网页
  static bool _isWeb() {
    return kIsWeb == true;
  }

  /// 是否移动设备
  static bool _isMobile() {
    return (_isAndroid() || _isIOS() || _isFuchsia());
  }

  /// 是否桌面系统
  static bool _isDesktop() {
    return (_isLinux() || _isMacOS() || _isWindows());
  }

  static bool _isAndroid() {
    return _isWeb() ? false : Platform.isAndroid;
  }

  static bool _isIOS() {
    return _isWeb() ? false : Platform.isIOS;
  }

  static bool _isMacOS() {
    return _isWeb() ? false : Platform.isMacOS;
  }

  static bool _isWindows() {
    return _isWeb() ? false : Platform.isWindows;
  }

  static bool _isFuchsia() {
    return _isWeb() ? false : Platform.isFuchsia;
  }

  static bool _isLinux() {
    return _isWeb() ? false : Platform.isLinux;
  }

  static bool _isOHOS() {
    return _isWeb() ? false : Platform.operatingSystem == "ohos";
  }

  static String _getPlatform() {
    if (kIsWeb) return 'web';
    return Platform.operatingSystem;
  }

  static String _getOperatingSystem() {
    if (_isWindows()) {
      return 'windows';
    } else if (_isLinux()) {
      return 'linux';
    } else if (_isAndroid()) {
      return 'android';
    } else if (_isMacOS()) {
      return 'macos';
    } else if (_isIOS()) {
      return 'ios';
    } else if (_isFuchsia()) {
      return 'fuchsia';
    } else if (_isOHOS()) {
      return 'ohos';
    } else {
      return 'unknown';
    }
  }
}
