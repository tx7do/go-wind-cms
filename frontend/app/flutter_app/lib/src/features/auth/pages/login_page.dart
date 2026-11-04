import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_query/cached_query.dart' show MutationSuccess;

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/layout/width_class.dart';
import 'package:flutter_app/src/core/themes/tokens.dart';
import 'package:flutter_app/src/features/auth/services/authentication_service.dart';
import 'package:flutter_app/generated/api/app/service/v1/index.dart'
    show AuthenticationServiceV1LoginResponse;

/// 登录页（全屏，不带壳层导航）
class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();
  final _authService = AuthenticationService();
  bool _obscurePassword = true;

  @override
  void dispose() {
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    if (!_formKey.currentState!.validate()) return;

    final mutation = _authService.loginMutation();
    try {
      final result = await mutation.mutate(
        LoginParams(
          username: _usernameController.text.trim(),
          password: _passwordController.text,
        ),
      );

      if (!mounted) return;

      if (result is MutationSuccess<AuthenticationServiceV1LoginResponse?> &&
          result.data != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(S.of(context).loginSuccess),
            backgroundColor: Theme.of(context).palette.success,
          ),
        );
        context.go('/');
      } else {
        _showLoginFailed();
      }
    } catch (e) {
      // 不吞错：失败原因留痕（服务端 401 之外的意外错误要能从日志追）
      debugPrint('[LoginPage] login failed: $e');
      if (!mounted) return;
      _showLoginFailed();
    }
  }

  void _showLoginFailed() {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(S.of(context).loginFailed),
        backgroundColor: Theme.of(context).palette.danger,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final p = theme.palette;
    final space = theme.space;
    final compact = context.widthClass.isCompact;
    final loc = S.of(context);

    return Scaffold(
      backgroundColor: p.background,
      body: SafeArea(
        child: Stack(
          children: [
            // 返回首页
            Positioned(
              top: compact ? 8 : 16,
              left: compact ? 8 : 16,
              child: IconButton(
                icon: const Icon(Icons.arrow_back_ios_new, size: 20),
                onPressed: () => context.go('/'),
                tooltip: loc.backToHome,
                style: IconButton.styleFrom(
                  backgroundColor: p.primaryTint,
                ),
              ),
            ),

            // 居中卡片
            Center(
              child: SingleChildScrollView(
                padding: EdgeInsets.symmetric(
                  horizontal: compact ? space.xxl : 48,
                  vertical: compact ? space.xxl : 48,
                ),
                child: Container(
                  constraints: BoxConstraints(
                    maxWidth: compact ? double.infinity : 420,
                  ),
                  decoration: BoxDecoration(
                    color: p.surface,
                    borderRadius: theme.radius.lgAll,
                    border: Border.all(
                      color: p.border,
                      width: theme.borderWidth.thin,
                    ),
                    boxShadow: theme.shadow.overlay,
                  ),
                  padding: EdgeInsets.symmetric(
                    horizontal: compact ? space.xxl : 40,
                    vertical: compact ? space.xxl + 8 : space.block,
                  ),
                  child: Form(
                    key: _formKey,
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        // 品牌位
                        Container(
                          width: 64,
                          height: 64,
                          decoration: BoxDecoration(
                            color: p.primarySoft,
                            borderRadius: theme.radius.lgAll,
                          ),
                          child: Icon(
                            Icons.article_outlined,
                            size: 32,
                            color: p.primary,
                          ),
                        ),
                        SizedBox(height: space.xl),
                        Text(
                          loc.appName,
                          style: theme.textTheme.headlineMedium,
                        ),
                        const SizedBox(height: 6),
                        Text(
                          loc.loginForMore,
                          style: theme.textTheme.bodySmall,
                        ),
                        SizedBox(height: space.xxl + space.lg),

                        // 用户名
                        TextFormField(
                          controller: _usernameController,
                          decoration: InputDecoration(
                            labelText: loc.username,
                            hintText: loc.usernameHint,
                            prefixIcon: const Icon(Icons.person_outline),
                          ),
                          validator: (value) {
                            if (value == null || value.trim().isEmpty) {
                              return loc.usernameHint;
                            }
                            return null;
                          },
                        ),
                        SizedBox(height: space.lg + 2),

                        // 密码
                        TextFormField(
                          controller: _passwordController,
                          obscureText: _obscurePassword,
                          decoration: InputDecoration(
                            labelText: loc.password,
                            hintText: loc.passwordHint,
                            prefixIcon: const Icon(Icons.lock_outline),
                            suffixIcon: IconButton(
                              icon: Icon(
                                _obscurePassword
                                    ? Icons.visibility_off_outlined
                                    : Icons.visibility_outlined,
                              ),
                              onPressed: () {
                                setState(() {
                                  _obscurePassword = !_obscurePassword;
                                });
                              },
                            ),
                          ),
                          validator: (value) {
                            if (value == null || value.isEmpty) {
                              return loc.passwordHint;
                            }
                            return null;
                          },
                          onFieldSubmitted: (_) => _handleLogin(),
                        ),
                        SizedBox(height: space.xxl + space.sm),

                        // 登录按钮
                        SizedBox(
                          width: double.infinity,
                          height: 48,
                          child: FilledButton(
                            onPressed: _handleLogin,
                            child: Text(
                              loc.loginButton,
                              style: theme.textTheme.labelLarge?.copyWith(
                                fontSize: 16,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
