import 'dart:async';

import 'package:get_it/get_it.dart';

import 'http/index.dart';

/// 初始化
Future<void> init() async {
  final getIt = GetIt.instance;
  getIt.registerLazySingleton<Dio>(() => createDio());
}
