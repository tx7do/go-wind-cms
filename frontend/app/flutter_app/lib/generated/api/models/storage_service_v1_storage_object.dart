// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'storage_service_v1_storage_object.g.dart';

/// 对象存储对象
@JsonSerializable()
class StorageServiceV1StorageObject {
  const StorageServiceV1StorageObject({
    this.bucketName,
    this.fileDirectory,
    this.objectName,
  });
  
  factory StorageServiceV1StorageObject.fromJson(Map<String, Object?> json) => _$StorageServiceV1StorageObjectFromJson(json);
  
  /// 目标文件桶名称
  final String? bucketName;

  /// 存储目录（如 'user/1001/avatar/'），服务端会在此目录下生成唯一文件名
  final String? fileDirectory;

  /// OSS 对象键（完整路径，如 'user/1001/avatar.jpg'）。若未提供，服务端将自动生成。
  final String? objectName;

  Map<String, Object?> toJson() => _$StorageServiceV1StorageObjectToJson(this);
}
