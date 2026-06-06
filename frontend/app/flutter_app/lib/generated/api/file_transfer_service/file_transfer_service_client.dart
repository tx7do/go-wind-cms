// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, unused_import, invalid_annotation_target, unnecessary_import

import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';

import '../models/storage_service_v1_download_file_response.dart';
import '../models/storage_service_v1_upload_file_request.dart';
import '../models/storage_service_v1_upload_file_response.dart';

part 'file_transfer_service_client.g.dart';

@RestApi()
abstract class FileTransferServiceClient {
  factory FileTransferServiceClient(Dio dio, {String? baseUrl}) = _FileTransferServiceClient;

  /// 下载文件.
  ///
  /// [rangeStart] - 可选的分段下载范围（闭区间 start，开区间 end；为 0/0 表示全量）.
  @GET('/app/v1/file/download')
  Future<StorageServiceV1DownloadFileResponse> fileTransferServiceDownloadFile({
    @Query('fileId') int? fileId,
    @Query('storageObject.bucketName') String? object0,
    @Query('storageObject.fileDirectory') String? object1,
    @Query('storageObject.objectName') String? object2,
    @Query('downloadUrl') String? downloadUrl,
    @Query('rangeStart') String? rangeStart,
    @Query('rangeEnd') String? rangeEnd,
    @Query('preferPresignedUrl') bool? preferPresignedUrl,
    @Query('presignExpireSeconds') int? presignExpireSeconds,
    @Query('disposition') String? disposition,
    @Query('acceptMime') String? acceptMime,
  });

  /// 上传文件 PUT 方式
  @PUT('/app/v1/file/upload')
  Future<StorageServiceV1UploadFileResponse> fileTransferServicePutUploadFile({
    @Body() required StorageServiceV1UploadFileRequest body,
  });

  /// 上传文件 POST 方式
  @POST('/app/v1/file/upload')
  Future<StorageServiceV1UploadFileResponse> fileTransferServicePostUploadFile({
    @Body() required StorageServiceV1UploadFileRequest body,
  });
}
