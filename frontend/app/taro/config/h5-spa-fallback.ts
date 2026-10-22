import path from 'path'
import type { IncomingMessage } from 'http'
import type { Plugin } from 'vite'

/**
 * 仅 H5 dev 生效：把「浏览器顶层导航」统一回落 index.html。
 *
 * Taro H5 browser 路由下，未配 customRoutes 的页面其浏览器路径就是 Taro 内部路径
 * （/pages/index/index）。Vite 的 transformMiddleware 排在 SPA history fallback 之前，
 * 且把无扩展名 URL 当作候选模块去 root（src/）下解析，于是命中
 * src/pages/index/index.tsx，dev server 回一份 application/javascript，
 * 浏览器把编译产物当正文打印出来。
 */
export function h5SpaFallback(): Plugin {
  return {
    name: 'taro-h5-spa-fallback',
    enforce: 'pre',
    configureServer(server) {
      server.middlewares.use((req: IncomingMessage, _res, next) => {
        if (isHtmlNavigation(req)) req.url = '/index.html'
        next()
      })
    },
  }
}

function isHtmlNavigation(req: IncomingMessage): boolean {
  const method = req.method ?? 'GET'
  if (method !== 'GET' && method !== 'HEAD') return false
  // 模块请求的 Accept 是 */*，只有文档导航才带 text/html
  if (!req.headers.accept?.includes('text/html')) return false

  const pathname = (req.url ?? '/').split('?')[0]
  if (pathname === '/' || pathname === '/index.html') return false
  if (pathname.startsWith('/@') || pathname.startsWith('/node_modules')) return false
  return path.extname(pathname) === ''
}
