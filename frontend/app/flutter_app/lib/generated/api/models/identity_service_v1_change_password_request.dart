// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'identity_service_v1_change_password_request.g.dart';

/// 修改用户密码（需要验证旧密码） - 请求
@JsonSerializable()
class IdentityServiceV1ChangePasswordRequest {
  const IdentityServiceV1ChangePasswordRequest({
    this.oldPassword,
    this.newPassword,
  });
  
  factory IdentityServiceV1ChangePasswordRequest.fromJson(Map<String, Object?> json) => _$IdentityServiceV1ChangePasswordRequestFromJson(json);
  
  /// 旧密码
  final String? oldPassword;

  /// 新密码
  final String? newPassword;

  Map<String, Object?> toJson() => _$IdentityServiceV1ChangePasswordRequestToJson(this);
}
