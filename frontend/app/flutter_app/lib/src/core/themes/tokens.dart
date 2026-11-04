import 'package:flutter/material.dart';

import 'palette.dart';

/// 圆角分级。刻意比 Material 默认更小、更成体系：
/// 大圆角（20+）是 MD3 观感的主要来源，中式界面靠小圆角 + 细线显「利落」。
@immutable
class RadiusTokens {
  const RadiusTokens({
    required this.xs,
    required this.sm,
    required this.md,
    required this.lg,
    required this.xl,
    required this.full,
  });

  static const RadiusTokens light = RadiusTokens(
    xs: 2,
    sm: 6,
    md: 9,
    lg: 12,
    xl: 16,
    full: 999,
  );

  /// 数值（dp）
  final double xs;
  final double sm;
  final double md;
  final double lg;
  final double xl;
  final double full;

  BorderRadius get xsAll => BorderRadius.circular(xs);
  BorderRadius get smAll => BorderRadius.circular(sm);
  BorderRadius get mdAll => BorderRadius.circular(md);
  BorderRadius get lgAll => BorderRadius.circular(lg);
  BorderRadius get xlAll => BorderRadius.circular(xl);
  BorderRadius get fullAll => BorderRadius.circular(full);

  /// 底部抽屉：只圆上沿
  BorderRadius get lgTop => BorderRadius.vertical(top: Radius.circular(lg));

  /// 题图/媒体块：上方圆角跟卡片，下方直角
  BorderRadius mediaTop(double cardRadius) => BorderRadius.vertical(
    top: Radius.circular(cardRadius),
  );
}

/// 4px 基数的间距阶（doc §4）。
/// 语义命名优先于数值命名：写 `space.card` 而不是 `space.xl`，
/// 页面才不会各自发明「随手 18」。
@immutable
class SpaceTokens {
  const SpaceTokens({
    required this.xxs,
    required this.xs,
    required this.sm,
    required this.md,
    required this.lg,
    required this.xl,
    required this.xxl,
    required this.section,
    required this.block,
    required this.cardPadding,
    required this.listRowGap,
    required this.gridGap,
    required this.chipHeight,
    required this.controlHeight,
  });

  static const SpaceTokens light = SpaceTokens(
    xxs: 2,
    xs: 4,
    sm: 8,
    md: 12,
    lg: 16,
    xl: 20,
    xxl: 24,
    section: 48,
    block: 64,
    // 卡片内边距 20（doc §4 的 p-5）
    cardPadding: 20,
    listRowGap: 12,
    gridGap: 16,
    chipHeight: 28,
    // 桌面 40、移动 44（doc §7 按钮高度）
    controlHeight: 40,
  );

  final double xxs;
  final double xs;
  final double sm;
  final double md;
  final double lg;
  final double xl;
  final double xxl;

  /// 区块纵向节奏
  final double section;
  final double block;

  final double cardPadding;
  final double listRowGap;
  final double gridGap;
  final double chipHeight;
  final double controlHeight;
}

/// 线宽：细线是这套设计的主要分层手段。
@immutable
class BorderTokens {
  const BorderTokens({
    required this.hairline,
    required this.thin,
    required this.emphasis,
  });

  static const BorderTokens light = BorderTokens(
    hairline: 0.5,
    thin: 1,
    emphasis: 1.5,
  );

  final double hairline;
  final double thin;
  final double emphasis;
}

@immutable
class MotionTokens {
  const MotionTokens({
    required this.micro,
    required this.fast,
    required this.base,
    required this.slow,
    required this.curve,
    required this.lift,
  });

  static const MotionTokens light = MotionTokens(
    micro: Duration(milliseconds: 150),
    fast: Duration(milliseconds: 200),
    base: Duration(milliseconds: 300),
    slow: Duration(milliseconds: 450),
    curve: Curves.easeOutCubic,
    // hover 位移上限（doc §6：≤4px）
    lift: 3,
  );

  final Duration micro;
  final Duration fast;
  final Duration base;
  final Duration slow;
  final Curve curve;
  final double lift;
}

@immutable
class ShadowTokens {
  const ShadowTokens({
    required this.card,
    required this.hover,
    required this.overlay,
  });

  final List<BoxShadow> card;
  final List<BoxShadow> hover;
  final List<BoxShadow> overlay;

  factory ShadowTokens.fromPalette(UiPalette p) {
    // 深色下投影几乎不可见，浮起改由底色阶梯表达（见 palette.kNightPalette 注释）
    final cardBlur = p.isDark ? 8.0 : 12.0;
    return ShadowTokens(
      card: [
        BoxShadow(
          color: p.cardShadow,
          blurRadius: cardBlur,
          offset: Offset(0, p.isDark ? 2 : 4),
        ),
      ],
      hover: [
        BoxShadow(
          color: p.hoverShadow,
          blurRadius: p.isDark ? 14 : 30,
          offset: Offset(0, p.isDark ? 4 : 12),
        ),
      ],
      overlay: [
        BoxShadow(
          color: p.overlayShadow,
          blurRadius: p.isDark ? 24 : 20,
          offset: Offset(0, p.isDark ? 8 : 8),
        ),
      ],
    );
  }
}

/// 挂在 ThemeData 上的整套设计 token。
///
/// 页面只允许从这里取值，不允许写死色值/圆角/间距——这是把 MD3 观感
/// 拉回中式克制的唯一着力点。
@immutable
class AppThemeExtension extends ThemeExtension<AppThemeExtension> {
  const AppThemeExtension({
    required this.palette,
    required this.radius,
    required this.space,
    required this.border,
    required this.motion,
    required this.shadow,
  });

  final UiPalette palette;
  final RadiusTokens radius;
  final SpaceTokens space;
  final BorderTokens border;
  final MotionTokens motion;
  final ShadowTokens shadow;

  static AppThemeExtension of(BuildContext context) =>
      Theme.of(context).appTokens;

  /// 主题切换时 ThemeData 会 lerp 整套 extension。
  /// 色值真插值（避免切换瞬间闪色），离散数值过半切换。
  @override
  AppThemeExtension lerp(AppThemeExtension? other, double t) {
    if (other == null) return this;
    if (t < 0.5) return this;
    return other;
  }

  @override
  AppThemeExtension copyWith({
    UiPalette? palette,
    RadiusTokens? radius,
    SpaceTokens? space,
    BorderTokens? border,
    MotionTokens? motion,
    ShadowTokens? shadow,
  }) {
    return AppThemeExtension(
      palette: palette ?? this.palette,
      radius: radius ?? this.radius,
      space: space ?? this.space,
      border: border ?? this.border,
      motion: motion ?? this.motion,
      shadow: shadow ?? this.shadow,
    );
  }
}

/// `Theme.of(context).palette` / `.radius` / `.space` / `.border` / `.motion` / `.shadow`
extension ThemeDataTokens on ThemeData {
  AppThemeExtension get appTokens =>
      extension<AppThemeExtension>() ??
      // 兜底：极少数测试上下文里没有挂 extension，宁可给默认值也别崩
      const AppThemeExtension(
        palette: kDayPalette,
        radius: RadiusTokens.light,
        space: SpaceTokens.light,
        border: BorderTokens.light,
        motion: MotionTokens.light,
        shadow: ShadowTokens(
          card: [],
          hover: [],
          overlay: [],
        ),
      );

  UiPalette get palette => appTokens.palette;
  RadiusTokens get radius => appTokens.radius;
  SpaceTokens get space => appTokens.space;
  BorderTokens get borderWidth => appTokens.border;
  MotionTokens get motion => appTokens.motion;
  ShadowTokens get shadow => appTokens.shadow;

  /// 卡片描边（hover 前）
  BorderSide get cardSide => BorderSide(width: 1, color: palette.border);
}
