import 'package:dio/dio.dart';
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/src/core/config/environments.dart';

/// 配置选项
void _configureOptions(Dio dio) {
  // debug('_configureOptions');

  dio.options.baseUrl = Environments.apiBaseUrl;
  dio.options.connectTimeout = Environments.connectionTimeout;
  dio.options.receiveTimeout = Environments.receiveTimeout;
  dio.options.responseType = ResponseType.json;
  dio.options.contentType = Headers.jsonContentType;
}

/// 注册默认拦截器
void _configureInterceptors(Dio dio) {
  // 日志（调试时取消注释）
  // dio.interceptors.add(LogInterceptor());
  // // jwt认证
  // dio.interceptors.add(AuthenticationInterceptor(
  //     authService: authService, autoRefreshToken: false));
}

/// 注册拦截器
void registerInterceptor(Interceptor interceptor) {
  final dio = GetIt.instance<Dio>();
  dio.interceptors.add(interceptor);
}

Dio createDio() {
  final Dio dio = Dio();

  _configureOptions(dio);
  _configureInterceptors(dio);

  return dio;
}
