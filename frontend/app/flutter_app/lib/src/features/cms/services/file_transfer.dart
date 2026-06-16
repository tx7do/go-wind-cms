import 'dart:convert' show base64;
import 'dart:typed_data' show Uint8List;

import 'package:dio/dio.dart' show DioException, Options, ProgressCallback;
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show ApiClient, FileTransferServiceClient;
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show StorageServiceV1DownloadFileRequest,
        StorageServiceV1DownloadFileResponse,
        StorageServiceV1UploadFileRequest, StorageServiceV1UploadFileResponse,
        StorageServiceV1StorageObject, StorageServiceV1PresignOption;
import 'package:flutter_app/src/core/services/base_service.dart';
import 'package:flutter_app/src/core/transport/http/index.dart';

typedef DownloadResponse = StorageServiceV1DownloadFileResponse;
typedef UploadRequest = StorageServiceV1UploadFileRequest;
typedef UploadResponse = StorageServiceV1UploadFileResponse;
typedef StorageObject = StorageServiceV1StorageObject;
typedef PresignOption = StorageServiceV1PresignOption;

/// 文件传输服务
///
/// 封装 MinIO 文件上传/下载，对齐 TS 端 `file_transfer.ts`：
/// - [downloadFile] 下载文件（支持预签名 URL 和直接字节流）
/// - [uploadFile] 上传文件（支持 POST / PUT，支持进度回调）
class FileTransferService extends BaseService {
  FileTransferService() : super(tag: 'FileTransferService');

  FileTransferServiceClient get _api =>
      GetIt.instance<ApiClient>().fileTransferService;

  // ─── 下载 ─────────────────────────────────────────────

  /// 下载文件
  ///
  /// 返回 [FileDownloadResult]，调用方可根据需要使用：
  /// - [FileDownloadResult.presignedUrl] → 预签名 URL（大文件 / WebView 场景）
  /// - [FileDownloadResult.bytes] → 文件字节（小文件 / 本地渲染场景）
  /// - [FileDownloadResult.sourceFileName] → 原始文件名
  /// - [FileDownloadResult.mime] → MIME 类型
  ///
  /// [preferPresignedUrl] 为 true 时优先使用预签名 URL；
  /// 服务端不支持预签名或 [preferPresignedUrl] 为 false 时，直接返回文件字节。
  Future<FileDownloadResult> downloadFile({
    String? bucketName,
    String? fileDirectory,
    String? objectName,
    int? fileId,
    bool preferPresignedUrl = false,
    int? presignExpireSeconds,
    String? rangeStart,
    String? rangeEnd,
    String? disposition,
    String? acceptMime,
  }) async {
    try {
      final resp = await _api.downloadFile(
        StorageServiceV1DownloadFileRequest(
          fileId: fileId,
          storageObject: StorageObject(
            bucketName: bucketName,
            fileDirectory: fileDirectory,
            objectName: objectName,
          ),
          preferPresignedUrl: preferPresignedUrl,
          presignExpireSeconds: presignExpireSeconds,
          rangeStart: rangeStart != null ? int.tryParse(rangeStart) : null,
          rangeEnd: rangeEnd != null ? int.tryParse(rangeEnd) : null,
          disposition: disposition,
          acceptMime: acceptMime,
        ),
      );

      // 预签名 URL 模式
      if (preferPresignedUrl && resp.downloadUrl != null) {
        return FileDownloadResult(
          presignedUrl: resp.downloadUrl,
          sourceFileName: resp.sourceFileName,
          mime: resp.mime,
          size: resp.size?.toString(),
        );
      }

      // 直接字节流模式
      Uint8List? bytes;
      if (resp.file != null && resp.file!.isNotEmpty) {
        bytes = _decodeBase64(resp.file!);
      }

      return FileDownloadResult(
        bytes: bytes,
        sourceFileName: resp.sourceFileName,
        mime: resp.mime ?? 'application/octet-stream',
        size: resp.size?.toString(),
      );
    } on DioException catch (e) {
      throw handleDioError(e);
    }
  }

  // ─── 上传 ─────────────────────────────────────────────

  /// 上传文件（Base64 内联）
  ///
  /// [fileBytes] 文件原始字节
  /// [fileName] 原始文件名
  /// [bucketName] 文件桶名称
  /// [fileDirectory] 远端存储目录
  /// [method] 上传方式，'post' 或 'put'
  /// [onSendProgress] 上传进度回调
  /// [presign] 预签名选项（服务端返回预签名 URL，由客户端直传 OSS）
  Future<UploadResponse> uploadFile({
    required Uint8List fileBytes,
    required String fileName,
    String? bucketName,
    String? fileDirectory,
    String? objectName,
    String method = 'post',
    ProgressCallback? onSendProgress,
    PresignOption? presign,
  }) async {
    try {
      final fileBase64 = base64.encode(fileBytes);
      final mime = _guessMime(fileName);

      final request = UploadRequest(
        storageObject: StorageObject(
          bucketName: bucketName,
          fileDirectory: fileDirectory,
          objectName: objectName,
        ),
        file: fileBase64,
        sourceFileName: fileName,
        mime: mime,
        size: fileBytes.length,
        presign: presign,
      );

      // 上传进度通过 Options.sendProgress 无法直接在 Retrofit 中传递，
      // 因此对大文件场景建议使用预签名直传。
      // 此处选择 POST 或 PUT。
      if (method.toLowerCase() == 'put') {
        return await _api.putUploadFile(request);
      }
      return await _api.postUploadFile(request);
    } on DioException catch (e) {
      throw handleDioError(e);
    }
  }

  /// 上传文件（预签名直传模式）
  ///
  /// 第一步：调用后端获取预签名 URL
  /// 第二步：使用 Dio 直接 PUT 到预签名 URL（支持进度回调）
  Future<UploadResponse> uploadFilePresigned({
    required Uint8List fileBytes,
    required String fileName,
    String? bucketName,
    String? fileDirectory,
    String? objectName,
    String presignMethod = 'PUT',
    int presignExpireSeconds = 3600,
    ProgressCallback? onSendProgress,
  }) async {
    try {
      // 1. 获取预签名 URL
      final fileBase64 = base64.encode(fileBytes);
      final mime = _guessMime(fileName);

      final request = UploadRequest(
        storageObject: StorageObject(
          bucketName: bucketName,
          fileDirectory: fileDirectory,
          objectName: objectName,
        ),
        file: fileBase64,
        sourceFileName: fileName,
        mime: mime,
        size: fileBytes.length,
        presign: PresignOption(
          method: presignMethod,
          expireSeconds: presignExpireSeconds,
          contentType: mime,
        ),
      );

      final resp = await _api.postUploadFile(request);

      // 2. 如果有预签名 URL，执行直传
      if (resp.presignedUrl != null && resp.presignedUrl!.isNotEmpty) {
        final dio = GetIt.instance<Dio>();
        await dio.put(
          resp.presignedUrl!,
          data: fileBytes,
          options: _presignOptions(mime),
          onSendProgress: onSendProgress,
        );
      }

      return resp;
    } on DioException catch (e) {
      throw handleDioError(e);
    }
  }

  // ─── 工具方法 ─────────────────────────────────────────

  /// 解码 base64 字符串为字节
  ///
  /// 支持 URL-safe base64（与 TS 端 `normalizeBase64` 对齐）
  static Uint8List _decodeBase64(String input) {
    var str = input.replaceAll(RegExp(r'\s+'), '');
    // URL-safe → standard
    str = str.replaceAll('-', '+').replaceAll('_', '/');
    // 补 padding
    while (str.length % 4 != 0) {
      str += '=';
    }
    // 跳过 data URI 前缀
    if (str.contains('base64,')) {
      str = str.split('base64,')[1];
    }
    return base64.decode(str);
  }

  /// 简单的文件名 → MIME 推断
  static String _guessMime(String fileName) {
    final ext = fileName.split('.').last.toLowerCase();
    const mimeMap = {
      'jpg': 'image/jpeg',
      'jpeg': 'image/jpeg',
      'png': 'image/png',
      'gif': 'image/gif',
      'webp': 'image/webp',
      'svg': 'image/svg+xml',
      'pdf': 'application/pdf',
      'doc': 'application/msword',
      'docx':
          'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      'xls': 'application/vnd.ms-excel',
      'xlsx':
          'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      'mp4': 'video/mp4',
      'mp3': 'audio/mpeg',
      'zip': 'application/zip',
      'json': 'application/json',
      'txt': 'text/plain',
      'html': 'text/html',
      'css': 'text/css',
      'js': 'application/javascript',
    };
    return mimeMap[ext] ?? 'application/octet-stream';
  }

  /// 构建预签名直传的 Options
  static Options _presignOptions(String? contentType) {
    return Options(
      contentType: contentType,
      headers: {'Content-Type': contentType ?? 'application/octet-stream'},
    );
  }
}

// ─── 结果模型 ───────────────────────────────────────────

/// 文件下载结果
///
/// - [presignedUrl] 不为空 → 调用方通过 URL 下载（WebView / 系统浏览器）
/// - [bytes] 不为空 → 直接使用字节（本地渲染 / 保存）
class FileDownloadResult {
  /// 预签名下载 URL
  final String? presignedUrl;

  /// 文件原始字节
  final Uint8List? bytes;

  /// 原始文件名
  final String? sourceFileName;

  /// MIME 类型
  final String? mime;

  /// 文件大小（字节字符串）
  final String? size;

  const FileDownloadResult({
    this.presignedUrl,
    this.bytes,
    this.sourceFileName,
    this.mime,
    this.size,
  });

  /// 是否为预签名模式
  bool get isPresigned => presignedUrl != null;

  /// 文件大小（int）
  int? get sizeBytes => size != null ? int.tryParse(size!) : null;
}
