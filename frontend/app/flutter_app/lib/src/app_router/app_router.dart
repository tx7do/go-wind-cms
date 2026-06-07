import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/src/core/constants/index.dart' as constants;

import 'package:flutter_app/src/core/widgets/not_found_page.dart';
import 'package:flutter_app/src/app_router/route_names.dart';
import 'package:flutter_app/src/features/cms/pages/main_scaffold.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/post_detail_page.dart';
import 'package:flutter_app/src/features/cms/pages/post_list/post_list_page.dart';
import 'package:flutter_app/src/features/cms/pages/tag_feed/tag_feed_page.dart';
import 'package:flutter_app/src/features/cms/pages/tag_list/tag_list_page.dart';
import 'package:flutter_app/src/features/cms/pages/category_list/category_list_page.dart';
import 'package:flutter_app/src/features/cms/pages/search/search_page.dart';
import 'package:flutter_app/src/features/auth/pages/login_page.dart';

/// CMS 应用路由
class AppRouter {
  static const initial = constants.AppRoutePath.initial;

  static final router = GoRouter(
    initialLocation: initial,
    redirect: _guard,
    errorBuilder: (context, state) => const NotFoundPage(),
    routes: [
      // 主页 (底部导航)
      GoRoute(
        path: constants.AppRoutePath.initial,
        name: RouteNames.home,
        builder: (context, state) {
          return const CmsMainScaffold();
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
          // 文章列表（支持按分类/标签过滤）
          GoRoute(
            name: RouteNames.postList,
            path: 'posts',
            builder: (context, state) {
              final categoryId = int.tryParse(
                state.uri.queryParameters['categoryId'] ?? '',
              );
              final tagId = int.tryParse(
                state.uri.queryParameters['tagId'] ?? '',
              );
              return PostListPage(categoryId: categoryId, tagId: tagId);
            },
          ),
          // 分类列表
          GoRoute(
            name: RouteNames.categoryList,
            path: 'categories',
            builder: (context, state) {
              return const CategoryListPage();
            },
          ),
          // 标签列表
          GoRoute(
            name: RouteNames.tagList,
            path: 'tags',
            builder: (context, state) {
              return const TagListPage();
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
    ],
  );

  static FutureOr<String?> _guard(BuildContext context, GoRouterState state) {
    // CMS 内容展示端暂时不需要登录验证，直接放行
    return null;
  }
}
