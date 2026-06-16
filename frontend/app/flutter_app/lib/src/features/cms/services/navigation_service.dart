import 'package:flutter/material.dart' show IconData, Icons;
import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart' show Mutation, Query;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ApiClient, NavigationServiceClient;
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show SiteServiceV1Navigation, SiteServiceV1NavigationItem,
        SiteServiceV1NavigationItem$LinkType, SiteServiceV1Navigation$Location,
        SiteServiceV1CreateNavigationRequest, SiteServiceV1ListNavigationResponse,
        SiteServiceV1UpdateNavigationRequest, SiteServiceV1GetNavigationRequest,
        SiteServiceV1DeleteNavigationRequest;
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/services/base_service.dart';
import 'package:flutter_app/src/core/transport/http/index.dart';

typedef Navigation = SiteServiceV1Navigation;
typedef NavigationItem = SiteServiceV1NavigationItem;
typedef NavigationItemLinkType = SiteServiceV1NavigationItem$LinkType;
typedef NavigationLocation = SiteServiceV1Navigation$Location;
typedef CreateNavigationRequest = SiteServiceV1CreateNavigationRequest;
typedef UpdateNavigationRequest = SiteServiceV1UpdateNavigationRequest;
typedef ListNavigationResponse = SiteServiceV1ListNavigationResponse;

// ─── 导航图标映射 ───────────────────────────────────────

/// 图标名称 → IconData 映射
///
/// 服务端 NavigationItem.icon 字段存储字符串图标标识，
/// 此映射将其转换为 Flutter 的 IconData。
/// 当服务端传回未知图标名时，fallback 到 Icons.article。
const Map<String, IconData> navIconMap = {
  // 底部导航
  'home_outlined': Icons.home_outlined,
  'home': Icons.home,
  'explore_outlined': Icons.explore_outlined,
  'explore': Icons.explore,
  'bookmark_border': Icons.bookmark_border,
  'bookmark': Icons.bookmark,
  'person_outline': Icons.person_outline,
  'person': Icons.person,
  // 顶部导航
  'label_outlined': Icons.label_outlined,
  'label': Icons.label,
  // 通用
  'search': Icons.search,
  'settings_outlined': Icons.settings_outlined,
  'article': Icons.article,
  'folder_outlined': Icons.folder_outlined,
  'comment_outlined': Icons.comment_outlined,
};

/// 根据图标名称获取 IconData，找不到则返回 fallback
IconData resolveNavIcon(String? iconName, {IconData fallback = Icons.article}) {
  if (iconName == null || iconName.isEmpty) return fallback;
  return navIconMap[iconName] ?? fallback;
}

/// 需要路径参数（:id）的路由前缀
///
/// 只有带 / 的前缀（如 /post/、/page/、/tag/）才是真正的参数化路由，
/// /post、/tag 等现在是合法的列表页路由，不应被视为“不完整”。
const _parameterizedPrefixes = ['/post/', '/page/', '/tag/'];

/// 检查路径是否为需要参数的路由但没有带上参数
///
/// 例如 `/post/`（有前缀但没 ID）是无效的，但 `/post` 是合法的列表页路由。
bool _isIncompleteRoute(String url) {
  for (final prefix in _parameterizedPrefixes) {
    // 只有路径恰好是 /post/ 或 /tag/ （以 / 结尾但缺少 ID）才是不完整的
    if (url == prefix) return true;
  }
  return false;
}

/// 获取 NavigationItem 的路由路径
///
/// 路由以 item.url 为准（后端已填好完整路径）。
/// 仅当 url 为空时，才根据 linkType + objectId 自动生成回退路径。
/// - EXTERNAL: 始终使用 url 字段（外部链接，不走 Flutter 路由）
/// - POST/PAGE/CATEGORY: 优先 url，否则从 objectId 拼接
/// - CUSTOM/其他: 使用 url 字段
String? resolveNavRoute(NavigationItem item) {
  final linkType = item.linkType;

  // 外部链接：直接返回 url
  if (linkType == NavigationItemLinkType.linkTypeExternal) {
    return item.url;
  }

  // 优先使用后端填写的 url 字段
  if (item.url != null && item.url!.isNotEmpty) {
    if (!_isIncompleteRoute(item.url!)) return item.url;
  }

  // url 为空时根据 linkType + objectId 回退
  final String? url;
  switch (linkType) {
    case NavigationItemLinkType.linkTypePost:
      url = item.objectId != null ? '/post/${item.objectId}' : '/post';
    case NavigationItemLinkType.linkTypePage:
      url = item.objectId != null ? '/page/${item.objectId}' : '/post';
    case NavigationItemLinkType.linkTypeCategory:
      url = item.objectId != null
          ? '/post?categoryId=${item.objectId}'
          : '/category';
    default:
      url = null;
  }

  if (url == null || url.isEmpty) return null;
  if (_isIncompleteRoute(url)) return null;
  return url;
}

/// 判断该导航项是否为外部链接
bool isExternalLink(NavigationItem item) {
  return item.linkType == NavigationItemLinkType.linkTypeExternal;
}

/// 从 Navigation 列表中按位置筛选
///
/// [navigations] 服务端返回的完整导航列表
/// [location] 目标位置（HEADER / MOBILE / SIDEBAR / FOOTER）
/// [locale] 语言代码（可选）。传入后会优先匹配对应语言的导航；
///   如果该语言没有任何匹配，则 fallback 到 locale 为空的通用导航。
List<Navigation> filterNavigationsByLocation(
  List<Navigation> navigations,
  NavigationLocation location, {
  String? locale,
}) {
  final active = navigations.where((nav) {
    if (nav.isActive != true) return false;
    if (nav.location != location) return false;
    return true;
  }).toList();

  if (locale == null || locale.isEmpty) {
    // 未指定 locale，只取通用导航（locale 为 null）
    return active
        .where((nav) => nav.locale == null || nav.locale!.isEmpty)
        .toList();
  }

  // 优先找匹配 locale 的
  final matched = active.where((nav) => locale == nav.locale).toList();
  if (matched.isNotEmpty) return matched;

  // fallback 到通用导航（locale 为 null 或空）
  return active
      .where((nav) => nav.locale == null || nav.locale!.isEmpty)
      .toList();
}

/// 获取扁平化的导航项（按 sortOrder 排序）
///
/// 从匹配的 Navigation 中提取所有 NavigationItem，
/// 只取顶级项（parentId == 0 或未设置），按 sortOrder 排序。
List<NavigationItem> getFlatNavItems(
  List<Navigation> navigations,
  NavigationLocation location, {
  String? locale,
}) {
  final filtered = filterNavigationsByLocation(
    navigations,
    location,
    locale: locale,
  );
  final items = <NavigationItem>[];

  for (final nav in filtered) {
    for (final item in nav.items ?? <NavigationItem>[]) {
      if (item.parentId != null && item.parentId != 0) continue;
      if (item.isInvalid == true) continue;
      items.add(item);
    }
  }

  items.sort((a, b) => (a.sortOrder ?? 0).compareTo(b.sortOrder ?? 0));
  return items;
}

/// 导航服务
///
/// 使用 [Query] 缓存列表和详情查询，
/// 使用 [Mutation] 管理创建/更新/删除等写操作。
class NavigationService extends BaseService {
  NavigationService() : super(tag: 'NavigationService');

  NavigationServiceClient get _api =>
      GetIt.instance<ApiClient>().navigationService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取导航列表 Query
  ///
  /// [query] 分页查询参数，不传则全量加载
  Query<ListNavigationResponse> listQuery([PaginationQuery? query]) {
    final q = query ?? const PaginationQuery();
    return Query<ListNavigationResponse>(
      key: 'navigations',
      queryFn: () => _api.list(q.toPagingRequest()),
    );
  }

  /// 获取单个导航详情 Query
  Query<Navigation> getQuery({required int id}) {
    return Query<Navigation>(
      key: 'navigation-$id',
      queryFn: () => _api.get(SiteServiceV1GetNavigationRequest(id: id)),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 创建导航 Mutation
  Mutation<Navigation, Navigation> createMutation() {
    return Mutation<Navigation, Navigation>(
      mutationFn: (navigation) => _api.create(
        CreateNavigationRequest(data: navigation),
      ),
      invalidateQueries: ['navigations'],
    );
  }

  /// 更新导航 Mutation
  Mutation<Navigation, UpdateNavigationParams> updateMutation() {
    return Mutation<Navigation, UpdateNavigationParams>(
      mutationFn: (params) => _api.update(
        UpdateNavigationRequest(
          id: params.id,
          data: params.data,
          updateMask: params.updateMask,
          allowMissing: params.allowMissing,
        ),
      ),
      invalidateQueries: ['navigations'],
    );
  }

  /// 删除导航 Mutation
  Mutation<void, int> deleteMutation() {
    return Mutation<void, int>(
      mutationFn: (id) => _api
          .delete(SiteServiceV1DeleteNavigationRequest(id: id)),
      invalidateQueries: ['navigations'],
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取导航列表
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

  /// 获取单个导航
  Future<dynamic> get(int id) async {
    try {
      return await _api.get(SiteServiceV1GetNavigationRequest(id: id));
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 创建导航
  Future<dynamic> create(Navigation navigation) async {
    try {
      return await _api.create(
        CreateNavigationRequest(data: navigation),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 更新导航
  Future<dynamic> update(
    int id,
    Navigation data, {
    String? updateMask,
    bool? allowMissing,
  }) async {
    try {
      return await _api.update(
        UpdateNavigationRequest(
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

  /// 删除导航
  Future<dynamic> delete(int id) async {
    try {
      await _api.delete(SiteServiceV1DeleteNavigationRequest(id: id));
      return null;
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }
}

/// 更新导航参数
class UpdateNavigationParams {
  final int id;
  final Navigation data;
  final String? updateMask;
  final bool? allowMissing;

  const UpdateNavigationParams({
    required this.id,
    required this.data,
    this.updateMask,
    this.allowMissing,
  });
}
