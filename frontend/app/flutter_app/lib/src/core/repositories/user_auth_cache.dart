import 'package:get_it/get_it.dart' show GetIt;
import 'package:jose/jose.dart';
import 'package:shared_preferences/shared_preferences.dart'
    show SharedPreferences;

import 'package:flutter_app/src/core/utilities/uuid.dart' show UuidUtils;
import 'package:flutter_app/src/core/utilities/logger.dart' show fatal;

/// 存储键
class _StorageKey {
  /// 访问令牌
  static const String accessToken = 'ACCESS_TOKEN';

  /// 刷新令牌
  static const String refreshToken = 'REFRESH_TOKEN';

  /// 客户端ID
  static const String mqttClientId = 'MQTT_CLIENT_ID';

  /// 设备ID
  static const String deviceId = 'DEVICE_ID';
}

/// 用户认证信息缓存
class UserAuthCache {
  // final _storage = GetIt.instance<FlutterSecureStorage>();
  SharedPreferences get _prefs => GetIt.instance<SharedPreferences>();

  int? _userId;

  String _assessToken = '';
  String _refreshToken = '';

  String _deviceId = '';
  String _mqttClientId = '';

  int get userId {
    return _userId ?? 0;
  }

  /// 获取访问令牌
  String? get accessToken {
    return _assessToken;
  }

  /// 获取刷新令牌
  String? get refreshToken {
    return _refreshToken;
  }

  /// 是否已经登录
  bool get hasLogin {
    final at = accessToken;
    return at != '' && at != null;
  }

  /// 获取设备ID
  String get deviceId {
    return _deviceId;
  }

  /// 获取MQTT客户端ID
  String get mqttClientId {
    return _mqttClientId;
  }

  init() async {
    // debug('UserAuthCache init...');
    await _initDeviceId();
    await _initToken();
  }

  /// 存储认证信息
  saveAuthInfo(String accessToken, {String? refreshToken}) async {
    await saveAccessToken(accessToken);
    await saveRefreshToken(refreshToken);
  }

  /// 存储访问令牌
  Future<void> saveAccessToken(String accessToken) async {
    _assessToken = accessToken;
    _refreshUserId(accessToken);
    await _prefs.setString(_StorageKey.accessToken, _assessToken);
  }

  /// 存储刷新令牌
  Future<void> saveRefreshToken(String? refreshToken) async {
    _refreshToken = refreshToken ?? '';
    await _prefs.setString(_StorageKey.refreshToken, _refreshToken);
  }

  /// 清除访问令牌
  Future<void> clearAccessToken() async {
    _userId = null;
    await _prefs.remove(_StorageKey.accessToken);
  }

  /// 清除刷新令牌
  Future<void> clearRefreshToken() async {
    await _prefs.remove(_StorageKey.refreshToken);
  }

  /// 清除所有令牌
  Future<void> clearTokens() async {
    await clearAccessToken();
    await clearRefreshToken();
  }

  /// 获取访问令牌
  Future<String?> _readAccessToken() async {
    return _prefs.getString(_StorageKey.accessToken);
  }

  /// 获取刷新令牌
  Future<String?> _readRefreshToken() async {
    return _prefs.getString(_StorageKey.refreshToken);
  }

  /// 获取设备ID
  Future<String?> _readDeviceId() async {
    return _prefs.getString(_StorageKey.deviceId);
  }

  /// 存储设备ID
  Future<void> _saveDeviceId(String? id) async {
    await _prefs.setString(_StorageKey.deviceId, id ?? '');
  }

  /// 生成设备ID
  String _generateDeviceId() {
    final xidDeviceId = UuidUtils.xid;
    return xidDeviceId;
  }

  /// 从JWT令牌解析出UserID
  int? _getUserIdFromJwt(String jwtToken) {
    try {
      // 解析JWT
      final jwt = JsonWebToken.unverified(jwtToken);

      // debug('parse jwt payload: ${jwt.claims.toJson()}');

      int? userId = jwt.claims.getTyped('uid');
      if (userId == null) {
        fatal('UserID not found in JWT payload: $jwtToken');
        return null;
      }

      return userId;
    } catch (e) {
      fatal('parse jwt payload failed: $e');
      return null;
    }
  }

  /// 刷新UserID
  int _refreshUserId(String accessToken) {
    final userId = _getUserIdFromJwt(accessToken);
    if (userId == null) {
      return 0;
    }

    _userId = userId;

    return userId;
  }

  Future<String?> _refreshAccessToken() async {
    final accessToken = await _readAccessToken();
    if (accessToken != null && accessToken.isNotEmpty) {
      _assessToken = accessToken;
      _refreshUserId(accessToken);
    }
    return accessToken;
  }

  Future<String?> _refreshRefreshToken() async {
    final refreshToken = await _readRefreshToken();
    if (refreshToken != null && refreshToken.isNotEmpty) {
      _refreshToken = refreshToken;
    }
    return refreshToken;
  }

  String _refreshDeviceId() {
    final deviceId = _generateDeviceId();
    _saveDeviceId(deviceId);
    return deviceId;
  }

  Future<String> _initDeviceId() async {
    String deviceId = await _readDeviceId() ?? '';
    if (deviceId.isEmpty) {
      deviceId = _refreshDeviceId();
    }
    _deviceId = deviceId;
    return _deviceId;
  }

  _initToken() async {
    await _refreshAccessToken();
    await _refreshRefreshToken();
  }
}
