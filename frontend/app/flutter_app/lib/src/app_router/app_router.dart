import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/src/core/constants/index.dart' as constants;

import 'package:flutter_app/src/core/widgets/not_found_page.dart';
import 'package:flutter_app/src/app_router/route_names.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_page.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/post_detail_page.dart';
import 'package:flutter_app/src/features/cms/pages/post_list/post_list_page.dart';
import 'package:flutter_app/src/features/cms/pages/tag_feed/tag_feed_page.dart';
import 'package:flutter_app/src/features/cms/pages/tag_list/tag_list_page.dart';
import 'package:flutter_app/src/features/cms/pages/category_list/category_list_page.dart';
import 'package:flutter_app/src/features/cms/pages/search/search_page.dart';
import 'package:flutter_app/src/features/auth/pages/login_page.dart';
import 'package:flutter_app/src/core/utils/responsive_utils.dart';
import 'package:flutter_app/src/features/cms/pages/explore/explore_page.dart';
import 'package:flutter_app/src/features/cms/pages/profile/profile_page.dart';
import 'package:flutter_app/src/features/cms/pages/settings/settings_page.dart';
import 'package:flutter_app/src/features/cms/pages/my_comments/my_comments_page.dart';
import 'package:flutter_app/src/features/cms/pages/about/about_page.dart';
import 'package:flutter_app/src/features/cms/pages/legal/legal_page.dart';

/// CMS 应用路由
class AppRouter {
  static const initial = constants.AppRoutePath.initial;

  static final router = GoRouter(
    initialLocation: initial,
    redirect: _guard,
    errorBuilder: (context, state) => const NotFoundPage(),
    routes: [
      // 主页
      GoRoute(
        path: constants.AppRoutePath.initial,
        name: RouteNames.home,
        builder: (context, state) {
          return const HomePage();
        },
        routes: [
          // 文章详情
          GoRoute(
            name: RouteNames.postDetail,
            path: 'post/:id',
            builder: (context, state) {
              final postId =
                  int.tryParse(state.pathParameters['id'] ?? '0') ?? 0;
              return PostDetailPage(postId: postId);
            },
          ),
          // 标签文章列表
          GoRoute(
            name: RouteNames.tagFeed,
            path: 'tag/:id',
            builder: (context, state) {
              final tagId =
                  int.tryParse(state.pathParameters['id'] ?? '0') ?? 0;
              return TagFeedPage(tagId: tagId);
            },
          ),
          // 搜索
          GoRoute(
            name: RouteNames.search,
            path: 'search',
            builder: (context, state) {
              return const SearchPage();
            },
          ),
          // 页面详情（复用 PostDetailPage）
          GoRoute(
            name: RouteNames.pageDetail,
            path: 'page/:id',
            builder: (context, state) {
              final pageId =
                  int.tryParse(state.pathParameters['id'] ?? '0') ?? 0;
              return PostDetailPage(postId: pageId);
            },
          ),
        ],
      ),
      // 登录页
      GoRoute(
        name: RouteNames.login,
        path: constants.AppRoutePath.login,
        builder: (context, state) {
          return const LoginPage();
        },
      ),
      // 个人中心
      GoRoute(
        name: RouteNames.profile,
        path: constants.AppRoutePath.profile,
        builder: (context, state) {
          if (ResponsiveUtils.isMobile(context)) {
            return const HomePage(initialRoute: '/profile');
          }
          return const ProfilePage();
        },
      ),
      // 发现页
      GoRoute(
        name: RouteNames.explore,
        path: constants.AppRoutePath.explore,
        builder: (context, state) {
          if (ResponsiveUtils.isMobile(context)) {
            return const HomePage(initialRoute: '/explore');
          }
          return const ExplorePage();
        },
      ),
      // 关于我们
      GoRoute(
        name: RouteNames.about,
        path: constants.AppRoutePath.about,
        builder: (context, state) {
          return const AboutPage();
        },
      ),
      // 设置
      GoRoute(
        name: RouteNames.settings,
        path: constants.AppRoutePath.settings,
        builder: (context, state) {
          return const SettingsPage();
        },
      ),
      // 我的评论
      GoRoute(
        name: RouteNames.myComments,
        path: constants.AppRoutePath.myComments,
        builder: (context, state) {
          return const MyCommentsPage();
        },
      ),
      // 联系我们
      GoRoute(
        name: RouteNames.contact,
        path: constants.AppRoutePath.contact,
        builder: (context, state) {
          return const LegalPage(type: LegalPageType.contact);
        },
      ),
      // 免责条款
      GoRoute(
        name: RouteNames.disclaimer,
        path: constants.AppRoutePath.disclaimer,
        builder: (context, state) {
          return const LegalPage(type: LegalPageType.disclaimer);
        },
      ),
      // 隐私协议
      GoRoute(
        name: RouteNames.privacy,
        path: constants.AppRoutePath.privacy,
        builder: (context, state) {
          return const LegalPage(type: LegalPageType.privacy);
        },
      ),
      // 服务条款
      GoRoute(
        name: RouteNames.terms,
        path: constants.AppRoutePath.terms,
        builder: (context, state) {
          return const LegalPage(type: LegalPageType.terms);
        },
      ),
      // 文章列表（/post）
      GoRoute(
        name: RouteNames.postList,
        path: '/post',
        builder: (context, state) {
          final categoryId = int.tryParse(
            state.uri.queryParameters['categoryId'] ?? '',
          );
          final tagId = int.tryParse(state.uri.queryParameters['tagId'] ?? '');
          return PostListPage(categoryId: categoryId, tagId: tagId);
        },
      ),
      // 分类列表（/category）
      GoRoute(
        name: RouteNames.categoryList,
        path: '/category',
        builder: (context, state) {
          return const CategoryListPage();
        },
      ),
      // 标签列表（/tag）
      GoRoute(
        name: RouteNames.tagList,
        path: '/tag',
        builder: (context, state) {
          return const TagListPage();
        },
      ),
    ],
  );

  static FutureOr<String?> _guard(BuildContext context, GoRouterState state) {
    // CMS 内容展示端暂时不需要登录验证，直接放行
    return null;
  }
}
