import createMiddleware from 'next-intl/middleware';

import {routing} from '@/app/[locale]/routing';

export default createMiddleware(routing);

export const config = {
    // 只匹配国际化路径
    matcher: ['/', '/(zh-CN|en-US)/:path*'],
};
