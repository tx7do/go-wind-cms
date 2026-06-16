import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart' show Mutation, Query;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ApiClient, CommentServiceClient;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show CommentServiceV1Comment, CommentServiceV1CreateCommentRequest,
    CommentServiceV1ListCommentResponse, CommentServiceV1UpdateCommentRequest,
    CommentServiceV1GetCommentRequest, CommentServiceV1DeleteCommentRequest;
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

  CommentServiceClient get _api => GetIt.instance<ApiClient>().commentService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取评论列表 Query
  ///
  /// [query] 分页查询参数，不传则全量加载
  Query<ListCommentResponse> listQuery([PaginationQuery? query]) {
    final q = query ?? const PaginationQuery(skipLocale: true);
    return Query<ListCommentResponse>(
      key: 'comments',
      queryFn: () => _api.list(q.toPagingRequest()),
    );
  }

  /// 获取单个评论详情 Query
  Query<Comment> getQuery({required int id}) {
    return Query<Comment>(
      key: 'comment-$id',
      queryFn: () => _api.get(CommentServiceV1GetCommentRequest(id: id)),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 创建评论 Mutation
  Mutation<Comment, Comment> createMutation() {
    return Mutation<Comment, Comment>(
      mutationFn: (comment) =>
          _api.create(CreateCommentRequest(data: comment)),
      invalidateQueries: ['comments'],
    );
  }

  /// 更新评论 Mutation
  Mutation<Comment, UpdateCommentParams> updateMutation() {
    return Mutation<Comment, UpdateCommentParams>(
      mutationFn: (params) => _api.update(
        UpdateCommentRequest(
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
      mutationFn: (id) =>
          _api.delete(CommentServiceV1DeleteCommentRequest(id: id)),
      invalidateQueries: ['comments'],
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取评论列表
  ///
  /// [query] 分页查询参数，不传则全量加载
  Future<dynamic> list([PaginationQuery? query]) async {
    final q = query ?? const PaginationQuery(skipLocale: true);
    try {
      return await _api.list(q.toPagingRequest());
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取单个评论
  Future<dynamic> get(int id) async {
    try {
      return await _api.get(CommentServiceV1GetCommentRequest(id: id));
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 创建评论
  Future<dynamic> create(Comment comment) async {
    try {
      return await _api.create(
        CreateCommentRequest(data: comment),
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
      return await _api.update(
        UpdateCommentRequest(
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
      await _api.delete(CommentServiceV1DeleteCommentRequest(id: id));
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
