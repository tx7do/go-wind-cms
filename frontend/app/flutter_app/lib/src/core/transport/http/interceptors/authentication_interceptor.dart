import 'dart:async';

import 'package:dio/dio.dart';

import 'package:flutter_app/src/core/constants/index.dart' show AppRoutePath;
import 'package:flutter_app/src/core/services/base_service.dart'
    show BaseService;
import 'package:flutter_app/src/core/utilities/logger.dart' show fatal;

/// 认证服务接口，定义了获取和刷新令牌的方法
abstract class AuthService extends BaseService {
  /// 获取当前访问令牌
  String? getAccessToken();

  /// 获取刷新令牌
  String? getRefreshToken();

  /// 刷新访问令牌
  /// 返回新的访问令牌，失败时返回null
  Future<String?> refreshToken();

  /// 认证失败处理
  authenticationFailed();
}

/// 认证拦截器
class AuthenticationInterceptor extends Interceptor {
  final AuthService _authService;
  final bool _autoRefreshToken;
  late Completer _refreshLock = Completer();

  /// 创建认证拦截器实例
  /// [authService] - 认证服务实现
  /// [autoRefreshToken] - 是否自动刷新令牌，默认为true
  AuthenticationInterceptor({
    required AuthService authService,
    bool autoRefreshToken = true,
  }) : _authService = authService,
       _autoRefreshToken = autoRefreshToken;

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    if (options.path != AppRoutePath.login) {
      final token = _authService.getAccessToken();

      if (token != null && token.isNotEmpty) {
        options.headers['Authorization'] = _makeBearerToken(accessToken: token);
      }
    }

    return handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) async {
    return handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response == null) {
      return;
    }

    // 如果不需要自动刷新令牌或不是401错误，则直接传递错误
    if (!_autoRefreshToken || err.response?.statusCode != 401) {
      return handler.next(err);
    }

    try {
      // 尝试刷新令牌并重新发送请求
      final newToken = await _refreshToken();
      if (newToken == null) {
        // 刷新令牌失败，清除令牌并返回错误
        // await _authService.authenticationFailed();
        return handler.next(err);
      }

      // 使用新令牌重试请求
      final options = err.requestOptions;
      options.headers['Authorization'] = _makeBearerToken(
        accessToken: newToken,
      );

      // 清除错误响应
      // err.response = null;
      final newException = DioException(
        requestOptions: err.requestOptions,
        error: err.error,
        type: err.type,
        // 不提供 response 参数，默认值为 null
      );
      err = newException;

      // 重新发送请求
      final response = await Dio().fetch(options);
      return handler.resolve(response);
    } catch (e) {
      fatal('Error refreshing token: $e');
      // 发生异常时清除令牌并传递错误
      await _authService.authenticationFailed();
      return handler.next(err);
    }
  }

  /// 刷新令牌的方法，使用锁机制防止并发刷新
  Future<String?> _refreshToken() async {
    // 如果已经有刷新请求在进行，等待它完成
    if (!_refreshLock.isCompleted) {
      await _refreshLock.future;
      return _authService.getAccessToken();
    }

    // 创建新的锁
    final Completer newLock = Completer();
    _refreshLock.complete();
    _refreshLock = newLock;

    try {
      // 尝试刷新令牌
      final newToken = await _authService.refreshToken();
      if (newToken != null) {
        _refreshLock.complete();
        return newToken;
      }

      _refreshLock.complete();
      return null;
    } catch (e) {
      _refreshLock.complete();
      rethrow;
    }
  }

  /// 创建Bearer令牌
  _makeBearerToken({String? accessToken}) {
    return "Bearer ${accessToken ?? ""}";
  }
}
