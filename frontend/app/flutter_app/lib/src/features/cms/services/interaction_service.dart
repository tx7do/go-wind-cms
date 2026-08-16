import 'package:dio/dio.dart' show DioException;
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ApiClient, InteractionServiceClient;
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show InteractionServiceV1LikeRequest,
        InteractionServiceV1LikeResponse,
        InteractionServiceV1WatchRequest,
        InteractionServiceV1WatchResponse,
        InteractionServiceV1GetInteractionStatusRequest,
        InteractionServiceV1GetInteractionStatusResponse,
        InteractionServiceV1TargetType,
        InteractionServiceV1GetCountsRequest,
        InteractionServiceV1GetCountsResponse,
        InteractionServiceV1CounterMetric,
        ContentServiceV1ListPostResponse;
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/services/base_service.dart';

/// 交互服务（点赞 / 收藏）
///
/// 作为 InteractionService 的客户端封装，对应后端点赞/收藏计数内聚子系统。
/// viewer 身份由后端鉴权上下文确定，前端无需传 user_id。
/// 点赞计数统一存于 interaction_counter 表（post/comment 表的散落计数列已移除），
/// 列表/详情渲染计数改读 [getCounts]；点赞/取消后由 LikeResponse.likeCount 乐观更新。
class InteractionService extends BaseService {
  InteractionService() : super(tag: 'InteractionService');

  InteractionServiceClient get _api =>
      GetIt.instance<ApiClient>().interactionService;

  // ─── 点赞 ──────────────────────────────────────────

  /// 点赞（post 或 comment）。幂等。返回操作后的点赞状态与最新计数。
  Future<InteractionServiceV1LikeResponse?> like(
      InteractionServiceV1TargetType targetType, int targetId) async {
    try {
      return await _api.like(
          InteractionServiceV1LikeRequest(targetType: targetType, targetId: targetId));
    } on DioException catch (e) {
      handleDioError(e);
      return null;
    }
  }

  /// 取消点赞。幂等。
  Future<InteractionServiceV1LikeResponse?> unlike(
      InteractionServiceV1TargetType targetType, int targetId) async {
    try {
      return await _api.unlike(
          InteractionServiceV1LikeRequest(targetType: targetType, targetId: targetId));
    } on DioException catch (e) {
      handleDioError(e);
      return null;
    }
  }

  // ─── 收藏 ──────────────────────────────────────────

  /// 收藏 post。幂等。
  Future<InteractionServiceV1WatchResponse?> watch(int postId) async {
    try {
      return await _api
          .watch(InteractionServiceV1WatchRequest(postId: postId));
    } on DioException catch (e) {
      handleDioError(e);
      return null;
    }
  }

  /// 取消收藏 post。幂等。
  Future<InteractionServiceV1WatchResponse?> unwatch(int postId) async {
    try {
      return await _api
          .unwatch(InteractionServiceV1WatchRequest(postId: postId));
    } on DioException catch (e) {
      handleDioError(e);
      return null;
    }
  }

  // ─── 状态查询 ──────────────────────────────────────

  /// 批量查询当前 viewer 对指定目标的交互状态。
  Future<InteractionServiceV1GetInteractionStatusResponse?> getStatus(
      InteractionServiceV1TargetType targetType, List<int> targetIds) async {
    try {
      return await _api.getInteractionStatus(
          InteractionServiceV1GetInteractionStatusRequest(
              targetType: targetType, targetIds: targetIds));
    } on DioException catch (e) {
      handleDioError(e);
      return null;
    }
  }

  // ─── 计数查询 ──────────────────────────────────────

  /// 批量查询指定目标的计数（如点赞数）。
  ///
  /// 计数统一存于 interaction_counter 表（post/comment 表的散落计数列已移除），
  /// 列表/详情渲染计数改读本接口。响应中未出现的 (target, metric) 组合按 0 兜底，
  /// 由 [extractCount] / [extractCountMap] 处理。
  Future<InteractionServiceV1GetCountsResponse?> getCounts(
      InteractionServiceV1TargetType targetType,
      List<int> targetIds,
      List<InteractionServiceV1CounterMetric> metrics) async {
    try {
      return await _api.getCounts(InteractionServiceV1GetCountsRequest(
          targetType: targetType,
          targetIds: targetIds,
          metrics: metrics));
    } on DioException catch (e) {
      handleDioError(e);
      return null;
    }
  }

  /// 从 [GetCountsResponse] 中提取单个 (targetId, metric) 的计数。
  /// 未出现则返回 0。逻辑对应 React 版 `extractCount`。
  static int extractCount(
      InteractionServiceV1GetCountsResponse? resp,
      int targetId,
      InteractionServiceV1CounterMetric metric) {
    final counts = resp?.counts;
    if (counts == null) return 0;
    final cm = counts[targetId.toString()];
    final mcList = cm?.counts;
    if (mcList == null) return 0;
    for (final mc in mcList) {
      if (mc.metric == metric && mc.count != null) {
        return mc.count!;
      }
    }
    return 0;
  }

  /// 批量构建 targetId → 计数 的映射，供列表页一次性渲染。
  /// 未出现在响应中的目标按 0 兜底。
  static Map<int, int> extractCountMap(
      InteractionServiceV1GetCountsResponse? resp,
      List<int> targetIds,
      InteractionServiceV1CounterMetric metric) {
    final m = <int, int>{};
    for (final id in targetIds) {
      m[id] = extractCount(resp, id, metric);
    }
    return m;
  }

  // ─── 收藏列表 ──────────────────────────────────────

  /// 列出当前 viewer 收藏的 post。
  Future<ContentServiceV1ListPostResponse?> listWatchedPosts(
      [PaginationQuery? query]) async {
    final q = query ?? const PaginationQuery();
    try {
      return await _api.listWatchedPosts(q.toPagingRequest());
    } on DioException catch (e) {
      handleDioError(e);
      return null;
    }
  }
}
