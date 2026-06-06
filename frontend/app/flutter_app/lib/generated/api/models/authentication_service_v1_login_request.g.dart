// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'authentication_service_v1_login_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AuthenticationServiceV1LoginRequest
_$AuthenticationServiceV1LoginRequestFromJson(Map<String, dynamic> json) =>
    AuthenticationServiceV1LoginRequest(
      clientId: json['client_id'] as String?,
      clientSecret: json['client_secret'] as String?,
      scope: json['scope'] as String?,
      redirectUri: json['redirect_uri'] as String?,
      userId: (json['user_id'] as num?)?.toInt(),
      username: json['username'] as String?,
      email: json['email'] as String?,
      mobile: json['mobile'] as String?,
      password: json['password'] as String?,
      refreshToken: json['refresh_token'] as String?,
      code: json['code'] as String?,
      clientType: json['client_type'] == null
          ? null
          : AuthenticationServiceV1LoginRequestClientType.fromJson(
              json['client_type'] as String,
            ),
      deviceId: json['device_id'] as String?,
      jti: json['jti'] as String?,
      grantType: json['grant_type'] == null
          ? AuthenticationServiceV1LoginRequestGrantType.password
          : AuthenticationServiceV1LoginRequestGrantType.fromJson(
              json['grant_type'] as String,
            ),
    );

Map<String, dynamic> _$AuthenticationServiceV1LoginRequestToJson(
  AuthenticationServiceV1LoginRequest instance,
) => <String, dynamic>{
  'grant_type':
      _$AuthenticationServiceV1LoginRequestGrantTypeEnumMap[instance
          .grantType]!,
  'client_id': instance.clientId,
  'client_secret': instance.clientSecret,
  'scope': instance.scope,
  'redirect_uri': instance.redirectUri,
  'user_id': instance.userId,
  'username': instance.username,
  'email': instance.email,
  'mobile': instance.mobile,
  'password': instance.password,
  'refresh_token': instance.refreshToken,
  'code': instance.code,
  'client_type':
      _$AuthenticationServiceV1LoginRequestClientTypeEnumMap[instance
          .clientType],
  'device_id': instance.deviceId,
  'jti': instance.jti,
};

const _$AuthenticationServiceV1LoginRequestGrantTypeEnumMap = {
  AuthenticationServiceV1LoginRequestGrantType.password: 'password',
  AuthenticationServiceV1LoginRequestGrantType.clientCredentials:
      'client_credentials',
  AuthenticationServiceV1LoginRequestGrantType.authorizationCode:
      'authorization_code',
  AuthenticationServiceV1LoginRequestGrantType.refreshToken: 'refresh_token',
  AuthenticationServiceV1LoginRequestGrantType.implicit: 'implicit',
  AuthenticationServiceV1LoginRequestGrantType.$unknown: r'$unknown',
};

const _$AuthenticationServiceV1LoginRequestClientTypeEnumMap = {
  AuthenticationServiceV1LoginRequestClientType.admin: 'admin',
  AuthenticationServiceV1LoginRequestClientType.app: 'app',
  AuthenticationServiceV1LoginRequestClientType.$unknown: r'$unknown',
};
