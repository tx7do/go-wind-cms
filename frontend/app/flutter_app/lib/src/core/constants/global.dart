import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/utilities/logger.dart';

/// 全局变量
class Global {
  Global._(); // Private constructor to prevent instantiation
  ///
  static final GlobalKey<NavigatorState> applicationKey = GlobalKey();

  ///
  static final GlobalKey<ScaffoldMessengerState> scaffoldMessengerKey =
      GlobalKey();

  /// 日志记录器
  static final Logger logger = Logger('im');
}
