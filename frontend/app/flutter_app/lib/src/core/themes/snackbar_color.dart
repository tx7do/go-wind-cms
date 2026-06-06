import 'dart:ui';

class SnackBarColor {
  final Color? error;
  final Color? info;
  final Color? warning;
  final Color? success;

  const SnackBarColor({
    required this.error,
    required this.info,
    required this.warning,
    required this.success,
  });

  SnackBarColor lerp(SnackBarColor? other, double t) {
    if (other is! SnackBarColor) {
      return this;
    }
    return SnackBarColor(
      error: Color.lerp(error, other.error, t),
      info: Color.lerp(info, other.info, t),
      warning: Color.lerp(warning, other.warning, t),
      success: Color.lerp(success, other.success, t),
    );
  }

  @override
  String toString() =>
      'SnackBarColor(error: $error, info: $info, warning: $warning, success: $success)';
}
