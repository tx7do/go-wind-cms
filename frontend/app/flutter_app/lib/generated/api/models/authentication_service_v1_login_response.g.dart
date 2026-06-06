// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'authentication_service_v1_login_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AuthenticationServiceV1LoginResponse
_$AuthenticationServiceV1LoginResponseFromJson(Map<String, dynamic> json) =>
    AuthenticationServiceV1LoginResponse(
      accessToken: json['access_token'] as String?,
      expiresIn: json['expires_in'] as String?,
      refreshToken: json['refresh_token'] as String?,
      scope: json['scope'] as String?,
      refreshExpiresIn: json['refresh_expires_in'] as String?,
      idToken: json['id_token'] as String?,
      tokenType: json['token_type'] == null
          ? AuthenticationServiceV1LoginResponseTokenType.bearer
          : AuthenticationServiceV1LoginResponseTokenType.fromJson(
              json['token_type'] as String,
            ),
    );

Map<String, dynamic> _$AuthenticationServiceV1LoginResponseToJson(
  AuthenticationServiceV1LoginResponse instance,
) => <String, dynamic>{
  'token_type':
      _$AuthenticationServiceV1LoginResponseTokenTypeEnumMap[instance
          .tokenType]!,
  'access_token': instance.accessToken,
  'expires_in': instance.expiresIn,
  'refresh_token': instance.refreshToken,
  'scope': instance.scope,
  'refresh_expires_in': instance.refreshExpiresIn,
  'id_token': instance.idToken,
};

const _$AuthenticationServiceV1LoginResponseTokenTypeEnumMap = {
  AuthenticationServiceV1LoginResponseTokenType.bearer: 'bearer',
  AuthenticationServiceV1LoginResponseTokenType.mac: 'mac',
  AuthenticationServiceV1LoginResponseTokenType.$unknown: r'$unknown',
};
