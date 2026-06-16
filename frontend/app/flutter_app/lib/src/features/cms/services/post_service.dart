import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart' show Mutation, Query;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ApiClient, PostServiceClient;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Post, ContentServiceV1PostTranslation,
    ContentServiceV1CreatePostRequest, ContentServiceV1ListPostResponse,
    ContentServiceV1UpdatePostRequest, ContentServiceV1GetPostRequest,
    ContentServiceV1DeletePostRequest;
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/services/base_service.dart';
import 'package:flutter_app/src/core/transport/http/index.dart';

typedef Post = ContentServiceV1Post;
typedef PostTranslation = ContentServiceV1PostTranslation;
typedef CreatePostRequest = ContentServiceV1CreatePostRequest;
typedef UpdatePostRequest = ContentServiceV1UpdatePostRequest;
typedef ListPostResponse = ContentServiceV1ListPostResponse;

/// 帖子服务
///
/// 使用 [Query] 缓存列表和详情查询，
/// 使用 [Mutation] 管理创建/更新/删除等写操作。
class PostService extends BaseService {
  PostService() : super(tag: 'PostService');

  PostServiceClient get _api => GetIt.instance<ApiClient>().postService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取帖子列表 Query
  ///
  /// [query] 分页查询参数，不传则全量加载
  Query<ListPostResponse> listQuery([PaginationQuery? query]) {
    final q = query ?? const PaginationQuery();
    return Query<ListPostResponse>(
      key: 'posts',
      queryFn: () => _api.list(q.toPagingRequest()),
    );
  }

  /// 获取单个帖子详情 Query
  Query<Post> getQuery({required int id, String? code, String? locale}) {
    return Query<Post>(
      key: 'post-$id',
      queryFn: () => _api.get(
        ContentServiceV1GetPostRequest(id: id, code: code, locale: locale),
      ),
    );
  }

  /// 获取帖子翻译 Query
  Query<PostTranslation> getTranslationQuery({
    required int id,
    String? code,
    String? locale,
  }) {
    return Query<PostTranslation>(
      key: 'post-$id-translation',
      queryFn: () => _api.getTranslation(
        ContentServiceV1GetPostRequest(id: id, code: code, locale: locale),
      ),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 创建帖子 Mutation
  Mutation<Post, Post> createMutation() {
    return Mutation<Post, Post>(
      mutationFn: (post) =>
          _api.create(CreatePostRequest(data: post)),
      invalidateQueries: ['posts'],
    );
  }

  /// 更新帖子 Mutation
  Mutation<Post, UpdatePostParams> updateMutation() {
    return Mutation<Post, UpdatePostParams>(
      mutationFn: (params) => _api.update(
        UpdatePostRequest(
          id: params.id,
          data: params.data,
          updateMask: params.updateMask,
          allowMissing: params.allowMissing,
        ),
      ),
      invalidateQueries: ['posts'],
    );
  }

  /// 删除帖子 Mutation
  Mutation<void, int> deleteMutation() {
    return Mutation<void, int>(
      mutationFn: (id) =>
          _api.delete(ContentServiceV1DeletePostRequest(id: id)),
      invalidateQueries: ['posts'],
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取帖子列表
  ///
  /// [query] 分页查询参数，不传则全量加载
  Future<dynamic> list([PaginationQuery? query]) async {
    final q = query ?? const PaginationQuery();
    try {
      return await _api.list(q.toPagingRequest());
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取帖子分页列表（便捷方法）
  ///
  /// [page] 页码（从 1 开始）
  /// [pageSize] 每页条数
  /// [query] 可传入额外的查询条件，page/pageSize 会覆盖
  Future<ListPostResponse?> listPaged({
    required int page,
    int pageSize = 10,
    PaginationQuery? query,
  }) async {
    final q = (query ?? const PaginationQuery()).copyWith(
      page: page,
      pageSize: pageSize,
    );
    try {
      return await _api.list(q.toPagingRequest());
    } on DioException catch (e) {
      handleDioError(e);
      return null;
    }
  }

  /// 获取单个帖子
  Future<dynamic> get(int id, {String? code, String? locale}) async {
    try {
      return await _api.get(
        ContentServiceV1GetPostRequest(id: id, code: code, locale: locale),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 创建帖子
  Future<dynamic> create(Post post) async {
    try {
      return await _api.create(CreatePostRequest(data: post));
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 更新帖子
  Future<dynamic> update(
    int id,
    Post data, {
    String? updateMask,
    bool? allowMissing,
  }) async {
    try {
      return await _api.update(
        UpdatePostRequest(
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

  /// 删除帖子
  Future<dynamic> delete(int id) async {
    try {
      await _api.delete(ContentServiceV1DeletePostRequest(id: id));
      return null;
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取翻译数据
  Future<dynamic> getTranslation(int id, {String? code, String? locale}) async {
    try {
      return await _api.getTranslation(
        ContentServiceV1GetPostRequest(id: id, code: code, locale: locale),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }
}

/// 更新帖子参数
class UpdatePostParams {
  final int id;
  final Post data;
  final String? updateMask;
  final bool? allowMissing;

  const UpdatePostParams({
    required this.id,
    required this.data,
    this.updateMask,
    this.allowMissing,
  });
}
