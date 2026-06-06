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

  const AppThemeInitial({this.themeMode = ThemeMode.light, this.seedColor});

  @override
  List<Object?> get props => [themeMode, seedColor];
}

// 修改状态
class AppThemeModified extends AppThemeState {
  final ThemeMode themeMode;
  final Color? seedColor;

  const AppThemeModified({this.themeMode = ThemeMode.light, this.seedColor});

  @override
  List<Object?> get props => [themeMode, seedColor];

  AppThemeModified copyWith({ThemeMode? themeMode, Color? seedColor}) => AppThemeModified(
        themeMode: themeMode ?? this.themeMode,
        seedColor: seedColor ?? this.seedColor,
      );
}
