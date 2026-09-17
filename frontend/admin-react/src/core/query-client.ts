import {QueryClient, defaultShouldDehydrateQuery} from '@tanstack/react-query'

// 单独导出，方便测试 & 维护
export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            // 管理后台不缓存查询响应：每次组件挂载都向后端重新请求
            staleTime: 0,

            // 查询在无观察者（组件卸载）后立即回收，跨挂载不残留缓存
            gcTime: 0,

            // 窗口聚焦时自动重新验证
            refetchOnWindowFocus: false,

            // 重试策略：非生产环境重试3次，生产环境仅网络错误重试
            retry: (failureCount, error) => {
                if (import.meta.env.PROD) {
                    return failureCount < 1 && error.message?.includes('Network');
                }
                return failureCount < 3;
            },

            // 重试延迟：指数退避
            retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
        },
        mutations: {
            // 突变失败后自动失效相关查询，触发重新获取
            onError: (error) => {
                console.error('Mutation error:', error);
            },
        },
        dehydrate: {
            // SSR 时只反序列化成功的数据
            shouldDehydrateQuery: (query) =>
                defaultShouldDehydrateQuery(query) ||
                query.state.status === 'pending',
        },
    },
})
