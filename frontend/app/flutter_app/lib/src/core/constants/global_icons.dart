import 'package:flutter/material.dart';

// 全局图标
class GlobalIcons {
  GlobalIcons._(); // Private constructor to prevent instantiation

  /// 向右箭头
  static const Icon navigateNext = Icon(Icons.navigate_next);

  /// 向左箭头
  static const Icon navigateBefore = Icon(Icons.navigate_before);

  /// 错误
  static Icon error({double? size = 50.0, Color? color = Colors.red}) {
    return Icon(Icons.error_outline, color: color, size: size);
  }

  /// 横向三个点图标
  static Icon moreHorizIcon(BuildContext context, {double? size}) {
    return Icon(
      Icons.more_horiz,
      color: Theme.of(context).colorScheme.onPrimary,
      size: size,
    );
  }
}
