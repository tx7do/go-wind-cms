import 'package:flutter/foundation.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

import 'package:flutter_app/src/core/utilities/convert.dart';

/// 环境变量
class Environments {
  Environments._(); // Private constructor to prevent instantiation

  /// 初始化
  static init() async {
    const String fileName = kDebugMode ? '.dev.env' : '.env';
    await dotenv.load(fileName: fileName);
  }

  /// NTP远程服务地址
  static String get ntpHost => dotenv.get('NTP_HOST');

  /// API基本链接地址
  static String get apiBaseUrl => dotenv.get('API_BASE_URL');

  /// SSE基本链接地址
  static String get sseUrl => dotenv.get('SSE_URL');

  /// 连接超时时间
  static Duration get connectionTimeout =>
      secondStringToDuration(dotenv.get('CONNECTION_TIMEOUT'));

  /// 接收数据超时时间
  static Duration get receiveTimeout =>
      secondStringToDuration(dotenv.get('RECEIVE_TIMEOUT'));

  /// IOS应用ID
  static String get iosAppId => dotenv.get('IOS_APP_ID');

  /// AES密钥
  static String get aesKey => dotenv.get('AES_KEY');

  /// Sentry异常监控DSN
  static String get sentryDSN => dotenv.get('SENTRY_DSN');

  /// 能否创建一个简单的WebView服务器
  static bool get enableInAppLocalhostServer =>
      parseBool(dotenv.get('ENABLE_IN_APP_LOCALHOST_SERVER'));
}
