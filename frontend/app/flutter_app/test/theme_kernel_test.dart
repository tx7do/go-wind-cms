import 'package:flutter/foundation.dart' show TargetPlatform, debugDefaultTargetPlatformOverride;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:flutter_app/src/core/themes/index.dart';

/// 设计系统内核契约测试。
///
/// 这些断言钉的是「中式克制」这条方向的硬指标：一旦被改回 MD3 默认观感
/// （大圆角、药丸指示器、elevation tint、灰粉容器），测试会红。
void main() {
  group('晴风 / 夜风 双态', () {
    test('主色固定为品牌蓝，且亮暗两态各自命中 token', () {
      final light = getLightTheme();
      final dark = getDarkTheme();

      expect(light.palette.primary, const Color(0xFF006BE6));
      expect(dark.palette.primary, const Color(0xFF2E96FF));
      expect(light.brightness, Brightness.light);
      expect(dark.brightness, Brightness.dark);
      // 夜风 filled 按钮底字是深墨而非白（#2E96FF 配白字只有 2.9:1）
      expect(dark.palette.onPrimary, const Color(0xFF041426));
    });

    test('elevation tint 关闭：分层只靠底色阶梯，不靠 MD3 泛色', () {
      for (final theme in [getLightTheme(), getDarkTheme()]) {
        expect(theme.colorScheme.surfaceTint, Colors.transparent);
        expect(theme.cardTheme.elevation, 0);
        expect(theme.appBarTheme.scrolledUnderElevation, 0);
      }
    });

    test('圆角收成小圆角体系，卡片描边而非投影', () {
      final theme = getLightTheme();
      expect(theme.radius.md, 9);
      expect(theme.cardTheme.shape, isA<RoundedRectangleBorder>());
      final shape = theme.cardTheme.shape! as RoundedRectangleBorder;
      expect(shape.borderRadius, theme.radius.mdAll);
      expect(shape.side.color, theme.palette.border);
    });

    test('TabBar 用下划线而非 MD3 药丸指示器', () {
      final theme = getLightTheme();
      expect(theme.tabBarTheme.indicator, isA<UnderlineTabIndicator>());
      final indicator = theme.tabBarTheme.indicator! as UnderlineTabIndicator;
      expect(indicator.borderSide.width, 2);
      expect(indicator.borderSide.color, theme.palette.primary);
    });

    test('排印：中文正文行高 ≥1.75，标题收字距', () {
      final theme = getLightTheme();
      expect(theme.textTheme.bodyLarge!.height! >= 1.75, isTrue);
      expect(theme.textTheme.displayLarge!.letterSpacing! <= 0, isTrue);
      expect(theme.textTheme.displayLarge!.fontStyle, isNot(FontStyle.italic));
    });

    test('字族：系统优先，内置思源退到 fallback', () {
      // 3.47 的 ThemeData 不再暴露 fontFamily getter，只能从 textTheme 读实际落地值
      String? bodyFamilyOf(ThemeData t) => t.textTheme.bodyMedium!.fontFamily;
      List<String>? fallbackOf(ThemeData t) =>
          t.textTheme.bodyMedium!.fontFamilyFallback;

      // 测试环境 defaultTargetPlatform 默认 android，逐平台覆盖验证
      debugDefaultTargetPlatformOverride = TargetPlatform.macOS;
      expect(bodyFamilyOf(getLightTheme()), 'PingFang SC');
      debugDefaultTargetPlatformOverride = TargetPlatform.windows;
      expect(bodyFamilyOf(getLightTheme()), 'Microsoft YaHei');
      debugDefaultTargetPlatformOverride = TargetPlatform.android;
      expect(bodyFamilyOf(getLightTheme()), AppFont.bundledFamily);
      debugDefaultTargetPlatformOverride = null;

      expect(fallbackOf(getLightTheme()), contains(AppFont.bundledFamily));
      // primaryTextTheme 不重染颜色：灰档必须逐角色保留
      final light = getLightTheme();
      expect(light.primaryTextTheme.bodySmall!.color, light.palette.inkSecondary);
      expect(light.textTheme.labelSmall!.color, light.palette.inkTertiary);
    });
  });

  testWidgets('主题渲染冒烟：页面可构建且 token 可达', (tester) async {
    late BuildContext captured;
    await tester.pumpWidget(
      MaterialApp(
        theme: getLightTheme(),
        home: Builder(
          builder: (context) {
            captured = context;
            return Scaffold(
              appBar: AppBar(title: const Text('晴风')),
              body: const Center(
                child: Card(child: Text('内容卡片')),
              ),
            );
          },
        ),
      ),
    );

    expect(find.text('晴风'), findsOneWidget);
    expect(find.text('内容卡片'), findsOneWidget);
    // ThemeDataTokens 扩展在真实上下文中解析出整套 extension
    expect(Theme.of(captured).palette.background, const Color(0xFFF4F7F9));
    expect(Theme.of(captured).space.cardPadding, 20);
  });
}
