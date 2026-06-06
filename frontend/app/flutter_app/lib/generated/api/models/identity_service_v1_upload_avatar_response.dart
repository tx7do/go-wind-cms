// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'identity_service_v1_upload_avatar_response.g.dart';

@JsonSerializable()
class IdentityServiceV1UploadAvatarResponse {
  const IdentityServiceV1UploadAvatarResponse({
    this.url,
  });
  
  factory IdentityServiceV1UploadAvatarResponse.fromJson(Map<String, Object?> json) => _$IdentityServiceV1UploadAvatarResponseFromJson(json);
  
  final String? url;

  Map<String, Object?> toJson() => _$IdentityServiceV1UploadAvatarResponseToJson(this);
}
