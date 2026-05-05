import {defineMock} from 'vite-plugin-mock-dev-server'

// 生成随机的 RFC 3339 格式时间戳
function generateRandomTimestamp(daysAgo: number = 0): string {
  const date = new Date(Date.now() - daysAgo * 86400000)
  return date.toISOString()
}

// 树形分类数据
const categories = [
  // 一级分类：技术分享
  {
    id: 1,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 1,
    isNav: true,
    icon: 'carbon:document',
    postCount: 45,
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
      {
        id: 101,
        categoryId: 1,
        languageCode: 'en-US',
        name: 'Tech Sharing',
        slug: 'tech',
        description: 'Share the latest technical articles and tutorials',
        thumbnail: 'https://picsum.photos/400/300?random=1',
        coverImage: 'https://picsum.photos/1200/400?random=1',
        createdAt: generateRandomTimestamp(30),
        updatedAt: generateRandomTimestamp(15),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 0,
    path: '/tech',
    parentId: null,
    children: [],
    createdAt: generateRandomTimestamp(30),
    updatedAt: generateRandomTimestamp(15),
  },
  // 二级分类：前端开发
  {
    id: 11,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 1,
    isNav: false,
    icon: 'carbon:code',
    postCount: 20,
    directPostCount: 20,
    translations: [
      {
        id: 11,
        categoryId: 11,
        languageCode: 'zh-CN',
        name: '前端开发',
        slug: 'frontend',
        description: '前端开发技术和框架',
        thumbnail: 'https://picsum.photos/400/300?random=11',
        coverImage: 'https://picsum.photos/1200/400?random=11',
        createdAt: generateRandomTimestamp(25),
        updatedAt: generateRandomTimestamp(10),
      },
      {
        id: 111,
        categoryId: 11,
        languageCode: 'en-US',
        name: 'Frontend Development',
        slug: 'frontend',
        description: 'Frontend development technologies and frameworks',
        thumbnail: 'https://picsum.photos/400/300?random=11',
        coverImage: 'https://picsum.photos/1200/400?random=11',
        createdAt: generateRandomTimestamp(25),
        updatedAt: generateRandomTimestamp(10),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 1,
    path: '/tech/frontend',
    parentId: 1,
    children: [],
    createdAt: generateRandomTimestamp(25),
    updatedAt: generateRandomTimestamp(10),
  },
  // 二级分类：后端开发
  {
    id: 12,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 2,
    isNav: false,
    icon: 'carbon:cloud',
    postCount: 15,
    directPostCount: 15,
    translations: [
      {
        id: 12,
        categoryId: 12,
        languageCode: 'zh-CN',
        name: '后端开发',
        slug: 'backend',
        description: '后端开发技术和架构',
        thumbnail: 'https://picsum.photos/400/300?random=12',
        coverImage: 'https://picsum.photos/1200/400?random=12',
        createdAt: generateRandomTimestamp(24),
        updatedAt: generateRandomTimestamp(9),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 1,
    path: '/tech/backend',
    parentId: 1,
    children: [],
    createdAt: generateRandomTimestamp(24),
    updatedAt: generateRandomTimestamp(9),
  },
  // 二级分类：移动开发
  {
    id: 13,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 3,
    isNav: false,
    icon: 'carbon:mobile',
    postCount: 10,
    directPostCount: 10,
    translations: [
      {
        id: 13,
        categoryId: 13,
        languageCode: 'zh-CN',
        name: '移动开发',
        slug: 'mobile',
        description: '移动端开发技术',
        thumbnail: 'https://picsum.photos/400/300?random=13',
        coverImage: 'https://picsum.photos/1200/400?random=13',
        createdAt: generateRandomTimestamp(23),
        updatedAt: generateRandomTimestamp(8),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 1,
    path: '/tech/mobile',
    parentId: 1,
    children: [],
    createdAt: generateRandomTimestamp(23),
    updatedAt: generateRandomTimestamp(8),
  },
  // 一级分类：生活随笔
  {
    id: 2,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 2,
    isNav: true,
    icon: 'carbon:blog',
    postCount: 30,
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
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 0,
    path: '/life',
    parentId: null,
    children: [],
    createdAt: generateRandomTimestamp(25),
    updatedAt: generateRandomTimestamp(12),
  },
  // 二级分类：旅行游记
  {
    id: 21,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 1,
    isNav: false,
    icon: 'carbon:map',
    postCount: 10,
    directPostCount: 10,
    translations: [
      {
        id: 21,
        categoryId: 21,
        languageCode: 'zh-CN',
        name: '旅行游记',
        slug: 'travel',
        description: '旅行见闻和游记',
        thumbnail: 'https://picsum.photos/400/300?random=21',
        coverImage: 'https://picsum.photos/1200/400?random=21',
        createdAt: generateRandomTimestamp(20),
        updatedAt: generateRandomTimestamp(7),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 1,
    path: '/life/travel',
    parentId: 2,
    children: [],
    createdAt: generateRandomTimestamp(20),
    updatedAt: generateRandomTimestamp(7),
  },
  // 二级分类：美食探店
  {
    id: 22,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 2,
    isNav: false,
    icon: 'carbon:favorite',
    postCount: 8,
    directPostCount: 8,
    translations: [
      {
        id: 22,
        categoryId: 22,
        languageCode: 'zh-CN',
        name: '美食探店',
        slug: 'food',
        description: '探索城市美食',
        thumbnail: 'https://picsum.photos/400/300?random=22',
        coverImage: 'https://picsum.photos/1200/400?random=22',
        createdAt: generateRandomTimestamp(19),
        updatedAt: generateRandomTimestamp(6),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 1,
    path: '/life/food',
    parentId: 2,
    children: [],
    createdAt: generateRandomTimestamp(19),
    updatedAt: generateRandomTimestamp(6),
  },
  // 一级分类：产品设计
  {
    id: 3,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 3,
    isNav: true,
    icon: 'carbon:chart-line',
    postCount: 25,
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
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 0,
    path: '/design',
    parentId: null,
    children: [],
    createdAt: generateRandomTimestamp(20),
    updatedAt: generateRandomTimestamp(10),
  },
  // 二级分类：UI 设计
  {
    id: 31,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 1,
    isNav: false,
    icon: 'carbon:color-switch',
    postCount: 10,
    directPostCount: 10,
    translations: [
      {
        id: 31,
        categoryId: 31,
        languageCode: 'zh-CN',
        name: 'UI 设计',
        slug: 'ui-design',
        description: '用户界面设计',
        thumbnail: 'https://picsum.photos/400/300?random=31',
        coverImage: 'https://picsum.photos/1200/400?random=31',
        createdAt: generateRandomTimestamp(18),
        updatedAt: generateRandomTimestamp(5),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 1,
    path: '/design/ui-design',
    parentId: 3,
    children: [],
    createdAt: generateRandomTimestamp(18),
    updatedAt: generateRandomTimestamp(5),
  },
  // 二级分类：UX 设计
  {
    id: 32,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 2,
    isNav: false,
    icon: 'carbon:user-profile',
    postCount: 7,
    directPostCount: 7,
    translations: [
      {
        id: 32,
        categoryId: 32,
        languageCode: 'zh-CN',
        name: 'UX 设计',
        slug: 'ux-design',
        description: '用户体验设计',
        thumbnail: 'https://picsum.photos/400/300?random=32',
        coverImage: 'https://picsum.photos/1200/400?random=32',
        createdAt: generateRandomTimestamp(17),
        updatedAt: generateRandomTimestamp(4),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 1,
    path: '/design/ux-design',
    parentId: 3,
    children: [],
    createdAt: generateRandomTimestamp(17),
    updatedAt: generateRandomTimestamp(4),
  },
  // 一级分类：创业思考
  {
    id: 4,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 4,
    isNav: true,
    icon: 'carbon:idea',
    postCount: 18,
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
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 0,
    path: '/startup',
    parentId: null,
    children: [],
    createdAt: generateRandomTimestamp(15),
    updatedAt: generateRandomTimestamp(8),
  },
  // 二级分类：团队管理
  {
    id: 41,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 1,
    isNav: false,
    icon: 'carbon:group',
    postCount: 5,
    directPostCount: 5,
    translations: [
      {
        id: 41,
        categoryId: 41,
        languageCode: 'zh-CN',
        name: '团队管理',
        slug: 'team-management',
        description: '团队建设和管理经验',
        thumbnail: 'https://picsum.photos/400/300?random=41',
        coverImage: 'https://picsum.photos/1200/400?random=41',
        createdAt: generateRandomTimestamp(14),
        updatedAt: generateRandomTimestamp(3),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 1,
    path: '/startup/team-management',
    parentId: 4,
    children: [],
    createdAt: generateRandomTimestamp(14),
    updatedAt: generateRandomTimestamp(3),
  },
  // 二级分类：产品思考
  {
    id: 42,
    status: 'CATEGORY_STATUS_ACTIVE',
    sortOrder: 2,
    isNav: false,
    icon: 'carbon:product',
    postCount: 3,
    directPostCount: 3,
    translations: [
      {
        id: 42,
        categoryId: 42,
        languageCode: 'zh-CN',
        name: '产品思考',
        slug: 'product-thinking',
        description: '产品规划和思考',
        thumbnail: 'https://picsum.photos/400/300?random=42',
        coverImage: 'https://picsum.photos/1200/400?random=42',
        createdAt: generateRandomTimestamp(13),
        updatedAt: generateRandomTimestamp(2),
      },
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    depth: 1,
    path: '/startup/product-thinking',
    parentId: 4,
    children: [],
    createdAt: generateRandomTimestamp(13),
    updatedAt: generateRandomTimestamp(2),
  },
]

// 构建父子关系
function buildCategoryTree() {
  const categoryMap: Record<number, any> = {}

  // 创建映射
  categories.forEach(category => {
    categoryMap[category.id] = category
    category.children = []
  })

  // 构建树形结构
  categories.forEach(category => {
    if (category.parentId) {
      const parent = categoryMap[category.parentId]
      if (parent) {
        parent.children.push(category)
        // 累加文章数量
        parent.postCount += category.postCount
      }
    }
  })

  // 返回根分类（没有父级的分类）
  return categories.filter(c => c.parentId === null)
}

const categoryTree = buildCategoryTree()

const categoryEnMap: Record<string, {
  name: string
  description: string
}> = {
  tech: { name: 'Tech Sharing', description: 'Share the latest technical articles and tutorials' },
  frontend: { name: 'Frontend Development', description: 'Frontend development technologies and frameworks' },
  backend: { name: 'Backend Development', description: 'Backend development technologies and architecture' },
  mobile: { name: 'Mobile Development', description: 'Mobile development technologies' },
  life: { name: 'Life Notes', description: 'Record moments and thoughts from daily life' },
  travel: { name: 'Travel', description: 'Travel experiences and journals' },
  food: { name: 'Food Exploration', description: 'Explore city delicacies' },
  design: { name: 'Product Design', description: 'Product design concepts and practices' },
  uidesign: { name: 'UI Design', description: 'User Interface Design' },
  uxdesign: { name: 'UX Design', description: 'User Experience Design' },
  startup: { name: 'Startup Insights', description: 'Reflections and summaries from startup journey' },
  teammanagement: { name: 'Team Management', description: 'Team building and management experience' },
  productthinking: { name: 'Product Thinking', description: 'Product planning and thinking' },
}

categories.forEach((category) => {
  const zh = category.translations.find(t => t.languageCode === 'zh-CN')
  if (!zh || category.translations.some(t => t.languageCode === 'en-US')) {
    return
  }

  const mapped = categoryEnMap[zh.slug.replace(/-/g, '')] || {
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
      let filteredCategories = categoryTree
      if (queryParams.locale) {
        filteredCategories = categoryTree.map(category => {
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

      // 树形结构不分页，返回所有根分类
      return {
        items: filteredCategories,
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
