// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'kratos_status.g.dart';

/// Kratos错误返回
@JsonSerializable()
class KratosStatus {
  const KratosStatus({
    this.code,
    this.message,
    this.reason,
    this.metadata,
  });
  
  factory KratosStatus.fromJson(Map<String, Object?> json) => _$KratosStatusFromJson(json);
  
  /// 错误码
  final num? code;

  /// 错误消息
  final String? message;

  /// 错误原因
  final String? reason;

  /// 元数据
  final dynamic metadata;

  Map<String, Object?> toJson() => _$KratosStatusToJson(this);
}
