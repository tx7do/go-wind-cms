import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/src/core/constants/index.dart' as constants;

import 'package:flutter_app/src/features/cms/pages/main_scaffold.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/post_detail_page.dart';
import 'package:flutter_app/src/features/cms/pages/tag_feed/tag_feed_page.dart';
import 'package:flutter_app/src/features/cms/pages/search/search_page.dart';

/// CMS 应用路由
class AppRouter {
  static const initial = constants.AppRoutePath.initial;

  static final router = GoRouter(
    initialLocation: initial,
    redirect: _guard,
    routes: [
      // 主页 (底部导航)
      GoRoute(
        path: constants.AppRoutePath.initial,
        name: 'home',
        builder: (context, state) {
          return const CmsMainScaffold();
        },
        routes: [
          // 文章详情
          GoRoute(
            name: 'post_detail',
            path: 'post/:id',
            builder: (context, state) {
              final postId =
                  int.tryParse(state.pathParameters['id'] ?? '0') ?? 0;
              return PostDetailPage(postId: postId);
            },
          ),
          // 标签文章列表
          GoRoute(
            name: 'tag_feed',
            path: 'tag/:id',
            builder: (context, state) {
              final tagId =
                  int.tryParse(state.pathParameters['id'] ?? '0') ?? 0;
              return TagFeedPage(tagId: tagId);
            },
          ),
          // 搜索
          GoRoute(
            name: 'search',
            path: 'search',
            builder: (context, state) {
              return const SearchPage();
            },
          ),
        ],
      ),
    ],
  );

  static FutureOr<String?> _guard(BuildContext context, GoRouterState state) {
    // CMS 内容展示端暂时不需要登录验证，直接放行
    return null;
  }
}
