import 'package:dio/dio.dart';
import 'package:get_it/get_it.dart';

import 'package:flutter_app/src/core/utilities/logger.dart';

abstract class BaseApiClient {
  final Dio _dio;
  final Logger _logger;

  BaseApiClient({Dio? dio, String tag = 'ApiClient'})
      : _dio = dio ?? GetIt.instance<Dio>(),
        _logger = Logger(tag);

  Dio get dio => _dio;

  Logger get logger => _logger;
}
