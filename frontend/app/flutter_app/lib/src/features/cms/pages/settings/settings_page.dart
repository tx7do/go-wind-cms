import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/layout/width_class.dart';
import 'package:flutter_app/src/core/themes/index.dart' as theme;
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';

/// 设置页
///
/// 外观设置：主题模式（跟随系统 / 晴风 / 夜风）与语言。
/// 主题色固定为品牌蓝，不再提供色板选择。
class SettingsPage extends StatelessWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context) {
    final compact = context.widthClass.isCompact;
    final themeData = Theme.of(context);
    final space = themeData.space;

    final body = SingleChildScrollView(
      padding: EdgeInsets.all(compact ? space.lg : space.xxl),
      child: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 600),
          child: Card(
            elevation: 0,
            margin: EdgeInsets.zero,
            child: Padding(
              padding: EdgeInsets.all(space.cardPadding),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 区标题
                  Row(
                    children: [
                      Icon(
                        Icons.palette_outlined,
                        size: 18,
                        color: themeData.palette.primary,
                      ),
                      SizedBox(width: space.sm),
                      Text(
                        S.of(context).appearance,
                        style: themeData.textTheme.titleMedium,
                      ),
                    ],
                  ),
                  SizedBox(height: space.lg),

                  // 主题模式
                  _buildThemeModeSection(context, compact),
                  SizedBox(height: space.xl),

                  // 语言切换
                  _buildLanguageSection(context, compact),
                ],
              ),
            ),
          ),
        ),
      ),
    );

    // 宽屏由壳层 AdaptiveShell 提供 Scaffold
    if (!compact) return body;

    return Scaffold(
      backgroundColor: themeData.scaffoldBackgroundColor,
      appBar: AppBar(
        leading: const AppBackButton(),
        title: Text(S.of(context).settings),
        centerTitle: true,
      ),
      body: body,
    );
  }

  Widget _buildThemeModeSection(BuildContext context, bool isMobile) {
    final cubit = context.read<theme.AppThemeCubit>();
    final currentThemeMode = cubit.currentValue;
    final themeData = Theme.of(context);

    return _SettingRow(
      icon: Icons.dark_mode_outlined,
      title: S.of(context).darkMode,
      isMobile: isMobile,
      trailing: Container(
        decoration: BoxDecoration(
          border: Border.all(
            color: themeData.palette.borderStrong,
            width: themeData.borderWidth.thin,
          ),
          borderRadius: BorderRadius.circular(20),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            _buildSegmentButton(
              context,
              S.of(context).light,
              currentThemeMode == ThemeMode.light,
              () => cubit.modify(ThemeMode.light),
              isMobile,
              isFirst: true,
            ),
            _buildSegmentButton(
              context,
              S.of(context).followSystem,
              currentThemeMode == ThemeMode.system,
              () => cubit.modify(ThemeMode.system),
              isMobile,
            ),
            _buildSegmentButton(
              context,
              S.of(context).dark,
              currentThemeMode == ThemeMode.dark,
              () => cubit.modify(ThemeMode.dark),
              isMobile,
              isLast: true,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSegmentButton(
    BuildContext context,
    String label,
    bool isSelected,
    VoidCallback onTap,
    bool isMobile, {
    bool isFirst = false,
    bool isLast = false,
  }) {
    final themeData = Theme.of(context);
    final p = themeData.palette;

    return GestureDetector(
      onTap: onTap,
      child: IntrinsicWidth(
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          decoration: BoxDecoration(
            color: isSelected ? p.primarySoft : Colors.transparent,
            borderRadius: BorderRadius.only(
              topLeft: isFirst ? const Radius.circular(19) : Radius.zero,
              bottomLeft: isFirst ? const Radius.circular(19) : Radius.zero,
              topRight: isLast ? const Radius.circular(19) : Radius.zero,
              bottomRight: isLast ? const Radius.circular(19) : Radius.zero,
            ),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (isSelected)
                Icon(Icons.check, size: 14, color: p.primaryOnSoft),
              if (isSelected) const SizedBox(width: 4),
              Text(
                label,
                style: themeData.textTheme.labelMedium?.copyWith(
                  fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                  color: isSelected ? p.primaryOnSoft : p.inkSecondary,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildLanguageSection(BuildContext context, bool isMobile) {
    final cubit = context.read<theme.AppThemeCubit>();
    final currentLocale = cubit.currentLocale;
    final supportedLocales = cubit.supportedLocales;
    final themeData = Theme.of(context);

    // 使用简短名称：中文(中国) / English
    final localeLabels = <Locale, String>{
      const Locale('zh', 'CN'): '中文(中国)',
      const Locale('en', 'US'): 'English',
    };

    return _SettingRow(
      icon: Icons.language,
      title: S.of(context).language,
      isMobile: isMobile,
      trailing: Container(
        decoration: BoxDecoration(
          border: Border.all(
            color: themeData.palette.borderStrong,
            width: themeData.borderWidth.thin,
          ),
          borderRadius: BorderRadius.circular(20),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: supportedLocales.asMap().entries.map((entry) {
            final index = entry.key;
            final locale = entry.value;
            final label = localeLabels[locale] ?? locale.toLanguageTag();
            final isSelected = currentLocale == locale;
            
            return _buildSegmentButton(
              context,
              label,
              isSelected,
              () => cubit.modifyLocale(locale),
              isMobile,
              isFirst: index == 0,
              isLast: index == supportedLocales.length - 1,
            );
          }).toList(),
        ),
      ),
    );
  }
}

/// 设置行
class _SettingRow extends StatelessWidget {
  final IconData icon;
  final String title;
  final bool isMobile;
  final Widget trailing;

  const _SettingRow({
    required this.icon,
    required this.title,
    required this.isMobile,
    required this.trailing,
  });

  @override
  Widget build(BuildContext context) {
    final themeData = Theme.of(context);
    final space = themeData.space;

    final iconAndTitle = Row(
      children: [
        Icon(icon, size: 20, color: themeData.palette.inkSecondary),
        SizedBox(width: space.sm + 2),
        Flexible(
          child: Text(
            title,
            style: themeData.textTheme.bodyMedium,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );

    if (isMobile) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          iconAndTitle,
          SizedBox(height: space.sm),
          SizedBox(width: double.infinity, child: trailing),
        ],
      );
    }

    return Row(
      children: [
        Expanded(child: iconAndTitle),
        SizedBox(width: space.lg),
        trailing,
      ],
    );
  }
}
