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
    final themeData = Theme.of(context);
    final isMobile = ResponsiveUtils.isMobile(context);
    final loc = S.of(context);

    return Scaffold(
      backgroundColor: themeData.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: themeData.colorScheme.surface,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, size: 22),
          onPressed: () => context.pop(),
        ),
        title: Text(
          loc.login,
          style: TextStyle(
            fontSize: isMobile ? 18.sp : 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        centerTitle: true,
      ),
      body: Center(
        child: SingleChildScrollView(
          padding: EdgeInsets.symmetric(
            horizontal: isMobile ? 24.w : 400,
            vertical: isMobile ? 32.h : 48,
          ),
          child: Form(
            key: _formKey,
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Logo / App Name
                Icon(
                  Icons.article_outlined,
                  size: isMobile ? 64.sp : 64,
                  color: themeData.colorScheme.primary,
                ),
                SizedBox(height: isMobile ? 16.h : 16),
                Text(
                  loc.appName,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: isMobile ? 24.sp : 24,
                    fontWeight: FontWeight.bold,
                    color: themeData.colorScheme.primary,
                  ),
                ),
                SizedBox(height: isMobile ? 48.h : 48),

                // Username
                TextFormField(
                  controller: _usernameController,
                  decoration: InputDecoration(
                    labelText: loc.username,
                    hintText: loc.usernameHint,
                    prefixIcon: const Icon(Icons.person_outline),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
                    ),
                  ),
                  validator: (value) {
                    if (value == null || value.trim().isEmpty) {
                      return loc.usernameHint;
                    }
                    return null;
                  },
                ),
                SizedBox(height: isMobile ? 16.h : 16),

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
                      borderRadius: BorderRadius.circular(isMobile ? 12.r : 12),
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
                SizedBox(height: isMobile ? 32.h : 32),

                // Login Button
                SizedBox(
                  height: isMobile ? 48.h : 48,
                  child: ElevatedButton(
                    onPressed: _handleLogin,
                    style: ElevatedButton.styleFrom(
                      backgroundColor: themeData.colorScheme.primary,
                      foregroundColor: themeData.colorScheme.onPrimary,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(
                          isMobile ? 12.r : 12,
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
    );
  }
}
