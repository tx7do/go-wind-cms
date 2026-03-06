import {defineMock} from 'vite-plugin-mock-dev-server'

// 生成随机的 RFC 3339 格式时间戳
function generateRandomTimestamp(daysAgo: number = 0): string {
  const date = new Date(Date.now() - daysAgo * 86400000)
  return date.toISOString()
}

const categories = [
  {
    id: 1,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 1,
    isNav: true,
    icon: 'carbon:document',
    postCount: 15,
    directPostCount: 15,
    translations: [
      {
        id: 1,
        categoryId: 1,
        languageCode: 'zh-CN',
        name: '技术分享',
        slug: 'tech',
        description: '分享最新的技术文章和教程',
        thumbnail: 'https://picsum.photos/400/300?random=1',
        coverImage: 'https://picsum.photos/1200/400?random=1',
        createdAt: generateRandomTimestamp(30),
        updatedAt: generateRandomTimestamp(15),
      },
    ],
    availableLanguages: ['zh-CN'],
    createdBy: 1,
    children: [],
    depth: 0,
    path: '/tech',
    createdAt: generateRandomTimestamp(30),
    updatedAt: generateRandomTimestamp(15),
  },
  {
    id: 2,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 2,
    isNav: true,
    icon: 'carbon:blog',
    postCount: 12,
    directPostCount: 12,
    translations: [
      {
        id: 2,
        categoryId: 2,
        languageCode: 'zh-CN',
        name: '生活随笔',
        slug: 'life',
        description: '记录生活中的点点滴滴',
        thumbnail: 'https://picsum.photos/400/300?random=2',
        coverImage: 'https://picsum.photos/1200/400?random=2',
        createdAt: generateRandomTimestamp(25),
        updatedAt: generateRandomTimestamp(12),
      },
    ],
    availableLanguages: ['zh-CN'],
    createdBy: 1,
    children: [],
    depth: 0,
    path: '/life',
    createdAt: generateRandomTimestamp(25),
    updatedAt: generateRandomTimestamp(12),
  },
  {
    id: 3,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 3,
    isNav: true,
    icon: 'carbon:chart-line',
    postCount: 8,
    directPostCount: 8,
    translations: [
      {
        id: 3,
        categoryId: 3,
        languageCode: 'zh-CN',
        name: '产品设计',
        slug: 'design',
        description: '产品设计理念与实践',
        thumbnail: 'https://picsum.photos/400/300?random=3',
        coverImage: 'https://picsum.photos/1200/400?random=3',
        createdAt: generateRandomTimestamp(20),
        updatedAt: generateRandomTimestamp(10),
      },
    ],
    availableLanguages: ['zh-CN'],
    createdBy: 1,
    children: [],
    depth: 0,
    path: '/design',
    createdAt: generateRandomTimestamp(20),
    updatedAt: generateRandomTimestamp(10),
  },
  {
    id: 4,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 4,
    isNav: true,
    icon: 'carbon:idea',
    postCount: 10,
    directPostCount: 10,
    translations: [
      {
        id: 4,
        categoryId: 4,
        languageCode: 'zh-CN',
        name: '创业思考',
        slug: 'startup',
        description: '创业路上的思考与总结',
        thumbnail: 'https://picsum.photos/400/300?random=4',
        coverImage: 'https://picsum.photos/1200/400?random=4',
        createdAt: generateRandomTimestamp(15),
        updatedAt: generateRandomTimestamp(8),
      },
    ],
    availableLanguages: ['zh-CN'],
    createdBy: 1,
    children: [],
    depth: 0,
    path: '/startup',
    createdAt: generateRandomTimestamp(15),
    updatedAt: generateRandomTimestamp(8),
  },
]

const categoryEnMap: Record<string, {
  name: string
  description: string
}> = {
  tech: {
    name: 'Tech Sharing',
    description: 'Share the latest technical articles and tutorials',
  },
  life: {
    name: 'Life Notes',
    description: 'Record moments and thoughts from daily life',
  },
  design: {
    name: 'Product Design',
    description: 'Product design concepts and practices',
  },
  startup: {
    name: 'Startup Insights',
    description: 'Reflections and summaries from startup journey',
  },
}

categories.forEach((category) => {
  const zh = category.translations.find(t => t.languageCode === 'zh-CN')
  if (!zh || category.translations.some(t => t.languageCode === 'en-US')) {
    return
  }

  const mapped = categoryEnMap[zh.slug] || {
    name: zh.name,
    description: zh.description,
  }

  category.translations.push({
    id: zh.id + 100,
    categoryId: category.id,
    languageCode: 'en-US',
    name: mapped.name,
    slug: zh.slug,
    description: mapped.description,
    thumbnail: zh.thumbnail,
    coverImage: zh.coverImage,
    createdAt: zh.createdAt,
    updatedAt: zh.updatedAt,
  })

  category.availableLanguages = ['zh-CN', 'en-US']
})

export default defineMock([
  {
    url: '/app/v1/categories',
    method: 'GET',
    body: ({query}: any) => {
      const page = parseInt(query.page || '1')
      const pageSize = parseInt(query.pageSize || '10')

      // 解析 query 参数
      let queryParams: any = {}
      if (query.query) {
        try {
          queryParams = typeof query.query === 'string'
            ? JSON.parse(query.query)
            : query.query
        } catch (e) {
          console.error('Failed to parse query:', e)
        }
      }

      // 根据 locale 过滤分类数据
      let filteredCategories = categories
      if (queryParams.locale) {
        filteredCategories = categories.map(category => {
          const filteredTranslations = category.translations.filter(
            t => t.languageCode === queryParams.locale
          )

          if (filteredTranslations.length > 0) {
            return {
              ...category,
              translations: filteredTranslations
            }
          }

          return {
            ...category,
            translations: [category.translations[0]]
          }
        })
      }

      const start = (page - 1) * pageSize
      const end = start + pageSize

      return {
        items: filteredCategories.slice(start, end),
        total: filteredCategories.length,
      }
    },
  },
  {
    url: '/app/v1/categories/:id',
    method: 'GET',
    body: ({params, query}: any) => {
      const id = parseInt(params.id)
      let category = categories.find(c => c.id === id)

      if (!category) {
        return {error: 'Category not found'}
      }

      // ✅ 根据 languageCode 参数过滤翻译数据
      if (query.languageCode) {
        const filteredTranslations = category.translations.filter(
          t => t.languageCode === query.languageCode
        )

        category = {
          ...category,
          translations: filteredTranslations.length > 0
            ? filteredTranslations
            : [category.translations[0]]
        }
      }

      return category
    },
  },
])
