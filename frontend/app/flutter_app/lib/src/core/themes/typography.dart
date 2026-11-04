import 'package:flutter/material.dart';

import 'palette.dart';

/// 排印阶——docs/design-language.md §3 的 Dart 落版权威。
///
/// 中文三条纪律，写死在这里而不是各页面自行发挥：
/// 1. 强调只用字重与主色，不用斜体（中文斜体是排版事故）；
/// 2. 标题收字距（-0.01em 量级），正文不收；
/// 3. 计数与时间用等宽数字（tabularFigures），列表右对齐才不会跳。
abstract final class AppTypography {
  /// 数字等宽特性
  static const List<FontFeature> _tabular = [FontFeature.tabularFigures()];

  static TextTheme build(UiPalette p) {
    // 只给公共字段，字号/字重/行高逐档覆写，避免继承链里的隐式默认
    TextStyle ts({
      required double size,
      required FontWeight weight,
      double height = 1.5,
      Color? color,
      double spacing = 0,
      bool tabular = false,
    }) {
      return TextStyle(
        fontSize: size,
        fontWeight: weight,
        height: height,
        letterSpacing: spacing,
        color: color ?? p.ink,
        fontFeatures: tabular ? _tabular : null,
      );
    }

    return TextTheme(
      // 首页 hero：30–36 / 700（doc 写的 800 在苹方上会被合成加粗，压到 700）
      displayLarge: ts(
        size: 30,
        weight: FontWeight.w700,
        height: 1.28,
        spacing: -0.3,
        color: p.inkStrong,
      ),
      displayMedium: ts(
        size: 26,
        weight: FontWeight.w700,
        height: 1.3,
        spacing: -0.25,
        color: p.inkStrong,
      ),
      displaySmall: ts(
        size: 22,
        weight: FontWeight.w700,
        height: 1.35,
        spacing: -0.2,
        color: p.inkStrong,
      ),

      // 页面标题 / 文章题图标题
      headlineLarge: ts(
        size: 24,
        weight: FontWeight.w700,
        height: 1.38,
        spacing: -0.2,
        color: p.inkStrong,
      ),
      headlineMedium: ts(
        size: 20,
        weight: FontWeight.w700,
        height: 1.42,
        spacing: -0.15,
        color: p.inkStrong,
      ),
      headlineSmall: ts(
        size: 18,
        weight: FontWeight.w700,
        height: 1.45,
        color: p.inkStrong,
      ),

      // 区块标题（配 3px 主色左条）与卡片标题
      titleLarge: ts(
        size: 17,
        weight: FontWeight.w600,
        height: 1.5,
        color: p.inkStrong,
      ),
      titleMedium: ts(
        size: 15,
        weight: FontWeight.w600,
        height: 1.5,
        color: p.inkStrong,
      ),
      titleSmall: ts(
        size: 13,
        weight: FontWeight.w600,
        height: 1.5,
        color: p.ink,
      ),

      // 文章正文 16–17 / 1.8；常规正文 14–15 / 1.7
      bodyLarge: ts(size: 16, weight: FontWeight.w400, height: 1.8),
      bodyMedium: ts(size: 15, weight: FontWeight.w400, height: 1.75),
      bodySmall: ts(
        size: 13,
        weight: FontWeight.w400,
        height: 1.65,
        color: p.inkSecondary,
      ),

      // 控件文字
      labelLarge: ts(size: 14, weight: FontWeight.w500, height: 1.4),
      labelMedium: ts(
        size: 13,
        weight: FontWeight.w500,
        height: 1.4,
        color: p.inkSecondary,
        tabular: true,
      ),
      // 元信息：时间、计数、序号
      labelSmall: ts(
        size: 12,
        weight: FontWeight.w400,
        height: 1.5,
        color: p.inkTertiary,
        tabular: true,
      ),
    );
  }
}
