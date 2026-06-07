import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

/// 统一的返回按钮
///
/// 智能返回策略：
/// - 如果导航栈中有上一页（canPop），执行 pop
/// - 否则 go 到首页 '/'
class AppBackButton extends StatelessWidget {
  final Color? color;
  final double size;

  const AppBackButton({super.key, this.color, this.size = 22});

  @override
  Widget build(BuildContext context) {
    return IconButton(
      icon: Icon(Icons.arrow_back, size: size, color: color),
      onPressed: () => _goBack(context),
    );
  }

  void _goBack(BuildContext context) {
    if (context.canPop()) {
      context.pop();
    } else {
      context.go('/');
    }
  }
}
