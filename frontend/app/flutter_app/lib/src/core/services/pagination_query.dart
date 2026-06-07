import 'dart:convert' show jsonEncode;

import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/src/core/preference/user_preference_cache.dart';

/// 通用分页查询类
///
/// 参考 TS 端的 PaginationQuery，统一封装所有 List API 的查询参数。
/// locale 自动从 [UserPreferenceCache] 获取（可覆盖）。
///
/// 用法：
/// ```dart
/// // 全量加载
/// _postService.list(PaginationQuery());
///
/// // 分页加载
/// _postService.list(PaginationQuery(page: 1, pageSize: 10));
///
/// // 带过滤条件
/// _postService.list(PaginationQuery(
///   page: 1,
///   pageSize: 10,
///   formValues: {'category_id': 5, 'status': 'published'},
///   orderBy: ['-created_at'],
/// ));
/// ```
class PaginationQuery {
  final int? page;
  final int? pageSize;

  /// 过滤条件（会被序列化为 query JSON 字符串）
  final Map<String, dynamic>? formValues;

  /// 字段掩码（SELECT 字段白名单）
  /// 支持 String（逗号分隔）或 List<String>
  final Object? fieldMask;

  /// 排序条件（默认 ['-created_at']）
  final List<String>? orderBy;

  /// 语言代码（如 'zh-CN'）。不传则自动从 UserPreferenceCache 获取
  final String? locale;

  /// 是否需要清理租户字段
  final bool isTenantUser;

  /// 是否跳过自动注入 locale
  ///
  /// 某些 API（如 Comment）不支持 locale 参数，设为 true 可跳过。
  final bool skipLocale;

  const PaginationQuery({
    this.page,
    this.pageSize,
    this.formValues,
    this.fieldMask,
    this.orderBy,
    this.locale,
    this.isTenantUser = false,
    this.skipLocale = false,
  });

  // ─── 私有工具方法 ──────────────────────────────────────

  /// 移除 Map 中的 null、undefined、空字符串值
  static Map<String, dynamic> _removeNullUndefined(Map<String, dynamic> obj) {
    return Map.fromEntries(
      obj.entries.where(
        (entry) => entry.value != null && entry.value.toString().isNotEmpty,
      ),
    );
  }

  /// 获取当前用户语言偏好，转换为 API 格式 (zh_CN → zh-CN)
  static String? get _currentLocale {
    try {
      final lang = GetIt.instance<UserPreferenceCache>().language;
      if (lang.isEmpty) return null;
      return lang.replaceAll('_', '-');
    } catch (_) {
      return null;
    }
  }

  // ─── 计算属性 ──────────────────────────────────────────

  /// 是否不分页
  ///
  /// 参考 TS: get noPaging() { return !this.paging?.page && !this.paging?.pageSize; }
  /// 有 page 或 pageSize → 分页
  /// 都没有 → 全量
  bool get noPaging => page == null && pageSize == null;

  /// 构建 query JSON 字符串
  ///
  /// 将 formValues + locale 合并为 JSON 字符串。
  /// 参考 TS: makeQueryString()
  String? get queryString {
    final Map<String, dynamic> values = {};

    if (formValues != null) {
      final cleaned = _removeNullUndefined(formValues!);
      values.addAll(cleaned);
    }

    if (!skipLocale) {
      final resolvedLocale = locale ?? _currentLocale;
      if (resolvedLocale != null && resolvedLocale.isNotEmpty) {
        values['locale'] = resolvedLocale;
      }
    }

    if (values.isEmpty) return null;

    if (isTenantUser) {
      values.remove('tenant_id');
      values.remove('tenantId');
      if (values.isEmpty) return null;
    }

    return jsonEncode(values);
  }

  /// 构建 orderBy JSON 字符串
  ///
  /// 默认 ['-created_at']（按创建时间倒序）
  String? get orderByString {
    final list = orderBy ?? const ['-created_at'];
    return jsonEncode(list);
  }

  /// 格式化 fieldMask
  ///
  /// List → 逗号分隔字符串；String → 原样返回
  String? get formattedFieldMask {
    if (fieldMask == null) return null;

    if (fieldMask is List) {
      return (fieldMask as List)
          .whereType<String>()
          .where((s) => s.isNotEmpty)
          .join(',');
    }

    final str = fieldMask.toString().trim();
    return str.isEmpty ? null : str;
  }

  // ─── 生成请求参数 ──────────────────────────────────────

  /// 生成后端 API 需要的参数 Map
  ///
  /// 可直接展开到 Retrofit 方法调用中：
  /// ```dart
  /// _api.postServiceList(...query.toRawParams());
  /// ```
  ///
  /// 注意：Retrofit 方法使用命名参数，所以需要逐个传：
  /// ```dart
  /// _api.postServiceList(
  ///   page: query.page,
  ///   pageSize: query.pageSize,
  ///   noPaging: query.noPaging,
  ///   query: query.queryString,
  ///   orderBy: query.orderByString,
  ///   fieldMask: query.formattedFieldMask,
  /// );
  /// ```
  Map<String, dynamic> toRawParams() {
    return {
      'page': page,
      'pageSize': pageSize,
      'noPaging': noPaging,
      'fieldMask': formattedFieldMask,
      'orderBy': orderByString,
      'query': queryString,
      // 以下参数预留，当前未使用
      'sorting': null,
      'offset': null,
      'limit': null,
      'token': null,
      'filter': null,
    };
  }

  /// 复制并修改部分字段
  PaginationQuery copyWith({
    int? page,
    int? pageSize,
    Map<String, dynamic>? formValues,
    Object? fieldMask,
    List<String>? orderBy,
    String? locale,
    bool? isTenantUser,
    bool? skipLocale,
    bool clearPage = false,
    bool clearFormValues = false,
  }) {
    return PaginationQuery(
      page: clearPage ? null : (page ?? this.page),
      pageSize: pageSize ?? this.pageSize,
      formValues: clearFormValues ? null : (formValues ?? this.formValues),
      fieldMask: fieldMask ?? this.fieldMask,
      orderBy: orderBy ?? this.orderBy,
      locale: locale ?? this.locale,
      isTenantUser: isTenantUser ?? this.isTenantUser,
      skipLocale: skipLocale ?? this.skipLocale,
    );
  }

  /// 下一页查询
  PaginationQuery nextPage() {
    return copyWith(page: (page ?? 1) + 1);
  }
}
