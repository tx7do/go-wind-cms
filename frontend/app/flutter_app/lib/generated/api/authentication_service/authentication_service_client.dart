// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';

import '../models/authentication_service_v1_login_request.dart';
import '../models/authentication_service_v1_login_response.dart';

part 'authentication_service_client.g.dart';

@RestApi()
abstract class AuthenticationServiceClient {
  factory AuthenticationServiceClient(Dio dio, {String? baseUrl}) = _AuthenticationServiceClient;

  /// 登录
  @POST('/app/v1/login')
  Future<AuthenticationServiceV1LoginResponse> authenticationServiceLogin({
    @Body() required AuthenticationServiceV1LoginRequest body,
  });

  /// 登出
  @POST('/app/v1/logout')
  Future<void> authenticationServiceLogout();

  /// 刷新认证令牌
  @POST('/app/v1/refresh-token')
  Future<AuthenticationServiceV1LoginResponse> authenticationServiceRefreshToken({
    @Body() required AuthenticationServiceV1LoginRequest body,
  });
}
