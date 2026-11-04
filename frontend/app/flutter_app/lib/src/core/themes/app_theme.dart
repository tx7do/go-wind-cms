import 'package:flutter/cupertino.dart' show CupertinoPageTransitionsBuilder;
import 'package:flutter/foundation.dart' show defaultTargetPlatform;
import 'package:flutter/gestures.dart' show PointerDeviceKind;
import 'package:flutter/material.dart';
import 'package:flutter/services.dart' show SystemUiOverlayStyle;

import 'fonts.dart';
import 'palette.dart';
import 'tokens.dart';
import 'typography.dart';

/// 亮/暗共用一套组件规则，只有 [UiPalette] 分档。
///
/// 这样做的直接原因：历史上 light/dark 两个文件各写各的（例如 NavigationBar
/// 只在暗色里设了背景色），两态越漂越远。观感差异一律以 token 差异表达。
ThemeData buildAppTheme({required UiPalette p, required Brightness brightness}) {
  final colorScheme = _buildColorScheme(p, brightness);
  // ThemeData 构造里 `defaultTextTheme.merge(textTheme)` 是逐角色替换：
  // 自定义 TextTheme 覆盖了全部 15 个角色，故 fontFamily 必须在自己的
  // TextTheme 上 apply，只传 ThemeData 的 fontFamily 参数会被整体丢掉。
  final textTheme = AppTypography.build(p).apply(
    fontFamily: AppFont.primaryFamily,
    fontFamilyFallback: AppFont.fallback,
  );
  final radius = RadiusTokens.light;
  final space = SpaceTokens.light;
  final border = BorderTokens.light;
  final motion = MotionTokens.light;
  final shadow = ShadowTokens.fromPalette(p);

  return ThemeData(
    // 不钉 platform：让它跟随真实系统，桌面用 mouse 拖拽、Apple 用回弹，
    // 都由 AppScrollBehavior 统一裁决（历史代码钉死 iOS 导致全站回弹）。
    colorScheme: colorScheme,
    scaffoldBackgroundColor: p.background,
    canvasColor: p.background,
    fontFamily: AppFont.primaryFamily,
    fontFamilyFallback: AppFont.fallback,
    textTheme: textTheme,
    // primaryTextTheme 不再重染颜色：AppTypography 已按 palette 逐档上色，
    // 用 apply(bodyColor/displayColor) 会把 bodySmall/labelSmall 的灰档冲掉。
    primaryTextTheme: textTheme,
    iconTheme: IconThemeData(color: p.inkSecondary, size: 20),
    primaryIconTheme: IconThemeData(color: p.onPrimary),
    splashFactory: InkRipple.splashFactory,
    splashColor: p.primary.withAlpha(0x14),
    highlightColor: p.primary.withAlpha(0x0A),
    visualDensity: VisualDensity.standard,
    materialTapTargetSize: MaterialTapTargetSize.padded,
    pageTransitionsTheme: const PageTransitionsTheme(
      builders: <TargetPlatform, PageTransitionsBuilder>{
        TargetPlatform.android: FadeForwardsPageTransitionsBuilder(),
        TargetPlatform.fuchsia: FadeForwardsPageTransitionsBuilder(),
        TargetPlatform.iOS: CupertinoPageTransitionsBuilder(),
        TargetPlatform.linux: FadeForwardsPageTransitionsBuilder(),
        TargetPlatform.macOS: CupertinoPageTransitionsBuilder(),
        TargetPlatform.windows: FadeForwardsPageTransitionsBuilder(),
      },
    ),
    // scrollBehavior 在 Flutter 3.47 已从 ThemeData 移到 MaterialApp，
    // 由 app.dart 传入。
    extensions: <ThemeExtension<dynamic>>[
      AppThemeExtension(
        palette: p,
        radius: radius,
        space: space,
        border: border,
        motion: motion,
        shadow: shadow,
      ),
    ],

    appBarTheme: AppBarTheme(
      backgroundColor: p.surface,
      foregroundColor: p.inkStrong,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
      scrolledUnderElevation: 0,
      shadowColor: Colors.transparent,
      centerTitle: false,
      titleSpacing: space.lg,
      // Border 是 ShapeBorder：只画底边那条 0.5px 细线，不要阴影
      shape: Border(bottom: BorderSide(color: p.border, width: border.hairline)),
      titleTextStyle: textTheme.headlineSmall,
      iconTheme: IconThemeData(color: p.ink, size: 22),
      actionsIconTheme: IconThemeData(color: p.inkSecondary, size: 22),
      systemOverlayStyle: p.isDark
          ? SystemUiOverlayStyle.light
          : SystemUiOverlayStyle.dark,
    ),

    cardTheme: CardThemeData(
      color: p.surface,
      elevation: 0,
      margin: EdgeInsets.zero,
      clipBehavior: Clip.antiAlias,
      shape: RoundedRectangleBorder(
        borderRadius: radius.mdAll,
        side: BorderSide(color: p.border, width: border.thin),
      ),
    ),

    dividerTheme: DividerThemeData(
      color: p.border,
      thickness: border.hairline,
      space: border.hairline,
    ),

    listTileTheme: ListTileThemeData(
      iconColor: p.inkSecondary,
      textColor: p.ink,
      titleTextStyle: textTheme.bodyMedium,
      subtitleTextStyle: textTheme.bodySmall,
      contentPadding: EdgeInsets.symmetric(horizontal: space.lg),
      shape: RoundedRectangleBorder(borderRadius: radius.smAll),
    ),

    filledButtonTheme: FilledButtonThemeData(
      style: ButtonStyle(
        backgroundColor: WidgetStateColor.resolveWith(
          (s) => s.contains(WidgetState.disabled) ? p.surfaceTop : p.primary,
        ),
        foregroundColor: WidgetStateColor.resolveWith(
          (s) => s.contains(WidgetState.disabled) ? p.inkDisabled : p.onPrimary,
        ),
        textStyle: WidgetStatePropertyAll(textTheme.labelLarge?.copyWith(
          fontWeight: FontWeight.w600,
        )),
        minimumSize: WidgetStatePropertyAll(Size(0, space.controlHeight)),
        padding: WidgetStatePropertyAll(
          EdgeInsets.symmetric(horizontal: space.xl),
        ),
        shape: WidgetStatePropertyAll(
          RoundedRectangleBorder(borderRadius: radius.smAll),
        ),
        elevation: const WidgetStatePropertyAll(0),
      ),
    ),
    // 页面里既有的 ElevatedButton 一律并入 filled 规格（同一观感，两种写法）
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ButtonStyle(
        backgroundColor: WidgetStateColor.resolveWith(
          (s) => s.contains(WidgetState.disabled) ? p.surfaceTop : p.primary,
        ),
        foregroundColor: WidgetStateColor.resolveWith(
          (s) => s.contains(WidgetState.disabled) ? p.inkDisabled : p.onPrimary,
        ),
        textStyle: WidgetStatePropertyAll(textTheme.labelLarge?.copyWith(
          fontWeight: FontWeight.w600,
        )),
        minimumSize: WidgetStatePropertyAll(Size(0, space.controlHeight)),
        padding: WidgetStatePropertyAll(
          EdgeInsets.symmetric(horizontal: space.xl),
        ),
        shape: WidgetStatePropertyAll(
          RoundedRectangleBorder(borderRadius: radius.smAll),
        ),
        elevation: const WidgetStatePropertyAll(0),
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: ButtonStyle(
        backgroundColor: WidgetStateColor.resolveWith(
          (s) => s.contains(WidgetState.hovered) ? p.primaryTint : p.surface,
        ),
        foregroundColor: WidgetStateColor.resolveWith(
          (s) => s.contains(WidgetState.disabled) ? p.inkDisabled : p.ink,
        ),
        side: WidgetStateBorderSide.resolveWith(
          (s) => BorderSide(
            width: border.thin,
            color: s.contains(WidgetState.hovered)
                ? p.primaryBorder
                : p.borderStrong,
          ),
        ),
        textStyle: WidgetStatePropertyAll(textTheme.labelLarge),
        minimumSize: WidgetStatePropertyAll(Size(0, space.controlHeight)),
        padding: WidgetStatePropertyAll(
          EdgeInsets.symmetric(horizontal: space.xl),
        ),
        shape: WidgetStatePropertyAll(
          RoundedRectangleBorder(borderRadius: radius.smAll),
        ),
      ),
    ),
    textButtonTheme: TextButtonThemeData(
      style: ButtonStyle(
        foregroundColor: WidgetStateColor.resolveWith(
          (s) => s.contains(WidgetState.disabled)
              ? p.inkDisabled
              : (s.contains(WidgetState.hovered) ? p.primaryDeep : p.primary),
        ),
        textStyle: WidgetStatePropertyAll(textTheme.labelLarge),
        minimumSize: WidgetStatePropertyAll(Size(0, space.controlHeight - 8)),
        padding: WidgetStatePropertyAll(
          EdgeInsets.symmetric(horizontal: space.md),
        ),
        shape: WidgetStatePropertyAll(
          RoundedRectangleBorder(borderRadius: radius.smAll),
        ),
        overlayColor: WidgetStatePropertyAll(p.primary.withAlpha(0x0F)),
      ),
    ),
    iconButtonTheme: IconButtonThemeData(
      style: ButtonStyle(
        iconColor: WidgetStateColor.resolveWith(
          (s) => s.contains(WidgetState.disabled)
              ? p.inkDisabled
              : (s.contains(WidgetState.hovered) ? p.ink : p.inkSecondary),
        ),
        iconSize: const WidgetStatePropertyAll(20),
        minimumSize: WidgetStatePropertyAll(Size.square(space.xxl)),
        shape: WidgetStatePropertyAll(
          RoundedRectangleBorder(borderRadius: radius.smAll),
        ),
        overlayColor: WidgetStatePropertyAll(p.primary.withAlpha(0x0F)),
      ),
    ),
    floatingActionButtonTheme: FloatingActionButtonThemeData(
      backgroundColor: p.primary,
      foregroundColor: p.onPrimary,
      elevation: 0,
      highlightElevation: 0,
      shape: RoundedRectangleBorder(borderRadius: radius.mdAll),
    ),

    // 3.47 起 ThemeData 收 Object?，数据类是 InputDecorationThemeData
    inputDecorationTheme: InputDecorationThemeData(
      filled: true,
      fillColor: p.surfaceLow,
      contentPadding: EdgeInsets.symmetric(
        horizontal: space.md,
        vertical: space.md,
      ),
      hintStyle: textTheme.bodyMedium?.copyWith(color: p.inkTertiary),
      labelStyle: textTheme.bodyMedium?.copyWith(color: p.inkSecondary),
      floatingLabelStyle: TextStyle(color: p.primary),
      errorStyle: textTheme.labelSmall?.copyWith(color: p.dangerText),
      helperStyle: textTheme.labelSmall?.copyWith(color: p.inkTertiary),
      iconColor: p.inkTertiary,
      prefixIconColor: p.inkTertiary,
      suffixIconColor: p.inkTertiary,
      border: _inputBorder(p.borderStrong, border.thin, radius),
      enabledBorder: _inputBorder(p.borderStrong, border.thin, radius),
      focusedBorder: _inputBorder(p.primary, border.emphasis, radius),
      errorBorder: _inputBorder(p.danger, border.thin, radius),
      focusedErrorBorder: _inputBorder(p.danger, border.emphasis, radius),
      disabledBorder: _inputBorder(p.border, border.thin, radius),
    ),

    chipTheme: ChipThemeData(
      backgroundColor: p.primarySoft,
      selectedColor: p.primary,
      disabledColor: p.surfaceTop,
      secondarySelectedColor: p.primary,
      labelStyle: textTheme.labelMedium?.copyWith(
        color: p.primaryOnSoft,
        fontWeight: FontWeight.w500,
      ),
      secondaryLabelStyle: textTheme.labelMedium?.copyWith(color: p.onPrimary),
      side: BorderSide.none,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.all(Radius.circular(999)),
      ),
      padding: EdgeInsets.symmetric(horizontal: space.sm, vertical: space.xxs),
      showCheckmark: false,
    ),

    tabBarTheme: TabBarThemeData(
      labelColor: p.inkStrong,
      unselectedLabelColor: p.inkSecondary,
      labelStyle: textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600),
      unselectedLabelStyle: textTheme.titleSmall,
      indicatorSize: TabBarIndicatorSize.tab,
      indicator: UnderlineTabIndicator(
        borderSide: BorderSide(color: p.primary, width: 2),
        insets: EdgeInsets.zero,
      ),
      dividerColor: Colors.transparent,
      tabAlignment: TabAlignment.start,
      overlayColor: WidgetStatePropertyAll(p.primary.withAlpha(0x0F)),
      splashBorderRadius: radius.smAll,
    ),

    navigationBarTheme: NavigationBarThemeData(
      backgroundColor: p.surface,
      indicatorColor: p.primaryTint,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
      height: 62,
      labelBehavior: NavigationDestinationLabelBehavior.alwaysShow,
      labelTextStyle: WidgetStateProperty.resolveWith(
        (s) => (s.contains(WidgetState.selected)
                ? textTheme.labelSmall?.copyWith(
                    color: p.primary, fontWeight: FontWeight.w600)
                : textTheme.labelSmall?.copyWith(color: p.inkTertiary))!
            .copyWith(fontSize: 11),
      ),
      iconTheme: WidgetStateProperty.resolveWith(
        (s) => IconThemeData(
          size: 22,
          color: s.contains(WidgetState.selected) ? p.primary : p.inkSecondary,
        ),
      ),
      overlayColor: WidgetStatePropertyAll(p.primary.withAlpha(0x0F)),
    ),
    bottomNavigationBarTheme: BottomNavigationBarThemeData(
      backgroundColor: p.surface,
      selectedItemColor: p.primary,
      unselectedItemColor: p.inkTertiary,
      selectedLabelStyle: textTheme.labelSmall?.copyWith(
        fontWeight: FontWeight.w600,
        fontSize: 11,
      ),
      unselectedLabelStyle: textTheme.labelSmall?.copyWith(fontSize: 11),
      elevation: 0,
      type: BottomNavigationBarType.fixed,
      showSelectedLabels: true,
      showUnselectedLabels: true,
    ),

    searchBarTheme: SearchBarThemeData(
      backgroundColor: WidgetStatePropertyAll(p.surfaceLow),
      surfaceTintColor: WidgetStatePropertyAll(Colors.transparent),
      elevation: WidgetStatePropertyAll(0),
      shadowColor: WidgetStatePropertyAll(Colors.transparent),
      padding: WidgetStatePropertyAll(
        EdgeInsets.symmetric(horizontal: space.md, vertical: space.sm),
      ),
      textStyle: WidgetStatePropertyAll(textTheme.bodyMedium),
      hintStyle: WidgetStatePropertyAll(
        textTheme.bodyMedium?.copyWith(color: p.inkTertiary),
      ),
      shape: WidgetStatePropertyAll(
        RoundedRectangleBorder(borderRadius: radius.smAll),
      ),
      side: WidgetStatePropertyAll(
        BorderSide(color: p.borderStrong, width: border.thin),
      ),
    ),

    dialogTheme: DialogThemeData(
      backgroundColor: p.overlay,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
      shadowColor: shadow.overlay.isNotEmpty
          ? p.overlayShadow
          : Colors.transparent,
      shape: RoundedRectangleBorder(
        borderRadius: radius.lgAll,
        side: BorderSide(color: p.border, width: border.thin),
      ),
      titleTextStyle: textTheme.headlineSmall,
      contentTextStyle: textTheme.bodyMedium,
    ),
    bottomSheetTheme: BottomSheetThemeData(
      backgroundColor: p.surface,
      surfaceTintColor: Colors.transparent,
      modalBarrierColor: p.scrim,
      elevation: 0,
      showDragHandle: true,
      dragHandleColor: p.inkDisabled,
      shape: RoundedRectangleBorder(borderRadius: radius.lgTop),
    ),
    popupMenuTheme: PopupMenuThemeData(
      color: p.overlay,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
      shadowColor: p.overlayShadow,
      iconColor: p.inkSecondary,
      textStyle: textTheme.bodyMedium,
      shape: RoundedRectangleBorder(
        borderRadius: radius.mdAll,
        side: BorderSide(color: p.border, width: border.thin),
      ),
    ),
    snackBarTheme: SnackBarThemeData(
      backgroundColor: p.inverseSurface,
      contentTextStyle: textTheme.bodySmall?.copyWith(color: p.inverseInk),
      actionTextColor: p.primary,
      elevation: 0,
      behavior: SnackBarBehavior.floating,
      shape: RoundedRectangleBorder(borderRadius: radius.mdAll),
    ),
    tooltipTheme: TooltipThemeData(
      decoration: BoxDecoration(
        color: p.inverseSurface,
        borderRadius: radius.xsAll,
      ),
      textStyle: textTheme.labelSmall?.copyWith(color: p.inverseInk),
      padding: EdgeInsets.symmetric(horizontal: space.sm, vertical: space.xs),
      waitDuration: motion.base,
    ),
    badgeTheme: BadgeThemeData(
      backgroundColor: p.danger,
      textColor: p.onPrimary,
      smallSize: 7,
      largeSize: 16,
      textStyle: textTheme.labelSmall?.copyWith(
        color: p.onPrimary,
        fontWeight: FontWeight.w600,
      ),
    ),
    progressIndicatorTheme: ProgressIndicatorThemeData(
      color: p.primary,
      linearTrackColor: p.surfaceTop,
      circularTrackColor: p.surfaceTop,
      linearMinHeight: 3,
    ),
    checkboxTheme: CheckboxThemeData(
      fillColor: WidgetStateColor.resolveWith(
        (s) => s.contains(WidgetState.selected) ? p.primary : Colors.transparent,
      ),
      checkColor: WidgetStatePropertyAll(p.onPrimary),
      side: BorderSide(color: p.borderStrong, width: border.emphasis),
      shape: RoundedRectangleBorder(borderRadius: radius.xsAll),
      splashRadius: 0,
    ),
    radioTheme: RadioThemeData(
      fillColor: WidgetStateColor.resolveWith(
        (s) => s.contains(WidgetState.selected) ? p.primary : p.borderStrong,
      ),
    ),
    switchTheme: SwitchThemeData(
      thumbColor: WidgetStateColor.resolveWith(
        (s) => s.contains(WidgetState.selected) ? p.primary : p.surface,
      ),
      trackColor: WidgetStateColor.resolveWith(
        (s) => s.contains(WidgetState.selected)
            ? p.primarySoft
            : p.surfaceTop,
      ),
      trackOutlineColor: WidgetStateColor.resolveWith(
        (s) => s.contains(WidgetState.selected) ? p.primary : p.borderStrong,
      ),
    ),
    sliderTheme: SliderThemeData(
      activeTrackColor: p.primary,
      inactiveTrackColor: p.surfaceTop,
      thumbColor: p.primary,
      overlayColor: p.primary.withAlpha(0x1A),
      trackHeight: 3,
    ),
    expansionTileTheme: ExpansionTileThemeData(
      iconColor: p.inkSecondary,
      collapsedIconColor: p.inkTertiary,
      textColor: p.inkStrong,
      childrenPadding: EdgeInsets.only(bottom: space.md),
    ),
    textSelectionTheme: TextSelectionThemeData(
      cursorColor: p.primary,
      selectionColor: p.primarySoft,
      selectionHandleColor: p.primary,
    ),
    scrollbarTheme: ScrollbarThemeData(
      thumbColor: WidgetStateColor.resolveWith(
        (s) => s.contains(WidgetState.hovered)
            ? p.inkTertiary.withAlpha(0xB3)
            : p.inkTertiary.withAlpha(0x66),
      ),
      radius: const Radius.circular(4),
      thickness: WidgetStateProperty.resolveWith(
        (s) => s.contains(WidgetState.hovered) ? 8 : 6,
      ),
      mainAxisMargin: 6,
      crossAxisMargin: 2,
      thumbVisibility: const WidgetStatePropertyAll(false),
    ),
  );
}

OutlineInputBorder _inputBorder(Color color, double width, RadiusTokens r) {
  return OutlineInputBorder(
    borderRadius: r.smAll,
    borderSide: BorderSide(color: color, width: width),
  );
}

ColorScheme _buildColorScheme(UiPalette p, Brightness brightness) {
  return ColorScheme(
    brightness: brightness,
    primary: p.primary,
    onPrimary: p.onPrimary,
    primaryContainer: p.primarySoft,
    onPrimaryContainer: p.primaryOnSoft,
    secondary: p.secondary,
    onSecondary: p.inkOnSolid,
    secondaryContainer: p.surfaceTop,
    onSecondaryContainer: p.inkStrong,
    // 第三色相不参与装饰，只借用主色的深浅两档（doc §1 铁律 1）
    tertiary: p.primary,
    onTertiary: p.onPrimary,
    tertiaryContainer: p.primaryTint,
    onTertiaryContainer: p.primaryOnSoft,
    error: p.danger,
    onError: p.inkOnSolid,
    errorContainer: p.dangerSoft,
    onErrorContainer: p.dangerText,
    surface: p.surface,
    onSurface: p.ink,
    onSurfaceVariant: p.inkSecondary,
    surfaceContainerLowest: p.surfaceLow,
    surfaceContainerLow: p.surface,
    surfaceContainer: p.isDark ? p.surfaceHigh : p.background,
    surfaceContainerHigh: p.surfaceHigh,
    surfaceContainerHighest: p.surfaceTop,
    surfaceBright: p.isDark ? p.surfaceHigh : p.surface,
    surfaceDim: p.isDark ? p.background : p.surfaceTop,
    outline: p.inkTertiary,
    outlineVariant: p.border,
    inverseSurface: p.inverseSurface,
    onInverseSurface: p.inverseInk,
    inversePrimary: p.primaryDeep,
    scrim: p.scrim,
    shadow: p.shadowInk,
    // elevation tint 关掉：分层由 surface 阶梯负责，否则深底会泛出 seed 味
    surfaceTint: Colors.transparent,
  );
}

/// 滚动观感：Apple 回弹、其余平台钳位；关掉 Material 的溢出光晕。
class AppScrollBehavior extends MaterialScrollBehavior {
  const AppScrollBehavior();

  @override
  Set<PointerDeviceKind> get dragDevices => const <PointerDeviceKind>{
    PointerDeviceKind.touch,
    PointerDeviceKind.mouse,
    PointerDeviceKind.trackpad,
    PointerDeviceKind.stylus,
  };

  @override
  ScrollPhysics getScrollPhysics(BuildContext context) {
    return switch (defaultTargetPlatform) {
      TargetPlatform.iOS || TargetPlatform.macOS =>
        const BouncingScrollPhysics(parent: AlwaysScrollableScrollPhysics()),
      _ => const ClampingScrollPhysics(parent: AlwaysScrollableScrollPhysics()),
    };
  }

  @override
  Widget buildOverscrollIndicator(
    BuildContext context,
    Widget child,
    ScrollableDetails details,
  ) => child;
}
