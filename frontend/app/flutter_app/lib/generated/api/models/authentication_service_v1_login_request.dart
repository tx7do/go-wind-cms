// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'authentication_service_v1_login_request_client_type.dart';
import 'authentication_service_v1_login_request_grant_type.dart';

part 'authentication_service_v1_login_request.g.dart';

/// 用户登录 - 请求
@JsonSerializable()
class AuthenticationServiceV1LoginRequest {
  const AuthenticationServiceV1LoginRequest({
    this.clientId,
    this.clientSecret,
    this.scope,
    this.redirectUri,
    this.userId,
    this.username,
    this.email,
    this.mobile,
    this.password,
    this.refreshToken,
    this.code,
    this.clientType,
    this.deviceId,
    this.jti,
    this.grantType = AuthenticationServiceV1LoginRequestGrantType.password,
  });
  
  factory AuthenticationServiceV1LoginRequest.fromJson(Map<String, Object?> json) => _$AuthenticationServiceV1LoginRequestFromJson(json);
  
  /// 授权类型，此处的值固定为"password"，必选项。
  @JsonKey(name: 'grant_type')
  final AuthenticationServiceV1LoginRequestGrantType grantType;

  /// 客户端ID
  @JsonKey(name: 'client_id')
  final String? clientId;

  /// 客户端密钥
  @JsonKey(name: 'client_secret')
  final String? clientSecret;

  /// 以空格分隔的用户授予范围列表。如果未提供，scope则授权任何范围，默认为空列表。
  final String? scope;

  /// 跳转链接
  @JsonKey(name: 'redirect_uri')
  final String? redirectUri;

  /// 用户ID
  @JsonKey(name: 'user_id')
  final int? userId;

  /// 用户名
  final String? username;

  /// 电子邮件地址
  final String? email;

  /// 手机号
  final String? mobile;

  /// 用户的密码
  final String? password;

  /// 更新令牌，用来获取下一次的访问令牌，可选项。如果访问令牌将过期，则返回刷新令牌很有用，应用程序可以使用该刷新令牌来获取另一个访问令牌。但是，通过隐式授予颁发的令牌不能颁发刷新令牌。
  @JsonKey(name: 'refresh_token')
  final String? refreshToken;

  /// 授权请求中收到的一次性验证/认证码。(当使用授权码模式时)
  final String? code;

  /// 客户端类型
  @JsonKey(name: 'client_type')
  final AuthenticationServiceV1LoginRequestClientType? clientType;

  /// 设备唯一标识（可选），用于设备绑定、推送、风控等
  @JsonKey(name: 'device_id')
  final String? deviceId;

  /// 建议客户端生成并提供 jti（JWT ID）作为唯一标识，服务端可据此防止重放攻击
  final String? jti;

  Map<String, Object?> toJson() => _$AuthenticationServiceV1LoginRequestToJson(this);
}
