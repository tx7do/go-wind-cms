// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_upload_avatar_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1UploadAvatarRequest
_$IdentityServiceV1UploadAvatarRequestFromJson(Map<String, dynamic> json) =>
    IdentityServiceV1UploadAvatarRequest(
      imageBase64: json['imageBase64'] as String?,
      imageUrl: json['imageUrl'] as String?,
    );

Map<String, dynamic> _$IdentityServiceV1UploadAvatarRequestToJson(
  IdentityServiceV1UploadAvatarRequest instance,
) => <String, dynamic>{
  'imageBase64': instance.imageBase64,
  'imageUrl': instance.imageUrl,
};
