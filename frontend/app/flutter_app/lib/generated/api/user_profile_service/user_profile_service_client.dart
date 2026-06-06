// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';

import '../models/identity_service_v1_bind_contact_request.dart';
import '../models/identity_service_v1_change_password_request.dart';
import '../models/identity_service_v1_update_user_request.dart';
import '../models/identity_service_v1_upload_avatar_request.dart';
import '../models/identity_service_v1_upload_avatar_response.dart';
import '../models/identity_service_v1_user.dart';
import '../models/identity_service_v1_verify_contact_request.dart';

part 'user_profile_service_client.g.dart';

@RestApi()
abstract class UserProfileServiceClient {
  factory UserProfileServiceClient(Dio dio, {String? baseUrl}) = _UserProfileServiceClient;

  /// 获取用户资料
  @GET('/app/v1/me')
  Future<IdentityServiceV1User> userProfileServiceGetUser();

  /// 更新用户资料
  @PUT('/app/v1/me')
  Future<void> userProfileServiceUpdateUser({
    @Body() required IdentityServiceV1UpdateUserRequest body,
  });

  /// 上传头像
  @POST('/app/v1/me/avatar')
  Future<IdentityServiceV1UploadAvatarResponse> userProfileServiceUploadAvatar({
    @Body() required IdentityServiceV1UploadAvatarRequest body,
  });

  /// 删除头像
  @DELETE('/app/v1/me/avatar')
  Future<void> userProfileServiceDeleteAvatar();

  /// 绑定手机号码/邮箱
  @POST('/app/v1/me/contact')
  Future<void> userProfileServiceBindContact({
    @Body() required IdentityServiceV1BindContactRequest body,
  });

  /// 验证手机号码/邮箱
  @POST('/app/v1/me/contact/verify')
  Future<void> userProfileServiceVerifyContact({
    @Body() required IdentityServiceV1VerifyContactRequest body,
  });

  /// 修改用户密码
  @POST('/app/v1/me/password')
  Future<void> userProfileServiceChangePassword({
    @Body() required IdentityServiceV1ChangePasswordRequest body,
  });
}
