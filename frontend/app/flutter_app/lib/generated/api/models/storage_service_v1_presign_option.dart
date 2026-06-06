// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'storage_service_v1_presign_option.g.dart';

/// 预签名选项
@JsonSerializable()
class StorageServiceV1PresignOption {
  const StorageServiceV1PresignOption({
    this.method,
    this.expireSeconds,
    this.contentType,
  });
  
  factory StorageServiceV1PresignOption.fromJson(Map<String, Object?> json) => _$StorageServiceV1PresignOptionFromJson(json);
  
  /// 预签名使用的 HTTP 方法，例如 PUT 或 POST
  final String? method;

  /// 预签名有效期（秒）
  final int? expireSeconds;

  /// 预签名约束的 Content-Type（可选）
  final String? contentType;

  Map<String, Object?> toJson() => _$StorageServiceV1PresignOptionToJson(this);
}
