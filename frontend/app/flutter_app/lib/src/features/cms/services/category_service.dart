import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart'
    show Mutation, Query;

import 'package:flutter_app/generated/api/category_service/category_service_client.dart'
    show CategoryServiceClient;
import 'package:flutter_app/generated/api/rest_client.dart' show RestClient;

import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category_translation.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_create_category_request.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_list_category_response.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_update_category_request.dart';
import 'package:flutter_app/src/core/services/base_service.dart';
import 'package:flutter_app/src/core/transport/http/index.dart';

typedef Category = ContentServiceV1Category;
typedef CategoryTranslation = ContentServiceV1CategoryTranslation;
typedef CreateCategoryRequest = ContentServiceV1CreateCategoryRequest;
typedef UpdateCategoryRequest = ContentServiceV1UpdateCategoryRequest;
typedef ListCategoryResponse = ContentServiceV1ListCategoryResponse;

/// 类别服务
///
/// 使用 [Query] 缓存列表和详情查询，
/// 使用 [Mutation] 管理创建/更新/删除等写操作。
class CategoryService extends BaseService {
  CategoryService() : super(tag: 'CategoryService');

  CategoryServiceClient get _api =>
      GetIt.instance<RestClient>().categoryService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取类别列表 Query
  ///
  /// 用法：
  /// ```dart
  /// final query = categoryService.listQuery();
  /// // 监听：QueryBuilder<ListQueryParams, ...>
  /// // 触发：query.fetch(params);
  /// ```
  Query<ListCategoryResponse> listQuery({
    int? page,
    int? pageSize,
    String? orderBy,
    bool noPaging = true,
  }) {
    return Query<ListCategoryResponse>(
      key: 'categories',
      queryFn: () => _api.categoryServiceList(
        page: page,
        pageSize: pageSize,
        orderBy: orderBy,
        noPaging: noPaging,
      ),
    );
  }

  /// 获取单个类别详情 Query
  Query<Category> getQuery({
    required int id,
    String? locale,
  }) {
    return Query<Category>(
      key: 'category-$id',
      queryFn: () => _api.categoryServiceGet(id: id, locale: locale),
    );
  }

  /// 获取类别翻译 Query
  Query<CategoryTranslation> getTranslationQuery({
    required int id,
    String? code,
    String? locale,
  }) {
    return Query<CategoryTranslation>(
      key: 'category-$id-translation',
      queryFn: () => _api.categoryServiceGetTranslation(
        id: id,
        code: code,
        locale: locale,
      ),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 创建类别 Mutation
  Mutation<Category, Category>
      createMutation() {
    return Mutation<Category, Category>(
      mutationFn: (category) => _api.categoryServiceCreate(
        body: CreateCategoryRequest(data: category),
      ),
      invalidateQueries: ['categories'],
    );
  }

  /// 更新类别 Mutation
  Mutation<Category, UpdateCategoryParams>
      updateMutation() {
    return Mutation<Category, UpdateCategoryParams>(
      mutationFn: (params) => _api.categoryServiceUpdate(
        id: params.id,
        body: UpdateCategoryRequest(
          id: params.id,
          data: params.data,
          updateMask: params.updateMask,
          allowMissing: params.allowMissing,
        ),
      ),
      invalidateQueries: ['categories'],
    );
  }

  /// 删除类别 Mutation
  Mutation<void, int> deleteMutation() {
    return Mutation<void, int>(
      mutationFn: (id) => _api.categoryServiceDelete(id: id),
      invalidateQueries: ['categories'],
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取类别列表
  Future<dynamic> list({
    int? page,
    int? pageSize,
    String? orderBy,
    bool noPaging = true,
  }) async {
    try {
      return await _api.categoryServiceList(
        page: page,
        pageSize: pageSize,
        orderBy: orderBy,
        noPaging: noPaging,
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取单个类别
  Future<dynamic> get(int id, {String? locale}) async {
    try {
      return await _api.categoryServiceGet(id: id, locale: locale);
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 创建类别
  Future<dynamic> create(Category category) async {
    try {
      return await _api.categoryServiceCreate(
        body: CreateCategoryRequest(data: category),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 更新类别
  Future<dynamic> update(
    int id,
    Category data, {
    String? updateMask,
    bool? allowMissing,
  }) async {
    try {
      return await _api.categoryServiceUpdate(
        id: id,
        body: UpdateCategoryRequest(
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

  /// 删除类别
  Future<dynamic> delete(int id) async {
    try {
      await _api.categoryServiceDelete(id: id);
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
      return await _api.categoryServiceGetTranslation(
        id: id,
        code: code,
        locale: locale,
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }
}

/// 更新类别参数
class UpdateCategoryParams {
  final int id;
  final Category data;
  final String? updateMask;
  final bool? allowMissing;

  const UpdateCategoryParams({
    required this.id,
    required this.data,
    this.updateMask,
    this.allowMissing,
  });
}
