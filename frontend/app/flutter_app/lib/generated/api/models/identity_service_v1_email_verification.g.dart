// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_email_verification.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1EmailVerification _$IdentityServiceV1EmailVerificationFromJson(
  Map<String, dynamic> json,
) => IdentityServiceV1EmailVerification(
  code: json['code'] as String,
  email: json['email'] as String?,
);

Map<String, dynamic> _$IdentityServiceV1EmailVerificationToJson(
  IdentityServiceV1EmailVerification instance,
) => <String, dynamic>{'email': instance.email, 'code': instance.code};
