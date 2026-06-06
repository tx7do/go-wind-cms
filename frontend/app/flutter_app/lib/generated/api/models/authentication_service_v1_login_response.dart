// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'authentication_service_v1_login_response_token_type.dart';

part 'authentication_service_v1_login_response.g.dart';

/// 用户登录 - 回应
@JsonSerializable()
class AuthenticationServiceV1LoginResponse {
  const AuthenticationServiceV1LoginResponse({
    this.accessToken,
    this.expiresIn,
    this.refreshToken,
    this.scope,
    this.refreshExpiresIn,
    this.idToken,
    this.tokenType = AuthenticationServiceV1LoginResponseTokenType.bearer,
  });
  
  factory AuthenticationServiceV1LoginResponse.fromJson(Map<String, Object?> json) => _$AuthenticationServiceV1LoginResponseFromJson(json);
  
  /// 令牌的类型，该值大小写不敏感，必选项，可以是bearer类型或mac类型，通常只是字符串“Bearer”。
  @JsonKey(name: 'token_type')
  final AuthenticationServiceV1LoginResponseTokenType tokenType;

  /// 访问令牌，必选项。授权服务器颁发的访问令牌字符串。
  @JsonKey(name: 'access_token')
  final String? accessToken;

  /// 访问令牌过期时间（秒）
  @JsonKey(name: 'expires_in')
  final String? expiresIn;

  /// 更新令牌，用来获取下一次的访问令牌，可选项。如果访问令牌将过期，则返回刷新令牌很有用，应用程序可以使用该刷新令牌来获取另一个访问令牌。但是，通过隐式授予颁发的令牌不能颁发刷新令牌。
  @JsonKey(name: 'refresh_token')
  final String? refreshToken;

  /// 以空格分隔的用户授予范围列表。如果未提供，scope则授权任何范围，默认为空列表。
  final String? scope;

  /// 刷新令牌过期时间（秒）
  @JsonKey(name: 'refresh_expires_in')
  final String? refreshExpiresIn;

  /// ID 令牌，OpenID Connect 扩展中定义的 JWT 格式令牌
  @JsonKey(name: 'id_token')
  final String? idToken;

  Map<String, Object?> toJson() => _$AuthenticationServiceV1LoginResponseToJson(this);
}
