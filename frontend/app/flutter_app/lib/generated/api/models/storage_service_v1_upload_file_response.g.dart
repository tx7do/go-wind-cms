// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'storage_service_v1_upload_file_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

StorageServiceV1UploadFileResponse _$StorageServiceV1UploadFileResponseFromJson(
  Map<String, dynamic> json,
) => StorageServiceV1UploadFileResponse(
  objectName: json['objectName'] as String?,
  presignedUrl: json['presignedUrl'] as String?,
);

Map<String, dynamic> _$StorageServiceV1UploadFileResponseToJson(
  StorageServiceV1UploadFileResponse instance,
) => <String, dynamic>{
  'objectName': instance.objectName,
  'presignedUrl': instance.presignedUrl,
};
