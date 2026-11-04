import 'package:flutter_test/flutter_test.dart';

import 'package:flutter_app/src/core/layout/width_class.dart';

/// 分档与栅格反推的契约测试。
///
/// 钉的是「宽窄屏适配」的两条硬指标：边界不重叠、列数不随宽度倒退。
void main() {
  group('WidthClass 分档', () {
    test('边界值归入更高一档，四档无缝隙', () {
      expect(WidthClassOf.fromWidth(375), WidthClass.compact);
      expect(WidthClassOf.fromWidth(599.9), WidthClass.compact);
      expect(WidthClassOf.fromWidth(600), WidthClass.medium);
      expect(WidthClassOf.fromWidth(899.9), WidthClass.medium);
      expect(WidthClassOf.fromWidth(900), WidthClass.expanded);
      expect(WidthClassOf.fromWidth(1199.9), WidthClass.expanded);
      expect(WidthClassOf.fromWidth(1200), WidthClass.wide);
      expect(WidthClassOf.fromWidth(2560), WidthClass.wide);
    });

    test('侧边栏只在 expanded/wide 出现，导航形态只在 compact 走底栏', () {
      for (final wc in WidthClass.values) {
        expect(wc.showSidebar, wc == WidthClass.expanded || wc == WidthClass.wide);
        expect(wc.useBottomNav, wc == WidthClass.compact);
      }
    });

    test('内容限宽：compact 满宽不限，其余三档单调不减', () {
      expect(WidthClass.compact.contentMaxWidth, double.infinity);
      var previous = 0.0;
      for (final wc in WidthClass.values.skip(1)) {
        expect(wc.contentMaxWidth >= previous, isTrue, reason: '$wc 限宽倒退');
        previous = wc.contentMaxWidth;
      }
      // medium 档至少要能容纳两张 300dp 卡片 + 沟槽，否则双栏是假双栏
      expect(WidthClass.medium.contentMaxWidth >= 600, isTrue);
    });
  });

  group('AdaptiveGrid 列数反推', () {
    test('卡片列数随宽度单调不减（窗口拖宽不会突然变少一列）', () {
      var previous = 0;
      for (var w = 300.0; w <= 2000.0; w += 1) {
        final c = AdaptiveGrid.columns(
          availableWidth: w,
          minItemWidth: AdaptiveGrid.minPostCardWidth,
          maxColumns: 3,
        );
        expect(c >= previous, isTrue, reason: '$w dp 处列数倒退：$c < $previous');
        expect(c >= 1 && c <= 3, isTrue);
        previous = c;
      }
    });

    test('列数不牺牲最小卡片宽', () {
      // 900dp 主列：300*3 + 16*2 = 932 > 900，必须落回 2 列
      expect(
        AdaptiveGrid.columns(
          availableWidth: 900,
          minItemWidth: AdaptiveGrid.minPostCardWidth,
          maxColumns: 3,
        ),
        2,
      );
      expect(
        AdaptiveGrid.columns(
          availableWidth: 932,
          minItemWidth: AdaptiveGrid.minPostCardWidth,
          maxColumns: 3,
        ),
        3,
      );
    });

    test('异常输入不炸（无限宽 / 零宽退化为单列）', () {
      expect(
        AdaptiveGrid.columns(
          availableWidth: double.infinity,
          minItemWidth: 300,
        ),
        1,
      );
      expect(AdaptiveGrid.columns(availableWidth: 0, minItemWidth: 300), 1);
    });
  });

  test('正文行长上限恒定 760，不随档位放宽', () {
    expect(kArticleLineWidth, 760);
  });
}
