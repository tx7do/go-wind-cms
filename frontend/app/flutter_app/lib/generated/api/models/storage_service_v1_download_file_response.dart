// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'storage_service_v1_download_file_response.g.dart';

/// 文件下载响应
@JsonSerializable()
class StorageServiceV1DownloadFileResponse {
  const StorageServiceV1DownloadFileResponse({
    this.file,
    this.downloadUrl,
    this.sourceFileName,
    this.mime,
    this.size,
    this.checksum,
    this.storagePath,
    this.updatedAt,
  });
  
  factory StorageServiceV1DownloadFileResponse.fromJson(Map<String, Object?> json) => _$StorageServiceV1DownloadFileResponseFromJson(json);
  
  /// 直接返回文件字节（小文件或同步场景）
  final String? file;

  /// 预签名 URL（大文件或异步下载）
  final String? downloadUrl;

  /// 原始文件名，客户端用于保存或显示
  final String? sourceFileName;

  /// MIME 类型，默认 application/octet-stream
  final String? mime;

  /// 文件大小（字节），可用于进度或校验
  final String? size;

  /// 可选，文件校验（如 MD5 或 SHA256）
  final String? checksum;

  /// 可选，后端存储位置或 key
  final String? storagePath;

  /// 可选，最后修改时间
  final DateTime? updatedAt;

  Map<String, Object?> toJson() => _$StorageServiceV1DownloadFileResponseToJson(this);
}
