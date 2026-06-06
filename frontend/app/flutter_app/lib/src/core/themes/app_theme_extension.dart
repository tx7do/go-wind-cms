import 'package:flutter/material.dart';

import 'snackbar_color.dart' show SnackBarColor;

@immutable
class AppThemeExtension extends ThemeExtension<AppThemeExtension> {
  const AppThemeExtension({
    required this.snackBarColor,
  });

  final SnackBarColor snackBarColor;

  @override
  AppThemeExtension copyWith({SnackBarColor? snackBarColor}) {
    return AppThemeExtension(
      snackBarColor: snackBarColor ?? this.snackBarColor,
    );
  }

  @override
  AppThemeExtension lerp(AppThemeExtension? other, double t) {
    if (other is! AppThemeExtension) {
      return this;
    }
    return AppThemeExtension(
      snackBarColor: snackBarColor.lerp(other.snackBarColor, t),
    );
  }

  // Optional
  @override
  String toString() =>
      'AppThemeExtension(snackBarColor: ${snackBarColor.toString()})';
}
