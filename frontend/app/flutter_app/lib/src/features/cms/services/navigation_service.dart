import 'package:flutter/material.dart' show IconData, Icons;
import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart'
    show Mutation, Query, InfiniteQuery;

import 'package:flutter_app/generated/api/navigation_service/navigation_service_client.dart'
    show NavigationServiceClient;
import 'package:flutter_app/generated/api/rest_client.dart' show RestClient;

import 'package:flutter_app/generated/api/models/site_service_v1_navigation.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation_item.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation_item_link_type.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_navigation_location.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_create_navigation_request.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_list_navigation_response.dart';
import 'package:flutter_app/generated/api/models/site_service_v1_update_navigation_request.dart';
import 'package:flutter_app/src/core/services/base_service.dart';
import 'package:flutter_app/src/core/transport/http/index.dart';

typedef Navigation = SiteServiceV1Navigation;
typedef NavigationItem = SiteServiceV1NavigationItem;
typedef NavigationItemLinkType = SiteServiceV1NavigationItemLinkType;
typedef NavigationLocation = SiteServiceV1NavigationLocation;
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

/// 获取 NavigationItem 的路由路径
///
/// 根据 linkType 和 url/objectId 生成 Flutter 路由：
/// - CUSTOM: 使用 url 字段
/// - POST: /post/{objectId}
/// - PAGE: /page/{objectId}
/// - CATEGORY: /category/{objectId}
/// - EXTERNAL: 使用 url 字段（外部链接，不走 Flutter 路由）
String? resolveNavRoute(NavigationItem item) {
  final linkType = item.linkType;
  if (linkType == null) return item.url;
  switch (linkType) {
    case NavigationItemLinkType.linkTypeCustom:
      return item.url != null && item.url!.isNotEmpty ? item.url : null;
    case NavigationItemLinkType.linkTypePost:
      return item.objectId != null ? '/post/${item.objectId}' : null;
    case NavigationItemLinkType.linkTypePage:
      return item.objectId != null ? '/page/${item.objectId}' : null;
    case NavigationItemLinkType.linkTypeCategory:
      return item.objectId != null ? '/category/${item.objectId}' : null;
    case NavigationItemLinkType.linkTypeExternal:
      return item.url != null && item.url!.isNotEmpty ? item.url : null;
    default:
      return item.url != null && item.url!.isNotEmpty ? item.url : null;
  }
}

/// 判断该导航项是否为外部链接
bool isExternalLink(NavigationItem item) {
  return item.linkType == NavigationItemLinkType.linkTypeExternal;
}

/// 从 Navigation 列表中按位置筛选
///
/// [navigations] 服务端返回的完整导航列表
/// [location] 目标位置（HEADER / MOBILE / SIDEBAR / FOOTER）
/// [locale] 语言代码（可选，空则不过滤）
List<Navigation> filterNavigationsByLocation(
  List<Navigation> navigations,
  NavigationLocation location, {
  String? locale,
}) {
  return navigations.where((nav) {
    if (nav.isActive != true) return false;
    if (nav.location != location) return false;
    if (locale != null) {
      final navLocale = nav.locale;
      if (navLocale != null && navLocale != locale) return false;
    }
    return true;
  }).toList();
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
  final filtered = filterNavigationsByLocation(navigations, location, locale: locale);
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
      GetIt.instance<RestClient>().navigationService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取导航列表 Query
  Query<ListNavigationResponse> listQuery({
    int? page,
    int? pageSize,
    String? orderBy,
    bool noPaging = true,
  }) {
    return Query<ListNavigationResponse>(
      key: 'navigations',
      queryFn: () => _api.navigationServiceList(
        page: page,
        pageSize: pageSize,
        orderBy: orderBy,
        noPaging: noPaging,
      ),
    );
  }

  /// 获取单个导航详情 Query
  Query<Navigation> getQuery({
    required int id,
  }) {
    return Query<Navigation>(
      key: 'navigation-$id',
      queryFn: () => _api.navigationServiceGet(id: id),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 创建导航 Mutation
  Mutation<Navigation, Navigation> createMutation() {
    return Mutation<Navigation, Navigation>(
      mutationFn: (navigation) => _api.navigationServiceCreate(
        body: CreateNavigationRequest(data: navigation),
      ),
      invalidateQueries: ['navigations'],
    );
  }

  /// 更新导航 Mutation
  Mutation<Navigation, UpdateNavigationParams> updateMutation() {
    return Mutation<Navigation, UpdateNavigationParams>(
      mutationFn: (params) => _api.navigationServiceUpdate(
        id: params.id,
        body: UpdateNavigationRequest(
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
      mutationFn: (id) => _api.navigationServiceDelete(id: id),
      invalidateQueries: ['navigations'],
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取导航列表
  Future<dynamic> list({
    int? page,
    int? pageSize,
    String? orderBy,
    bool noPaging = true,
  }) async {
    try {
      return await _api.navigationServiceList(
        page: page,
        pageSize: pageSize,
        orderBy: orderBy,
        noPaging: noPaging,
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 获取单个导航
  Future<dynamic> get(int id) async {
    try {
      return await _api.navigationServiceGet(id: id);
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 创建导航
  Future<dynamic> create(Navigation navigation) async {
    try {
      return await _api.navigationServiceCreate(
        body: CreateNavigationRequest(data: navigation),
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
      return await _api.navigationServiceUpdate(
        id: id,
        body: UpdateNavigationRequest(
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
      await _api.navigationServiceDelete(id: id);
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
