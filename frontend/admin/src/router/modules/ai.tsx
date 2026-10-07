import type { AppRouteObject } from '@/core/router/types';
import { createLazyRoute } from '@/core/router';

/**
 * AI 路由配置：一期仅模型提供商管理（内容生成走各表单的 AI 生成按钮，
 * 无独立页面；对话/知识库/问数待后端能力补齐后按需追加）。
 */
export const aiRoutes: AppRouteObject[] = [
  {
    name: 'ai',
    path: 'ai',
    meta: {
      title: 'routes:ai',
      icon: 'lucide:sparkles',
      order: 2010,
      keepAlive: true,
    },
    children: [
      {
        name: 'ai-chat',
        path: 'chat',
        element: createLazyRoute(() => import('@/pages/app/ai/chat')),
        meta: {
          title: 'routes:ai-chat',
          icon: 'lucide:message-circle',
          order: 1,
        },
      },
      {
        name: 'ai-query',
        path: 'query',
        element: createLazyRoute(() => import('@/pages/app/ai/query')),
        meta: {
          title: 'routes:ai-query',
          icon: 'lucide:search-code',
          order: 2,
        },
      },
      {
        name: 'ai-knowledge',
        path: 'knowledge',
        element: createLazyRoute(() => import('@/pages/app/ai/knowledge')),
        meta: {
          title: 'routes:ai-knowledge',
          icon: 'lucide:book-open',
          order: 3,
        },
      },
      {
        name: 'ai-usage',
        path: 'usage',
        element: createLazyRoute(() => import('@/pages/app/ai/usage')),
        meta: {
          title: 'routes:ai-usage',
          icon: 'lucide:bar-chart-3',
          order: 4,
        },
      },
      {
        name: 'ai-providers',
        path: 'providers',
        element: createLazyRoute(() => import('@/pages/app/ai/provider')),
        meta: {
          title: 'routes:ai-providers',
          icon: 'lucide:cpu',
          order: 5,
        },
      },
    ],
  },
];
export default aiRoutes;
