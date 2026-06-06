// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_bind_email_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1BindEmailRequest _$IdentityServiceV1BindEmailRequestFromJson(
  Map<String, dynamic> json,
) => IdentityServiceV1BindEmailRequest(
  email: json['email'] as String?,
  verificationCode: json['verificationCode'] as String?,
);

Map<String, dynamic> _$IdentityServiceV1BindEmailRequestToJson(
  IdentityServiceV1BindEmailRequest instance,
) => <String, dynamic>{
  'email': instance.email,
  'verificationCode': instance.verificationCode,
};
