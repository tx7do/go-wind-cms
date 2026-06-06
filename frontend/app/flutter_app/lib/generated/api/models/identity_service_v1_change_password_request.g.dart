// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_change_password_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1ChangePasswordRequest
_$IdentityServiceV1ChangePasswordRequestFromJson(Map<String, dynamic> json) =>
    IdentityServiceV1ChangePasswordRequest(
      oldPassword: json['oldPassword'] as String?,
      newPassword: json['newPassword'] as String?,
    );

Map<String, dynamic> _$IdentityServiceV1ChangePasswordRequestToJson(
  IdentityServiceV1ChangePasswordRequest instance,
) => <String, dynamic>{
  'oldPassword': instance.oldPassword,
  'newPassword': instance.newPassword,
};
