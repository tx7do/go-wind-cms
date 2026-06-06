// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_update_user_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1UpdateUserRequest _$IdentityServiceV1UpdateUserRequestFromJson(
  Map<String, dynamic> json,
) => IdentityServiceV1UpdateUserRequest(
  data: IdentityServiceV1User.fromJson(json['data'] as Map<String, dynamic>),
  id: (json['id'] as num?)?.toInt(),
  password: json['password'] as String?,
  updateMask: json['updateMask'] as String?,
  allowMissing: json['allowMissing'] as bool?,
);

Map<String, dynamic> _$IdentityServiceV1UpdateUserRequestToJson(
  IdentityServiceV1UpdateUserRequest instance,
) => <String, dynamic>{
  'id': instance.id,
  'data': instance.data,
  'password': instance.password,
  'updateMask': instance.updateMask,
  'allowMissing': instance.allowMissing,
};
