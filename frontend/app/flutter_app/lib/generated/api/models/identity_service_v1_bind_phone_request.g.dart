// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_bind_phone_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1BindPhoneRequest _$IdentityServiceV1BindPhoneRequestFromJson(
  Map<String, dynamic> json,
) => IdentityServiceV1BindPhoneRequest(
  phone: json['phone'] as String?,
  code: json['code'] as String?,
);

Map<String, dynamic> _$IdentityServiceV1BindPhoneRequestToJson(
  IdentityServiceV1BindPhoneRequest instance,
) => <String, dynamic>{'phone': instance.phone, 'code': instance.code};
