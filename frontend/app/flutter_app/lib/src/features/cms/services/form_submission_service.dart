import 'package:dio/dio.dart' show DioException;
import 'package:flutter/foundation.dart' show debugPrint;
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ApiClient, FormSubmissionServiceClient;
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1SubmitFormSubmissionRequest;
import 'package:flutter_app/src/core/services/base_service.dart';

/// 表单提交服务（页面 FORM 区块的访客提交）。
///
/// 匿名可调用；校验（页面已发布/区块 FORM/字段集约束）在 core 侧执行。
/// 失败经 [handleDioError] 归一为 Status 供调用方区分错误/成功。
class FormSubmissionService extends BaseService {
  FormSubmissionService() : super(tag: 'FormSubmissionService');

  FormSubmissionServiceClient get _api =>
      GetIt.instance<ApiClient>().formSubmissionService;

  /// 提交表单。成功返回提交记录，失败返回 Status（BaseService 错误归一）。
  Future<dynamic> submit({
    required int pageId,
    required int sectionId,
    required Map<String, String> data,
  }) async {
    try {
      return await _api.submit(ContentServiceV1SubmitFormSubmissionRequest(
        pageId: pageId,
        sectionId: sectionId,
        data: data,
      ));
    } on DioException catch (e) {
      debugPrint('form submit failed: ${e.message}');
      return handleDioError(e);
    }
  }
}
