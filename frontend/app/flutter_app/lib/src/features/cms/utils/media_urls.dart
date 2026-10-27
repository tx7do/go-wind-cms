import 'package:flutter_app/src/core/config/environments.dart';

/// 媒体引用工具（即时图片变换接入）：
/// customFields 的 IMAGE 值约定为字符串——`media_asset:<id>` 引用或 http 直链。
/// 变换端点只缩小不放大；gif/svg 等不可变换源返回 415、资产缺失 404，由调用方兜底。

/// `media_asset:<id>` 引用值形态（与 admin 写侧约定对齐）
final RegExp mediaRefPattern = RegExp(r'^media_asset:(\d+)$');

/// 解析 `media_asset:<id>` 引用，非引用返回 null
int? parseMediaRef(String value) {
  final m = mediaRefPattern.firstMatch(value.trim());
  return m != null ? int.parse(m.group(1)!) : null;
}

/// 即时图片变换端点 URL（免鉴权；w=1280 覆盖详情列宽 2x 余量、q=80）
String mediaTransformUrl(int id) =>
    '${Environments.apiBaseUrl}/app/v1/media-assets/$id/content?w=1280&q=80';

/// IMAGE 值是否可按图片渲染（http 直链或 media_asset 引用）
bool isMediaImageValue(String value) =>
    parseMediaRef(value) != null || value.startsWith('http');
