import 'package:flutter/material.dart';

import 'package:flutter_app/src/core/widgets/responsive_layout.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_content_view.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_web_view.dart';

/// CMS 首页
///
/// 壳层（顶栏/底栏）归 AdaptiveShell；这里只决定正文内容形态：
/// - 手机（< 600dp）：单栏瀑布流（NestedScrollView + TabBarView）
/// - 平板/Web（>= 600dp）：宽屏视图（左右双栏 + 侧栏）
///
/// 注意：这里仍用 ResponsiveLayout 双树切换（P4 首页重写时会合并成单树），
/// 宽窄切换会重建 State——这是已知过渡态，不再承担 tab 路由职责。
class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return ResponsiveLayout(
      mobileBody: const HomeContentView(),
      webBody: const HomeWebView(),
    );
  }
}
