import createMiddleware from 'next-intl/middleware';

import {routing} from '@/app/[locale]/routing';

// next-intl 请求处理函数：由 src/proxy.ts（Next 16 入口约定）挂载，
// matcher 配置在入口文件内静态声明（Next 要求 config 不可 re-export）。
export default createMiddleware(routing);
