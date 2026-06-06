// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'storage_service_v1_upload_file_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

StorageServiceV1UploadFileRequest _$StorageServiceV1UploadFileRequestFromJson(
  Map<String, dynamic> json,
) => StorageServiceV1UploadFileRequest(
  storageObject: json['storageObject'] == null
      ? null
      : StorageServiceV1StorageObject.fromJson(
          json['storageObject'] as Map<String, dynamic>,
        ),
  file: json['file'] as String?,
  presign: json['presign'] == null
      ? null
      : StorageServiceV1PresignOption.fromJson(
          json['presign'] as Map<String, dynamic>,
        ),
  sourceFileName: json['sourceFileName'] as String?,
  mime: json['mime'] as String?,
  size: json['size'] as String?,
  tenantId: (json['tenantId'] as num?)?.toInt(),
  userId: (json['userId'] as num?)?.toInt(),
);

Map<String, dynamic> _$StorageServiceV1UploadFileRequestToJson(
  StorageServiceV1UploadFileRequest instance,
) => <String, dynamic>{
  'storageObject': instance.storageObject,
  'file': instance.file,
  'presign': instance.presign,
  'sourceFileName': instance.sourceFileName,
  'mime': instance.mime,
  'size': instance.size,
  'tenantId': instance.tenantId,
  'userId': instance.userId,
};
