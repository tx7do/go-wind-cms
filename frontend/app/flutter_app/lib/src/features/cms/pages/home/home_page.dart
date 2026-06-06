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
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return const ResponsiveLayout(
      mobileBody: HomeMobileView(),
      webBody: HomeWebView(),
    );
  }
}
