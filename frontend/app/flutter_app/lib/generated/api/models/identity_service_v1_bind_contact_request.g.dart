// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_bind_contact_request.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1BindContactRequest
_$IdentityServiceV1BindContactRequestFromJson(Map<String, dynamic> json) =>
    IdentityServiceV1BindContactRequest(
      phone: json['phone'] == null
          ? null
          : IdentityServiceV1BindPhoneRequest.fromJson(
              json['phone'] as Map<String, dynamic>,
            ),
      email: json['email'] == null
          ? null
          : IdentityServiceV1BindEmailRequest.fromJson(
              json['email'] as Map<String, dynamic>,
            ),
    );

Map<String, dynamic> _$IdentityServiceV1BindContactRequestToJson(
  IdentityServiceV1BindContactRequest instance,
) => <String, dynamic>{'phone': instance.phone, 'email': instance.email};
