// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'identity_service_v1_upload_avatar_request.g.dart';

@JsonSerializable()
class IdentityServiceV1UploadAvatarRequest {
  const IdentityServiceV1UploadAvatarRequest({
    this.imageBase64,
    this.imageUrl,
  });
  
  factory IdentityServiceV1UploadAvatarRequest.fromJson(Map<String, Object?> json) => _$IdentityServiceV1UploadAvatarRequestFromJson(json);
  
  /// 图片的Base64编码
  final String? imageBase64;

  /// 图片的URL地址
  final String? imageUrl;

  Map<String, Object?> toJson() => _$IdentityServiceV1UploadAvatarRequestToJson(this);
}
