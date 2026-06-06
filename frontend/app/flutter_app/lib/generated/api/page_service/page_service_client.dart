// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';

import '../models/content_service_v1_create_page_request.dart';
import '../models/content_service_v1_list_page_response.dart';
import '../models/content_service_v1_page.dart';
import '../models/content_service_v1_page_translation.dart';
import '../models/content_service_v1_update_page_request.dart';
import '../models/enum3.dart';

part 'page_service_client.g.dart';

@RestApi()
abstract class PageServiceClient {
  factory PageServiceClient(Dio dio, {String? baseUrl}) = _PageServiceClient;

  /// 获取页面列表.
  ///
  /// [page] - 当前页码（从1开始，默认1）.
  ///
  /// [pageSize] - 每页条数（默认10，建议设置上限如100）.
  ///
  /// [offset] - 跳过的记录数（从0开始，默认0）.
  ///
  /// [limit] - 最多返回的记录数（默认10，建议设置上限如100）.
  ///
  /// [token] - 上一页最后一条记录的游标（如ID/时间戳+ID，首次请求为空）.
  ///
  /// [noPaging] - 是否不分页，如果为true，则page和pageSize参数无效。.
  ///
  /// [query] - JSON字符串过滤条件，基础语法：{"field1":"val1", "field2___icontains":"val2"}，具体请参见：https://github.com/tx7do/go-crud/tree/main/pagination/filter/README.md.
  ///
  /// [filter] - Google AIP规范字符串过滤条件.
  ///
  /// [filterExprType] - 过滤表达式类型.
  ///
  /// [orderBy] - 排序条件.
  ///
  /// [fieldMask] - 字段掩码，其作用为SELECT中的字段，其语法为使用逗号分隔字段名，例如：id,realName,userName。如果为空则选中所有字段，即SELECT *。.
  @GET('/app/v1/pages')
  Future<ContentServiceV1ListPageResponse> pageServiceList({
    @Query('page') int? page,
    @Query('pageSize') int? pageSize,
    @Query('offset') String? offset,
    @Query('limit') int? limit,
    @Query('token') String? token,
    @Query('noPaging') bool? noPaging,
    @Query('query') String? query,
    @Query('filter') String? filter,
    @Query('filterExpr.type') Enum3? enum3,
    @Query('orderBy') String? orderBy,
    @Query('fieldMask') String? fieldMask,
  });

  /// 创建页面
  @POST('/app/v1/pages')
  Future<ContentServiceV1Page> pageServiceCreate({
    @Body() required ContentServiceV1CreatePageRequest body,
  });

  /// 获取页面数据
  @GET('/app/v1/pages/{id}')
  Future<ContentServiceV1Page> pageServiceGet({
    @Path('id') required int id,
    @Query('slug') String? slug,
    @Query('locale') String? locale,
    @Query('viewMask') String? viewMask,
  });

  /// 更新页面
  @PUT('/app/v1/pages/{id}')
  Future<ContentServiceV1Page> pageServiceUpdate({
    @Path('id') required int id,
    @Body() required ContentServiceV1UpdatePageRequest body,
  });

  /// 删除页面
  @DELETE('/app/v1/pages/{id}')
  Future<void> pageServiceDelete({
    @Path('id') required int id,
  });

  /// 获取翻译数据
  @GET('/app/v1/pages/{id}/translation')
  Future<ContentServiceV1PageTranslation> pageServiceGetTranslation({
    @Path('id') required int id,
    @Query('slug') String? slug,
    @Query('locale') String? locale,
    @Query('viewMask') String? viewMask,
  });
}
