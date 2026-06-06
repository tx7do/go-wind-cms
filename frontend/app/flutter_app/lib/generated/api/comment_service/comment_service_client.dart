// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';

import '../models/comment_service_v1_comment.dart';
import '../models/comment_service_v1_create_comment_request.dart';
import '../models/comment_service_v1_list_comment_response.dart';
import '../models/comment_service_v1_update_comment_request.dart';
import '../models/enum1.dart';

part 'comment_service_client.g.dart';

@RestApi()
abstract class CommentServiceClient {
  factory CommentServiceClient(Dio dio, {String? baseUrl}) = _CommentServiceClient;

  /// 获取评论列表.
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
  @GET('/app/v1/comments')
  Future<CommentServiceV1ListCommentResponse> commentServiceList({
    @Query('page') int? page,
    @Query('pageSize') int? pageSize,
    @Query('offset') String? offset,
    @Query('limit') int? limit,
    @Query('token') String? token,
    @Query('noPaging') bool? noPaging,
    @Query('query') String? query,
    @Query('filter') String? filter,
    @Query('filterExpr.type') Enum1? enum1,
    @Query('orderBy') String? orderBy,
    @Query('fieldMask') String? fieldMask,
  });

  /// 创建评论
  @POST('/app/v1/comments')
  Future<CommentServiceV1Comment> commentServiceCreate({
    @Body() required CommentServiceV1CreateCommentRequest body,
  });

  /// 获取评论数据
  @GET('/app/v1/comments/{id}')
  Future<CommentServiceV1Comment> commentServiceGet({
    @Path('id') required int id,
    @Query('viewMask') String? viewMask,
  });

  /// 更新评论
  @PUT('/app/v1/comments/{id}')
  Future<CommentServiceV1Comment> commentServiceUpdate({
    @Path('id') required int id,
    @Body() required CommentServiceV1UpdateCommentRequest body,
  });

  /// 删除评论
  @DELETE('/app/v1/comments/{id}')
  Future<void> commentServiceDelete({
    @Path('id') required int id,
  });
}
