// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'storage_service_v1_presign_option.dart';
import 'storage_service_v1_storage_object.dart';

part 'storage_service_v1_upload_file_request.g.dart';

@JsonSerializable()
class StorageServiceV1UploadFileRequest {
  const StorageServiceV1UploadFileRequest({
    this.storageObject,
    this.file,
    this.presign,
    this.sourceFileName,
    this.mime,
    this.size,
    this.tenantId,
    this.userId,
  });
  
  factory StorageServiceV1UploadFileRequest.fromJson(Map<String, Object?> json) => _$StorageServiceV1UploadFileRequestFromJson(json);
  
  final StorageServiceV1StorageObject? storageObject;

  /// 直接上传的文件内容（内联字节）
  final String? file;
  final StorageServiceV1PresignOption? presign;

  /// 原文件文件名
  final String? sourceFileName;

  /// 文件的MIME类型
  final String? mime;

  /// 文件大小（字节）
  final String? size;

  /// 租户ID，0代表系统全局角色
  final int? tenantId;

  /// 用户ID
  final int? userId;

  Map<String, Object?> toJson() => _$StorageServiceV1UploadFileRequestToJson(this);
}
