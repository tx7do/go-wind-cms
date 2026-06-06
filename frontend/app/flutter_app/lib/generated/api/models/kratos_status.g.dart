// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'kratos_status.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

KratosStatus _$KratosStatusFromJson(Map<String, dynamic> json) => KratosStatus(
  code: json['code'] as num?,
  message: json['message'] as String?,
  reason: json['reason'] as String?,
  metadata: json['metadata'],
);

Map<String, dynamic> _$KratosStatusToJson(KratosStatus instance) =>
    <String, dynamic>{
      'code': instance.code,
      'message': instance.message,
      'reason': instance.reason,
      'metadata': instance.metadata,
    };
