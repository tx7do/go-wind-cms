import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart'
    show Mutation, Query;

import 'package:flutter_app/generated/api/comment_service/comment_service_client.dart'
    show CommentServiceClient;
import 'package:flutter_app/generated/api/rest_client.dart' show RestClient;

import 'package:flutter_app/generated/api/models/comment_service_v1_comment.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_create_comment_request.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_list_comment_response.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_update_comment_request.dart';
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/services/base_service.dart';
import 'package:flutter_app/src/core/transport/http/index.dart';

typedef Comment = CommentServiceV1Comment;
typedef CreateCommentRequest = CommentServiceV1CreateCommentRequest;
typedef UpdateCommentRequest = CommentServiceV1UpdateCommentRequest;
typedef ListCommentResponse = CommentServiceV1ListCommentResponse;

/// 评论服务
///
/// 使用 [Query] 缓存列表和详情查询，
/// 使用 [Mutation] 管理创建/更新/删除等写操作。
class CommentService extends BaseService {
  CommentService() : super(tag: 'CommentService');

  CommentServiceClient get _api =>
      GetIt.instance<RestClient>().commentService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取评论列表 Query
  ///
  /// [query] 分页查询参数，不传则全量加载
  Query<ListCommentResponse> listQuery([PaginationQuery? query]) {
    final q = query ?? const PaginationQuery();
    return Query<ListCommentResponse>(
      key: 'comments',
      queryFn: () => _api.commentServiceList(
        page: q.page,
        pageSize: q.pageSize,
        noPaging: q.noPaging,
        orderBy: q.orderByString,
        query: q.queryString,
      ),
    );
  }

  /// 获取单个评论详情 Query
  Query<Comment> getQuery({
    required int id,
  }) {
    return Query<Comment>(
      key: 'comment-$id',
      queryFn: () => _api.commentServiceGet(id: id),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 创建评论 Mutation
  Mutation<Comment, Comment> createMutation() {
    return Mutation<Comment, Comment>(
      mutationFn: (comment) => _api.commentServiceCreate(
        body: CreateCommentRequest(data: comment),
      ),
      invalidateQueries: ['comments'],
    );
  }

  /// 更新评论 Mutation
  Mutation<Comment, UpdateCommentParams> updateMutation() {
    return Mutation<Comment, UpdateCommentParams>(
      mutationFn: (params) => _api.commentServiceUpdate(
        id: params.id,
        body: UpdateCommentRequest(
          id: params.id,
          data: params.data,
          updateMask: params.updateMask,
          allowMissing: params.allowMissing,
        ),
      ),
      invalidateQueries: ['comments'],
    );
  }

  /// 删除评论 Mutation
  Mutation<void, int> deleteMutation() {
    return Mutation<void, int>(
      mutationFn: (id) => _api.commentServiceDelete(id: id),
      invalidateQueries: ['comments'],
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取评论列表
  ///
  /// [query] 分页查询参数，不传则全量加载
  Future<dynamic> list([PaginationQuery? query]) async {
    final q = query ?? const PaginationQuery();
    try {
      return await _api.commentServiceList(
        page: q.page,
        pageSize: q.pageSize,
        noPaging: q.noPaging,
        orderBy: q.orderByString,
        query: q.queryString,
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取单个评论
  Future<dynamic> get(int id) async {
    try {
      return await _api.commentServiceGet(id: id);
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 创建评论
  Future<dynamic> create(Comment comment) async {
    try {
      return await _api.commentServiceCreate(
        body: CreateCommentRequest(data: comment),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 更新评论
  Future<dynamic> update(
    int id,
    Comment data, {
    String? updateMask,
    bool? allowMissing,
  }) async {
    try {
      return await _api.commentServiceUpdate(
        id: id,
        body: UpdateCommentRequest(
          id: id,
          data: data,
          updateMask: updateMask,
          allowMissing: allowMissing,
        ),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 删除评论
  Future<dynamic> delete(int id) async {
    try {
      await _api.commentServiceDelete(id: id);
      return null;
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }
}

/// 更新评论参数
class UpdateCommentParams {
  final int id;
  final Comment data;
  final String? updateMask;
  final bool? allowMissing;

  const UpdateCommentParams({
    required this.id,
    required this.data,
    this.updateMask,
    this.allowMissing,
  });
}
