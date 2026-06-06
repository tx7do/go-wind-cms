import 'package:get_it/get_it.dart' show GetIt;
import 'package:shared_preferences/shared_preferences.dart';

/// 存储键
class _StorageKey {
  /// 最后一个登录的账号名
  static const String lastLoginAccount = 'LAST_LOGIN_ACCOUNT';

  /// 最后一个登录的密码
  static const String lastLoginPassword = 'LAST_LOGIN_PASSWORD';
}

/// 本地存储
class UserLoginCache {
  /// 本地存储
  SharedPreferences get _prefs => GetIt.instance<SharedPreferences>();

  /// 最后一个登录的账号名
  String get lastLoginAccount {
    return _prefs.getString(_StorageKey.lastLoginAccount) ?? "";
  }

  Future<void> setLastLoginAccount(String username) async {
    await _prefs.setString(_StorageKey.lastLoginAccount, username);
  }

  /// 最后一个登录的账号名
  String get lastLoginPassword {
    return _prefs.getString(_StorageKey.lastLoginPassword) ?? "";
  }

  Future<void> setLastLoginPassword(String password) async {
    await _prefs.setString(_StorageKey.lastLoginPassword, password);
  }
}
