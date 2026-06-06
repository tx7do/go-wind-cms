// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';

import '../models/content_service_v1_create_tag_request.dart';
import '../models/content_service_v1_list_tag_response.dart';
import '../models/content_service_v1_tag.dart';
import '../models/content_service_v1_tag_translation.dart';
import '../models/content_service_v1_update_tag_request.dart';
import '../models/enum5.dart';

part 'tag_service_client.g.dart';

@RestApi()
abstract class TagServiceClient {
  factory TagServiceClient(Dio dio, {String? baseUrl}) = _TagServiceClient;

  /// 获取标签列表.
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
  @GET('/app/v1/tags')
  Future<ContentServiceV1ListTagResponse> tagServiceList({
    @Query('page') int? page,
    @Query('pageSize') int? pageSize,
    @Query('offset') String? offset,
    @Query('limit') int? limit,
    @Query('token') String? token,
    @Query('noPaging') bool? noPaging,
    @Query('query') String? query,
    @Query('filter') String? filter,
    @Query('filterExpr.type') Enum5? enum5,
    @Query('orderBy') String? orderBy,
    @Query('fieldMask') String? fieldMask,
  });

  /// 创建标签
  @POST('/app/v1/tags')
  Future<ContentServiceV1Tag> tagServiceCreate({
    @Body() required ContentServiceV1CreateTagRequest body,
  });

  /// 获取标签数据
  @GET('/app/v1/tags/{id}')
  Future<ContentServiceV1Tag> tagServiceGet({
    @Path('id') required int id,
    @Query('code') String? code,
    @Query('locale') String? locale,
    @Query('viewMask') String? viewMask,
  });

  /// 更新标签
  @PUT('/app/v1/tags/{id}')
  Future<ContentServiceV1Tag> tagServiceUpdate({
    @Path('id') required int id,
    @Body() required ContentServiceV1UpdateTagRequest body,
  });

  /// 删除标签
  @DELETE('/app/v1/tags/{id}')
  Future<void> tagServiceDelete({
    @Path('id') required int id,
  });

  /// 获取翻译数据
  @GET('/app/v1/tags/{id}/translation')
  Future<ContentServiceV1TagTranslation> tagServiceGetTranslation({
    @Path('id') required int id,
    @Query('code') String? code,
    @Query('locale') String? locale,
    @Query('viewMask') String? viewMask,
  });
}
