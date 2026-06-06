// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'identity_service_v1_bind_email_request.dart';
import 'identity_service_v1_bind_phone_request.dart';

part 'identity_service_v1_bind_contact_request.g.dart';

@JsonSerializable()
class IdentityServiceV1BindContactRequest {
  const IdentityServiceV1BindContactRequest({
    this.phone,
    this.email,
  });
  
  factory IdentityServiceV1BindContactRequest.fromJson(Map<String, Object?> json) => _$IdentityServiceV1BindContactRequestFromJson(json);
  
  final IdentityServiceV1BindPhoneRequest? phone;
  final IdentityServiceV1BindEmailRequest? email;

  Map<String, Object?> toJson() => _$IdentityServiceV1BindContactRequestToJson(this);
}
