/**
 * 静态构建后处理：生成根路径 index.html 重定向到默认语言
 * 这样无论 nginx 如何配置，访问 / 都能正确跳转
 */
import { writeFileSync } from 'node:fs'
import { resolve, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const outputDir = resolve(__dirname, '..', '.output', 'public')

const html = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta http-equiv="refresh" content="0;url=/zh-CN/">
  <link rel="canonical" href="/zh-CN/">
  <title>Redirecting...</title>
</head>
<body>
  <noscript>
    <meta http-equiv="refresh" content="0;url=/zh-CN/">
  </noscript>
  <p>Redirecting to <a href="/zh-CN/">/zh-CN/</a>...</p>
  <script>window.location.replace('/zh-CN/' + window.location.search + window.location.hash)</script>
</body>
</html>
`

writeFileSync(resolve(outputDir, 'index.html'), html)
console.log('✔ Generated .output/public/index.html (redirects to /zh-CN/)')
