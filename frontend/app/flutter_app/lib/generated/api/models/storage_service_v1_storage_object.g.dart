// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'storage_service_v1_storage_object.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

StorageServiceV1StorageObject _$StorageServiceV1StorageObjectFromJson(
  Map<String, dynamic> json,
) => StorageServiceV1StorageObject(
  bucketName: json['bucketName'] as String?,
  fileDirectory: json['fileDirectory'] as String?,
  objectName: json['objectName'] as String?,
);

Map<String, dynamic> _$StorageServiceV1StorageObjectToJson(
  StorageServiceV1StorageObject instance,
) => <String, dynamic>{
  'bucketName': instance.bucketName,
  'fileDirectory': instance.fileDirectory,
  'objectName': instance.objectName,
};
