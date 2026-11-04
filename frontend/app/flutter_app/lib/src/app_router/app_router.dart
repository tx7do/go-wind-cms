import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'package:flutter_app/src/core/constants/index.dart' as constants;
import 'package:flutter_app/src/core/constants/router_paths.dart';

import 'package:flutter_app/src/core/widgets/adaptive_shell.dart';
import 'package:flutter_app/src/core/widgets/not_found_page.dart';
import 'package:flutter_app/src/app_router/route_names.dart';
import 'package:flutter_app/src/features/cms/pages/home/home_page.dart';
import 'package:flutter_app/src/features/cms/pages/post_detail/post_detail_page.dart';
import 'package:flutter_app/src/features/cms/pages/post_list/post_list_page.dart';
import 'package:flutter_app/src/features/cms/pages/tag_feed/tag_feed_page.dart';
import 'package:flutter_app/src/features/cms/pages/tag_list/tag_list_page.dart';
import 'package:flutter_app/src/features/cms/pages/category_list/category_list_page.dart';
import 'package:flutter_app/src/features/cms/pages/category_detail/category_detail_page.dart';
import 'package:flutter_app/src/features/cms/pages/search/search_page.dart';
import 'package:flutter_app/src/features/auth/pages/login_page.dart';
import 'package:flutter_app/src/features/cms/pages/explore/explore_page.dart';
import 'package:flutter_app/src/features/cms/pages/bookmarks/bookmarks_page.dart';
import 'package:flutter_app/src/features/cms/pages/page_detail/page_detail_page.dart';
import 'package:flutter_app/src/features/cms/pages/profile/profile_page.dart';
import 'package:flutter_app/src/features/cms/pages/settings/settings_page.dart';
import 'package:flutter_app/src/features/cms/pages/my_comments/my_comments_page.dart';
import 'package:flutter_app/src/features/cms/pages/about/about_page.dart';
import 'package:flutter_app/src/features/cms/pages/legal/legal_page.dart';

/// 底栏四支的根路径，顺序必须与 [StatefulShellRoute] 的 branches 一致
/// （AdaptiveShell 靠它判断「当前是否停在 tab 根位置」）。
const _branchRoots = <String>[
  AppRoutePath.initial,
  AppRoutePath.explore,
  AppRoutePath.bookmarks,
  AppRoutePath.profile,
];

/// CMS 应用路由
class AppRouter {
  static const initial = constants.AppRoutePath.initial;

  static final router = GoRouter(
    initialLocation: initial,
    redirect: _guard,
    errorBuilder: (context, state) => const NotFoundPage(),
    routes: [
      // 自适应壳层：窄屏底栏 + 宽屏顶栏，四支各自持有 Navigator
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) => AdaptiveShell(
          navigationShell: navigationShell,
          branchRoutes: _branchRoots,
          currentPath: state.uri.path,
        ),
        branches: [
          // 0 首页支：详情/列表/搜索都挂在里面，返回时首页滚动位置不丢
          StatefulShellBranch(
            initialLocation: AppRoutePath.initial,
            routes: [
              GoRoute(
                path: AppRoutePath.initial,
                name: RouteNames.home,
                builder: (context, state) => const HomePage(),
                routes: [
                  // 文章列表（/post）
                  GoRoute(
                    path: 'post',
                    name: RouteNames.postList,
                    builder: (context, state) {
                      final categoryId = int.tryParse(
                        state.uri.queryParameters['categoryId'] ?? '',
                      );
                      final tagId = int.tryParse(
                        state.uri.queryParameters['tagId'] ?? '',
                      );
                      return PostListPage(
                        categoryId: categoryId,
                        tagId: tagId,
                      );
                    },
                  ),
                  // 文章详情（/post/:id）——与列表平级，避免嵌套叠页
                  GoRoute(
                    path: 'post/:id',
                    name: RouteNames.postDetail,
                    builder: (context, state) {
                      final postId =
                          int.tryParse(state.pathParameters['id'] ?? '0') ?? 0;
                      return PostDetailPage(postId: postId);
                    },
                  ),
                  // 标签列表（/tag）与标签文章流（/tag/:id）
                  GoRoute(
                    path: 'tag',
                    name: RouteNames.tagList,
                    builder: (context, state) => const TagListPage(),
                  ),
                  GoRoute(
                    path: 'tag/:id',
                    name: RouteNames.tagFeed,
                    builder: (context, state) {
                      final tagId =
                          int.tryParse(state.pathParameters['id'] ?? '0') ?? 0;
                      return TagFeedPage(tagId: tagId);
                    },
                  ),
                  // 分类列表（/category）与分类详情（/category/:id）
                  GoRoute(
                    path: 'category',
                    name: RouteNames.categoryList,
                    builder: (context, state) => const CategoryListPage(),
                  ),
                  GoRoute(
                    path: 'category/:id',
                    name: RouteNames.categoryDetail,
                    builder: (context, state) {
                      final categoryId = int.tryParse(
                            state.pathParameters['id'] ?? '0',
                          ) ??
                          0;
                      return CategoryDetailPage(categoryId: categoryId);
                    },
                  ),
                  GoRoute(
                    path: 'search',
                    name: RouteNames.search,
                    builder: (context, state) => const SearchPage(),
                  ),
                  // CMS 页面详情（Page + Section 区块渲染）
                  GoRoute(
                    path: 'page/:id',
                    name: RouteNames.pageDetail,
                    builder: (context, state) {
                      final pageId =
                          int.tryParse(state.pathParameters['id'] ?? '0') ?? 0;
                      return PageDetailPage(pageId: pageId);
                    },
                    routes: [
                      GoRoute(
                        path: 'slug/:slug',
                        builder: (context, state) => PageDetailPage(
                          slug: state.pathParameters['slug'] ?? '',
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ],
          ),
          // 1 发现
          StatefulShellBranch(
            initialLocation: AppRoutePath.explore,
            routes: [
              GoRoute(
                path: AppRoutePath.explore,
                name: RouteNames.explore,
                builder: (context, state) => const ExplorePage(),
              ),
            ],
          ),
          // 2 收藏（此前只有服务端导航能进，现在可深链）
          StatefulShellBranch(
            initialLocation: AppRoutePath.bookmarks,
            routes: [
              GoRoute(
                path: AppRoutePath.bookmarks,
                name: RouteNames.bookmarks,
                builder: (context, state) => const BookmarksPage(),
              ),
            ],
          ),
          // 3 我的：设置/评论/关于/法务都在这一支里，来回切不丢分支状态
          StatefulShellBranch(
            initialLocation: AppRoutePath.profile,
            routes: [
              GoRoute(
                path: AppRoutePath.profile,
                name: RouteNames.profile,
                builder: (context, state) => const ProfilePage(),
              ),
              GoRoute(
                path: AppRoutePath.settings,
                name: RouteNames.settings,
                builder: (context, state) => const SettingsPage(),
              ),
              GoRoute(
                path: AppRoutePath.myComments,
                name: RouteNames.myComments,
                builder: (context, state) => const MyCommentsPage(),
              ),
              GoRoute(
                path: AppRoutePath.about,
                name: RouteNames.about,
                builder: (context, state) => const AboutPage(),
              ),
              GoRoute(
                path: AppRoutePath.contact,
                name: RouteNames.contact,
                builder: (context, state) =>
                    const LegalPage(type: LegalPageType.contact),
              ),
              GoRoute(
                path: AppRoutePath.disclaimer,
                name: RouteNames.disclaimer,
                builder: (context, state) =>
                    const LegalPage(type: LegalPageType.disclaimer),
              ),
              GoRoute(
                path: AppRoutePath.privacy,
                name: RouteNames.privacy,
                builder: (context, state) =>
                    const LegalPage(type: LegalPageType.privacy),
              ),
              GoRoute(
                path: AppRoutePath.terms,
                name: RouteNames.terms,
                builder: (context, state) =>
                    const LegalPage(type: LegalPageType.terms),
              ),
            ],
          ),
        ],
      ),
      // 登录页（不在壳内：全屏，不带任何导航外框）
      GoRoute(
        name: RouteNames.login,
        path: AppRoutePath.login,
        builder: (context, state) => const LoginPage(),
      ),
    ],
  );

  static FutureOr<String?> _guard(BuildContext context, GoRouterState state) {
    // CMS 内容展示端暂时不需要登录验证，直接放行
    return null;
  }
}
