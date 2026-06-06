import 'package:dio/dio.dart' show DioException;
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/src/core/repositories/user_auth_cache.dart';
import 'package:flutter_app/src/core/transport/http/status.dart';
import 'package:flutter_app/src/core/utilities/logger.dart';

abstract class BaseService {
  final Logger _logger;
  final bool _loadFromLocal = false;

  BaseService({String tag = 'Service'}) : _logger = Logger(tag);

  Logger get logger => _logger;

  bool get loadFromLocal => _loadFromLocal;

  int get currentUserId {
    return GetIt.instance<UserAuthCache>().userId;
  }

  String get mqttClientId {
    return GetIt.instance<UserAuthCache>().mqttClientId;
  }

  /// 将 DioException 统一转换为 Status 对象
  ///
  /// 所有继承 BaseService 的 service 都可直接使用。
  /// Retrofit/Dio 在服务端返回非 2xx 时会抛出 DioException，
  /// 通过此方法可将其转为业务层可判断的 Status。
  ///
  /// 用法：
  /// ```dart
  /// try {
  ///   final response = await _api.someMethod(body: request);
  ///   return response;
  /// } on DioException catch (e) {
  ///   return handleDioError(e);
  /// }
  /// ```
  Status handleDioError(DioException e) {
    final data = e.response?.data;
    if (data is Map<String, dynamic>) {
      return Status(
        code: e.response?.statusCode,
        reason: data['reason'] as String?,
        message: data['message'] as String? ?? e.message,
        metadata: data['metadata'] != null
            ? Map<String, String>.from(data['metadata'] as Map)
            : null,
      );
    }

    return Status(
      code: e.response?.statusCode,
      reason: e.type.name,
      message: e.message,
    );
  }
}
