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
  final Color? seedColor;
  final Locale? locale;

  const AppThemeInitial({this.themeMode = ThemeMode.light, this.seedColor, this.locale});

  @override
  List<Object?> get props => [themeMode, seedColor, locale];
}

// 修改状态
class AppThemeModified extends AppThemeState {
  final ThemeMode themeMode;
  final Color? seedColor;
  final Locale? locale;

  const AppThemeModified({this.themeMode = ThemeMode.light, this.seedColor, this.locale});

  @override
  List<Object?> get props => [themeMode, seedColor, locale];

  AppThemeModified copyWith({ThemeMode? themeMode, Color? seedColor, Locale? locale}) => AppThemeModified(
        themeMode: themeMode ?? this.themeMode,
        seedColor: seedColor ?? this.seedColor,
        locale: locale ?? this.locale,
      );
}
