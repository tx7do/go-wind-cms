import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart' show Mutation, Query;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ApiClient, PageServiceClient;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Page, ContentServiceV1PageTranslation,
    ContentServiceV1CreatePageRequest, ContentServiceV1ListPageResponse,
    ContentServiceV1UpdatePageRequest, ContentServiceV1GetPageRequest,
    ContentServiceV1DeletePageRequest;
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/services/base_service.dart';
import 'package:flutter_app/src/core/transport/http/index.dart';

typedef PageModel = ContentServiceV1Page;
typedef PageTranslation = ContentServiceV1PageTranslation;
typedef CreatePageRequest = ContentServiceV1CreatePageRequest;
typedef UpdatePageRequest = ContentServiceV1UpdatePageRequest;
typedef ListPageResponse = ContentServiceV1ListPageResponse;

/// 页面服务
///
/// 使用 [Query] 缓存列表和详情查询，
/// 使用 [Mutation] 管理创建/更新/删除等写操作。
class PageService extends BaseService {
  PageService() : super(tag: 'PageService');

  PageServiceClient get _api => GetIt.instance<ApiClient>().pageService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取页面列表 Query
  ///
  /// [query] 分页查询参数，不传则全量加载
  Query<ListPageResponse> listQuery([PaginationQuery? query]) {
    final q = query ?? const PaginationQuery();
    return Query<ListPageResponse>(
      key: 'pages',
      queryFn: () => _api.list(q.toPagingRequest()),
    );
  }

  /// 获取单个页面详情 Query
  Query<PageModel> getQuery({required int id, String? slug, String? locale}) {
    return Query<PageModel>(
      key: 'page-$id',
      queryFn: () => _api.get(
        ContentServiceV1GetPageRequest(id: id, slug: slug, locale: locale),
      ),
    );
  }

  /// 获取页面翻译 Query
  Query<PageTranslation> getTranslationQuery({
    required int id,
    String? slug,
    String? locale,
  }) {
    return Query<PageTranslation>(
      key: 'page-$id-translation',
      queryFn: () => _api.getTranslation(
        ContentServiceV1GetPageRequest(id: id, slug: slug, locale: locale),
      ),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 创建页面 Mutation
  Mutation<PageModel, PageModel> createMutation() {
    return Mutation<PageModel, PageModel>(
      mutationFn: (page) =>
          _api.create(CreatePageRequest(data: page)),
      invalidateQueries: ['pages'],
    );
  }

  /// 更新页面 Mutation
  Mutation<PageModel, UpdatePageParams> updateMutation() {
    return Mutation<PageModel, UpdatePageParams>(
      mutationFn: (params) => _api.update(
        UpdatePageRequest(
          id: params.id,
          data: params.data,
          updateMask: params.updateMask,
          allowMissing: params.allowMissing,
        ),
      ),
      invalidateQueries: ['pages'],
    );
  }

  /// 删除页面 Mutation
  Mutation<void, int> deleteMutation() {
    return Mutation<void, int>(
      mutationFn: (id) =>
          _api.delete(ContentServiceV1DeletePageRequest(id: id)),
      invalidateQueries: ['pages'],
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取页面列表
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

  /// 获取单个页面
  Future<dynamic> get(int id, {String? slug, String? locale}) async {
    try {
      return await _api.get(
        ContentServiceV1GetPageRequest(id: id, slug: slug, locale: locale),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 创建页面
  Future<dynamic> create(PageModel page) async {
    try {
      return await _api.create(CreatePageRequest(data: page));
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 更新页面
  Future<dynamic> update(
    int id,
    PageModel data, {
    String? updateMask,
    bool? allowMissing,
  }) async {
    try {
      return await _api.update(
        UpdatePageRequest(
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

  /// 删除页面
  Future<dynamic> delete(int id) async {
    try {
      await _api.delete(ContentServiceV1DeletePageRequest(id: id));
      return null;
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取翻译数据
  Future<dynamic> getTranslation(int id, {String? slug, String? locale}) async {
    try {
      return await _api.getTranslation(
        ContentServiceV1GetPageRequest(id: id, slug: slug, locale: locale),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }
}

/// 更新页面参数
class UpdatePageParams {
  final int id;
  final PageModel data;
  final String? updateMask;
  final bool? allowMissing;

  const UpdatePageParams({
    required this.id,
    required this.data,
    this.updateMask,
    this.allowMissing,
  });
}
