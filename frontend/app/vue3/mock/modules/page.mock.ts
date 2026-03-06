import { defineMock } from 'vite-plugin-mock-dev-server'

// 生成随机的 RFC 3339 格式时间戳
function generateRandomTimestamp(daysAgo: number = 0): string {
  const date = new Date(Date.now() - daysAgo * 86400000)
  return date.toISOString()
}

const pages = [
  {
    id: 1,
    status: 'PAGE_STATUS_PUBLISHED',
    type: 'PAGE_TYPE_DEFAULT',
    editorType: 'EDITOR_TYPE_MARKDOWN',
    slug: 'about',
    authorId: 1,
    authorName: '张三',
    showInNavigation: true,
    sortOrder: 1,
    visits: 456,
    customFields: {},
    translations: [
      {
        id: 1,
        pageId: 1,
        languageCode: 'zh-CN',
        title: '关于我们',
        slug: 'about',
        summary: '了解我们的团队和使命',
        content: `# 关于 GoWind CMS

GoWind CMS 是一个强大的多租户 Headless CMS 平台，为现代团队提供灵活的内容管理解决方案。

## 我们的使命

让内容管理变得更简单、更高效、更灵活。

## 核心特性

### 多租户架构
完整的租户隔离与共享基础设施，支持无限数量的独立项目。

### 灵活的内容管理
强大的工具用于管理所有类型的内容，支持多语言、自定义字段、版本控制等。

### 企业级安全
- 基于角色的访问控制
- 数据加密
- 审计日志
- 定期安全审查

### 高级分析
深入了解您的内容表现，包括访问统计、用户行为分析等。

## 联系我们

如果您有任何问题或建议，欢迎联系我们：

- 邮箱：support@gowindcms.com
- 电话：+86 123 4567 8900
- 地址：中国北京市朝阳区xxx路xx号`,
        thumbnail: 'https://picsum.photos/800/450?random=100',
        fullPath: '/about',
        wordCount: 200,
      },
    ],
    availableLanguages: ['zh-CN'],
    createdBy: 1,
    children: [],
    depth: 0,
    path: '/about',
    createdAt: generateRandomTimestamp(30),
    updatedAt: generateRandomTimestamp(10),
  },
]

const pageEnMap: Record<string, {
  title: string
  summary: string
  content: string
  seoTitle?: string
  metaDescription?: string
}> = {
  about: {
    title: 'About Us',
    summary: 'Learn about our team and mission',
    content: `# About GoWind CMS

GoWind CMS is a powerful multi-tenant headless CMS platform that provides flexible content management for modern teams.

## Our Mission

Make content management simpler, more efficient, and more flexible.

## Core Features

### Multi-tenant Architecture
Complete tenant isolation with shared infrastructure, supporting unlimited independent projects.

### Flexible Content Management
Powerful tools to manage all content types, including multilingual support, custom fields, and version control.

### Enterprise-grade Security
- Role-based access control
- Data encryption
- Audit logs
- Regular security reviews

### Advanced Analytics
Gain insights into content performance, including traffic statistics and user behavior analysis.

## Contact Us

If you have any questions or suggestions, feel free to contact us:

- Email: support@gowindcms.com
- Phone: +86 123 4567 8900
- Address: Chaoyang District, Beijing, China`,
  },
}

pages.forEach((page) => {
  const zh = page.translations.find(t => t.languageCode === 'zh-CN')
  if (!zh || page.translations.some(t => t.languageCode === 'en-US')) {
    return
  }

  const mapped = pageEnMap[zh.slug] || {
    title: zh.title,
    summary: zh.summary,
    content: zh.content,
  }

  page.translations.push({
    id: zh.id + 100,
    pageId: page.id,
    languageCode: 'en-US',
    title: mapped.title,
    slug: zh.slug,
    summary: mapped.summary,
    content: mapped.content,
    thumbnail: zh.thumbnail,
    fullPath: zh.fullPath,
    wordCount: zh.wordCount,
  })

  page.availableLanguages = ['zh-CN', 'en-US']
})

export default defineMock([
  {
    url: '/app/v1/pages',
    method: 'GET',
    body: ({ query }: any) => {
      const page = parseInt(query.page || '1')
      const pageSize = parseInt(query.pageSize || '10')

      let filteredPages = [...pages]

      // 按 slug 筛选
      if (query.query) {
        try {
          const queryObj = JSON.parse(query.query)
          if (queryObj.slug) {
            filteredPages = filteredPages.filter(p =>
              p.translations[0].slug === queryObj.slug
            )
          }
          if (queryObj.type) {
            filteredPages = filteredPages.filter(p => p.type === queryObj.type)
          }

          // 根据 locale 过滤翻译数据
          if (queryObj.locale) {
            filteredPages = filteredPages.map(page => {
              const filteredTranslations = page.translations.filter(
                t => t.languageCode === queryObj.locale
              )

              if (filteredTranslations.length > 0) {
                return {
                  ...page,
                  translations: filteredTranslations
                }
              }

              return {
                ...page,
                translations: [page.translations[0]]
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
        items: filteredPages.slice(start, end),
        total: filteredPages.length,
      }
    },
  },
  {
    url: '/app/v1/pages/:id',
    method: 'GET',
    body: ({ params, query }: any) => {
      const id = parseInt(params.id)
      let page = pages.find(p => p.id === id)

      if (!page) {
        return { error: 'Page not found' }
      }

      // ✅ 根据 languageCode 参数过滤翻译数据
      if (query.languageCode) {
        const filteredTranslations = page.translations.filter(
          t => t.languageCode === query.languageCode
        )

        page = {
          ...page,
          translations: filteredTranslations.length > 0
            ? filteredTranslations
            : [page.translations[0]]
        }
      }

      return page
    },
  },
])
