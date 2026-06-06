// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

part 'identity_service_v1_bind_phone_request.g.dart';

@JsonSerializable()
class IdentityServiceV1BindPhoneRequest {
  const IdentityServiceV1BindPhoneRequest({
    this.phone,
    this.code,
  });
  
  factory IdentityServiceV1BindPhoneRequest.fromJson(Map<String, Object?> json) => _$IdentityServiceV1BindPhoneRequestFromJson(json);
  
  /// 手机号码
  final String? phone;

  /// 验证码
  final String? code;

  Map<String, Object?> toJson() => _$IdentityServiceV1BindPhoneRequestToJson(this);
}
