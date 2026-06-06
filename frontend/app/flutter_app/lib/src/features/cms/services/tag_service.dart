import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart'
    show Mutation, Query;

import 'package:flutter_app/generated/api/tag_service/tag_service_client.dart'
    show TagServiceClient;
import 'package:flutter_app/generated/api/rest_client.dart' show RestClient;

import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag_translation.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_create_tag_request.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_tag_response.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_update_tag_request.dart';
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

  TagServiceClient get _api =>
      GetIt.instance<RestClient>().tagService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取标签列表 Query
  Query<ListTagResponse> listQuery({
    int? page,
    int? pageSize,
    String? orderBy,
    bool noPaging = true,
  }) {
    return Query<ListTagResponse>(
      key: 'tags',
      queryFn: () => _api.tagServiceList(
        page: page,
        pageSize: pageSize,
        orderBy: orderBy,
        noPaging: noPaging,
      ),
    );
  }

  /// 获取单个标签详情 Query
  Query<Tag> getQuery({
    required int id,
    String? code,
    String? locale,
  }) {
    return Query<Tag>(
      key: 'tag-$id',
      queryFn: () => _api.tagServiceGet(id: id, code: code, locale: locale),
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
      queryFn: () => _api.tagServiceGetTranslation(
        id: id,
        code: code,
        locale: locale,
      ),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 创建标签 Mutation
  Mutation<Tag, Tag> createMutation() {
    return Mutation<Tag, Tag>(
      mutationFn: (tag) => _api.tagServiceCreate(
        body: CreateTagRequest(data: tag),
      ),
      invalidateQueries: ['tags'],
    );
  }

  /// 更新标签 Mutation
  Mutation<Tag, UpdateTagParams> updateMutation() {
    return Mutation<Tag, UpdateTagParams>(
      mutationFn: (params) => _api.tagServiceUpdate(
        id: params.id,
        body: UpdateTagRequest(
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
      mutationFn: (id) => _api.tagServiceDelete(id: id),
      invalidateQueries: ['tags'],
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取标签列表
  Future<dynamic> list({
    int? page,
    int? pageSize,
    String? orderBy,
    bool noPaging = true,
  }) async {
    try {
      return await _api.tagServiceList(
        page: page,
        pageSize: pageSize,
        orderBy: orderBy,
        noPaging: noPaging,
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取单个标签
  Future<dynamic> get(int id, {String? code, String? locale}) async {
    try {
      return await _api.tagServiceGet(id: id, code: code, locale: locale);
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 创建标签
  Future<dynamic> create(Tag tag) async {
    try {
      return await _api.tagServiceCreate(
        body: CreateTagRequest(data: tag),
      );
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
      return await _api.tagServiceUpdate(
        id: id,
        body: UpdateTagRequest(
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
      await _api.tagServiceDelete(id: id);
      return null;
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取翻译数据
  Future<dynamic> getTranslation(
    int id, {
    String? code,
    String? locale,
  }) async {
    try {
      return await _api.tagServiceGetTranslation(
        id: id,
        code: code,
        locale: locale,
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
