import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart' show Mutation;

import 'package:flutter_app/generated/api/authentication_service/authentication_service_client.dart'
    show AuthenticationServiceClient;
import 'package:flutter_app/generated/api/rest_client.dart' show RestClient;
import 'package:flutter_app/src/app_router/app_router.dart' show AppRouter;

import 'package:flutter_app/generated/api/models/authentication_service_v1_login_request.dart';
import 'package:flutter_app/generated/api/models/authentication_service_v1_login_request_grant_type.dart';
import 'package:flutter_app/generated/api/models/authentication_service_v1_login_response.dart';
import 'package:flutter_app/src/core/repositories/user_auth_cache.dart'
    show UserAuthCache;

import 'package:flutter_app/src/core/transport/http/index.dart';

export 'package:get_it/get_it.dart';

/// 认证服务
///
/// 使用 Mutation 管理登录/登出/刷新令牌等写操作，
/// 自动获得加载状态、错误处理、缓存失效等能力。
class AuthenticationService extends AuthService {
  AuthenticationServiceClient get _api =>
      GetIt.instance<RestClient>().authenticationService;

  /// 获取当前访问令牌
  @override
  String? getAccessToken() {
    return GetIt.instance<UserAuthCache>().accessToken;
  }

  /// 获取刷新令牌
  @override
  String? getRefreshToken() {
    return GetIt.instance<UserAuthCache>().refreshToken;
  }

  /// 刷新访问令牌（内部方法，供 AuthenticationInterceptor 调用）
  @override
  Future<String?> refreshToken() async {
    final token = getRefreshToken();
    if (token == null || token.isEmpty) {
      doAuthenticationFailed();
      return null;
    }

    final result = await _doRefreshToken(token);
    if (result is AuthenticationServiceV1LoginResponse) {
      return result.accessToken;
    }

    doAuthenticationFailed();
    return null;
  }

  @override
  authenticationFailed() {
    doAuthenticationFailed();
  }

  /// 清除所有令牌
  clearTokens() async {
    await GetIt.instance<UserAuthCache>().clearTokens();
  }

  /// 保存新的令牌
  saveTokens(String accessToken, {String? refreshToken}) async {
    await GetIt.instance<UserAuthCache>().saveAuthInfo(
      accessToken,
      refreshToken: refreshToken,
    );
  }

  // ─── Mutations ───────────────────────────────────────

  /// 用户登录 Mutation
  ///
  /// 用法：
  /// ```dart
  /// final mutation = authService.loginMutation();
  /// // 在 UI 中通过 MutationBuilder 监听状态
  /// // 触发：mutation.mutate(LoginParams(username: 'xxx', password: 'xxx'));
  /// ```
  Mutation<AuthenticationServiceV1LoginResponse, LoginParams> loginMutation() {
    return Mutation<AuthenticationServiceV1LoginResponse, LoginParams>(
      mutationFn: (params) async {
        final request = AuthenticationServiceV1LoginRequest(
          grantType: AuthenticationServiceV1LoginRequestGrantType.password,
          username: params.username,
          password: params.password,
        );
        return _api.authenticationServiceLogin(body: request);
      },
      onSuccess: (response, _) => doLoginSuccess(response),
    );
  }

  /// 用户登出 Mutation
  Mutation<void, void> logoutMutation() {
    return Mutation<void, void>(
      mutationFn: (_) async {
        await _api.authenticationServiceLogout();
        await GetIt.instance<UserAuthCache>().clearTokens();
      },
    );
  }

  /// 刷新令牌 Mutation
  Mutation<AuthenticationServiceV1LoginResponse, String>
  refreshTokenMutation() {
    return Mutation<AuthenticationServiceV1LoginResponse, String>(
      mutationFn: (refreshToken) async {
        final request = AuthenticationServiceV1LoginRequest(
          grantType: AuthenticationServiceV1LoginRequestGrantType.refreshToken,
          refreshToken: refreshToken,
        );
        return _api.authenticationServiceRefreshToken(body: request);
      },
      onSuccess: (response, _) => doLoginSuccess(response),
    );
  }

  /// 直接调用刷新令牌（供 AuthenticationInterceptor 使用）
  Future<dynamic> _doRefreshToken(String refreshToken) async {
    try {
      final request = AuthenticationServiceV1LoginRequest(
        grantType: AuthenticationServiceV1LoginRequestGrantType.refreshToken,
        refreshToken: refreshToken,
      );
      final response = await _api.authenticationServiceRefreshToken(
        body: request,
      );
      doLoginSuccess(response);
      return response;
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  // ─── 内部方法 ────────────────────────────────────────

  /// 认证成功处理
  doLoginSuccess(AuthenticationServiceV1LoginResponse msg) {
    GetIt.instance<UserAuthCache>().saveAuthInfo(
      msg.accessToken ?? '',
      refreshToken: msg.refreshToken,
    );
  }

  /// 认证失败处理
  doAuthenticationFailed() {
    GetIt.instance<UserAuthCache>().clearTokens();
    AppRouter.router.goNamed('home');
  }
}

/// 登录参数
class LoginParams {
  final String username;
  final String password;

  const LoginParams({required this.username, required this.password});
}
