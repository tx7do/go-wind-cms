// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'identity_service_v1_email_verification.g.dart';

/// 邮箱验证
@JsonSerializable()
class IdentityServiceV1EmailVerification {
  const IdentityServiceV1EmailVerification({
    required this.code,
    this.email,
  });
  
  factory IdentityServiceV1EmailVerification.fromJson(Map<String, Object?> json) => _$IdentityServiceV1EmailVerificationFromJson(json);
  
  /// 邮箱地址
  final String? email;

  /// 邮箱验证码
  final String code;

  Map<String, Object?> toJson() => _$IdentityServiceV1EmailVerificationToJson(this);
}
