import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/generated/l10n.dart' as l10n;
import 'package:flutter_app/src/core/preference/preference.dart';
import 'package:flutter_app/src/core/themes/light_theme.dart' show kDefaultSeedColor;

part 'app_theme_state.dart';

class AppThemeCubit extends Cubit<AppThemeState> {
  UserPreferenceCache get _cache => GetIt.instance<UserPreferenceCache>();

  ThemeMode? _originalThemeMode;

  AppThemeCubit() : super(const AppThemeInitial());

  ThemeMode get themeMode => _cache.themeMode;

  /// 当前主题色
  Color get seedColor => _cache.seedColorValue ?? kDefaultSeedColor;

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

  /// 获取当前 seedColor
  Color get currentSeedColor {
    if (state is AppThemeInitial) {
      return (state as AppThemeInitial).seedColor ?? seedColor;
    } else if (state is AppThemeModified) {
      return (state as AppThemeModified).seedColor ?? seedColor;
    }
    return kDefaultSeedColor;
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
    final color = _cache.seedColorValue ?? kDefaultSeedColor;
    final locale = _localeFromString(_cache.language);
    emit(AppThemeInitial(themeMode: currentMode, seedColor: color, locale: locale));
  }

  /// 设置初始化值
  void setInitValue(ThemeMode newMode) {
    _originalThemeMode ??= newMode;
    emit(AppThemeInitial(themeMode: newMode, seedColor: seedColor, locale: currentLocale));
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
  modify(ThemeMode newMode) async {
    final current = currentValue;

    if (current == newMode) {
      return;
    }

    await _cache.setMaterialThemeMode(newMode);
    _setEasyLoadingTheme(newMode);

    emit(AppThemeModified(themeMode: newMode, seedColor: currentSeedColor, locale: currentLocale));
  }

  /// 修改主题色
  void modifySeedColor(Color color) async {
    if (color == currentSeedColor) return;

    await _cache.setSeedColor(color);

    emit(AppThemeModified(themeMode: currentValue, seedColor: color, locale: currentLocale));
  }

  /// 强制修改
  void forceModify(ThemeMode newMode) {
    _cache.setMaterialThemeMode(newMode);
    emit(AppThemeModified(themeMode: newMode, seedColor: currentSeedColor, locale: currentLocale));
  }

  _setEasyLoadingTheme(ThemeMode mode) {
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
      seedColor: currentSeedColor,
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
