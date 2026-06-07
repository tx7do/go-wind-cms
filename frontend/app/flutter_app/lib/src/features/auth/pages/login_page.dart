import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_query/cached_query.dart' show MutationSuccess;

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/features/auth/services/authentication_service.dart';
import 'package:flutter_app/generated/api/models/authentication_service_v1_login_response.dart';

/// 登录页面
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
            backgroundColor: Colors.green,
          ),
        );
        context.go('/');
      } else {
        _showLoginFailed();
      }
    } catch (e) {
      if (!mounted) return;
      _showLoginFailed();
    }
  }

  void _showLoginFailed() {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(S.of(context).loginFailed),
        backgroundColor: Colors.red,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isMobile = ResponsiveUtils.isMobile(context);
    final loc = S.of(context);
    final size = MediaQuery.of(context).size;

    return Scaffold(
      backgroundColor: theme.colorScheme.surface,
      body: SafeArea(
        child: Stack(
          children: [
            // 返回按钮
            Positioned(
              top: isMobile ? 8 : 16,
              left: isMobile ? 8 : 16,
              child: IconButton(
                icon: const Icon(Icons.arrow_back_ios_new, size: 20),
                onPressed: () => context.go('/'),
                tooltip: loc.backToHome,
                style: IconButton.styleFrom(
                  backgroundColor: theme.colorScheme.onSurface.withAlpha(15),
                ),
              ),
            ),

            // 主内容：居中卡片
            Center(
              child: SingleChildScrollView(
                padding: EdgeInsets.symmetric(
                  horizontal: isMobile ? 24.w : 48,
                  vertical: isMobile ? 32.h : 48,
                ),
                child: Container(
                  constraints: BoxConstraints(
                    maxWidth: isMobile ? double.infinity : 420,
                  ),
                  decoration: BoxDecoration(
                    color: theme.scaffoldBackgroundColor,
                    borderRadius: BorderRadius.circular(isMobile ? 20.r : 24),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.black.withAlpha((0.06 * 255).round()),
                        blurRadius: 32,
                        offset: const Offset(0, 8),
                      ),
                    ],
                  ),
                  child: Padding(
                    padding: EdgeInsets.symmetric(
                      horizontal: isMobile ? 28.w : 40,
                      vertical: isMobile ? 32.h : 48,
                    ),
                    child: Form(
                      key: _formKey,
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          // Logo
                          Container(
                            width: isMobile ? 64.w : 72,
                            height: isMobile ? 64.w : 72,
                            decoration: BoxDecoration(
                              color: theme.colorScheme.primaryContainer,
                              borderRadius: BorderRadius.circular(
                                isMobile ? 18.r : 20,
                              ),
                            ),
                            child: Icon(
                              Icons.article_outlined,
                              size: isMobile ? 32.sp : 36,
                              color: theme.colorScheme.primary,
                            ),
                          ),
                          SizedBox(height: isMobile ? 20.h : 24),
                          Text(
                            loc.appName,
                            style: TextStyle(
                              fontSize: isMobile ? 22.sp : 24,
                              fontWeight: FontWeight.bold,
                              color: theme.colorScheme.onSurface,
                            ),
                          ),
                          SizedBox(height: isMobile ? 4.h : 6),
                          Text(
                            loc.loginForMore,
                            style: TextStyle(
                              fontSize: isMobile ? 13.sp : 14,
                              color: theme.colorScheme.onSurface.withAlpha(140),
                            ),
                          ),
                          SizedBox(height: isMobile ? 32.h : 40),

                          // Username
                          TextFormField(
                            controller: _usernameController,
                            decoration: InputDecoration(
                              labelText: loc.username,
                              hintText: loc.usernameHint,
                              prefixIcon: const Icon(Icons.person_outline),
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(
                                  isMobile ? 12.r : 14,
                                ),
                              ),
                            ),
                            validator: (value) {
                              if (value == null || value.trim().isEmpty) {
                                return loc.usernameHint;
                              }
                              return null;
                            },
                          ),
                          SizedBox(height: isMobile ? 16.h : 18),

                          // Password
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
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(
                                  isMobile ? 12.r : 14,
                                ),
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
                          SizedBox(height: isMobile ? 28.h : 32),

                          // Login Button
                          SizedBox(
                            width: double.infinity,
                            height: isMobile ? 48.h : 50,
                            child: FilledButton(
                              onPressed: _handleLogin,
                              style: FilledButton.styleFrom(
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(
                                    isMobile ? 12.r : 14,
                                  ),
                                ),
                              ),
                              child: Text(
                                loc.loginButton,
                                style: TextStyle(
                                  fontSize: isMobile ? 16.sp : 16,
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
            ),
          ],
        ),
      ),
    );
  }
}
