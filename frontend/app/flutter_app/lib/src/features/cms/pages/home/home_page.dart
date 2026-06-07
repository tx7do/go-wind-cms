import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_mobile_view.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_web_view.dart';

/// CMS 首页
///
/// 根据屏幕宽度自动切换：
/// - 手机（< 600dp）：单栏瀑布流（NestedScrollView + TabBarView）
/// - 平板/Web（>= 600dp）：宽屏视图（顶部导航 + 左右双栏）
class HomePage extends StatelessWidget {
  /// Mobile 模式下初始显示的路由 tab（如 '/profile'）
  final String? initialRoute;

  const HomePage({super.key, this.initialRoute});

  @override
  Widget build(BuildContext context) {
    return ResponsiveLayout(
      mobileBody: HomeMobileView(initialRoute: initialRoute),
      webBody: const HomeWebView(),
    );
  }
}
