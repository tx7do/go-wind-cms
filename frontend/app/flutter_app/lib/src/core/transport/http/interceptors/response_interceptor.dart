import 'package:dio/dio.dart';

class ResponseInterceptor extends Interceptor {
  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    // 处理通用的响应格式
    if (response.data is Map) {
      final data = response.data as Map<String, dynamic>;

      // 检查业务状态码
      if (data.containsKey('code') && data['code'] != 200) {
        // 业务错误，创建自定义错误
        return handler.reject(
          DioException(
            requestOptions: response.requestOptions,
            response: response,
            type: DioExceptionType.badResponse,
            error: '业务错误: ${data['message']}',
          ),
        );
      }

      // 提取实际数据
      if (data.containsKey('data')) {
        response.data = data['data'];
      }
    }

    // 继续响应流程
    return handler.next(response);
  }
}
