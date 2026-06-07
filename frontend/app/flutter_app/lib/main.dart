import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_native_splash/flutter_native_splash.dart';

import 'package:flutter_app/src/app.dart' show CMSApp;
import 'package:flutter_app/src/init.dart' show init;
import 'package:flutter_app/src/core/themes/index.dart' show AppThemeCubit;

/// 程序入口
Future<void> main() async {
  final widgetsBinding = WidgetsFlutterBinding.ensureInitialized();
  FlutterNativeSplash.preserve(widgetsBinding: widgetsBinding);

  await init();

  // 初始化完成后移除原生闪屏
  FlutterNativeSplash.remove();

  run();
}

void run() {
  // ignore: missing_provider_scope
  runApp(
    MultiBlocProvider(
      providers: [BlocProvider(create: (_) => AppThemeCubit())],
      child: const CMSApp(),
    ),
  );
}
