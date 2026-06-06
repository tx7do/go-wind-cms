// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'identity_service_v1_user.dart';

part 'identity_service_v1_update_user_request.g.dart';

/// 更新用户 - 请求
@JsonSerializable()
class IdentityServiceV1UpdateUserRequest {
  const IdentityServiceV1UpdateUserRequest({
    required this.data,
    this.id,
    this.password,
    this.updateMask,
    this.allowMissing,
  });
  
  factory IdentityServiceV1UpdateUserRequest.fromJson(Map<String, Object?> json) => _$IdentityServiceV1UpdateUserRequestFromJson(json);
  
  final int? id;
  final IdentityServiceV1User data;

  /// 用户登录密码
  final String? password;

  /// 要更新的字段列表
  final String? updateMask;

  /// 如果设置为true的时候，资源不存在则会新增(插入)，并且在这种情况下`updateMask`字段将会被忽略。
  final bool? allowMissing;

  Map<String, Object?> toJson() => _$IdentityServiceV1UpdateUserRequestToJson(this);
}
