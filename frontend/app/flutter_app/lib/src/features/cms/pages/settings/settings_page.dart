import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/themes/index.dart' as theme;
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/widgets/app_back_button.dart';

/// 设置页
///
/// 包含主题色选择、深色模式切换等外观设置。
class SettingsPage extends StatelessWidget {
  const SettingsPage({super.key});

  /// 预设主题色列表
  static const List<Color> _presetColors = [
    Color(0xFF3A7CA5), // 默认蓝
    Color(0xFF6750A4), // 紫色
    Color(0xFF006B5E), // 墨绿
    Color(0xFFB7262E), // 中国红
    Color(0xFFF57C00), // 橙色
    Color(0xFF1565C0), // 深蓝
    Color(0xFFAD1457), // 玫红
    Color(0xFF2E7D32), // 绿色
  ];

  @override
  Widget build(BuildContext context) {
    final isMobile = ResponsiveUtils.isMobile(context);
    final themeData = Theme.of(context);

    return Scaffold(
      backgroundColor: themeData.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: themeData.colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        leading: const AppBackButton(),
        title: Text(
          S.of(context).settings,
          style: TextStyle(
            fontSize: isMobile ? 18.sp : 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        centerTitle: true,
      ),
      body: Theme(
        data: themeData.copyWith(
          splashFactory: NoSplash.splashFactory,
          highlightColor: Colors.transparent,
          splashColor: Colors.transparent,
        ),
        child: SingleChildScrollView(
          padding: EdgeInsets.all(isMobile ? 16.w : 24),
          child: Center(
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 600),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 外观设置卡片
                  Card(
                    elevation: 0,
                    shape: RoundedRectangleBorder(
                      borderRadius:
                          BorderRadius.circular(isMobile ? 16.r : 16),
                    ),
                    color: themeData.colorScheme.surface,
                    child: Padding(
                      padding: EdgeInsets.all(isMobile ? 16.w : 20),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // 区标题
                          Row(
                            children: [
                              Icon(
                                Icons.palette_outlined,
                                size: 18,
                                color: themeData.colorScheme.primary,
                              ),
                              SizedBox(width: isMobile ? 8.w : 8),
                              Text(
                                S.of(context).appearance,
                                style: TextStyle(
                                  fontSize: isMobile ? 15.sp : 15,
                                  fontWeight: FontWeight.w600,
                                  color: themeData.colorScheme.onSurface,
                                ),
                              ),
                            ],
                          ),
                          SizedBox(height: isMobile ? 16.h : 16),

                          // 主题色
                          _buildThemeColorSection(context, isMobile),
                          SizedBox(height: isMobile ? 20.h : 20),

                          // 深色模式
                          _buildThemeModeSection(context, isMobile),
                          SizedBox(height: isMobile ? 20.h : 20),

                          // 语言切换
                          _buildLanguageSection(context, isMobile),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildThemeColorSection(BuildContext context, bool isMobile) {
    final themeData = Theme.of(context);
    final cubit = context.read<theme.AppThemeCubit>();
    final currentSeedColor = cubit.currentSeedColor;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          S.of(context).themeColor,
          style: TextStyle(
            fontSize: isMobile ? 13.sp : 13,
            color: themeData.colorScheme.onSurface.withAlpha(160),
          ),
        ),
        SizedBox(height: isMobile ? 10.h : 10),
        Wrap(
          spacing: 12,
          runSpacing: 8,
          children: _presetColors.map((color) {
            final isSelected =
                color.toARGB32() == currentSeedColor.toARGB32();
            return MouseRegion(
              cursor: SystemMouseCursors.click,
              child: GestureDetector(
                onTap: () => cubit.modifySeedColor(color),
                child: AnimatedContainer(
                  duration: const Duration(milliseconds: 200),
                  width: isMobile ? 40.w : 40,
                  height: isMobile ? 40.w : 40,
                  decoration: BoxDecoration(
                    color: color,
                    shape: BoxShape.circle,
                    border: isSelected
                        ? Border.all(
                            color: themeData.colorScheme.onSurface,
                            width: 3,
                          )
                        : null,
                    boxShadow: isSelected
                        ? [
                            BoxShadow(
                              color: color.withAlpha((0.4 * 255).round()),
                              blurRadius: 8,
                              offset: const Offset(0, 2),
                            ),
                          ]
                        : null,
                  ),
                  child: isSelected
                      ? Icon(
                          Icons.check,
                          size: isMobile ? 20.sp : 20,
                          color: Colors.white,
                        )
                      : null,
                ),
              ),
            );
          }).toList(),
        ),
      ],
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
            color: themeData.colorScheme.outline.withAlpha(100),
            width: 1,
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
    final fontSize = isMobile ? 12.sp : 12.0;

    return GestureDetector(
      onTap: onTap,
      child: IntrinsicWidth(
        child: Container(
          padding: EdgeInsets.symmetric(
            horizontal: isMobile ? 16.w : 16,
            vertical: isMobile ? 8.h : 8,
          ),
          decoration: BoxDecoration(
            color: isSelected
                ? themeData.colorScheme.primaryContainer
                : Colors.transparent,
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
                Icon(
                  Icons.check,
                  size: isMobile ? 14.sp : 14,
                  color: themeData.colorScheme.onSurface,
                ),
              if (isSelected) SizedBox(width: isMobile ? 4.w : 4),
              Text(
                label,
                style: TextStyle(
                  fontSize: fontSize,
                  fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
                  color: themeData.colorScheme.onSurface,
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
            color: themeData.colorScheme.outline.withAlpha(100),
            width: 1,
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

    if (isMobile) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(
                icon,
                size: 20.sp,
                color: themeData.colorScheme.onSurface.withAlpha(180),
              ),
              SizedBox(width: 10.w),
              Text(
                title,
                style: TextStyle(
                  fontSize: 14.sp,
                  color: themeData.colorScheme.onSurface,
                ),
              ),
            ],
          ),
          SizedBox(height: 8.h),
          SizedBox(
            width: double.infinity,
            child: trailing,
          ),
        ],
      );
    }

    return Row(
      children: [
        Icon(
          icon,
          size: 20,
          color: themeData.colorScheme.onSurface.withAlpha(180),
        ),
        const SizedBox(width: 10),
        Flexible(
          child: Text(
            title,
            style: TextStyle(
              fontSize: 14,
              color: themeData.colorScheme.onSurface,
            ),
            overflow: TextOverflow.ellipsis,
          ),
        ),
        const Spacer(),
        trailing,
      ],
    );
  }
}
