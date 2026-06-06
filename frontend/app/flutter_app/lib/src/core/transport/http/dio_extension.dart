import 'package:dio/dio.dart';
import 'package:protobuf/protobuf.dart';

import 'package:flutter_app/src/core/extensions/protobuf.dart';

import 'status.dart';

extension DioExtension on Dio {
  /// GET
  Future<Response<dynamic>?> getEx<T extends GeneratedMessage>(
    String path, {
    GeneratedMessage? queryParameters,
    Options? options,
    CancelToken? cancelToken,
    ProgressCallback? onReceiveProgress,
    required T? Function() responseCreator,
  }) async {
    try {
      final ret = await get(
        path,
        queryParameters: queryParameters?.writeToProto3JsonMap(),
        options: options,
        cancelToken: cancelToken,
        onReceiveProgress: onReceiveProgress,
      );

      if (ret.statusCode == 200) {
        final response = responseCreator();
        if (response != null) {
          response.mergeFromProto3Json(ret.data);
        }
        ret.data = response;
        return ret;
      }
    } on DioException catch (e) {
      return _handleException(e);
    }

    return null;
  }

  /// DELETE
  Future<Response<dynamic>?> deleteEx<T extends GeneratedMessage>(
    String path, {
    GeneratedMessage? queryParameters,
    Options? options,
    CancelToken? cancelToken,
    required T? Function() responseCreator,
  }) async {
    try {
      final ret = await delete(
        path,
        queryParameters: queryParameters?.writeToProto3JsonMap(),
        options: options,
        cancelToken: cancelToken,
      );

      if (ret.statusCode == 200) {
        final response = responseCreator();
        if (response != null) {
          response.mergeFromProto3Json(ret.data);
        }
        ret.data = response;
        return ret;
      }
    } on DioException catch (e) {
      return _handleException(e);
    }

    return null;
  }

  /// POST
  Future<Response<dynamic>?> postEx<T extends GeneratedMessage>(
    String path, {
    GeneratedMessage? data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
    ProgressCallback? onSendProgress,
    ProgressCallback? onReceiveProgress,
    required T? Function() responseCreator,
  }) async {
    try {
      final ret = await post(
        path,
        data: data?.toProto3Json(),
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
        onSendProgress: onSendProgress,
        onReceiveProgress: onReceiveProgress,
      );

      if (ret.statusCode == 200) {
        final response = responseCreator();
        if (response != null) {
          response.mergeFromProto3Json(ret.data);
        }
        ret.data = response;
        return ret;
      }
    } on DioException catch (e) {
      return _handleException(e);
    }

    return null;
  }

  /// PUT
  Future<Response<dynamic>?> putEx<T extends GeneratedMessage>(
    String path, {
    GeneratedMessage? data,
    Map<String, dynamic>? queryParameters,
    Options? options,
    CancelToken? cancelToken,
    ProgressCallback? onSendProgress,
    ProgressCallback? onReceiveProgress,
    required T? Function() responseCreator,
  }) async {
    try {
      final ret = await put(
        path,
        data: data?.toProto3Json(),
        queryParameters: queryParameters,
        options: options,
        cancelToken: cancelToken,
        onSendProgress: onSendProgress,
        onReceiveProgress: onReceiveProgress,
      );

      if (ret.statusCode == 200) {
        final response = responseCreator();
        if (response != null) {
          response.mergeFromProto3Json(ret.data);
        }
        ret.data = response;
        return ret;
      }
    } on DioException catch (e) {
      return _handleException(e);
    }

    return null;
  }

  /// 创建Status的返回数据
  _makeStatus(DioException e) {
    if (e.response != null && e.response?.data != null) {
      e.response?.data = Status.fromJson(e.response?.data);
      return e.response;
    }
    return null;
  }

  /// 处理DIO的异常信息
  Future<Response<dynamic>?> _handleException(DioException e) {
    final status = _makeStatus(e);
    return status;
  }
}
