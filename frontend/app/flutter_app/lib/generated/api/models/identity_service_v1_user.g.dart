// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'identity_service_v1_user.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

IdentityServiceV1User _$IdentityServiceV1UserFromJson(
  Map<String, dynamic> json,
) => IdentityServiceV1User(
  roles: (json['roles'] as List<dynamic>?)?.map((e) => e as String).toList(),
  tenantId: (json['tenantId'] as num?)?.toInt(),
  tenantName: json['tenantName'] as String?,
  orgUnitId: (json['orgUnitId'] as num?)?.toInt(),
  orgUnitIds: (json['orgUnitIds'] as List<dynamic>?)
      ?.map((e) => (e as num).toInt())
      .toList(),
  orgUnitName: json['orgUnitName'] as String?,
  orgUnitNames: (json['orgUnitNames'] as List<dynamic>?)
      ?.map((e) => e as String)
      .toList(),
  positionId: (json['positionId'] as num?)?.toInt(),
  positionIds: (json['positionIds'] as List<dynamic>?)
      ?.map((e) => (e as num).toInt())
      .toList(),
  positionName: json['positionName'] as String?,
  positionNames: (json['positionNames'] as List<dynamic>?)
      ?.map((e) => e as String)
      .toList(),
  roleId: (json['roleId'] as num?)?.toInt(),
  roleIds: (json['roleIds'] as List<dynamic>?)
      ?.map((e) => (e as num).toInt())
      .toList(),
  id: (json['id'] as num?)?.toInt(),
  roleNames: (json['roleNames'] as List<dynamic>?)
      ?.map((e) => e as String)
      .toList(),
  username: json['username'] as String?,
  nickname: json['nickname'] as String?,
  realname: json['realname'] as String?,
  avatar: json['avatar'] as String?,
  email: json['email'] as String?,
  mobile: json['mobile'] as String?,
  telephone: json['telephone'] as String?,
  gender: json['gender'] == null
      ? null
      : IdentityServiceV1UserGender.fromJson(json['gender'] as String),
  address: json['address'] as String?,
  region: json['region'] as String?,
  description: json['description'] as String?,
  remark: json['remark'] as String?,
  deletedAt: json['deletedAt'] == null
      ? null
      : DateTime.parse(json['deletedAt'] as String),
  lastLoginIp: json['lastLoginIp'] as String?,
  status: json['status'] == null
      ? null
      : IdentityServiceV1UserStatus.fromJson(json['status'] as String),
  lockedUntil: json['lockedUntil'] == null
      ? null
      : DateTime.parse(json['lockedUntil'] as String),
  followers: json['followers'] as String?,
  following: json['following'] as String?,
  postCount: json['postCount'] as String?,
  commentCount: json['commentCount'] as String?,
  likeCount: json['likeCount'] as String?,
  createdBy: (json['createdBy'] as num?)?.toInt(),
  updatedBy: (json['updatedBy'] as num?)?.toInt(),
  deletedBy: (json['deletedBy'] as num?)?.toInt(),
  createdAt: json['createdAt'] == null
      ? null
      : DateTime.parse(json['createdAt'] as String),
  updatedAt: json['updatedAt'] == null
      ? null
      : DateTime.parse(json['updatedAt'] as String),
  lastLoginAt: json['lastLoginAt'] == null
      ? null
      : DateTime.parse(json['lastLoginAt'] as String),
);

Map<String, dynamic> _$IdentityServiceV1UserToJson(
  IdentityServiceV1User instance,
) => <String, dynamic>{
  'id': instance.id,
  'tenantId': instance.tenantId,
  'tenantName': instance.tenantName,
  'orgUnitId': instance.orgUnitId,
  'orgUnitIds': instance.orgUnitIds,
  'orgUnitName': instance.orgUnitName,
  'orgUnitNames': instance.orgUnitNames,
  'positionId': instance.positionId,
  'positionIds': instance.positionIds,
  'positionName': instance.positionName,
  'positionNames': instance.positionNames,
  'roleId': instance.roleId,
  'roleIds': instance.roleIds,
  'roles': instance.roles,
  'roleNames': instance.roleNames,
  'username': instance.username,
  'nickname': instance.nickname,
  'realname': instance.realname,
  'avatar': instance.avatar,
  'email': instance.email,
  'mobile': instance.mobile,
  'telephone': instance.telephone,
  'gender': _$IdentityServiceV1UserGenderEnumMap[instance.gender],
  'address': instance.address,
  'region': instance.region,
  'description': instance.description,
  'remark': instance.remark,
  'lastLoginAt': instance.lastLoginAt?.toIso8601String(),
  'lastLoginIp': instance.lastLoginIp,
  'status': _$IdentityServiceV1UserStatusEnumMap[instance.status],
  'lockedUntil': instance.lockedUntil?.toIso8601String(),
  'followers': instance.followers,
  'following': instance.following,
  'postCount': instance.postCount,
  'commentCount': instance.commentCount,
  'likeCount': instance.likeCount,
  'createdBy': instance.createdBy,
  'updatedBy': instance.updatedBy,
  'deletedBy': instance.deletedBy,
  'createdAt': instance.createdAt?.toIso8601String(),
  'updatedAt': instance.updatedAt?.toIso8601String(),
  'deletedAt': instance.deletedAt?.toIso8601String(),
};

const _$IdentityServiceV1UserGenderEnumMap = {
  IdentityServiceV1UserGender.secret: 'SECRET',
  IdentityServiceV1UserGender.male: 'MALE',
  IdentityServiceV1UserGender.female: 'FEMALE',
  IdentityServiceV1UserGender.$unknown: r'$unknown',
};

const _$IdentityServiceV1UserStatusEnumMap = {
  IdentityServiceV1UserStatus.disabled: 'DISABLED',
  IdentityServiceV1UserStatus.normal: 'NORMAL',
  IdentityServiceV1UserStatus.pending: 'PENDING',
  IdentityServiceV1UserStatus.locked: 'LOCKED',
  IdentityServiceV1UserStatus.expired: 'EXPIRED',
  IdentityServiceV1UserStatus.closed: 'CLOSED',
  IdentityServiceV1UserStatus.$unknown: r'$unknown',
};
