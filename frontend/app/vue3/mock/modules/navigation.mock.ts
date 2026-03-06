import { defineMock } from 'vite-plugin-mock-dev-server'

// 生成随机的 RFC 3339 格式时间戳
function generateRandomTimestamp(daysAgo: number = 0, hoursOffset: number = 0): string {
  const date = new Date(Date.now() - daysAgo * 86400000 - hoursOffset * 3600000)
  return date.toISOString()
}

const navigations = [
  {
    id: 1,
    name: 'Main Navigation',
    location: 'header',
    isActive: true,
    locale: 'zh-CN',
    items: [
      {
        id: 1,
        navigationId: 1,
        title: '首页',
        url: '/',
        icon: 'home',
        description: '返回首页',
        linkType: 'LINK_TYPE_CUSTOM',
        sortOrder: 1,
        isOpenNewTab: false,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [],
        createdBy: 1,
        createdAt: generateRandomTimestamp(30),
        updatedAt: generateRandomTimestamp(15),
      },
      {
        id: 2,
        navigationId: 1,
        title: '文章',
        url: '/post',
        icon: 'document',
        description: '浏览所有文章',
        linkType: 'LINK_TYPE_CUSTOM',
        sortOrder: 2,
        isOpenNewTab: false,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [],
        createdBy: 1,
        createdAt: generateRandomTimestamp(30),
        updatedAt: generateRandomTimestamp(14),
      },
      {
        id: 3,
        navigationId: 1,
        title: '分类',
        url: '/category',
        icon: 'folder',
        description: '浏览所有分类',
        linkType: 'LINK_TYPE_CUSTOM',
        sortOrder: 3,
        isOpenNewTab: false,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [
          {
            id: 4,
            navigationId: 1,
            title: '技术分享',
            url: '/category/1',
            icon: 'code',
            description: '技术文章分类',
            linkType: 'LINK_TYPE_CATEGORY',
            objectId: 1,
            sortOrder: 1,
            isOpenNewTab: false,
            isInvalid: false,
            cssClass: '',
            requiredPermission: '',
            parentId: 3,
            children: [],
            createdBy: 1,
            createdAt: generateRandomTimestamp(29),
            updatedAt: generateRandomTimestamp(13),
          },
          {
            id: 5,
            navigationId: 1,
            title: '生活随笔',
            url: '/category/2',
            icon: 'blog',
            description: '生活文章分类',
            linkType: 'LINK_TYPE_CATEGORY',
            objectId: 2,
            sortOrder: 2,
            isOpenNewTab: false,
            isInvalid: false,
            cssClass: '',
            requiredPermission: '',
            parentId: 3,
            children: [],
            createdBy: 1,
            createdAt: generateRandomTimestamp(28),
            updatedAt: generateRandomTimestamp(12),
          },
        ],
        createdBy: 1,
        createdAt: generateRandomTimestamp(30),
        updatedAt: generateRandomTimestamp(13),
      },
      {
        id: 6,
        navigationId: 1,
        title: '关于',
        url: '/about',
        icon: 'information',
        description: '关于我们',
        linkType: 'LINK_TYPE_PAGE',
        objectId: 1,
        sortOrder: 4,
        isOpenNewTab: false,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [],
        createdBy: 1,
        createdAt: generateRandomTimestamp(27),
        updatedAt: generateRandomTimestamp(11),
      },
    ],
    createdBy: 1,
    createdAt: generateRandomTimestamp(30),
    updatedAt: generateRandomTimestamp(12),
  },
  {
    id: 2,
    name: 'Footer Navigation',
    location: 'footer',
    isActive: true,
    locale: 'zh-CN',
    items: [
      {
        id: 7,
        navigationId: 2,
        title: '联系我们',
        url: '/contact',
        icon: 'email',
        description: '联系我们',
        linkType: 'LINK_TYPE_PAGE',
        objectId: 2,
        sortOrder: 1,
        isOpenNewTab: false,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [],
        createdBy: 1,
        createdAt: generateRandomTimestamp(25),
        updatedAt: generateRandomTimestamp(10),
      },
      {
        id: 8,
        navigationId: 2,
        title: '隐私政策',
        url: '/privacy',
        icon: 'shield-checkmark',
        description: '隐私政策',
        linkType: 'LINK_TYPE_PAGE',
        objectId: 3,
        sortOrder: 2,
        isOpenNewTab: false,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [],
        createdBy: 1,
        createdAt: generateRandomTimestamp(24),
        updatedAt: generateRandomTimestamp(9),
      },
      {
        id: 9,
        navigationId: 2,
        title: '服务条款',
        url: '/terms',
        icon: 'document',
        description: '服务条款',
        linkType: 'LINK_TYPE_PAGE',
        objectId: 4,
        sortOrder: 3,
        isOpenNewTab: false,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [],
        createdBy: 1,
        createdAt: generateRandomTimestamp(23),
        updatedAt: generateRandomTimestamp(8),
      },
      {
        id: 10,
        navigationId: 2,
        title: 'GitHub',
        url: 'https://github.com',
        icon: 'logo-github',
        description: '访问我们的GitHub',
        linkType: 'LINK_TYPE_EXTERNAL',
        sortOrder: 4,
        isOpenNewTab: true,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [],
        createdBy: 1,
        createdAt: generateRandomTimestamp(22),
        updatedAt: generateRandomTimestamp(7),
      },
    ],
    createdBy: 1,
    createdAt: generateRandomTimestamp(25),
    updatedAt: generateRandomTimestamp(8),
  },
  {
    id: 3,
    name: 'Sidebar Navigation',
    location: 'sidebar',
    isActive: true,
    locale: 'zh-CN',
    items: [
      {
        id: 11,
        navigationId: 3,
        title: '热门标签',
        url: '/tag',
        icon: 'pricetag',
        description: '浏览热门标签',
        linkType: 'LINK_TYPE_CUSTOM',
        sortOrder: 1,
        isOpenNewTab: false,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [],
        createdBy: 1,
        createdAt: generateRandomTimestamp(20),
        updatedAt: generateRandomTimestamp(6),
      },
      {
        id: 12,
        navigationId: 3,
        title: '归档',
        url: '/archive',
        icon: 'archive',
        description: '文章归档',
        linkType: 'LINK_TYPE_CUSTOM',
        sortOrder: 2,
        isOpenNewTab: false,
        isInvalid: false,
        cssClass: '',
        requiredPermission: '',
        parentId: 0,
        children: [],
        createdBy: 1,
        createdAt: generateRandomTimestamp(19),
        updatedAt: generateRandomTimestamp(5),
      },
    ],
    createdBy: 1,
    createdAt: generateRandomTimestamp(20),
    updatedAt: generateRandomTimestamp(4),
  },
]

const navigationNameEnMap: Record<string, string> = {
  'Main Navigation': 'Main Navigation',
  'Footer Navigation': 'Footer Navigation',
  'Sidebar Navigation': 'Sidebar Navigation',
}

const navigationItemTitleEnMap: Record<string, string> = {
  '首页': 'Home',
  '文章': 'Posts',
  '分类': 'Categories',
  '技术分享': 'Tech Sharing',
  '生活随笔': 'Life Notes',
  '关于': 'About',
  '联系我们': 'Contact',
  '隐私政策': 'Privacy Policy',
  '服务条款': 'Terms of Service',
  '热门标签': 'Popular Tags',
  '归档': 'Archive',
}

const navigationItemDescEnMap: Record<string, string> = {
  '返回首页': 'Back to homepage',
  '浏览所有文章': 'Browse all posts',
  '浏览所有分类': 'Browse all categories',
  '技术文章分类': 'Tech article category',
  '生活文章分类': 'Life article category',
  '关于我们': 'About us',
  '联系我们': 'Contact us',
  '隐私政策': 'Privacy policy',
  '服务条款': 'Terms of service',
  '浏览热门标签': 'Browse popular tags',
  '文章归档': 'Post archive',
}

const navigationItemDescEnById: Record<number, string> = {
  10: 'Visit our GitHub',
}

function toEnglishItems(items: any[], navigationId: number): any[] {
  return items.map(item => ({
    ...item,
    id: item.id + 100,
    navigationId,
    parentId: item.parentId === 0 ? 0 : item.parentId + 100,
    title: navigationItemTitleEnMap[item.title] || item.title,
    description: navigationItemDescEnById[item.id] || navigationItemDescEnMap[item.description] || item.description,
    children: toEnglishItems(item.children || [], navigationId),
  }))
}

const enNavigations = navigations.map(nav => {
  const newNavigationId = nav.id + 100
  return {
    ...nav,
    id: newNavigationId,
    name: navigationNameEnMap[nav.name] || nav.name,
    locale: 'en-US',
    items: toEnglishItems(nav.items || [], newNavigationId),
  }
})

navigations.push(...enNavigations)

export default defineMock([
  {
    url: '/app/v1/navigations',
    method: 'GET',
    body: ({ query }) => {
      const page = Number(query.page) || 1
      const pageSize = Number(query.pageSize) || 10
      const noPaging = query.noPaging === 'true'

      let items = [...navigations]

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

      // 根据 location 过滤
      if (queryParams.location) {
        items = items.filter(nav => nav.location === queryParams.location)
      }

      // 根据 locale 过滤
      if (queryParams.locale) {
        items = items.filter(nav => nav.locale === queryParams.locale)
      }

      // 根据 isActive 过滤
      if (queryParams.isActive !== undefined) {
        items = items.filter(nav => nav.isActive === queryParams.isActive)
      }

      const total = items.length

      if (noPaging) {
        return {
          items,
          total,
        }
      }

      const start = (page - 1) * pageSize
      const end = start + pageSize
      const paginatedItems = items.slice(start, end)


      return {
        items: paginatedItems,
        total,
      }
    },
  },
  {
    url: '/app/v1/navigations/:id',
    method: 'GET',
    body: ({ params }) => {
      const id = Number(params.id)
      const navigation = navigations.find(nav => nav.id === id)

      if (!navigation) {
        return {
          code: 404,
          message: 'Navigation not found',
        }
      }

      return navigation
    },
  },
  {
    url: '/navigations',
    method: 'POST',
    body: ({ body }) => {
      const newId = Math.max(...navigations.map(nav => nav.id)) + 1
      const now = generateRandomTimestamp(0, 0)

      const newNavigation = {
        id: newId,
        name: body.data.name || '',
        location: body.data.location || '',
        isActive: body.data.isActive !== undefined ? body.data.isActive : true,
        locale: body.data.locale || 'zh-CN',
        items: body.data.items || [],
        createdBy: body.data.createdBy || 1,
        createdAt: now,
        updatedAt: now,
      }

      navigations.push(newNavigation)

      return newNavigation
    },
  },
  {
    url: '/navigations/:id',
    method: 'PUT',
    body: ({ params, body }) => {
      const id = Number(params.id)
      const index = navigations.findIndex(nav => nav.id === id)

      if (index === -1) {
        return {
          code: 404,
          message: 'Navigation not found',
        }
      }

      const now = generateRandomTimestamp(0, 0)

      navigations[index] = {
        ...navigations[index],
        ...body.data,
        id,
        updatedAt: now,
      }

      return navigations[index]
    },
  },
  {
    url: '/navigations/:id',
    method: 'DELETE',
    body: ({ params }) => {
      const id = Number(params.id)
      const index = navigations.findIndex(nav => nav.id === id)

      if (index === -1) {
        return {
          code: 404,
          message: 'Navigation not found',
        }
      }

      navigations.splice(index, 1)

      return {}
    },
  },
])

