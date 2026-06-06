// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_phone_verification.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1PhoneVerification _$IdentityServiceV1PhoneVerificationFromJson(
  Map<String, dynamic> json,
) => IdentityServiceV1PhoneVerification(
  code: json['code'] as String,
  phone: json['phone'] as String?,
);

Map<String, dynamic> _$IdentityServiceV1PhoneVerificationToJson(
  IdentityServiceV1PhoneVerification instance,
) => <String, dynamic>{'phone': instance.phone, 'code': instance.code};
