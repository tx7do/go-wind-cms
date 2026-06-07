import 'package:flutter_app/src/features/cms/services/navigation_service.dart'
    show Navigation, NavigationItem, NavigationItemLinkType, NavigationLocation;

import 'package:flutter_app/generated/api/models/content_service_v1_category.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_category_translation.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_post_translation.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_comment.dart';
import 'package:flutter_app/generated/api/models/comment_service_v1_comment_content_type.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag.dart';
import 'package:flutter_app/generated/api/models/content_service_v1_tag_translation.dart';

typedef Category = ContentServiceV1Category;
typedef CategoryTranslation = ContentServiceV1CategoryTranslation;
typedef Post = ContentServiceV1Post;
typedef PostTranslation = ContentServiceV1PostTranslation;
typedef Comment = CommentServiceV1Comment;
typedef CommentContentType = CommentServiceV1CommentContentType;
typedef Tag = ContentServiceV1Tag;
typedef TagTranslation = ContentServiceV1TagTranslation;


/// 辅助：创建 DateTime
DateTime _ts(
  int year, [
  int month = 1,
  int day = 1,
  int hour = 0,
  int minute = 0,
]) {
  return DateTime(year, month, day, hour, minute);
}

/// Mock 分类数据
final mockCategories = [
  Category(
    id: 1,
    sortOrder: 0,
    isNav: true,
    postCount: 120,
    translations: [
      ContentServiceV1CategoryTranslation(languageCode: 'zh', name: '推荐', slug: 'recommend'),
    ],
  ),
  Category(
    id: 2,
    sortOrder: 1,
    isNav: true,
    postCount: 85,
    translations: [
      CategoryTranslation(
        languageCode: 'zh',
        name: '科技',
        slug: 'technology',
        description: '科技前沿资讯',
      ),
    ],
  ),
  Category(
    id: 3,
    sortOrder: 2,
    isNav: true,
    postCount: 64,
    translations: [
      CategoryTranslation(
        languageCode: 'zh',
        name: '生活',
        slug: 'life',
        description: '生活感悟与分享',
      ),
    ],
  ),
  Category(
    id: 4,
    sortOrder: 3,
    isNav: true,
    postCount: 47,
    translations: [
      CategoryTranslation(
        languageCode: 'zh',
        name: '设计',
        slug: 'design',
        description: 'UI/UX 设计心得',
      ),
    ],
  ),
  Category(
    id: 5,
    sortOrder: 4,
    isNav: true,
    postCount: 93,
    translations: [
      CategoryTranslation(
        languageCode: 'zh',
        name: '开发',
        slug: 'development',
        description: '编程开发技术',
      ),
    ],
  ),
  Category(
    id: 6,
    sortOrder: 5,
    isNav: true,
    postCount: 38,
    translations: [
      CategoryTranslation(
        languageCode: 'zh',
        name: '美食',
        slug: 'food',
        description: '美食探店与食谱',
      ),
    ],
  ),
  Category(
    id: 7,
    sortOrder: 6,
    isNav: true,
    postCount: 52,
    translations: [
      CategoryTranslation(
        languageCode: 'zh',
        name: '旅行',
        slug: 'travel',
        description: '旅行见闻',
      ),
    ],
  ),
];

final mockTags = <Tag>[
  Tag(
    id: 1,
    translations: [
      TagTranslation(languageCode: 'zh', name: 'Flutter', slug: 'flutter'),
    ],
  ),
  Tag(
    id: 2,
    translations: [
      TagTranslation(languageCode: 'zh', name: 'Dart', slug: 'dart'),
    ],
  ),
  Tag(
    id: 3,
    translations: [
      TagTranslation(languageCode: 'zh', name: 'CMS', slug: 'cms'),
    ],
  ),
  Tag(
    id: 4,
    translations: [
      TagTranslation(languageCode: 'zh', name: 'UI设计', slug: 'ui-design'),
    ],
  ),
  Tag(
    id: 5,
    translations: [
      TagTranslation(languageCode: 'zh', name: '前端', slug: 'frontend'),
    ],
  ),
  Tag(
    id: 6,
    translations: [
      TagTranslation(languageCode: 'zh', name: '后端', slug: 'backend'),
    ],
  ),
  Tag(
    id: 7,
    translations: [
      TagTranslation(languageCode: 'zh', name: '人工智能', slug: 'ai'),
    ],
  ),
  Tag(
    id: 8,
    translations: [
      TagTranslation(languageCode: 'zh', name: '摄影', slug: 'photography'),
    ],
  ),
  Tag(
    id: 9,
    translations: [
      TagTranslation(languageCode: 'zh', name: '美食探店', slug: 'food'),
    ],
  ),
  Tag(
    id: 10,
    translations: [
      TagTranslation(languageCode: 'zh', name: '极简生活', slug: 'minimalist'),
    ],
  ),
];

/// Mock 文章数据
final mockPosts = [
  Post(
    id: 1,
    isFeatured: true,
    visits: 3420,
    likes: 186,
    commentCount: 24,
    authorName: '张明',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: 'Flutter 3.x 新特性完全指南',
        slug: 'flutter-3x-guide',
        summary:
            '深入探索 Flutter 3.x 带来的革命性变化，包括 Impeller 渲染引擎、全新的 Material 3 组件库以及平台特定功能的增强支持。',
        content: _longContent,
      ),
    ],
    categoryIds: [5],
    tagIds: [1, 2, 5],
    publishTime: _ts(2025, 12, 28),
    createdAt: _ts(2025, 12, 28),
  ),
  Post(
    id: 2,
    isFeatured: true,
    visits: 2180,
    likes: 145,
    commentCount: 18,
    authorName: '李华',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: '现代 CMS 架构设计：从单体到微服务',
        slug: 'cms-architecture',
        summary: '探讨内容管理系统从传统单体架构向微服务架构演进的最佳实践，涵盖 API 设计、缓存策略和内容分发等核心话题。',
        content: _longContent,
      ),
    ],
    categoryIds: [2],
    tagIds: [3, 6, 7],
    publishTime: _ts(2025, 12, 27),
    createdAt: _ts(2025, 12, 27),
  ),
  Post(
    id: 3,
    isFeatured: true,
    visits: 5100,
    likes: 267,
    commentCount: 32,
    authorName: '王芳',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: '用 AI 提升你的设计工作流',
        slug: 'ai-design-workflow',
        summary: '如何利用人工智能工具优化设计流程，从灵感收集到最终交付，AI 正在改变设计师的工作方式。',
        content: _longContent,
      ),
    ],
    categoryIds: [4],
    tagIds: [7, 4],
    publishTime: _ts(2025, 12, 26),
    createdAt: _ts(2025, 12, 26),
  ),
  Post(
    id: 4,
    isFeatured: true,
    visits: 6780,
    likes: 312,
    commentCount: 45,
    authorName: '陈晨',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: '探索城市中的隐藏美食小店',
        slug: 'hidden-food-shops',
        summary: '城市角落里藏着许多不起眼的美食小店，跟随我们的脚步，发现那些令人惊喜的味道。',
        content: _longContent,
      ),
    ],
    categoryIds: [6],
    tagIds: [9],
    publishTime: _ts(2025, 12, 25),
    createdAt: _ts(2025, 12, 25),
  ),
  Post(
    id: 5,
    isFeatured: true,
    visits: 1560,
    likes: 98,
    commentCount: 12,
    authorName: '赵伟',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: 'Dart 异步编程：从 Future 到 Stream',
        slug: 'dart-async',
        summary: '全面解析 Dart 异步编程模型，从基础的 Future 到高级的 Stream 操作，帮助你写出更高效的异步代码。',
        content: _longContent,
      ),
    ],
    categoryIds: [5],
    tagIds: [2, 1],
    publishTime: _ts(2025, 12, 24),
    createdAt: _ts(2025, 12, 24),
  ),
  Post(
    id: 6,
    visits: 4320,
    likes: 201,
    commentCount: 28,
    authorName: '林小美',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: '极简主义生活实践指南',
        slug: 'minimalist-life',
        summary: '如何用更少的物品过更丰富的生活？分享极简主义实践中的心得与技巧。',
        content: _longContent,
      ),
    ],
    categoryIds: [3],
    tagIds: [10],
    publishTime: _ts(2025, 12, 23),
    createdAt: _ts(2025, 12, 23),
  ),
  Post(
    id: 7,
    visits: 8900,
    likes: 389,
    commentCount: 56,
    authorName: '张明',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: 'React vs Flutter：2025 跨平台框架对比',
        slug: 'react-vs-flutter',
        summary: '从性能、开发体验、生态系统等多个维度对比 React Native 和 Flutter 两大跨平台框架。',
        content: _longContent,
      ),
    ],
    categoryIds: [5],
    tagIds: [1, 5],
    publishTime: _ts(2025, 12, 22),
    createdAt: _ts(2025, 12, 22),
  ),
  Post(
    id: 8,
    visits: 3210,
    likes: 156,
    commentCount: 19,
    authorName: '周影',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: '手机摄影构图技巧：拍出电影感',
        slug: 'mobile-photography',
        summary: '不需要专业相机，掌握这些构图技巧，用手机也能拍出大片效果。',
        content: _longContent,
      ),
    ],
    categoryIds: [3],
    tagIds: [8],
    publishTime: _ts(2025, 12, 21),
    createdAt: _ts(2025, 12, 21),
  ),
  Post(
    id: 9,
    visits: 2340,
    likes: 112,
    commentCount: 14,
    authorName: '李华',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: 'Go 语言在后端开发中的最佳实践',
        slug: 'go-backend-best-practices',
        summary: '分享使用 Go 语言开发高性能后端服务的经验总结，涵盖并发模型、错误处理和项目结构。',
        content: _longContent,
      ),
    ],
    categoryIds: [5],
    tagIds: [6],
    publishTime: _ts(2025, 12, 20),
    createdAt: _ts(2025, 12, 20),
  ),
  Post(
    id: 10,
    visits: 9870,
    likes: 425,
    commentCount: 38,
    authorName: '行者无疆',
    translations: [
      PostTranslation(
        languageCode: 'zh',
        title: '一个人的旅行：冰岛环岛自驾手记',
        slug: 'iceland-road-trip',
        summary: '在冰岛的十天自驾之旅，冰川、火山、极光，记录一路上震撼人心的自然景观。',
        content: _longContent,
      ),
    ],
    categoryIds: [7],
    tagIds: [],
    publishTime: _ts(2025, 12, 19),
    createdAt: _ts(2025, 12, 19),
  ),
];

/// Mock 评论数据
final mockComments = [
  Comment(
    id: 1,
    contentType: CommentContentType.contentTypePost,
    objectId: 1,
    content: '写得太好了！Impeller 那部分分析得很到位，期待更多 Flutter 相关内容。',
    authorName: '小王',
    likeCount: 12,
    createdAt: _ts(2025, 12, 28, 14, 30),
  ),
  Comment(
    id: 2,
    contentType: CommentContentType.contentTypePost,
    objectId: 1,
    content: 'Material 3 的适配确实需要花一些功夫，不过效果还是很不错的。',
    authorName: '阿飞',
    likeCount: 8,
    createdAt: _ts(2025, 12, 28, 15, 45),
  ),
  Comment(
    id: 3,
    contentType: CommentContentType.contentTypePost,
    objectId: 1,
    content: '有没有关于 Flutter Web 的性能优化建议？',
    authorName: '前端新手',
    parentId: 2,
    likeCount: 3,
    createdAt: _ts(2025, 12, 28, 16, 20),
  ),
  Comment(
    id: 4,
    contentType: CommentContentType.contentTypePost,
    objectId: 2,
    content: '微服务架构确实更适合大型 CMS 系统，我们公司最近也在做类似的迁移。',
    authorName: '架构师A',
    likeCount: 15,
    createdAt: _ts(2025, 12, 27, 10, 0),
  ),
  Comment(
    id: 5,
    contentType: CommentContentType.contentTypePost,
    objectId: 2,
    content: '缓存策略那部分可以再详细讲讲吗？Redis + CDN 的方案如何？',
    authorName: '运维老李',
    likeCount: 7,
    createdAt: _ts(2025, 12, 27, 11, 30),
  ),
  Comment(
    id: 6,
    contentType: CommentContentType.contentTypePost,
    objectId: 4,
    content: '那家小店我也去过！老板人超好，推荐他们的招牌牛肉面。',
    authorName: '吃货一号',
    likeCount: 22,
    createdAt: _ts(2025, 12, 25, 19, 0),
  ),
];

const _longContent = '''
## 正文内容

这是一篇示例文章的正文内容。在实际项目中，这里会从后端 API 获取完整的 Markdown 或 HTML 格式的文章内容。

Flutter 是 Google 推出的跨平台 UI 框架，可以用一套代码库同时构建 iOS、Android、Web 和桌面应用。自发布以来，Flutter 已经成为最受欢迎的跨平台开发框架之一。

### 主要优势

1. **热重载** - 开发过程中可以实时看到代码修改的效果
2. **丰富的组件库** - Material Design 和 Cupertino 风格的组件应有尽有
3. **高性能** - 使用 Dart 语言编译为原生代码，性能接近原生应用
4. **活跃的社区** - 大量的开源包和插件可供使用

### 总结

无论你是新手还是经验丰富的开发者，Flutter 都值得一试。它不仅能提高开发效率，还能带来出色的用户体验。
''';

/// Mock 导航数据 - 移动端底部导航 (MOBILE)
final mockMobileNavigation = Navigation(
  id: 1,
  name: '底部导航',
  location: NavigationLocation.mobile,
  locale: 'zh',
  isActive: true,
  items: [
    NavigationItem(
      id: 1,
      title: '首页',
      icon: 'home_outlined',
      url: '/',
      linkType: NavigationItemLinkType.linkTypeCustom,
      sortOrder: 0,
    ),
    NavigationItem(
      id: 2,
      title: '发现',
      icon: 'explore_outlined',
      url: '/explore',
      linkType: NavigationItemLinkType.linkTypeCustom,
      sortOrder: 1,
    ),
    NavigationItem(
      id: 3,
      title: '收藏',
      icon: 'bookmark_border',
      url: '/bookmarks',
      linkType: NavigationItemLinkType.linkTypeCustom,
      sortOrder: 2,
    ),
    NavigationItem(
      id: 4,
      title: '我的',
      icon: 'person_outline',
      url: '/profile',
      linkType: NavigationItemLinkType.linkTypeCustom,
      sortOrder: 3,
    ),
  ],
);

/// Mock 导航数据 - Web端顶部导航 (HEADER)
final mockHeaderNavigation = Navigation(
  id: 2,
  name: '顶部导航',
  location: NavigationLocation.header,
  locale: 'zh',
  isActive: true,
  items: [
    NavigationItem(
      id: 10,
      title: '首页',
      icon: 'home_outlined',
      url: '/',
      linkType: NavigationItemLinkType.linkTypeCustom,
      sortOrder: 0,
    ),
    NavigationItem(
      id: 11,
      title: '发现',
      icon: 'explore_outlined',
      url: '/explore',
      linkType: NavigationItemLinkType.linkTypeCustom,
      sortOrder: 1,
    ),
    NavigationItem(
      id: 12,
      title: '标签',
      icon: 'label_outlined',
      url: '/tag',
      linkType: NavigationItemLinkType.linkTypeCustom,
      sortOrder: 2,
    ),
  ],
);

/// 所有 Mock 导航数据（模拟 ListNavigationResponse）
final mockNavigations = [
  mockMobileNavigation,
  mockHeaderNavigation,
];
