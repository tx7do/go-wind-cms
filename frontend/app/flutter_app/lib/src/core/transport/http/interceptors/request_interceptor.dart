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

    // locale 不在此全局注入：
    // OpenAPI 中只有 Get 单条/GetTranslation 接口接受 locale 参数，
    // List/Create/Update/Delete 等接口不接受，全局注入会导致无效参数。
    // locale 由各 service 层的 get()/getQuery() 按需传递。

    // 继续请求流程
    return handler.next(options);
  }
}
