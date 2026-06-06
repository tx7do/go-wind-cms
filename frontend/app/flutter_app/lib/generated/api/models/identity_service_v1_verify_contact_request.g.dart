// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_verify_contact_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1VerifyContactRequest
_$IdentityServiceV1VerifyContactRequestFromJson(Map<String, dynamic> json) =>
    IdentityServiceV1VerifyContactRequest(
      phone: json['phone'] == null
          ? null
          : IdentityServiceV1PhoneVerification.fromJson(
              json['phone'] as Map<String, dynamic>,
            ),
      email: json['email'] == null
          ? null
          : IdentityServiceV1EmailVerification.fromJson(
              json['email'] as Map<String, dynamic>,
            ),
      verificationId: json['verificationId'] as String?,
    );

Map<String, dynamic> _$IdentityServiceV1VerifyContactRequestToJson(
  IdentityServiceV1VerifyContactRequest instance,
) => <String, dynamic>{
  'phone': instance.phone,
  'email': instance.email,
  'verificationId': instance.verificationId,
};
