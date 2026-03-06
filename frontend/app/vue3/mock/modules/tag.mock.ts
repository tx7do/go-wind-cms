import { defineMock } from 'vite-plugin-mock-dev-server'

// 生成随机的 RFC 3339 格式时间戳
function generateRandomTimestamp(daysAgo: number = 0): string {
  const date = new Date(Date.now() - daysAgo * 86400000)
  return date.toISOString()
}

const tags = [
  {
    id: 1,
    status: 'TAG_STATUS_ACTIVE',
    color: '#3B82F6',
    icon: 'logo-vue',
    sortOrder: 1,
    isFeatured: true,
    postCount: 5,
    translations: [
      {
        id: 1,
        tagId: 1,
        languageCode: 'zh-CN',
        name: 'Vue.js',
        slug: 'vuejs',
        description: 'Vue.js 相关内容',
        seoTitle: 'Vue.js 教程和最佳实践',
        metaDescription: '学习 Vue.js 框架的最新教程、最佳实践和开发经验分享',
        metaKeywords: 'Vue.js,Vue3,前端框架',
        createdAt: generateRandomTimestamp(10),
        updatedAt: generateRandomTimestamp(8),
      },
      {
        id: 11,
        tagId: 1,
        languageCode: 'en-US',
        name: 'Vue.js',
        slug: 'vuejs',
        description: 'Vue.js related content',
        seoTitle: 'Vue.js Tutorial and Best Practices',
        metaDescription: 'Learn the latest Vue.js framework tutorials, best practices and development experience sharing',
        metaKeywords: 'Vue.js,Vue3,frontend framework',
        createdAt: generateRandomTimestamp(10),
        updatedAt: generateRandomTimestamp(8),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(10),
    updatedAt: generateRandomTimestamp(8),
  },
  {
    id: 2,
    status: 'TAG_STATUS_ACTIVE',
    color: '#10B981',
    icon: 'code',
    sortOrder: 2,
    isFeatured: true,
    postCount: 8,
    translations: [
      {
        id: 2,
        tagId: 2,
        languageCode: 'zh-CN',
        name: '前端开发',
        slug: 'frontend',
        description: '前端开发技术',
        seoTitle: '前端开发教程和实战',
        metaDescription: '前端开发技术博客，包含 HTML、CSS、JavaScript、框架等内容',
        metaKeywords: '前端,开发,HTML,CSS,JavaScript',
        createdAt: generateRandomTimestamp(9),
        updatedAt: generateRandomTimestamp(7),
      },
      {
        id: 12,
        tagId: 2,
        languageCode: 'en-US',
        name: 'Frontend Development',
        slug: 'frontend',
        description: 'Frontend development technology',
        seoTitle: 'Frontend Development Tutorial and Practice',
        metaDescription: 'Frontend development tech blog covering HTML, CSS, JavaScript, frameworks and more',
        metaKeywords: 'Frontend,Development,HTML,CSS,JavaScript',
        createdAt: generateRandomTimestamp(9),
        updatedAt: generateRandomTimestamp(7),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(9),
    updatedAt: generateRandomTimestamp(7),
  },
  {
    id: 3,
    status: 'TAG_STATUS_ACTIVE',
    color: '#F59E0B',
    icon: 'logo-typescript',
    sortOrder: 3,
    isFeatured: true,
    postCount: 6,
    translations: [
      {
        id: 3,
        tagId: 3,
        languageCode: 'zh-CN',
        name: 'TypeScript',
        slug: 'typescript',
        description: 'TypeScript 编程语言',
        seoTitle: 'TypeScript 完全指南',
        metaDescription: 'TypeScript 编程语言教程和最佳实践，提升代码质量',
        metaKeywords: 'TypeScript,类型系统,JavaScript',
        createdAt: generateRandomTimestamp(8),
        updatedAt: generateRandomTimestamp(6),
      },
      {
        id: 13,
        tagId: 3,
        languageCode: 'en-US',
        name: 'TypeScript',
        slug: 'typescript',
        description: 'TypeScript programming language',
        seoTitle: 'TypeScript Complete Guide',
        metaDescription: 'TypeScript programming language tutorials and best practices to improve code quality',
        metaKeywords: 'TypeScript,Type System,JavaScript',
        createdAt: generateRandomTimestamp(8),
        updatedAt: generateRandomTimestamp(6),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(8),
    updatedAt: generateRandomTimestamp(6),
  },
  {
    id: 4,
    status: 'TAG_STATUS_ACTIVE',
    color: '#EF4444',
    icon: 'content-management',
    sortOrder: 4,
    isFeatured: false,
    postCount: 4,
    translations: [
      {
        id: 4,
        tagId: 4,
        languageCode: 'zh-CN',
        name: 'CMS',
        slug: 'cms',
        description: '内容管理系统',
        seoTitle: 'CMS 系统设计和实现',
        metaDescription: '内容管理系统的设计原理、架构和实现方案',
        metaKeywords: 'CMS,内容管理,系统设计',
        createdAt: generateRandomTimestamp(7),
        updatedAt: generateRandomTimestamp(5),
      },
      {
        id: 14,
        tagId: 4,
        languageCode: 'en-US',
        name: 'CMS',
        slug: 'cms',
        description: 'Content Management System',
        seoTitle: 'CMS System Design and Implementation',
        metaDescription: 'Design principles, architecture, and implementation solutions for content management systems',
        metaKeywords: 'CMS,content management,system design',
        createdAt: generateRandomTimestamp(7),
        updatedAt: generateRandomTimestamp(5),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(7),
    updatedAt: generateRandomTimestamp(5),
  },
  {
    id: 5,
    status: 'TAG_STATUS_ACTIVE',
    color: '#8B5CF6',
    icon: 'architecture',
    sortOrder: 5,
    isFeatured: false,
    postCount: 3,
    translations: [
      {
        id: 5,
        tagId: 5,
        languageCode: 'zh-CN',
        name: '架构设计',
        slug: 'architecture',
        description: '软件架构设计',
        seoTitle: '软件架构设计最佳实践',
        metaDescription: '系统架构设计、分布式系统、微服务等架构相关知识',
        metaKeywords: '架构,设计,分布式,微服务',
        createdAt: generateRandomTimestamp(6),
        updatedAt: generateRandomTimestamp(4),
      },
      {
        id: 15,
        tagId: 5,
        languageCode: 'en-US',
        name: 'Architecture Design',
        slug: 'architecture',
        description: 'Software architecture design',
        seoTitle: 'Best Practices for Software Architecture Design',
        metaDescription: 'Architecture-related knowledge including system design, distributed systems, and microservices',
        metaKeywords: 'architecture,design,distributed,microservices',
        createdAt: generateRandomTimestamp(6),
        updatedAt: generateRandomTimestamp(4),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(6),
    updatedAt: generateRandomTimestamp(4),
  },
  {
    id: 6,
    status: 'TAG_STATUS_ACTIVE',
    color: '#EC4899',
    icon: 'user-circle',
    sortOrder: 6,
    isFeatured: false,
    postCount: 7,
    translations: [
      {
        id: 6,
        tagId: 6,
        languageCode: 'zh-CN',
        name: '生活',
        slug: 'life',
        description: '生活相关',
        seoTitle: '生活随笔和思考',
        metaDescription: '分享生活中的所思所想，包含个人随笔和生活感悟',
        metaKeywords: '生活,随笔,感悟,思考',
        createdAt: generateRandomTimestamp(5),
        updatedAt: generateRandomTimestamp(3),
      },
      {
        id: 16,
        tagId: 6,
        languageCode: 'en-US',
        name: 'Life',
        slug: 'life',
        description: 'Life related content',
        seoTitle: 'Life Essays and Reflections',
        metaDescription: 'Sharing thoughts and reflections from daily life, including personal essays and insights',
        metaKeywords: 'life,essay,insight,reflection',
        createdAt: generateRandomTimestamp(5),
        updatedAt: generateRandomTimestamp(3),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(5),
    updatedAt: generateRandomTimestamp(3),
  },
  {
    id: 7,
    status: 'TAG_STATUS_ACTIVE',
    color: '#06B6D4',
    icon: 'palette',
    sortOrder: 7,
    isFeatured: false,
    postCount: 5,
    translations: [
      {
        id: 7,
        tagId: 7,
        languageCode: 'zh-CN',
        name: 'UI设计',
        slug: 'ui-design',
        description: 'UI 界面设计',
        seoTitle: 'UI 设计原则和最佳实践',
        metaDescription: 'UI 界面设计教程，包含设计原则、色彩搭配、交互设计等',
        metaKeywords: 'UI设计,界面设计,用户界面',
        createdAt: generateRandomTimestamp(4),
        updatedAt: generateRandomTimestamp(2),
      },
      {
        id: 17,
        tagId: 7,
        languageCode: 'en-US',
        name: 'UI Design',
        slug: 'ui-design',
        description: 'UI interface design',
        seoTitle: 'UI Design Principles and Best Practices',
        metaDescription: 'UI design tutorials covering principles, color matching, and interaction design',
        metaKeywords: 'UI design,interface design,user interface',
        createdAt: generateRandomTimestamp(4),
        updatedAt: generateRandomTimestamp(2),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(4),
    updatedAt: generateRandomTimestamp(2),
  },
  {
    id: 8,
    status: 'TAG_STATUS_ACTIVE',
    color: '#14B8A6',
    icon: 'cursor-1',
    sortOrder: 8,
    isFeatured: false,
    postCount: 4,
    translations: [
      {
        id: 8,
        tagId: 8,
        languageCode: 'zh-CN',
        name: 'UX',
        slug: 'ux',
        description: '用户体验',
        seoTitle: 'UX 用户体验设计',
        metaDescription: '用户体验设计、用户研究、交互设计和可用性测试',
        metaKeywords: 'UX,用户体验,交互设计,可用性',
        createdAt: generateRandomTimestamp(3),
        updatedAt: generateRandomTimestamp(1),
      },
      {
        id: 18,
        tagId: 8,
        languageCode: 'en-US',
        name: 'UX',
        slug: 'ux',
        description: 'User Experience',
        seoTitle: 'UX User Experience Design',
        metaDescription: 'User experience design, user research, interaction design, and usability testing',
        metaKeywords: 'UX,user experience,interaction design,usability',
        createdAt: generateRandomTimestamp(3),
        updatedAt: generateRandomTimestamp(1),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(3),
    updatedAt: generateRandomTimestamp(1),
  },
  {
    id: 9,
    status: 'TAG_STATUS_ACTIVE',
    color: '#F97316',
    icon: 'rocket',
    sortOrder: 9,
    isFeatured: true,
    postCount: 6,
    translations: [
      {
        id: 9,
        tagId: 9,
        languageCode: 'zh-CN',
        name: '创业',
        slug: 'startup',
        description: '创业相关',
        seoTitle: '创业经验和创业思考',
        metaDescription: '分享创业过程中的经验、教训和思考，包含创业故事和建议',
        metaKeywords: '创业,企业,创业故事,创业建议',
        createdAt: generateRandomTimestamp(2),
        updatedAt: generateRandomTimestamp(0),
      },
      {
        id: 19,
        tagId: 9,
        languageCode: 'en-US',
        name: 'Startup',
        slug: 'startup',
        description: 'Startup related content',
        seoTitle: 'Startup Experience and Reflections',
        metaDescription: 'Sharing startup experiences, lessons learned, and thoughts, including stories and advice',
        metaKeywords: 'startup,business,startup story,startup advice',
        createdAt: generateRandomTimestamp(2),
        updatedAt: generateRandomTimestamp(0),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(2),
    updatedAt: generateRandomTimestamp(0),
  },
  {
    id: 10,
    status: 'TAG_STATUS_ACTIVE',
    color: '#0EA5E9',
    icon: 'team',
    sortOrder: 10,
    isFeatured: false,
    postCount: 3,
    translations: [
      {
        id: 10,
        tagId: 10,
        languageCode: 'zh-CN',
        name: '团队管理',
        slug: 'team-management',
        description: '团队建设与管理',
        seoTitle: '团队管理和领导力发展',
        metaDescription: '团队建设、管理方法、领导力培养和团队协作最佳实践',
        metaKeywords: '团队管理,领导力,团队协作,管理',
        createdAt: generateRandomTimestamp(1),
        updatedAt: generateRandomTimestamp(0),
      },
      {
        id: 20,
        tagId: 10,
        languageCode: 'en-US',
        name: 'Team Management',
        slug: 'team-management',
        description: 'Team building and management',
        seoTitle: 'Team Management and Leadership Development',
        metaDescription: 'Best practices for team building, management methods, leadership development, and collaboration',
        metaKeywords: 'team management,leadership,team collaboration,management',
        createdAt: generateRandomTimestamp(1),
        updatedAt: generateRandomTimestamp(0),
      }
    ],
    availableLanguages: ['zh-CN', 'en-US'],
    createdBy: 1,
    createdAt: generateRandomTimestamp(1),
    updatedAt: generateRandomTimestamp(0),
  },
]

export default defineMock([
  {
    url: '/app/v1/tags',
    method: 'GET',
    body: ({ query }: any) => {
      const page = parseInt(query.page || '1')
      const pageSize = parseInt(query.pageSize || '10')

      // 解析 query 参数（可能是 JSON 字符串）
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

      // 根据 locale 过滤标签数据
      let filteredTags = tags
      if (queryParams.locale) {
        filteredTags = tags.map(tag => {
          // 过滤 translations，只保留匹配 locale 的翻译
          const filteredTranslations = tag.translations.filter(
            t => t.languageCode === queryParams.locale
          )

          // 如果有匹配的翻译，返回带有过滤后翻译的标签
          if (filteredTranslations.length > 0) {
            return {
              ...tag,
              translations: filteredTranslations
            }
          }

          // 如果没有匹配的翻译，返回第一个翻译（回退）
          return {
            ...tag,
            translations: [tag.translations[0]]
          }
        })
      }

      const start = (page - 1) * pageSize
      const end = start + pageSize

      return {
        items: filteredTags.slice(start, end),
        total: filteredTags.length,
      }
    },
  },
  {
    url: '/app/v1/tags/:id',
    method: 'GET',
    body: ({ params, query }: any) => {
      const id = parseInt(params.id)
      let tag = tags.find(t => t.id === id)

      if (!tag) {
        return { error: 'Tag not found' }
      }

      // ✅ 根据 languageCode 参数过滤翻译数据
      if (query.languageCode) {
        const filteredTranslations = tag.translations.filter(
          t => t.languageCode === query.languageCode
        )

        tag = {
          ...tag,
          translations: filteredTranslations.length > 0
            ? filteredTranslations
            : [tag.translations[0]]
        }
      }

      return tag
    },
  },
])

