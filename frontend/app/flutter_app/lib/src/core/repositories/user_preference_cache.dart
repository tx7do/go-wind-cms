import 'dart:convert' show json;

import 'package:flutter/material.dart' show ThemeMode, Color;
import 'package:flutter_app/src/core/models/user_preference.dart'
    show UserPreference;
import 'package:get_it/get_it.dart' show GetIt;
import 'package:shared_preferences/shared_preferences.dart';

class ThemeModeMapper {
  static String convertToString(ThemeMode mode) {
    switch (mode) {
      case ThemeMode.light:
        return 'light';
      case ThemeMode.dark:
        return 'dark';
      case ThemeMode.system:
        return 'system';
    }
  }

  static int convertToInt(ThemeMode mode) {
    switch (mode) {
      case ThemeMode.light:
        return 0;
      case ThemeMode.dark:
        return 1;
      case ThemeMode.system:
        return 2;
    }
  }

  static ThemeMode intToThemeMode(int mode) {
    switch (mode) {
      case 0:
        return ThemeMode.light;
      case 1:
        return ThemeMode.dark;
      case 2:
        return ThemeMode.system;

      default:
        return ThemeMode.system; // 默认返回系统模式
    }
  }

  static ThemeMode convertToThemeMode(String mode) {
    switch (mode) {
      case 'light':
        return ThemeMode.light;
      case 'dark':
        return ThemeMode.dark;
      case 'system':
        return ThemeMode.system;
      default:
        return ThemeMode.system;
    }
  }
}

/// 存储键
class _StorageKey {
  static const String userPreference = 'USER_PREFERENCE';
  static const String seedColor = 'SEED_COLOR';
}

/// 用户偏好设置缓存
class UserPreferenceCache {
  UserPreference? _userPreference;

  SharedPreferences get _prefs => GetIt.instance<SharedPreferences>();

  init() async {
    // 初始化时可以加载用户偏好设置
    _userPreference = await loadUserPreference();
  }

  /// 主题风格
  Future<void> setMaterialThemeMode(ThemeMode mode) async {
    final up = await loadUserPreference();
    up.themeMode = ThemeModeMapper.convertToString(mode);
    setUserPreference(up);
    _saveUserPreference(up);
  }

  ThemeMode get themeMode {
    // final up = await getUserPreference();
    if (_userPreference == null) {
      return ThemeMode.system;
    }

    return ThemeModeMapper.convertToThemeMode(_userPreference!.themeMode);
  }

  bool get isDarkMode => themeMode == ThemeMode.dark;

  bool get isLightMode => themeMode == ThemeMode.light;

  /// 主题色（seedColor）
  Color? get seedColorValue {
    final value = _prefs.getInt(_StorageKey.seedColor);
    return value != null ? Color(value) : null;
  }

  Future<void> setSeedColor(Color color) async {
    await _prefs.setInt(_StorageKey.seedColor, color.toARGB32());
  }

  String get language {
    if (_userPreference == null) {
      return 'zh_CN';
    }

    return _userPreference!.language;
  }

  /// 语言
  Future<void> setLanguage(String lang, bool save) async {
    final up = await loadUserPreference();
    up.language = lang;
    // setUserPreference(up);
    if (save) {
      _saveUserPreference(up);
    }
  }

  setUserPreference(UserPreference? userPreference) {
    _userPreference = userPreference;
  }

  /// 保存数据
  _saveUserPreference(UserPreference? userPreference) async {
    if (userPreference == null) {
      return;
    }

    await _prefs.setString(
      _StorageKey.userPreference,
      json.encode(userPreference.toProto3Json()),
    );
  }

  /// 加载数据
  Future<UserPreference> loadUserPreference() async {
    final strJson = _prefs.getString(_StorageKey.userPreference) ?? "";
    if (strJson == '') {
      _userPreference = _makeDefaultUserPreference();
      return _userPreference!;
    } else {
      final info = UserPreference.fromJson(strJson);
      _userPreference = info;
      return info;
    }
  }

  /// 创建默认数据
  _makeDefaultUserPreference() {
    return UserPreference(
      language: 'zh_CN',
      themeMode: 'light',
      fontSize: 12,
      allowNotifications: true,
      soundEnabled: true,
      vibrationEnabled: true,
      badgeEnabled: true,
      anonymousMode: true,
      shareAnalytics: true,
      dataCollection: true,
      locationSharing: true,
      autoPlayMedia: true,
      showPreviews: true,
      useBiometric: true,
      hapticFeedback: true,
      rememberLogin: true,
    );
  }
}
