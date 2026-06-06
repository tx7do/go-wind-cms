import 'dart:async';

import 'package:dio/dio.dart' show Dio;
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_app/generated/api/rest_client.dart' show RestClient;

import 'package:flutter_app/src/core/config/environments.dart';
import 'package:flutter_app/src/core/repositories/init.dart' as repos;
import 'package:flutter_app/src/core/transport/init.dart' as transport;
import 'package:flutter_app/src/core/widgets/error_page.dart';
import 'package:get_it/get_it.dart' show GetIt;

import 'init_thirdparty_plugins.dart';

/// 应用初始化
Future<void> init() async {
  WidgetsFlutterBinding.ensureInitialized();

  await Environments.init();

  await initThirdPartyPlugins();

  _initTransport();

  await repos.init();

  _initErrorWidget();
}

/// 初始化传输层
void _initTransport() {
  transport.init();

  final getIt = GetIt.instance;
  getIt.registerLazySingleton<RestClient>(
    () => RestClient(GetIt.instance<Dio>()),
  );
}

/// 自定义报错页面
void _initErrorWidget() {
  ErrorWidget.builder = (FlutterErrorDetails details) {
    debugPrint(details.toString());

    if (kDebugMode) {
      return ErrorWidget(details.exception);
    }

    return CustomErrorWidget(errorMessage: details.exceptionAsString());
  };
}

/// 清理
void globalDispose() {}
