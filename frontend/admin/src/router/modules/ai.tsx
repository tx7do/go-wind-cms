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
          // 问数查全平台原生表，平台管理员专属（后端对租户 403）
          authority: ['sys:platform_admin'],
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
          // 提供商是平台级配置（且 api_key 敏感），租户不可见
          authority: ['sys:platform_admin'],
        },
      },
      {
        name: 'ai-agents',
        path: 'agents', // 相对路径，最终为 /ai/agents
        element: createLazyRoute(() => import('@/pages/app/ai/agent')),
        meta: {
          title: 'routes:ai-agents',
          icon: 'lucide:workflow',
          order: 6,
          // 权限：平台管理员与租户管理员（编排应用按租户隔离）
        },
      },
    ],
  },
];
export default aiRoutes;
