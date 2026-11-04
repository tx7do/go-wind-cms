import 'package:flutter/material.dart';

/// 双态原始色值——本项目视觉 token 的单一事实源。
///
/// 权威出处 docs/design-language.md §2.2（浅色「晴风」/ 深色「夜风」）。
/// 与文档的三处有意偏离，都写在字段注释里：
/// 1. 文字分四档而非文档的单 foreground（中式排印讲究标题重、正文退、元信息灰）；
/// 2. 深色模式 filled 按钮的底字改深墨（#2E96FF 配白字对比度只有 2.9:1）；
/// 3. 补了文档没定义的 surfaceLow/High/Top 三档灰阶（分层靠底色，不靠 MD3 elevation tint）。
@immutable
class UiPalette {
  const UiPalette({
    required this.background,
    required this.surface,
    required this.surfaceLow,
    required this.surfaceHigh,
    required this.surfaceTop,
    required this.overlay,
    required this.inkStrong,
    required this.ink,
    required this.inkSecondary,
    required this.inkTertiary,
    required this.inkDisabled,
    required this.border,
    required this.borderStrong,
    required this.primary,
    required this.primaryDeep,
    required this.primaryOnSoft,
    required this.primarySoft,
    required this.primaryTint,
    required this.onPrimary,
    required this.inverseSurface,
    required this.inverseInk,
    required this.success,
    required this.warning,
    required this.danger,
    required this.successText,
    required this.warningText,
    required this.dangerText,
    required this.successSoft,
    required this.warningSoft,
    required this.dangerSoft,
    required this.scrim,
    required this.onImage,
    required this.heroStart,
    required this.heroMid,
    required this.heroEnd,
    required this.cardShadow,
    required this.hoverShadow,
    required this.overlayShadow,
    required this.shadowInk,
    required this.isDark,
  });

  /// 页面底
  final Color background;

  /// 卡片 / 顶栏 / 底栏底
  final Color surface;

  /// 输入框等内凹底（比 surface 退一档）
  final Color surfaceLow;

  /// hover 浮起底（比 surface 进一档）
  final Color surfaceHigh;

  /// 选中底 / 骨架屏（再进一档）
  final Color surfaceTop;

  /// 弹层、对话框、下拉
  final Color overlay;

  /// 标题专用：比正文重一档
  final Color inkStrong;

  /// 正文
  final Color ink;

  /// 摘要、说明
  final Color inkSecondary;

  /// 元信息（时间、计数）——不可用于可点文字
  final Color inkTertiary;

  /// 禁用态文字与图标
  final Color inkDisabled;

  /// 卡片描边、分隔线
  final Color border;

  /// 输入框描边
  final Color borderStrong;

  /// 唯一行动色
  final Color primary;

  /// primary 的按下 / 悬停变体
  final Color primaryDeep;

  /// 浅底（primarySoft）上的主色文字
  final Color primaryOnSoft;

  /// 徽底、头像底、激活底色
  final Color primarySoft;

  /// 比 primarySoft 更淡一档（大面积 tint）
  final Color primaryTint;

  /// primary 实色底上的文字
  final Color onPrimary;

  /// 反白底（SnackBar、Tooltip）
  final Color inverseSurface;

  /// inverseSurface 上的文字
  final Color inverseInk;

  final Color success;
  final Color warning;
  final Color danger;

  /// 语义色的「文字档」：实色语义色压在白底上对比度不足 4.5:1
  /// （warning #F59E0B 只有 2.1:1），文字一律走这三档。
  final Color successText;
  final Color warningText;
  final Color dangerText;

  final Color successSoft;
  final Color warningSoft;
  final Color dangerSoft;

  /// 题图/轮播上的文字压暗层
  final Color scrim;

  /// 照片上的文字色——图片是主题无关的，配深色 scrim 两态都用白字
  final Color onImage;

  /// hero 渐变三端（仅大面积氛围层与品牌图形允许用渐变）
  final Color heroStart;
  final Color heroMid;
  final Color heroEnd;

  /// 阴影：静态卡片 / hover 浮起 / 弹层
  final Color cardShadow;
  final Color hoverShadow;
  final Color overlayShadow;

  /// ColorScheme.shadow 用的不透明投影基色（本项目 elevation 全为 0，
  /// 仅兜住少数默认取 shadowColor 的组件）。
  final Color shadowInk;

  final bool isDark;

  /// 主色描边（hover 时卡片边框转这一档，对应 doc §2.3 的 primary/30）
  Color get primaryBorder => primary.withAlpha(isDark ? 0x40 : 0x4D);

  /// 细线与骨架的通用弱描边
  Color get hairline => border.withAlpha(isDark ? 0xCC : 0x99);

  /// MD3 secondary 一律取灰蓝，不引入第三色相（doc §1 铁律 1）。
  Color get secondary => isDark ? const Color(0xFF8FA6BD) : const Color(0xFF475569);

  /// 亮态 surface 即纯白——用它代替 Colors.white，避免深色模式误用白字
  Color get inkOnSolid => isDark ? background : surface;

  /// hero 渐变（氛围层专用，见 doc §2）
  Gradient get heroGradient => LinearGradient(
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
        colors: [heroStart, heroMid, heroEnd],
      );
}

/// 晴风（浅色）
const UiPalette kDayPalette = UiPalette(
  background: Color(0xFFF4F7F9),
  surface: Color(0xFFFFFFFF),
  surfaceLow: Color(0xFFF8FAFC),
  surfaceHigh: Color(0xFFEFF3F8),
  surfaceTop: Color(0xFFE6EDF5),
  overlay: Color(0xFFFFFFFF),

  inkStrong: Color(0xFF0F172A),
  ink: Color(0xFF1E293B),
  inkSecondary: Color(0xFF64748B),
  inkTertiary: Color(0xFF94A3B8),
  inkDisabled: Color(0xFFCBD5E1),

  border: Color(0xFFE4E9F0),
  borderStrong: Color(0xFFD6DFEA),

  primary: Color(0xFF006BE6),
  primaryDeep: Color(0xFF005AC2),
  primaryOnSoft: Color(0xFF005AB5),
  primarySoft: Color(0xFFE6F1FE),
  primaryTint: Color(0xFFF0F6FF),
  onPrimary: Color(0xFFFFFFFF),

  inverseSurface: Color(0xFF0F172A),
  inverseInk: Color(0xFFF4F7F9),

  success: Color(0xFF16A34A),
  warning: Color(0xFFF59E0B),
  danger: Color(0xFFEF4444),
  successText: Color(0xFF15803D),
  warningText: Color(0xFFB45309),
  dangerText: Color(0xFFDC2626),
  successSoft: Color(0xFFE6F7EC),
  warningSoft: Color(0xFFFEF4E0),
  dangerSoft: Color(0xFFFEEBEB),

  scrim: Color(0x6B020817),
  onImage: Color(0xFFFFFFFF),
  heroStart: Color(0xFF70B3FF),
  heroMid: Color(0xFF0D74F2),
  heroEnd: Color(0xFF005AC2),

  cardShadow: Color(0x0A0F172A),
  hoverShadow: Color(0x14006BE6),
  overlayShadow: Color(0x1F0F172A),
  shadowInk: Color(0xFF0F172A),

  isDark: false,
);

/// 夜风（深色）：分层靠三级底色，不靠投影
const UiPalette kNightPalette = UiPalette(
  background: Color(0xFF050B14),
  surface: Color(0xFF0D1626),
  surfaceLow: Color(0xFF0A1220),
  surfaceHigh: Color(0xFF142038),
  surfaceTop: Color(0xFF1A2843),
  overlay: Color(0xFF142038),

  inkStrong: Color(0xFFF1F5FA),
  ink: Color(0xFFE9EEF5),
  inkSecondary: Color(0xFF94A3B8),
  inkTertiary: Color(0xFF6B7A8F),
  inkDisabled: Color(0xFF3E4C63),

  border: Color(0xFF1E293B),
  borderStrong: Color(0xFF2A3B54),

  primary: Color(0xFF2E96FF),
  primaryDeep: Color(0xFF66B4FF),
  primaryOnSoft: Color(0xFF8FC3FF),
  primarySoft: Color(0xFF172C45),
  primaryTint: Color(0xFF10233A),
  // 白字压在 #2E96FF 上只有 2.9:1，深墨底字才够 4.5:1
  onPrimary: Color(0xFF041426),

  inverseSurface: Color(0xFF1A2843),
  inverseInk: Color(0xFFE9EEF5),

  success: Color(0xFF22C55E),
  warning: Color(0xFFFBBF24),
  danger: Color(0xFFF87171),
  successText: Color(0xFF4ADE80),
  warningText: Color(0xFFFCD34D),
  dangerText: Color(0xFFF87171),
  successSoft: Color(0xFF0F2A1C),
  warningSoft: Color(0xFF2A2313),
  dangerSoft: Color(0xFF2C1618),

  scrim: Color(0x8C030912),
  onImage: Color(0xFFFFFFFF),
  heroStart: Color(0xFF123C6B),
  heroMid: Color(0xFF1E5F9E),
  heroEnd: Color(0xFF0B2A4A),

  cardShadow: Color(0x33000000),
  hoverShadow: Color(0x1A2E96FF),
  overlayShadow: Color(0x66000000),
  shadowInk: Color(0xFF000000),

  isDark: true,
);
