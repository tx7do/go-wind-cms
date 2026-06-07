import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

/// 统一的返回按钮
///
/// 智能返回策略：
/// - 如果导航栈中有上一页（canPop），执行 pop
/// - 否则 go 到上一个 tab 页面（/profile 或 /explore）
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
      // 无法 pop 时，返回到上一个 tab 页面
      // 根据当前路由判断应该返回到哪里
      final location = GoRouterState.of(context).uri.toString();
      if (location.startsWith('/profile') || 
          location.startsWith('/settings') ||
          location.startsWith('/my_comments') ||
          location.startsWith('/about')) {
        context.go('/profile');
      } else {
        context.go('/');
      }
    }
  }
}
