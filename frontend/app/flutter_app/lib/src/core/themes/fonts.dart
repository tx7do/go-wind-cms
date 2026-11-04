import 'package:flutter/foundation.dart'
    show TargetPlatform, defaultTargetPlatform;

/// 字体族解析。
///
/// 系统优先：苹方 / 微软雅黑 / 系统思源 是中文用户本机的默认观感，
/// 内置的 `Noto Sans SC`（VF 变量字体）退到兜底位——它在 Windows 上会盖掉
/// 雅黑、在 macOS 上会盖掉苹方，显「扁钝」，只在不该落空的平台（Android、
/// 缺字体的 Linux、拿不到本地字体的 Web）才启用。
abstract final class AppFont {
  /// pubspec.yaml 里注册的内置中文字体族名
  static const String bundledFamily = 'Noto Sans SC';

  /// 主题主字体族（按平台取系统字体）
  static String get primaryFamily {
    // defaultTargetPlatform 在 Web 上按浏览器 UA 给出，native 上按真实系统给出
    return switch (defaultTargetPlatform) {
      TargetPlatform.iOS || TargetPlatform.macOS => 'PingFang SC',
      TargetPlatform.windows => 'Microsoft YaHei',
      TargetPlatform.linux => 'Noto Sans CJK SC',
      // Android 系统 CJK 与内置字体同源，直接用内置档保证跨 ROM 一致
      TargetPlatform.android || TargetPlatform.fuchsia => bundledFamily,
    };
  }

  /// 逐字回退链：先补齐其他中文系字体，再落内置思源，最后 emoji。
  static List<String> get fallback {
    // Web 端 CanvasKit 拿不到本地字体时，内置思源保证中文不 tofu；
    // emoji 资产（NotoColorEmoji）已移除，Web 端 emoji 由浏览器字体兜底。
    return const <String>[
      'PingFang SC',
      'Hiragino Sans GB',
      'Microsoft YaHei',
      'Source Han Sans SC',
      'Noto Sans CJK SC',
      bundledFamily,
      'Apple Color Emoji',
      'Segoe UI Emoji',
      'Noto Color Emoji',
    ];
  }

  /// 等宽（代码块、日志、序列号）。没有注册等宽资产，靠逐字回退落到系统 mono。
  static const List<String> monoFallback = <String>[
    'SFMono-Regular',
    'Menlo',
    'Consolas',
    'Liberation Mono',
    'Courier New',
  ];
}
