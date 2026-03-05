import {defineMock} from 'vite-plugin-mock-dev-server'

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
      },
    ],
    availableLanguages: ['zh-CN'],
    customFields: {},
    children: [],
    depth: 0,
    path: '/tech',
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-03-01T00:00:00Z',
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
      },
    ],
    availableLanguages: ['zh-CN'],
    customFields: {},
    children: [],
    depth: 0,
    path: '/life',
    createdAt: '2026-01-02T00:00:00Z',
    updatedAt: '2026-03-02T00:00:00Z',
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
      },
    ],
    availableLanguages: ['zh-CN'],
    customFields: {},
    children: [],
    depth: 0,
    path: '/design',
    createdAt: '2026-01-03T00:00:00Z',
    updatedAt: '2026-03-03T00:00:00Z',
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
      },
    ],
    availableLanguages: ['zh-CN'],
    customFields: {},
    children: [],
    depth: 0,
    path: '/startup',
    createdAt: '2026-01-04T00:00:00Z',
    updatedAt: '2026-03-04T00:00:00Z',
  },
]

export default defineMock([
  {
    url: '/app/v1/categories',
    method: 'GET',
    body: ({query}: any) => {
      const page = parseInt(query.page || '1')
      const pageSize = parseInt(query.pageSize || '10')
      const start = (page - 1) * pageSize
      const end = start + pageSize

      return {
        items: categories.slice(start, end),
        total: categories.length,
      }
    },
  },
  {
    url: '/app/v1/categories/:id',
    method: 'GET',
    body: ({params}: any) => {
      const id = parseInt(params.id)
      const category = categories.find(c => c.id === id)
      return category || {error: 'Category not found'}
    },
  },
])

