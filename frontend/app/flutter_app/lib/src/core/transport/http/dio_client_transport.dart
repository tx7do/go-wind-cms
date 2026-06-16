import 'dart:async';

import 'package:dio/dio.dart';
import 'package:get_it/get_it.dart';

import 'package:flutter_app/generated/api/transport.dart';
import 'package:flutter_app/src/core/utilities/logger.dart';

/// 基于 Dio 的 [ClientTransport] 实现。
///
/// 将 protoc-gen-dart-http 生成的统一 transport 接口桥接到项目已有的 Dio 实例，
/// 复用 baseUrl、拦截器（认证、请求、响应）等配置。
class DioClientTransport implements ClientTransport {
  final Dio _dio;
  final Logger _logger = Logger('DioClientTransport');

  DioClientTransport({Dio? dio}) : _dio = dio ?? GetIt.instance<Dio>();

  @override
  Future<dynamic> unary(
    String path,
    String method,
    String? body,
    TransportMeta meta, {
    Map<String, String>? headers,
  }) async {
    _logger.d('unary ${meta.service}.${meta.method} $method $path');

    final options = Options(
      method: method,
      headers: headers,
      responseType: ResponseType.json,
    );

    Object? data;
    if (body != null) {
      data = body;
    }

    final response = await _dio.request<dynamic>(
      path,
      data: data,
      options: options,
    );

    return response.data;
  }

  @override
  Stream<Map<String, dynamic>> serverStream(
    String path,
    TransportMeta meta, {
    Map<String, String>? headers,
  }) {
    throw UnimplementedError(
      'serverStream is not supported by DioClientTransport',
    );
  }

  @override
  DuplexConnection duplexStream(
    String path,
    TransportMeta meta, {
    Map<String, String>? headers,
  }) {
    throw UnimplementedError(
      'duplexStream is not supported by DioClientTransport',
    );
  }
}
