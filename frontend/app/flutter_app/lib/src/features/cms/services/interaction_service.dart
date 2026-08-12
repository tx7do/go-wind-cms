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
        ContentServiceV1ListPostResponse;
import 'package:flutter_app/src/core/services/pagination_query.dart';
import 'package:flutter_app/src/core/services/base_service.dart';

/// 交互服务（点赞 / 收藏）
///
/// 作为 InteractionService 的客户端封装，对应后端点赞/收藏计数内聚子系统。
/// viewer 身份由后端鉴权上下文确定，前端无需传 user_id。
/// 计数（likes）由后端在 ledger 写入的同一事务内递增/递减并回读，前端直读返回值。
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
