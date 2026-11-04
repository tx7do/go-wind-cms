import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/generated/l10n.dart' as l10n;
import 'package:flutter_app/src/core/preference/preference.dart';

part 'app_theme_state.dart';

/// 外观偏好：主题模式（跟随系统 / 晴风 / 夜风）+ 语言。
///
/// 主题色不再可配——品牌蓝是 docs/design-language.md §1 的铁律，
/// 八个预设色板会把「唯一行动色」降级成用户口味。
class AppThemeCubit extends Cubit<AppThemeState> {
  UserPreferenceCache get _cache => GetIt.instance<UserPreferenceCache>();

  ThemeMode? _originalThemeMode;

  AppThemeCubit() : super(const AppThemeInitial());

  ThemeMode get themeMode => _cache.themeMode;

  bool get isDarkMode => _cache.isDarkMode;

  bool get isLightMode => _cache.isLightMode;

  /// 获取当前值
  ThemeMode get currentValue {
    if (state is AppThemeInitial) {
      return (state as AppThemeInitial).themeMode;
    } else if (state is AppThemeModified) {
      return (state as AppThemeModified).themeMode;
    }
    return ThemeMode.system;
  }

  /// 初始化值
  ThemeMode get initialValue {
    if (state is AppThemeInitial) {
      return (state as AppThemeInitial).themeMode;
    }
    return ThemeMode.system;
  }

  /// 是否已被修改
  bool get hasModify => state is AppThemeModified;

  /// 初始化
  void init() async {
    final currentMode = _cache.themeMode;
    final locale = _localeFromString(_cache.language);
    emit(AppThemeInitial(themeMode: currentMode, locale: locale));
  }

  /// 设置初始化值
  void setInitValue(ThemeMode newMode) {
    _originalThemeMode ??= newMode;
    emit(AppThemeInitial(themeMode: newMode, locale: currentLocale));
  }

  /// 重置
  void resetTheme() {
    if (_originalThemeMode == null) {
      return;
    }

    final current = currentValue;
    if (current != _originalThemeMode) {
      _cache.setMaterialThemeMode(_originalThemeMode!);
    }
  }

  /// 修改主题模式
  Future<void> modify(ThemeMode newMode) async {
    final current = currentValue;

    if (current == newMode) {
      return;
    }

    await _cache.setMaterialThemeMode(newMode);
    _setEasyLoadingTheme(newMode);

    emit(AppThemeModified(themeMode: newMode, locale: currentLocale));
  }

  /// 强制修改
  void forceModify(ThemeMode newMode) {
    _cache.setMaterialThemeMode(newMode);
    emit(AppThemeModified(themeMode: newMode, locale: currentLocale));
  }

  void _setEasyLoadingTheme(ThemeMode mode) {
    EasyLoadingStyle easyLoadingStyle = EasyLoadingStyle.dark;
    if (mode == ThemeMode.dark) {
      easyLoadingStyle = EasyLoadingStyle.light;
    } else if (mode == ThemeMode.system) {
      if (!isDarkMode) {
        easyLoadingStyle = EasyLoadingStyle.light;
      }
    }
    EasyLoading.instance.loadingStyle = easyLoadingStyle;
  }

  // =================== 语言管理 ===================

  /// 获取当前 locale
  Locale get currentLocale {
    if (state is AppThemeInitial) {
      return (state as AppThemeInitial).locale ?? _localeFromString(_cache.language);
    } else if (state is AppThemeModified) {
      return (state as AppThemeModified).locale ?? _localeFromString(_cache.language);
    }
    return _localeFromString(_cache.language);
  }

  /// 修改语言
  void modifyLocale(Locale newLocale) async {
    if (newLocale == currentLocale) return;

    final country = newLocale.countryCode;
    final langString = country != null && country.isNotEmpty
        ? '${newLocale.languageCode}_$country'
        : newLocale.languageCode;
    await _cache.setLanguage(langString, true);

    emit(AppThemeModified(
      themeMode: currentValue,
      locale: newLocale,
    ));
  }

  /// 支持的语言列表
  List<Locale> get supportedLocales => l10n.S.delegate.supportedLocales;

  /// 将字符串转为 Locale（zh_CN / zh-CN → Locale('zh', 'CN')）
  static Locale _localeFromString(String code) {
    // 标准化：zh-CN → zh_CN
    final normalized = code.replaceAll('-', '_');

    // 尝试匹配 supportedLocales
    final supported = l10n.S.delegate.supportedLocales;
    for (final locale in supported) {
      final str = locale.countryCode != null
          ? '${locale.languageCode}_${locale.countryCode}'
          : locale.languageCode;
      if (str == normalized) return locale;
      // 也匹配不带国家的
      if (locale.countryCode == null && locale.languageCode == normalized) return locale;
    }
    // fallback: 拆分
    final parts = normalized.split('_');
    if (parts.length == 2) {
      return Locale(parts[0], parts[1]);
    }
    return Locale(parts[0]);
  }
}
