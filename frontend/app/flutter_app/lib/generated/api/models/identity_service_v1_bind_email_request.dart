// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'identity_service_v1_bind_email_request.g.dart';

@JsonSerializable()
class IdentityServiceV1BindEmailRequest {
  const IdentityServiceV1BindEmailRequest({
    this.email,
    this.verificationCode,
  });
  
  factory IdentityServiceV1BindEmailRequest.fromJson(Map<String, Object?> json) => _$IdentityServiceV1BindEmailRequestFromJson(json);
  
  /// 邮箱地址
  final String? email;

  /// 邮箱验证码（可选）
  final String? verificationCode;

  Map<String, Object?> toJson() => _$IdentityServiceV1BindEmailRequestToJson(this);
}
