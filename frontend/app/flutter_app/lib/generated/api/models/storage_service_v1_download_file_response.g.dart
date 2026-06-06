// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'storage_service_v1_download_file_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

StorageServiceV1DownloadFileResponse
_$StorageServiceV1DownloadFileResponseFromJson(Map<String, dynamic> json) =>
    StorageServiceV1DownloadFileResponse(
      file: json['file'] as String?,
      downloadUrl: json['downloadUrl'] as String?,
      sourceFileName: json['sourceFileName'] as String?,
      mime: json['mime'] as String?,
      size: json['size'] as String?,
      checksum: json['checksum'] as String?,
      storagePath: json['storagePath'] as String?,
      updatedAt: json['updatedAt'] == null
          ? null
          : DateTime.parse(json['updatedAt'] as String),
    );

Map<String, dynamic> _$StorageServiceV1DownloadFileResponseToJson(
  StorageServiceV1DownloadFileResponse instance,
) => <String, dynamic>{
  'file': instance.file,
  'downloadUrl': instance.downloadUrl,
  'sourceFileName': instance.sourceFileName,
  'mime': instance.mime,
  'size': instance.size,
  'checksum': instance.checksum,
  'storagePath': instance.storagePath,
  'updatedAt': instance.updatedAt?.toIso8601String(),
};
