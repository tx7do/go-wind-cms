// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'identity_service_v1_phone_verification.g.dart';

/// 手机验证
@JsonSerializable()
class IdentityServiceV1PhoneVerification {
  const IdentityServiceV1PhoneVerification({
    required this.code,
    this.phone,
  });
  
  factory IdentityServiceV1PhoneVerification.fromJson(Map<String, Object?> json) => _$IdentityServiceV1PhoneVerificationFromJson(json);
  
  /// 手机号码，带国家码
  final String? phone;

  /// 短信验证码
  final String code;

  Map<String, Object?> toJson() => _$IdentityServiceV1PhoneVerificationToJson(this);
}
