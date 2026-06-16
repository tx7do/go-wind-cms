import 'package:get_it/get_it.dart' show GetIt;
import 'package:cached_query/cached_query.dart' show Mutation, Query;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ApiClient, UserProfileServiceClient;
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show IdentityServiceV1User, IdentityServiceV1UpdateUserRequest,
        IdentityServiceV1UploadAvatarRequest, IdentityServiceV1UploadAvatarResponse,
        IdentityServiceV1BindContactRequest, IdentityServiceV1VerifyContactRequest,
        IdentityServiceV1ChangePasswordRequest;
import 'package:flutter_app/src/core/services/base_service.dart';
import 'package:flutter_app/src/core/transport/http/index.dart';

typedef User = IdentityServiceV1User;
typedef UpdateUserRequest = IdentityServiceV1UpdateUserRequest;
typedef UploadAvatarRequest = IdentityServiceV1UploadAvatarRequest;
typedef UploadAvatarResponse = IdentityServiceV1UploadAvatarResponse;
typedef BindContactRequest = IdentityServiceV1BindContactRequest;
typedef VerifyContactRequest = IdentityServiceV1VerifyContactRequest;
typedef ChangePasswordRequest = IdentityServiceV1ChangePasswordRequest;

/// 用户资料服务
///
/// 使用 [Query] 缓存用户资料查询，
/// 使用 [Mutation] 管理更新/头像/联系方式/密码等写操作。
class UserProfileService extends BaseService {
  UserProfileService() : super(tag: 'UserProfileService');

  UserProfileServiceClient get _api =>
      GetIt.instance<ApiClient>().userProfileService;

  // ─── Queries ──────────────────────────────────────────

  /// 获取当前用户资料 Query
  Query<User> getQuery() {
    return Query<User>(
      key: 'user-profile',
      queryFn: () => _api.getUser({}),
    );
  }

  // ─── Mutations ────────────────────────────────────────

  /// 更新用户资料 Mutation
  Mutation<void, UpdateUserParams> updateMutation() {
    return Mutation<void, UpdateUserParams>(
      mutationFn: (params) => _api.updateUser(
        UpdateUserRequest(
          id: params.id,
          data: params.data,
          password: params.password,
          updateMask: params.updateMask,
          allowMissing: params.allowMissing,
        ),
      ),
      invalidateQueries: ['user-profile'],
    );
  }

  /// 上传头像 Mutation
  Mutation<UploadAvatarResponse, UploadAvatarParams> uploadAvatarMutation() {
    return Mutation<UploadAvatarResponse, UploadAvatarParams>(
      mutationFn: (params) => _api.uploadAvatar(
        UploadAvatarRequest(
          imageBase64: params.imageBase64,
          imageUrl: params.imageUrl,
        ),
      ),
      invalidateQueries: ['user-profile'],
    );
  }

  /// 删除头像 Mutation
  Mutation<void, void> deleteAvatarMutation() {
    return Mutation<void, void>(
      mutationFn: (_) => _api.deleteAvatar({}),
      invalidateQueries: ['user-profile'],
    );
  }

  /// 绑定联系方式 Mutation
  Mutation<void, BindContactRequest> bindContactMutation() {
    return Mutation<void, BindContactRequest>(
      mutationFn: (request) =>
          _api.bindContact(request),
    );
  }

  /// 验证联系方式 Mutation
  Mutation<void, VerifyContactRequest> verifyContactMutation() {
    return Mutation<void, VerifyContactRequest>(
      mutationFn: (request) =>
          _api.verifyContact(request),
    );
  }

  /// 修改密码 Mutation
  Mutation<void, ChangePasswordRequest> changePasswordMutation() {
    return Mutation<void, ChangePasswordRequest>(
      mutationFn: (request) =>
          _api.changePassword(request),
    );
  }

  // ─── 直接调用方法（非 Mutation，适合简单场景） ─────────

  /// 获取当前用户资料
  Future<dynamic> get() async {
    try {
      return await _api.getUser({});
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 更新用户资料
  Future<dynamic> update(
    User data, {
    int? id,
    String? password,
    String? updateMask,
    bool? allowMissing,
  }) async {
    try {
      return await _api.updateUser(
        UpdateUserRequest(
          id: id,
          data: data,
          password: password,
          updateMask: updateMask,
          allowMissing: allowMissing,
        ),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 上传头像（Base64 或 URL）
  Future<dynamic> uploadAvatar({String? imageBase64, String? imageUrl}) async {
    try {
      return await _api.uploadAvatar(
        UploadAvatarRequest(imageBase64: imageBase64, imageUrl: imageUrl),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 删除头像
  Future<dynamic> deleteAvatar() async {
    try {
      await _api.deleteAvatar({});
      return null;
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 绑定手机号码/邮箱
  Future<dynamic> bindContact(BindContactRequest request) async {
    try {
      return await _api.bindContact(request);
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 验证手机号码/邮箱
  Future<dynamic> verifyContact(VerifyContactRequest request) async {
    try {
      return await _api.verifyContact(request);
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }

  /// 修改用户密码
  Future<dynamic> changePassword({
    required String oldPassword,
    required String newPassword,
  }) async {
    try {
      return await _api.changePassword(
        ChangePasswordRequest(
          oldPassword: oldPassword,
          newPassword: newPassword,
        ),
      );
    } on DioException catch (e) {
      return handleDioError(e);
    }
  }
}

/// 更新用户参数
class UpdateUserParams {
  final int? id;
  final User data;
  final String? password;
  final String? updateMask;
  final bool? allowMissing;

  const UpdateUserParams({
    this.id,
    required this.data,
    this.password,
    this.updateMask,
    this.allowMissing,
  });
}

/// 上传头像参数
class UploadAvatarParams {
  final String? imageBase64;
  final String? imageUrl;

  const UploadAvatarParams({this.imageBase64, this.imageUrl});
}
