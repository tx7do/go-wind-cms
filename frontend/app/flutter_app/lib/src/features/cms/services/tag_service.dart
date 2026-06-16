import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart' show Mutation, Query;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ApiClient, TagServiceClient;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ContentServiceV1Tag, ContentServiceV1TagTranslation,
    ContentServiceV1CreateTagRequest, ContentServiceV1ListTagResponse,
    ContentServiceV1UpdateTagRequest, ContentServiceV1GetTagRequest,
    ContentServiceV1DeleteTagRequest;
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/services/base_service.dart';
import 'package:flutter_app/src/core/transport/http/index.dart';

typedef Tag = ContentServiceV1Tag;
typedef TagTranslation = ContentServiceV1TagTranslation;
typedef CreateTagRequest = ContentServiceV1CreateTagRequest;
typedef UpdateTagRequest = ContentServiceV1UpdateTagRequest;
typedef ListTagResponse = ContentServiceV1ListTagResponse;

/// 标签服务
///
/// 使用 [Query] 缓存列表和详情查询，
/// 使用 [Mutation] 管理创建/更新/删除等写操作。
class TagService extends BaseService {
  TagService() : super(tag: 'TagService');

  TagServiceClient get _api => GetIt.instance<ApiClient>().tagService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取标签列表 Query
  ///
  /// [query] 分页查询参数，不传则全量加载
  Query<ListTagResponse> listQuery([PaginationQuery? query]) {
    final q = query ?? const PaginationQuery();
    return Query<ListTagResponse>(
      key: 'tags',
      queryFn: () => _api.list(q.toPagingRequest()),
    );
  }

  /// 获取单个标签详情 Query
  Query<Tag> getQuery({required int id, String? code, String? locale}) {
    return Query<Tag>(
      key: 'tag-$id',
      queryFn: () => _api.get(
        ContentServiceV1GetTagRequest(id: id, code: code, locale: locale),
      ),
    );
  }

  /// 获取标签翻译 Query
  Query<TagTranslation> getTranslationQuery({
    required int id,
    String? code,
    String? locale,
  }) {
    return Query<TagTranslation>(
      key: 'tag-$id-translation',
      queryFn: () => _api.getTranslation(
        ContentServiceV1GetTagRequest(id: id, code: code, locale: locale),
      ),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 创建标签 Mutation
  Mutation<Tag, Tag> createMutation() {
    return Mutation<Tag, Tag>(
      mutationFn: (tag) =>
          _api.create(CreateTagRequest(data: tag)),
      invalidateQueries: ['tags'],
    );
  }

  /// 更新标签 Mutation
  Mutation<Tag, UpdateTagParams> updateMutation() {
    return Mutation<Tag, UpdateTagParams>(
      mutationFn: (params) => _api.update(
        UpdateTagRequest(
          id: params.id,
          data: params.data,
          updateMask: params.updateMask,
          allowMissing: params.allowMissing,
        ),
      ),
      invalidateQueries: ['tags'],
    );
  }

  /// 删除标签 Mutation
  Mutation<void, int> deleteMutation() {
    return Mutation<void, int>(
      mutationFn: (id) =>
          _api.delete(ContentServiceV1DeleteTagRequest(id: id)),
      invalidateQueries: ['tags'],
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取标签列表
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

  /// 获取单个标签
  Future<dynamic> get(int id, {String? code, String? locale}) async {
    try {
      return await _api.get(
        ContentServiceV1GetTagRequest(id: id, code: code, locale: locale),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 创建标签
  Future<dynamic> create(Tag tag) async {
    try {
      return await _api.create(CreateTagRequest(data: tag));
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 更新标签
  Future<dynamic> update(
    int id,
    Tag data, {
    String? updateMask,
    bool? allowMissing,
  }) async {
    try {
      return await _api.update(
        UpdateTagRequest(
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

  /// 删除标签
  Future<dynamic> delete(int id) async {
    try {
      await _api.delete(ContentServiceV1DeleteTagRequest(id: id));
      return null;
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取翻译数据
  Future<dynamic> getTranslation(int id, {String? code, String? locale}) async {
    try {
      return await _api.getTranslation(
        ContentServiceV1GetTagRequest(id: id, code: code, locale: locale),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }
}

/// 更新标签参数
class UpdateTagParams {
  final int id;
  final Tag data;
  final String? updateMask;
  final bool? allowMissing;

  const UpdateTagParams({
    required this.id,
    required this.data,
    this.updateMask,
    this.allowMissing,
  });
}
