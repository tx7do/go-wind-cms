// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:dio/dio.dart';

import 'category_service/category_service_client.dart';
import 'comment_service/comment_service_client.dart';
import 'file_transfer_service/file_transfer_service_client.dart';
import 'authentication_service/authentication_service_client.dart';
import 'user_profile_service/user_profile_service_client.dart';
import 'navigation_service/navigation_service_client.dart';
import 'page_service/page_service_client.dart';
import 'post_service/post_service_client.dart';
import 'tag_service/tag_service_client.dart';

/// GoWind Content Hub App API `v1.0`.
///
/// GoWind Content Hub App API.
class RestClient {
  RestClient(
    Dio dio, {
    String? baseUrl,
  })  : _dio = dio,
        _baseUrl = baseUrl;

  final Dio _dio;
  final String? _baseUrl;

  static String get version => '1.0';

  CategoryServiceClient? _categoryService;
  CommentServiceClient? _commentService;
  FileTransferServiceClient? _fileTransferService;
  AuthenticationServiceClient? _authenticationService;
  UserProfileServiceClient? _userProfileService;
  NavigationServiceClient? _navigationService;
  PageServiceClient? _pageService;
  PostServiceClient? _postService;
  TagServiceClient? _tagService;

  CategoryServiceClient get categoryService => _categoryService ??= CategoryServiceClient(_dio, baseUrl: _baseUrl);

  CommentServiceClient get commentService => _commentService ??= CommentServiceClient(_dio, baseUrl: _baseUrl);

  FileTransferServiceClient get fileTransferService => _fileTransferService ??= FileTransferServiceClient(_dio, baseUrl: _baseUrl);

  AuthenticationServiceClient get authenticationService => _authenticationService ??= AuthenticationServiceClient(_dio, baseUrl: _baseUrl);

  UserProfileServiceClient get userProfileService => _userProfileService ??= UserProfileServiceClient(_dio, baseUrl: _baseUrl);

  NavigationServiceClient get navigationService => _navigationService ??= NavigationServiceClient(_dio, baseUrl: _baseUrl);

  PageServiceClient get pageService => _pageService ??= PageServiceClient(_dio, baseUrl: _baseUrl);

  PostServiceClient get postService => _postService ??= PostServiceClient(_dio, baseUrl: _baseUrl);

  TagServiceClient get tagService => _tagService ??= TagServiceClient(_dio, baseUrl: _baseUrl);
}
