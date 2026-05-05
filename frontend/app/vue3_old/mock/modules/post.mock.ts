import { defineMock } from 'vite-plugin-mock-dev-server'

// 生成随机的 RFC 3339 格式时间戳
function generateRandomTimestamp(daysAgo: number = 0, hoursOffset: number = 0): string {
  const date = new Date(Date.now() - daysAgo * 86400000 - hoursOffset * 3600000)
  return date.toISOString()
}

const posts = [
  {
    id: 1,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'vue3-best-practices',
    disallowComment: false,
    isFeatured: true,
    sortOrder: 1,
    visits: 1234,
    likes: 89,
    commentCount: 15,
    authorId: 1,
    authorName: '张三',
    customFields: {},
    translations: [
      {
        id: 1,
        postId: 1,
        languageCode: 'zh-CN',
        title: 'Vue 3 最佳实践指南',
        slug: 'vue3-best-practices',
        summary: '深入探讨 Vue 3 开发中的最佳实践，包括 Composition API、响应式系统、性能优化等内容。',
        content: `# Vue 3 最佳实践指南

Vue 3 带来了许多令人兴奋的新特性，本文将介绍一些最佳实践。

## Composition API

Composition API 是 Vue 3 最重要的特性之一，它提供了更好的逻辑复用和代码组织方式。

\`\`\`typescript
import { ref, computed } from 'vue'

export function useCounter() {
  const count = ref(0)
  const double = computed(() => count.value * 2)

  function increment() {
    count.value++
  }

  return { count, double, increment }
}
\`\`\`

## 性能优化

1. 使用 v-memo 减少重渲染
2. 合理使用 computed 和 watch
3. 懒加载组件

更多内容请关注后续文章...`,
        thumbnail: 'https://picsum.photos/800/450?random=11',
        fullPath: '/post/1',
        wordCount: 500,
        metaKeywords: 'Vue3,最佳实践,前端开发',
        metaDescription: 'Vue 3 开发最佳实践完整指南',
        seoTitle: 'Vue 3 最佳实践指南 - 技术博客',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [1],
    tagIds: [1, 2],
    createdBy: 1,
    createdAt: generateRandomTimestamp(5, 6),
    updatedAt: generateRandomTimestamp(4, 8),
    publishTime: generateRandomTimestamp(5, 6),
  },
  {
    id: 2,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'typescript-advanced',
    disallowComment: false,
    isFeatured: true,
    sortOrder: 2,
    visits: 987,
    likes: 65,
    commentCount: 12,
    authorId: 1,
    authorName: '张三',
    customFields: {},
    translations: [
      {
        id: 2,
        postId: 2,
        languageCode: 'zh-CN',
        title: 'TypeScript 高级技巧',
        slug: 'typescript-advanced',
        summary: '掌握 TypeScript 的高级特性，包括泛型、类型推导、条件类型等，提升代码质量。',
        content: `# TypeScript 高级技巧

TypeScript 提供了强大的类型系统，本文将介绍一些高级用法。

## 泛型约束

\`\`\`typescript
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
  return obj[key]
}
\`\`\`

## 条件类型

条件类型可以根据类型关系进行类型选择。

## 工具类型

TypeScript 内置了许多实用的工具类型，如 Partial、Required、Pick 等。`,
        thumbnail: 'https://picsum.photos/800/450?random=12',
        fullPath: '/post/2',
        wordCount: 800,
        metaKeywords: 'TypeScript,类型系统,高级技巧',
        metaDescription: 'TypeScript 高级特性完全指南',
        seoTitle: 'TypeScript 高级技巧 - 技术博客',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [1],
    tagIds: [1, 3],
    createdBy: 1,
    createdAt: generateRandomTimestamp(8, 4),
    updatedAt: generateRandomTimestamp(7, 6),
    publishTime: generateRandomTimestamp(8, 4),
  },
  {
    id: 3,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'headless-content-hub-intro',
    disallowComment: false,
    isFeatured: true,
    sortOrder: 3,
    visits: 756,
    likes: 45,
    commentCount: 8,
    authorId: 2,
    authorName: '李四',
    customFields: {},
    translations: [
      {
        id: 3,
        postId: 3,
        languageCode: 'zh-CN',
        title: '内容中台架构设计',
        slug: 'headless-content-hub-intro',
        summary: '深入了解 Content Hub 的架构设计理念，以及如何构建一个高性能、可扩展的内容管理系统。',
        content: `# Content Hub 架构设计

Content Hub 是现代内容管理的趋势。

## 什么是 Content Hub

Content Hub 将内容管理后端与前端展示解耦，通过 API 提供内容服务。

## 核心优势

1. 多渠道内容分发
2. 技术栈灵活
3. 更好的性能
4. 易于扩展

## 架构组件

- RESTful API
- GraphQL 支持
- 多租户隔离
- 缓存策略`,
        thumbnail: 'https://picsum.photos/800/450?random=13',
        fullPath: '/post/3',
        wordCount: 1200,
        metaKeywords: 'Content Hub,架构设计,内容管理',
        metaDescription: 'Content Hub 架构设计完整指南',
        seoTitle: 'Content Hub 架构设计 - 技术博客',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [1],
    tagIds: [4, 5],
    createdBy: 2,
    createdAt: generateRandomTimestamp(10, 2),
    updatedAt: generateRandomTimestamp(9, 4),
    publishTime: generateRandomTimestamp(10, 2),
  },
  {
    id: 4,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'work-life-balance',
    disallowComment: false,
    isFeatured: false,
    sortOrder: 4,
    visits: 543,
    likes: 32,
    commentCount: 6,
    authorId: 1,
    authorName: '张三',
    customFields: {},
    translations: [
      {
        id: 4,
        postId: 4,
        languageCode: 'zh-CN',
        title: '工作与生活的平衡之道',
        slug: 'work-life-balance',
        summary: '作为程序员，如何在高强度的工作中保持生活的平衡？分享我的一些经验和思考。',
        content: `# 工作与生活的平衡之道

程序员的工作往往很繁忙，但生活同样重要。

## 时间管理

合理安排工作和休息时间，提高效率。

## 健康第一

1. 定期运动
2. 充足睡眠
3. 健康饮食

## 兴趣爱好

培养工作之外的兴趣爱好，让生活更加丰富多彩。`,
        thumbnail: 'https://picsum.photos/800/450?random=14',
        fullPath: '/post/4',
        wordCount: 600,
        metaKeywords: '工作生活平衡,程序员,健康',
        metaDescription: '程序员如何平衡工作与生活',
        seoTitle: '工作与生活的平衡之道 - 生活随笔',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [2],
    tagIds: [6],
    createdBy: 1,
    createdAt: generateRandomTimestamp(12, 3),
    updatedAt: generateRandomTimestamp(11, 5),
    publishTime: generateRandomTimestamp(12, 3),
  },
  {
    id: 5,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'ui-design-principles',
    disallowComment: false,
    isFeatured: false,
    sortOrder: 5,
    visits: 678,
    likes: 48,
    commentCount: 10,
    authorId: 2,
    authorName: '李四',
    customFields: {},
    translations: [
      {
        id: 5,
        postId: 5,
        languageCode: 'zh-CN',
        title: 'UI 设计的基本原则',
        slug: 'ui-design-principles',
        summary: '好的 UI 设计需要遵循一些基本原则，本文将介绍设计中的关键要素。',
        content: `# UI 设计的基本原则

优秀的 UI 设计能够提升用户体验。

## 一致性

保持设计元素的一致性，包括颜色、字体、间距等。

## 可用性

1. 清晰的信息层级
2. 直观的操作反馈
3. 符合用户习惯

## 美观性

在功能的基础上追求视觉美感。

## 响应式设计

适配不同设备和屏幕尺寸。`,
        thumbnail: 'https://picsum.photos/800/450?random=15',
        fullPath: '/post/5',
        wordCount: 700,
        metaKeywords: 'UI设计,用户体验,设计原则',
        metaDescription: 'UI 设计基本原则完整指南',
        seoTitle: 'UI 设计的基本原则 - 产品设计',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [3],
    tagIds: [7, 8],
    createdBy: 2,
    createdAt: generateRandomTimestamp(14, 1),
    updatedAt: generateRandomTimestamp(13, 3),
    publishTime: generateRandomTimestamp(14, 1),
  },
  {
    id: 6,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'startup-lessons',
    disallowComment: false,
    isFeatured: false,
    sortOrder: 6,
    visits: 892,
    likes: 56,
    commentCount: 14,
    authorId: 1,
    authorName: '张三',
    customFields: {},
    translations: [
      {
        id: 6,
        postId: 6,
        languageCode: 'zh-CN',
        title: '创业一年的经验总结',
        slug: 'startup-lessons',
        summary: '从零到一，创业路上的经验和教训，希望能给其他创业者一些启发。',
        content: `# 创业一年的经验总结

创业一年，收获颇多。

## 产品为王

好的产品是成功的基础，要深入了解用户需求。

## 团队建设

1. 找到志同道合的伙伴
2. 明确分工
3. 保持沟通

## 资金管理

合理规划资金使用，控制成本。

## 持续学习

市场在变化，要不断学习和调整。`,
        thumbnail: 'https://picsum.photos/800/450?random=16',
        fullPath: '/post/6',
        wordCount: 950,
        metaKeywords: '创业经验，团队建设，产品开发',
        metaDescription: '创业一年的经验和教训总结',
        seoTitle: '创业一年的经验总结 - 创业思考',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [4],
    tagIds: [9, 10],
    createdBy: 1,
    createdAt: generateRandomTimestamp(16, 2),
    updatedAt: generateRandomTimestamp(15, 4),
    publishTime: generateRandomTimestamp(16, 2),
  },
  // 前端开发 (category 11)
  {
    id: 7,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'react-hooks-guide',
    disallowComment: false,
    isFeatured: true,
    sortOrder: 7,
    visits: 1567,
    likes: 123,
    commentCount: 28,
    authorId: 1,
    authorName: '张三',
    customFields: {},
    translations: [
      {
        id: 7,
        postId: 7,
        languageCode: 'zh-CN',
        title: 'React Hooks 完全指南',
        slug: 'react-hooks-guide',
        summary: '深入理解 React Hooks 的使用，包括 useState、useEffect、useContext 等常用 Hook。',
        content: `# React Hooks 完全指南

React Hooks 是 React 16.8 引入的新特性。

## useState

管理组件状态。

## useEffect

处理副作用。

## useContext

简化状态管理。`,
        thumbnail: 'https://picsum.photos/800/450?random=17',
        fullPath: '/post/7',
        wordCount: 1200,
        metaKeywords: 'React,Hooks,前端开发',
        metaDescription: 'React Hooks 完整使用指南',
        seoTitle: 'React Hooks 完全指南 - 前端开发',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [11],
    tagIds: [1, 2],
    createdBy: 1,
    createdAt: generateRandomTimestamp(3, 2),
    updatedAt: generateRandomTimestamp(2, 4),
    publishTime: generateRandomTimestamp(3, 2),
  },
  // 后端开发 (category 12)
  {
    id: 8,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'nodejs-microservices',
    disallowComment: false,
    isFeatured: true,
    sortOrder: 8,
    visits: 1345,
    likes: 98,
    commentCount: 22,
    authorId: 2,
    authorName: '李四',
    customFields: {},
    translations: [
      {
        id: 8,
        postId: 8,
        languageCode: 'zh-CN',
        title: 'Node.js 微服务架构实践',
        slug: 'nodejs-microservices',
        summary: '使用 Node.js 构建微服务架构的完整实践，包括服务拆分、通信、监控等。',
        content: `# Node.js 微服务架构实践

微服务架构是现代后端开发的主流方案。

## 服务拆分

按业务领域进行合理拆分。

## 服务通信

RESTful API 和 gRPC。

## 监控与日志

集中式日志和监控体系。`,
        thumbnail: 'https://picsum.photos/800/450?random=18',
        fullPath: '/post/8',
        wordCount: 1500,
        metaKeywords: 'Node.js，微服务，后端开发',
        metaDescription: 'Node.js 微服务架构完整实践',
        seoTitle: 'Node.js 微服务架构实践 - 后端开发',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [12],
    tagIds: [3, 4],
    createdBy: 2,
    createdAt: generateRandomTimestamp(6, 1),
    updatedAt: generateRandomTimestamp(5, 3),
    publishTime: generateRandomTimestamp(6, 1),
  },
  // 移动开发 (category 13)
  {
    id: 9,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'flutter-cross-platform',
    disallowComment: false,
    isFeatured: false,
    sortOrder: 9,
    visits: 876,
    likes: 67,
    commentCount: 15,
    authorId: 1,
    authorName: '张三',
    customFields: {},
    translations: [
      {
        id: 9,
        postId: 9,
        languageCode: 'zh-CN',
        title: 'Flutter 跨平台开发实战',
        slug: 'flutter-cross-platform',
        summary: '使用 Flutter 一套代码构建 iOS 和 Android 应用，提升开发效率。',
        content: `# Flutter 跨平台开发实战

Flutter 是 Google 推出的跨平台框架。

## 为什么选择 Flutter

高性能、热重载、丰富的组件。

## Widget 系统

一切皆 Widget。

## 状态管理

Provider、Riverpod 等方案。`,
        thumbnail: 'https://picsum.photos/800/450?random=19',
        fullPath: '/post/9',
        wordCount: 1100,
        metaKeywords: 'Flutter，跨平台，移动开发',
        metaDescription: 'Flutter 跨平台开发完整指南',
        seoTitle: 'Flutter 跨平台开发实战 - 移动开发',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [13],
    tagIds: [5, 6],
    createdBy: 1,
    createdAt: generateRandomTimestamp(9, 3),
    updatedAt: generateRandomTimestamp(8, 5),
    publishTime: generateRandomTimestamp(9, 3),
  },
  // 旅行游记 (category 21)
  {
    id: 10,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'yunnan-travel',
    disallowComment: false,
    isFeatured: false,
    sortOrder: 10,
    visits: 654,
    likes: 45,
    commentCount: 12,
    authorId: 2,
    authorName: '李四',
    customFields: {},
    translations: [
      {
        id: 10,
        postId: 10,
        languageCode: 'zh-CN',
        title: '云南七日自由行',
        slug: 'yunnan-travel',
        summary: '大理、丽江、香格里拉，七日云南之旅的精彩见闻和攻略。',
        content: `# 云南七日自由行

云南是一个美丽的地方。

## 大理

苍山洱海，风花雪月。

## 丽江

古城夜景，纳西文化。

## 香格里拉

高原风光，藏族风情。`,
        thumbnail: 'https://picsum.photos/800/450?random=20',
        fullPath: '/post/10',
        wordCount: 800,
        metaKeywords: '云南，旅行，游记',
        metaDescription: '云南七日自由行完整攻略',
        seoTitle: '云南七日自由行 - 旅行游记',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [21],
    tagIds: [7, 8],
    createdBy: 2,
    createdAt: generateRandomTimestamp(11, 2),
    updatedAt: generateRandomTimestamp(10, 4),
    publishTime: generateRandomTimestamp(11, 2),
  },
  // 美食探店 (category 22)
  {
    id: 11,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'shanghai-food',
    disallowComment: false,
    isFeatured: false,
    sortOrder: 11,
    visits: 543,
    likes: 38,
    commentCount: 9,
    authorId: 1,
    authorName: '张三',
    customFields: {},
    translations: [
      {
        id: 11,
        postId: 11,
        languageCode: 'zh-CN',
        title: '上海美食探店指南',
        slug: 'shanghai-food',
        summary: '探索上海的地道美食，从本帮菜到国际料理，满足你的味蕾。',
        content: `# 上海美食探店指南

上海是美食的天堂。

## 本帮菜

浓油赤酱，原汁原味。

## 小笼包

皮薄馅大，汤汁鲜美。

## 国际美食

各国料理应有尽有。`,
        thumbnail: 'https://picsum.photos/800/450?random=21',
        fullPath: '/post/11',
        wordCount: 700,
        metaKeywords: '上海，美食，探店',
        metaDescription: '上海美食探店完整指南',
        seoTitle: '上海美食探店指南 - 美食探店',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [22],
    tagIds: [9, 10],
    createdBy: 1,
    createdAt: generateRandomTimestamp(13, 1),
    updatedAt: generateRandomTimestamp(12, 3),
    publishTime: generateRandomTimestamp(13, 1),
  },
  // UI 设计 (category 31)
  {
    id: 12,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'figma-design-system',
    disallowComment: false,
    isFeatured: true,
    sortOrder: 12,
    visits: 1234,
    likes: 89,
    commentCount: 20,
    authorId: 2,
    authorName: '李四',
    customFields: {},
    translations: [
      {
        id: 12,
        postId: 12,
        languageCode: 'zh-CN',
        title: 'Figma 设计系统搭建指南',
        slug: 'figma-design-system',
        summary: '使用 Figma 构建设计系统，提升设计效率和一致性。',
        content: `# Figma 设计系统搭建指南

设计系统是现代设计 workflow 的核心。

## 组件库

创建可复用的 UI 组件。

## 设计令牌

统一管理颜色、字体等。

## 协作流程

设计与开发的协作规范。`,
        thumbnail: 'https://picsum.photos/800/450?random=22',
        fullPath: '/post/12',
        wordCount: 950,
        metaKeywords: 'Figma，设计系统，UI 设计',
        metaDescription: 'Figma 设计系统完整指南',
        seoTitle: 'Figma 设计系统搭建指南 - UI 设计',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [31],
    tagIds: [11, 12],
    createdBy: 2,
    createdAt: generateRandomTimestamp(4, 2),
    updatedAt: generateRandomTimestamp(3, 4),
    publishTime: generateRandomTimestamp(4, 2),
  },
  // UX 设计 (category 32)
  {
    id: 13,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'ux-research-methods',
    disallowComment: false,
    isFeatured: false,
    sortOrder: 13,
    visits: 765,
    likes: 54,
    commentCount: 11,
    authorId: 1,
    authorName: '张三',
    customFields: {},
    translations: [
      {
        id: 13,
        postId: 13,
        languageCode: 'zh-CN',
        title: '用户体验研究方法大全',
        slug: 'ux-research-methods',
        summary: '掌握用户体验研究的各种方法，包括用户访谈、可用性测试等。',
        content: `# 用户体验研究方法大全

用户研究是优秀 UX 的基础。

## 用户访谈

深入了解用户需求。

## 可用性测试

发现产品设计中的问题。

## 数据分析

量化用户行为。`,
        thumbnail: 'https://picsum.photos/800/450?random=23',
        fullPath: '/post/13',
        wordCount: 1050,
        metaKeywords: 'UX 研究，用户体验，设计方法',
        metaDescription: '用户体验研究方法完整指南',
        seoTitle: '用户体验研究方法大全 - UX 设计',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [32],
    tagIds: [13, 14],
    createdBy: 1,
    createdAt: generateRandomTimestamp(7, 3),
    updatedAt: generateRandomTimestamp(6, 5),
    publishTime: generateRandomTimestamp(7, 3),
  },
  // 团队管理 (category 41)
  {
    id: 14,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'agile-team-management',
    disallowComment: false,
    isFeatured: false,
    sortOrder: 14,
    visits: 987,
    likes: 72,
    commentCount: 16,
    authorId: 2,
    authorName: '李四',
    customFields: {},
    translations: [
      {
        id: 14,
        postId: 14,
        languageCode: 'zh-CN',
        title: '敏捷团队管理实战',
        slug: 'agile-team-management',
        summary: '敏捷开发团队的管理经验，包括 Scrum、看板等方法论。',
        content: `# 敏捷团队管理实战

敏捷开发需要合适的管理方式。

## Scrum 框架

Sprint、站会、回顾会。

## 看板方法

可视化工作流，限制在制品。

## 团队文化

建立信任和透明度。`,
        thumbnail: 'https://picsum.photos/800/450?random=24',
        fullPath: '/post/14',
        wordCount: 1150,
        metaKeywords: '敏捷管理，团队建设，Scrum',
        metaDescription: '敏捷团队管理完整实战指南',
        seoTitle: '敏捷团队管理实战 - 团队管理',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [41],
    tagIds: [15, 16],
    createdBy: 2,
    createdAt: generateRandomTimestamp(8, 1),
    updatedAt: generateRandomTimestamp(7, 3),
    publishTime: generateRandomTimestamp(8, 1),
  },
  // 产品思考 (category 42)
  {
    id: 15,
    status: 'POST_STATUS_PUBLISHED',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'product-market-fit',
    disallowComment: false,
    isFeatured: true,
    sortOrder: 15,
    visits: 1456,
    likes: 112,
    commentCount: 25,
    authorId: 1,
    authorName: '张三',
    customFields: {},
    translations: [
      {
        id: 15,
        postId: 15,
        languageCode: 'zh-CN',
        title: '如何找到产品市场契合点',
        slug: 'product-market-fit',
        summary: '产品市场契合点是创业成功的关键，本文分享寻找 PMF 的方法和经验。',
        content: `# 如何找到产品市场契合点

PMF 是创业成功的核心要素。

## 理解 PMF

产品满足市场的真实需求。

## 验证方法

用户增长、留存率、口碑传播。

## 迭代优化

快速试错，持续改进。`,
        thumbnail: 'https://picsum.photos/800/450?random=25',
        fullPath: '/post/15',
        wordCount: 1300,
        metaKeywords: '产品市场契合，PMF，创业',
        metaDescription: '寻找产品市场契合点的完整指南',
        seoTitle: '如何找到产品市场契合点 - 产品思考',
      },
    ],
    availableLanguages: ['zh-CN'],
    categoryIds: [42],
    tagIds: [17, 18],
    createdBy: 1,
    createdAt: generateRandomTimestamp(2, 1),
    updatedAt: generateRandomTimestamp(1, 3),
    publishTime: generateRandomTimestamp(2, 1),
  },
]

// 不再需要自动生成更多文章，我们已经为每个 category 添加了专门的文章

const postEnMap: Record<string, {
  title: string
  summary: string
  content: string
  metaKeywords: string
  metaDescription: string
  seoTitle: string
}> = {
  'vue3-best-practices': {
    title: 'Vue 3 Best Practices Guide',
    summary: 'A deep dive into Vue 3 best practices, including Composition API, reactivity, and performance optimization.',
    content: `# Vue 3 Best Practices Guide

Vue 3 introduces many exciting features. This article shares practical best practices for real-world projects.

## Composition API

Composition API improves logic reuse and code organization.

## Performance Optimization

1. Use v-memo to reduce unnecessary renders
2. Use computed and watch appropriately
3. Lazy-load heavy components`,
    metaKeywords: 'Vue3,best practices,frontend development',
    metaDescription: 'A complete guide to Vue 3 development best practices',
    seoTitle: 'Vue 3 Best Practices Guide - Tech Blog',
  },
  'typescript-advanced': {
    title: 'Advanced TypeScript Techniques',
    summary: 'Master advanced TypeScript features such as generics, type inference, and conditional types.',
    content: `# Advanced TypeScript Techniques

TypeScript provides a powerful type system. This article introduces practical advanced usage.

## Generic Constraints

Use generic constraints to build safer reusable APIs.

## Conditional Types

Conditional types enable flexible type-level logic.

## Utility Types

Use Partial, Required, Pick, and others to improve maintainability.`,
    metaKeywords: 'TypeScript,type system,advanced techniques',
    metaDescription: 'Complete guide to advanced TypeScript features',
    seoTitle: 'Advanced TypeScript Techniques - Tech Blog',
  },
  'headless-content-hub-intro': {
    title: 'Content Hub Architecture Design',
    summary: 'Understand the architecture of Content Hub and how to build a high-performance, scalable content platform.',
    content: `# Content Hub Architecture Design

Content Hub is a major trend in modern content management.

## What is Content Hub

It decouples content management from frontend presentation and delivers content via APIs.

## Core Benefits

1. Multi-channel delivery
2. Flexible technology stack
3. Better performance
4. Easier scalability`,
    metaKeywords: 'Content Hub,architecture,content management',
    metaDescription: 'Complete guide to Content Hub architecture design',
    seoTitle: 'Content Hub Architecture Design - Tech Blog',
  },
  'work-life-balance': {
    title: 'How to Balance Work and Life',
    summary: 'How can developers keep balance in high-intensity work? Here are practical lessons and reflections.',
    content: `# How to Balance Work and Life

Developer work can be intense, but life matters equally.

## Time Management

Plan focused work and quality rest.

## Health First

1. Exercise regularly
2. Get enough sleep
3. Maintain healthy eating habits

## Build Hobbies

Interests outside of work enrich life and prevent burnout.`,
    metaKeywords: 'work-life balance,developer,health',
    metaDescription: 'How developers can balance work and life',
    seoTitle: 'How to Balance Work and Life - Life Notes',
  },
  'ui-design-principles': {
    title: 'Core Principles of UI Design',
    summary: 'Great UI design follows core principles. This article explains essential factors in practical design work.',
    content: `# Core Principles of UI Design

Excellent UI design improves user experience.

## Consistency

Keep colors, typography, and spacing consistent.

## Usability

1. Clear information hierarchy
2. Intuitive interaction feedback
3. Follow user habits

## Aesthetics

Deliver visual quality while preserving functionality.`,
    metaKeywords: 'UI design,user experience,design principles',
    metaDescription: 'Complete guide to UI design principles',
    seoTitle: 'Core Principles of UI Design - Product Design',
  },
  'startup-lessons': {
    title: 'Lessons Learned After One Year of Startup',
    summary: 'From zero to one: startup lessons and reflections that may help other founders.',
    content: `# Lessons Learned After One Year of Startup

After one year of building a startup, there are many lessons worth sharing.

## Product First

A strong product starts from deep user understanding.

## Team Building

1. Find aligned partners
2. Clarify responsibilities
3. Keep communication transparent

## Financial Discipline

Plan spending carefully and control cost risks.`,
    metaKeywords: 'startup lessons,team building,product development',
    metaDescription: 'Lessons and reflections from one year of startup',
    seoTitle: 'Lessons Learned After One Year of Startup - Startup Insights',
  },
}

posts.forEach((post) => {
  const zh = post.translations.find(t => t.languageCode === 'zh-CN')
  if (!zh || post.translations.some(t => t.languageCode === 'en-US')) {
    return
  }

  const mapped = postEnMap[post.slug]
  post.translations.push({
    id: zh.id + 100,
    postId: post.id,
    languageCode: 'en-US',
    title: mapped?.title ?? `Post Title ${post.id}`,
    slug: zh.slug,
    summary: mapped?.summary ?? `This is the summary of post ${post.id}, sharing useful ideas and insights.`,
    content: mapped?.content ?? `# Post Title ${post.id}

This is the content body for post ${post.id}.

## Section 1

Interesting details...

## Section 2

More practical content...`,
    thumbnail: zh.thumbnail,
    fullPath: zh.fullPath,
    wordCount: zh.wordCount,
    metaKeywords: mapped?.metaKeywords ?? `post-${post.id},blog,article`,
    metaDescription: mapped?.metaDescription ?? `Description for post ${post.id}`,
    seoTitle: mapped?.seoTitle ?? `Post Title ${post.id} - Blog`,
  })

  post.availableLanguages = ['zh-CN', 'en-US']
})

export default defineMock([
  {
    url: '/app/v1/posts',
    method: 'GET',
    body: ({ query }: any) => {
      const page = parseInt(query.page || '1')
      const pageSize = parseInt(query.pageSize || '10')

      let filteredPosts = [...posts]

      // 按状态筛选
      if (query.query) {
        try {
          const queryObj = JSON.parse(query.query)

          // 按分类筛选
          if (queryObj.categoryIds) {
            const categoryIds = Array.isArray(queryObj.categoryIds) ? queryObj.categoryIds : [queryObj.categoryIds]
            filteredPosts = filteredPosts.filter(post =>
              post.categoryIds.some(id => categoryIds.includes(id))
            )
          }

          // 按推荐筛选
          if (queryObj.isFeatured !== undefined) {
            filteredPosts = filteredPosts.filter(post => post.isFeatured === queryObj.isFeatured)
          }

          // 按标题搜索
          if (queryObj.title___icontains) {
            const keyword = queryObj.title___icontains.toLowerCase()
            filteredPosts = filteredPosts.filter(post =>
              post.translations[0].title.toLowerCase().includes(keyword)
            )
          }

          // 根据 locale 过滤翻译数据
          if (queryObj.locale) {
            filteredPosts = filteredPosts.map(post => {
              const filteredTranslations = post.translations.filter(
                t => t.languageCode === queryObj.locale
              )

              if (filteredTranslations.length > 0) {
                return {
                  ...post,
                  translations: filteredTranslations
                }
              }

              return {
                ...post,
                translations: [post.translations[0]]
              }
            })
          }
        } catch (e) {
          console.error('Parse query error:', e)
        }
      }

      const start = (page - 1) * pageSize
      const end = start + pageSize

      return {
        items: filteredPosts.slice(start, end),
        total: filteredPosts.length,
      }
    },
  },
  {
    url: '/app/v1/posts/:id',
    method: 'GET',
    body: ({ params, query }: any) => {
      const id = parseInt(params.id)
      let post = posts.find(p => p.id === id)

      if (!post) {
        return { error: 'Post not found' }
      }

      // ✅ 根据 languageCode 参数过滤翻译数据
      if (query.languageCode) {
        const filteredTranslations = post.translations.filter(
          t => t.languageCode === query.languageCode
        )

        post = {
          ...post,
          translations: filteredTranslations.length > 0
            ? filteredTranslations
            : [post.translations[0]]
        }
      }

      return post
    },
  },
])

