import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_app/src/features/auth/services/authentication_service.dart';

import 'package:flutter_app/generated/l10n.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/core/repositories/user_auth_cache.dart';

/// 个人中心页
class ProfilePage extends StatelessWidget {
  const ProfilePage({super.key});

  @override
  Widget build(BuildContext context) {
    final themeData = Theme.of(context);
    final isMobile = ResponsiveUtils.isMobile(context);

    return Scaffold(
      backgroundColor: themeData.scaffoldBackgroundColor,
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            pinned: true,
            floating: true,
            elevation: 0,
            backgroundColor: themeData.colorScheme.surface,
            surfaceTintColor: Colors.transparent,
            leading: !isMobile
                ? IconButton(
                    icon: const Icon(Icons.arrow_back),
                    onPressed: () => context.pop(),
                  )
                : null,
            title: Text(
              S.of(context).me,
              style: TextStyle(
                fontSize: isMobile ? 22.sp : 22,
                fontWeight: FontWeight.bold,
                color: themeData.colorScheme.onSurface,
              ),
            ),
          ),

          // 用户信息卡片
          SliverToBoxAdapter(
            child: Padding(
              padding: EdgeInsets.all(isMobile ? 16.w : 16),
              child: _UserProfileCard(isMobile: isMobile),
            ),
          ),

          // 功能菜单
          SliverToBoxAdapter(
            child: Padding(
              padding: EdgeInsets.symmetric(horizontal: isMobile ? 16.w : 16),
              child: _MenuSection(isMobile: isMobile),
            ),
          ),

          // 数据统计
          SliverToBoxAdapter(
            child: Padding(
              padding: EdgeInsets.fromLTRB(
                isMobile ? 16.w : 16,
                20,
                isMobile ? 16.w : 16,
                8,
              ),
              child: Text(
                S.of(context).readingStats,
                style: TextStyle(
                  fontSize: isMobile ? 16.sp : 16,
                  fontWeight: FontWeight.w600,
                  color: themeData.colorScheme.onSurface,
                ),
              ),
            ),
          ),
          SliverToBoxAdapter(
            child: Padding(
              padding: EdgeInsets.symmetric(horizontal: isMobile ? 16.w : 16),
              child: _StatsSection(isMobile: isMobile),
            ),
          ),

          SliverToBoxAdapter(child: SizedBox(height: isMobile ? 24.h : 24)),
        ],
      ),
    );
  }
}

class _UserProfileCard extends StatefulWidget {
  final bool isMobile;

  const _UserProfileCard({required this.isMobile});

  @override
  State<_UserProfileCard> createState() => _UserProfileCardState();
}

class _UserProfileCardState extends State<_UserProfileCard> {
  @override
  Widget build(BuildContext context) {
    final themeData = Theme.of(context);
    final loc = S.of(context);
    final isMobile = widget.isMobile;

    return ValueListenableBuilder<bool>(
      valueListenable: GetIt.instance<UserAuthCache>().loginStateNotifier,
      builder: (context, hasLogin, _) {
        return Card(
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(isMobile ? 16.r : 16),
          ),
          color: themeData.colorScheme.surface,
          child: Padding(
            padding: EdgeInsets.all(isMobile ? 20.w : 20),
            child: Row(
              children: [
                CircleAvatar(
                  radius: isMobile ? 32.r : 32,
                  backgroundColor: themeData.colorScheme.primaryContainer,
                  child: Icon(
                    Icons.person,
                    size: isMobile ? 32.sp : 32,
                    color: themeData.colorScheme.onSurface,
                  ),
                ),
                SizedBox(width: isMobile ? 16.w : 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        hasLogin ? loc.appName : loc.guestUser,
                        style: TextStyle(
                          fontSize: isMobile ? 18.sp : 18,
                          fontWeight: FontWeight.bold,
                          color: themeData.colorScheme.onSurface,
                        ),
                      ),
                      SizedBox(height: isMobile ? 4.h : 4),
                      Text(
                        hasLogin ? loc.welcomeBack : loc.loginForMore,
                        style: TextStyle(
                          fontSize: isMobile ? 13.sp : 13,
                          color: themeData.colorScheme.onSurface.withAlpha(140),
                        ),
                      ),
                    ],
                  ),
                ),
                if (hasLogin)
                  OutlinedButton(
                    onPressed: () => _handleLogout(context),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: Colors.red,
                      padding: EdgeInsets.symmetric(
                        horizontal: isMobile ? 20.w : 20,
                        vertical: isMobile ? 8.h : 8,
                      ),
                      shape: RoundedRectangleBorder(
                        borderRadius:
                            BorderRadius.circular(isMobile ? 20.r : 20),
                      ),
                    ),
                    child: Text(
                      loc.logout,
                      style: TextStyle(fontSize: isMobile ? 14.sp : 14),
                    ),
                  )
                else
                  OutlinedButton(
                    onPressed: () => context.go('/login'),
                    style: OutlinedButton.styleFrom(
                      padding: EdgeInsets.symmetric(
                        horizontal: isMobile ? 20.w : 20,
                        vertical: isMobile ? 8.h : 8,
                      ),
                      shape: RoundedRectangleBorder(
                        borderRadius:
                            BorderRadius.circular(isMobile ? 20.r : 20),
                      ),
                    ),
                    child: Text(
                      loc.login,
                      style: TextStyle(fontSize: isMobile ? 14.sp : 14),
                    ),
                  ),
              ],
            ),
          ),
        );
      },
    );
  }

  void _handleLogout(BuildContext context) {
    final loc = S.of(context);
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(loc.logout),
        content: Text(loc.logoutConfirm),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(),
            child: Text(loc.cancel),
          ),
          TextButton(
            onPressed: () async {
              Navigator.of(ctx).pop();
              await AuthenticationService().logoutMutation().mutate(null);
              if (mounted) {
                setState(() {});
                ScaffoldMessenger.of(
                  context,
                ).showSnackBar(SnackBar(content: Text(loc.loginButton)));
              }
            },
            child: Text(loc.confirm, style: const TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }
}

class _MenuSection extends StatelessWidget {
  final bool isMobile;

  const _MenuSection({required this.isMobile});

  @override
  Widget build(BuildContext context) {
    final themeData = Theme.of(context);

    final menuItems = [
      _MenuItem(
        Icons.history_outlined,
        S.of(context).browseHistory,
        S.of(context).viewReadingHistory,
      ),
      _MenuItem(
        Icons.comment_outlined,
        S.of(context).myComments,
        S.of(context).manageComments,
      ),
      _MenuItem(
        Icons.notifications_outlined,
        S.of(context).notifications,
        S.of(context).noNewMessages,
      ),
      _MenuItem(
        Icons.settings_outlined,
        S.of(context).settings,
        S.of(context).themeLanguagePrefs,
      ),
      _MenuItem(
        Icons.info_outline,
        S.of(context).about,
        S.of(context).versionInfo,
      ),
    ];

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(isMobile ? 16.r : 16),
      ),
      color: themeData.colorScheme.surface,
      child: Column(
        children: menuItems.map((item) {
          return ListTile(
            leading: Icon(
              item.icon,
              size: isMobile ? 22.sp : 22,
              color: themeData.colorScheme.onSurface.withAlpha(180),
            ),
            title: Text(
              item.title,
              style: TextStyle(
                fontSize: isMobile ? 15.sp : 15,
                color: themeData.colorScheme.onSurface,
              ),
            ),
            subtitle: Text(
              item.subtitle,
              style: TextStyle(
                fontSize: isMobile ? 12.sp : 12,
                color: themeData.colorScheme.onSurface.withAlpha(100),
              ),
            ),
            trailing: Icon(
              Icons.chevron_right,
              size: isMobile ? 20.sp : 20,
              color: themeData.colorScheme.onSurface.withAlpha(80),
            ),
            contentPadding: EdgeInsets.symmetric(
              horizontal: isMobile ? 16.w : 16,
              vertical: isMobile ? 4.h : 4,
            ),
            onTap: () {
              if (item.icon == Icons.settings_outlined) {
                context.push('/settings');
              } else if (item.icon == Icons.info_outline) {
                context.push('/about');
              } else if (item.icon == Icons.comment_outlined) {
                context.push('/my_comments');
              } else {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(S.of(context).featureNotAvailable),
                    duration: const Duration(seconds: 1),
                  ),
                );
              }
            },
          );
        }).toList(),
      ),
    );
  }
}

class _MenuItem {
  final IconData icon;
  final String title;
  final String subtitle;

  const _MenuItem(this.icon, this.title, this.subtitle);
}

class _StatsSection extends StatelessWidget {
  final bool isMobile;

  const _StatsSection({required this.isMobile});

  @override
  Widget build(BuildContext context) {
    final themeData = Theme.of(context);

    final stats = [
      _StatItem('12', S.of(context).readPosts),
      _StatItem('5', S.of(context).myComments),
      _StatItem('3', S.of(context).bookmarkedPostsLabel),
      _StatItem('2h', S.of(context).readingTime),
    ];

    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(isMobile ? 16.r : 16),
      ),
      color: themeData.colorScheme.surface,
      child: Padding(
        padding: EdgeInsets.symmetric(vertical: isMobile ? 16.h : 16),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: stats.map((stat) {
            return Column(
              children: [
                Text(
                  stat.value,
                  style: TextStyle(
                    fontSize: isMobile ? 20.sp : 20,
                    fontWeight: FontWeight.bold,
                    color: themeData.colorScheme.onSurface,
                  ),
                ),
                SizedBox(height: isMobile ? 4.h : 4),
                Text(
                  stat.label,
                  style: TextStyle(
                    fontSize: isMobile ? 12.sp : 12,
                    color: themeData.colorScheme.onSurface.withAlpha(120),
                  ),
                ),
              ],
            );
          }).toList(),
        ),
      ),
    );
  }
}

class _StatItem {
  final String value;
  final String label;

  const _StatItem(this.value, this.label);
}
