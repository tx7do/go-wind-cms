// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:json_annotation/json_annotation.dart';

import 'identity_service_v1_user_gender.dart';
import 'identity_service_v1_user_status.dart';

part 'identity_service_v1_user.g.dart';

/// 用户
@JsonSerializable()
class IdentityServiceV1User {
  const IdentityServiceV1User({
    this.roles,
    this.tenantId,
    this.tenantName,
    this.orgUnitId,
    this.orgUnitIds,
    this.orgUnitName,
    this.orgUnitNames,
    this.positionId,
    this.positionIds,
    this.positionName,
    this.positionNames,
    this.roleId,
    this.roleIds,
    this.id,
    this.roleNames,
    this.username,
    this.nickname,
    this.realname,
    this.avatar,
    this.email,
    this.mobile,
    this.telephone,
    this.gender,
    this.address,
    this.region,
    this.description,
    this.remark,
    this.deletedAt,
    this.lastLoginIp,
    this.status,
    this.lockedUntil,
    this.followers,
    this.following,
    this.postCount,
    this.commentCount,
    this.likeCount,
    this.createdBy,
    this.updatedBy,
    this.deletedBy,
    this.createdAt,
    this.updatedAt,
    this.lastLoginAt,
  });
  
  factory IdentityServiceV1User.fromJson(Map<String, Object?> json) => _$IdentityServiceV1UserFromJson(json);
  
  /// 用户ID
  final int? id;

  /// 租户ID
  final int? tenantId;

  /// 租户名称
  final String? tenantName;

  /// 组织ID
  final int? orgUnitId;

  /// 归属组织ID列表
  final List<int>? orgUnitIds;

  /// 组织名称
  final String? orgUnitName;

  /// 组织名称列表
  final List<String>? orgUnitNames;

  /// 职位ID
  final int? positionId;

  /// 职位列表
  final List<int>? positionIds;

  /// 职位名称
  final String? positionName;

  /// 职位名称列表
  final List<String>? positionNames;

  /// 角色ID
  final int? roleId;

  /// 角色ID列表
  final List<int>? roleIds;

  /// 角色码列表
  final List<String>? roles;

  /// 角色名称列表
  final List<String>? roleNames;

  /// 用户名
  final String? username;

  /// 昵称
  final String? nickname;

  /// 真实姓名
  final String? realname;

  /// 头像
  final String? avatar;

  /// 邮箱
  final String? email;

  /// 手机号
  final String? mobile;

  /// 座机号
  final String? telephone;

  /// 性别
  final IdentityServiceV1UserGender? gender;

  /// 住址
  final String? address;

  /// 国家地区
  final String? region;

  /// 个人描述
  final String? description;

  /// 备注
  final String? remark;

  /// 最后登录时间
  final DateTime? lastLoginAt;

  /// 最后登录IP
  final String? lastLoginIp;

  /// 状态
  final IdentityServiceV1UserStatus? status;

  /// 锁定截止时间
  final DateTime? lockedUntil;

  /// 粉丝数
  final String? followers;

  /// 关注数
  final String? following;

  /// 发帖数
  final String? postCount;

  /// 评论数
  final String? commentCount;

  /// 获赞数
  final String? likeCount;

  /// 创建者用户ID
  final int? createdBy;

  /// 更新者用户ID
  final int? updatedBy;

  /// 删除者用户ID
  final int? deletedBy;

  /// 创建时间
  final DateTime? createdAt;

  /// 更新时间
  final DateTime? updatedAt;

  /// 删除时间
  final DateTime? deletedAt;

  Map<String, Object?> toJson() => _$IdentityServiceV1UserToJson(this);
}
