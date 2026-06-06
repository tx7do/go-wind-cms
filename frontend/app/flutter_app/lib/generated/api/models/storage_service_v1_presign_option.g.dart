// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'storage_service_v1_presign_option.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

StorageServiceV1PresignOption _$StorageServiceV1PresignOptionFromJson(
  Map<String, dynamic> json,
) => StorageServiceV1PresignOption(
  method: json['method'] as String?,
  expireSeconds: (json['expireSeconds'] as num?)?.toInt(),
  contentType: json['contentType'] as String?,
);

Map<String, dynamic> _$StorageServiceV1PresignOptionToJson(
  StorageServiceV1PresignOption instance,
) => <String, dynamic>{
  'method': instance.method,
  'expireSeconds': instance.expireSeconds,
  'contentType': instance.contentType,
};
