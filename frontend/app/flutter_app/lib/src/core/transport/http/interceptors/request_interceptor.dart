import 'package:dio/dio.dart';

class RequestInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    // 添加通用请求头
    options.headers['Content-Type'] = 'application/json';
    options.headers['Accept'] = 'application/json';

    // 添加时间戳参数
    options.queryParameters['timestamp'] =
        DateTime.now().millisecondsSinceEpoch;

    // 继续请求流程
    return handler.next(options);
  }
}
