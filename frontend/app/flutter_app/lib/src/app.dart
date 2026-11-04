import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

import 'package:flutter_app/generated/l10n.dart' as l10n;

import 'package:flutter_app/src/core/extensions/app_localizations_context.dart';
import 'package:flutter_app/src/core/themes/index.dart' as theme;
import 'package:flutter_app/src/core/utilities/logger.dart' show debug;

import 'app_router/app_router.dart' as pages;

/// 应用程序
class CMSApp extends StatefulWidget {
  const CMSApp({super.key});

  @override
  State<CMSApp> createState() => _CMSAppState();
}

class _CMSAppState extends State<CMSApp> {
  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    // ScreenUtil 已移除：观感全部走 token + WidthClass，不再按 designSize 缩放
    return _buildMaterialApp(context);
  }

  Widget _buildMaterialApp(BuildContext context) {
    final delegate = l10n.S.delegate;
    final cubit = context.read<theme.AppThemeCubit>();
    cubit.init();

    return MaterialApp.router(
      debugShowCheckedModeBanner: false,
      onGenerateTitle: (context) => context.loc.appName,
      routerConfig: pages.AppRouter.router,
      theme: theme.getLightTheme(),
      darkTheme: theme.getDarkTheme(),
      // 3.47 起滚动观感归 MaterialApp（ThemeData 不再收 scrollBehavior）
      scrollBehavior: const theme.AppScrollBehavior(),
      themeMode: context.watch<theme.AppThemeCubit>().themeMode,
      locale: context.watch<theme.AppThemeCubit>().currentLocale,
      builder: (context, routerChild) {
        // EasyLoading 作为路由的外部包装
        final easyLoadingBuilder = EasyLoading.init();
        return easyLoadingBuilder(context, routerChild);
      },
      localizationsDelegates: [
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
        delegate,
      ],
      supportedLocales: l10n.S.delegate.supportedLocales,
      localeListResolutionCallback: (locales, supportedLocales) {
        debug('locale list: $locales');
        if (locales == null || locales.isEmpty) {
          return supportedLocales.first;
        }
        for (final locale in locales) {
          // 精确匹配
          if (supportedLocales.contains(locale)) return locale;
          // 语言代码匹配（忽略国家代码差异，如 zh-CN → zh_CN）
          final match = supportedLocales.where(
            (sl) => sl.languageCode == locale.languageCode,
          );
          if (match.isNotEmpty) return match.first;
        }
        return supportedLocales.first;
      },
    );
  }
}
