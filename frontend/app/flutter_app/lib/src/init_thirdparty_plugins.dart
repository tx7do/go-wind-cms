import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'package:get_it/get_it.dart' show GetIt;
import 'package:shared_preferences/shared_preferences.dart'
    show SharedPreferences;

import 'package:wakelock_plus/wakelock_plus.dart';

import 'package:flutter_app/src/core/utilities/platform.dart';

/// 注入插件
Future<void> injectPlugins() async {
  final getIt = GetIt.instance;

  final prefs = await SharedPreferences.getInstance();
  getIt.registerSingleton<SharedPreferences>(prefs);
}

/// 初始化第三方插件
Future<void> initThirdPartyPlugins() async {
  await injectPlugins();

  // 强制竖屏
  SystemChrome.setPreferredOrientations([DeviceOrientation.portraitUp]);

  // 启用和切换屏幕唤醒锁，以防止屏幕自动关闭。
  if (PlatformUtils.isMobile) {
    WakelockPlus.enable();
  }

  /// Android状态栏透明
  if (PlatformUtils.isAndroid) {
    SystemUiOverlayStyle systemUiOverlayStyle =
        const SystemUiOverlayStyle(statusBarColor: Colors.transparent);
    SystemChrome.setSystemUIOverlayStyle(systemUiOverlayStyle);
  }

  // 解决使用自签证书报错问题
  // io.HttpOverrides.global = GlobalHttpOverrides();

  // 注册视频播放器插件
  // fvp.registerWith();
}
