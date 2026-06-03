import {QueryClient, defaultShouldDehydrateQuery} from '@tanstack/react-query';

// 单独导出，方便测试 & 维护
export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            // 数据新鲜度：5分钟内不重复请求
            staleTime: 5 * 60 * 1000,

            // 缓存时间：30分钟后清除
            gcTime: 30 * 60 * 1000,

            // 窗口聚焦时自动重新验证
            refetchOnWindowFocus: false,

            // 重试策略：非生产环境重试3次，生产环境仅网络错误重试
            retry: (failureCount, error) => {
                if (process.env.NODE_ENV === 'production') {
                    return failureCount < 1 && !!error.message?.includes('Network');
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
});
