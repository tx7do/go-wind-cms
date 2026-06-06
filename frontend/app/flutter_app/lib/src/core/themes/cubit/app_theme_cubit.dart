import 'package:bloc/bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:flutter/material.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get_it/get_it.dart' show GetIt;

import 'package:flutter_app/src/core/repositories/index.dart';
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
    emit(AppThemeInitial(themeMode: currentMode, seedColor: color));
  }

  /// 设置初始化值
  void setInitValue(ThemeMode newMode) {
    _originalThemeMode ??= newMode;
    emit(AppThemeInitial(themeMode: newMode, seedColor: seedColor));
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

    emit(AppThemeModified(themeMode: newMode, seedColor: currentSeedColor));
  }

  /// 修改主题色
  void modifySeedColor(Color color) async {
    if (color == currentSeedColor) return;

    await _cache.setSeedColor(color);

    emit(AppThemeModified(themeMode: currentValue, seedColor: color));
  }

  /// 强制修改
  void forceModify(ThemeMode newMode) {
    _cache.setMaterialThemeMode(newMode);
    emit(AppThemeModified(themeMode: newMode, seedColor: currentSeedColor));
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
}
