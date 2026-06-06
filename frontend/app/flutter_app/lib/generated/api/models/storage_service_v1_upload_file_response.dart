// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'storage_service_v1_upload_file_response.g.dart';

@JsonSerializable()
class StorageServiceV1UploadFileResponse {
  const StorageServiceV1UploadFileResponse({
    this.objectName,
    this.presignedUrl,
  });
  
  factory StorageServiceV1UploadFileResponse.fromJson(Map<String, Object?> json) => _$StorageServiceV1UploadFileResponseFromJson(json);
  
  /// OSS 对象键
  final String? objectName;

  /// 预签名上传链接
  final String? presignedUrl;

  Map<String, Object?> toJson() => _$StorageServiceV1UploadFileResponseToJson(this);
}
