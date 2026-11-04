part of 'app_theme_cubit.dart';

@immutable
sealed class AppThemeState extends Equatable {
  const AppThemeState();

  @override
  List<Object?> get props => [];
}

// 初始状态
final class AppThemeInitial extends AppThemeState {
  final ThemeMode themeMode;
  final Locale? locale;

  const AppThemeInitial({this.themeMode = ThemeMode.light, this.locale});

  @override
  List<Object?> get props => [themeMode, locale];
}

// 修改状态
class AppThemeModified extends AppThemeState {
  final ThemeMode themeMode;
  final Locale? locale;

  const AppThemeModified({this.themeMode = ThemeMode.light, this.locale});

  @override
  List<Object?> get props => [themeMode, locale];

  AppThemeModified copyWith({ThemeMode? themeMode, Locale? locale}) =>
      AppThemeModified(
        themeMode: themeMode ?? this.themeMode,
        locale: locale ?? this.locale,
      );
}
