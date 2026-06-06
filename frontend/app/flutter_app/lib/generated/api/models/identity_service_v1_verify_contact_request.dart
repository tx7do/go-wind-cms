// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'identity_service_v1_email_verification.dart';
import 'identity_service_v1_phone_verification.dart';

part 'identity_service_v1_verify_contact_request.g.dart';

@JsonSerializable()
class IdentityServiceV1VerifyContactRequest {
  const IdentityServiceV1VerifyContactRequest({
    this.phone,
    this.email,
    this.verificationId,
  });
  
  factory IdentityServiceV1VerifyContactRequest.fromJson(Map<String, Object?> json) => _$IdentityServiceV1VerifyContactRequestFromJson(json);
  
  final IdentityServiceV1PhoneVerification? phone;
  final IdentityServiceV1EmailVerification? email;

  /// 服务端生成的验证码会话ID（可选）
  final String? verificationId;

  Map<String, Object?> toJson() => _$IdentityServiceV1VerifyContactRequestToJson(this);
}
