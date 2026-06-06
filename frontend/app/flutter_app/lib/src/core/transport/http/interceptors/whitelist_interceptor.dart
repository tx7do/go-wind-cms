import 'dart:async' show FutureOr;

import 'package:dio/dio.dart';

import 'package:flutter_app/src/core/utilities/logger.dart' show debug;

/// 白名单拦截器
/// 用于选择性地对请求应用拦截逻辑，跳过白名单中的路径
class WhitelistInterceptor extends Interceptor {
  /// 白名单路径列表
  final List<String> _whitelistPaths;

  /// 拦截器名称，用于日志记录
  final String _name;

  /// 拦截器执行函数
  final FutureOr<void> Function(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) _interceptor;

  /// 创建白名单拦截器
  /// [whitelistPaths] - 白名单路径列表
  /// [name] - 拦截器名称，用于日志记录
  /// [interceptor] - 拦截器执行函数
  WhitelistInterceptor({
    required List<String> whitelistPaths,
    String name = 'WhitelistInterceptor',
    required FutureOr<void> Function(
      RequestOptions options,
      RequestInterceptorHandler handler,
    ) interceptor,
  })  : _whitelistPaths = whitelistPaths,
        _name = name,
        _interceptor = interceptor;

  @override
  Future<void> onRequest(
      RequestOptions options, RequestInterceptorHandler handler) async {
    // 检查请求路径是否在白名单中
    if (_isPathInWhitelist(options.path)) {
      debug('[$_name] 跳过白名单路径: ${options.path}');
      return handler.next(options);
    }

    debug('[$_name] 处理非白名单路径: ${options.path}');

    // 执行拦截逻辑
    await _interceptor(options, handler);
  }

  /// 检查路径是否在白名单中
  bool _isPathInWhitelist(String path) {
    // 规范化路径，移除可能的查询参数
    final normalizedPath = path.split('?').first;

    // 检查是否匹配任何白名单路径
    return _whitelistPaths.any((whitelistPath) {
      // 精确匹配
      if (whitelistPath == normalizedPath) return true;

      // 支持通配符匹配（*）
      if (whitelistPath.contains('*')) {
        final pattern = whitelistPath.replaceAll('*', '.*');
        final regex = RegExp('^$pattern\$');
        return regex.hasMatch(normalizedPath);
      }

      return false;
    });
  }
}
