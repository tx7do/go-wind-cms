import {defineMock} from 'vite-plugin-mock-dev-server'

// 生成随机的 RFC 3339 格式时间戳
function generateRandomTimestamp(daysAgo: number = 0, hoursOffset: number = 0): string {
  const date = new Date(Date.now() - daysAgo * 86400000 - hoursOffset * 3600000)
  return date.toISOString()
}

// 树形评论数据 - 支持嵌套回复
const comments = [
  // ========== 文章 ID: 1 的评论 ==========
  {
    id: 1,
    contentType: 'CONTENT_TYPE_POST',
    objectId: 1,
    content: '非常好的文章，学到了很多！',
    authorId: 0,
    authorName: '王五',
    authorEmail: 'wangwu@example.com',
    authorUrl: '',
    authorType: 'AUTHOR_TYPE_GUEST',
    status: 'STATUS_APPROVED',
    likeCount: 5,
    dislikeCount: 0,
    replyCount: 2,
    ipAddress: '192.168.1.1',
    location: '北京',
    parentId: 0,
    children: [
      {
        id: 7,
        contentType: 'CONTENT_TYPE_POST',
        objectId: 1,
        content: '确实很不错，尤其是关于 Composition API 的部分',
        authorId: 0,
        authorName: '张三',
        authorEmail: 'zhangsan@example.com',
        authorUrl: '',
        authorType: 'AUTHOR_TYPE_GUEST',
        status: 'STATUS_APPROVED',
        likeCount: 2,
        dislikeCount: 0,
        replyCount: 0,
        ipAddress: '192.168.1.7',
        location: '广州',
        parentId: 1,
        replyToId: 1,
        children: [],
        createdAt: generateRandomTimestamp(2, 3),
        updatedAt: generateRandomTimestamp(2, 3),
      },
      {
        id: 8,
        contentType: 'CONTENT_TYPE_POST',
        objectId: 1,
        content: '@王五 同意！博主写得很用心',
        authorId: 0,
        authorName: '李四',
        authorEmail: 'lisi@example.com',
        authorUrl: '',
        authorType: 'AUTHOR_TYPE_GUEST',
        status: 'STATUS_APPROVED',
        likeCount: 1,
        dislikeCount: 0,
        replyCount: 1,
        ipAddress: '192.168.1.8',
        location: '武汉',
        parentId: 1,
        replyToId: 1,
        children: [
          {
            id: 11,
            contentType: 'CONTENT_TYPE_POST',
            objectId: 1,
            content: '@李四 谢谢支持！会继续努力的',
            authorId: 1,
            authorName: '博主',
            authorEmail: 'admin@example.com',
            authorUrl: 'https://example.com',
            authorType: 'AUTHOR_TYPE_ADMIN',
            status: 'STATUS_APPROVED',
            likeCount: 3,
            dislikeCount: 0,
            replyCount: 0,
            ipAddress: '192.168.1.100',
            location: '北京',
            parentId: 8,
            replyToId: 8,
            children: [],
            createdAt: generateRandomTimestamp(1, 10),
            updatedAt: generateRandomTimestamp(1, 10),
          },
        ],
        createdAt: generateRandomTimestamp(2, 2),
        updatedAt: generateRandomTimestamp(2, 2),
      },
    ],
    createdAt: generateRandomTimestamp(2, 4),
    updatedAt: generateRandomTimestamp(2, 4),
  },
  {
    id: 2,
    contentType: 'CONTENT_TYPE_POST',
    objectId: 1,
    content: 'Composition API 确实很强大，感谢分享！',
    authorId: 0,
    authorName: '赵六',
    authorEmail: 'zhaoliu@example.com',
    authorUrl: '',
    authorType: 'AUTHOR_TYPE_GUEST',
    status: 'STATUS_APPROVED',
    likeCount: 3,
    dislikeCount: 0,
    replyCount: 1,
    ipAddress: '192.168.1.2',
    location: '上海',
    parentId: 0,
    children: [
      {
        id: 9,
        contentType: 'CONTENT_TYPE_POST',
        objectId: 1,
        content: '@赵六 对，比 Options API 灵活多了',
        authorId: 0,
        authorName: '孙七',
        authorEmail: 'sunqi@example.com',
        authorUrl: '',
        authorType: 'AUTHOR_TYPE_GUEST',
        status: 'STATUS_APPROVED',
        likeCount: 1,
        dislikeCount: 0,
        replyCount: 0,
        ipAddress: '192.168.1.9',
        location: '南京',
        parentId: 2,
        replyToId: 2,
        children: [],
        createdAt: generateRandomTimestamp(1, 12),
        updatedAt: generateRandomTimestamp(1, 12),
      },
    ],
    createdAt: generateRandomTimestamp(2, 2),
    updatedAt: generateRandomTimestamp(2, 2),
  },
  {
    id: 3,
    contentType: 'CONTENT_TYPE_POST',
    objectId: 1,
    content: '期待更多关于 Vue 3 的文章',
    authorId: 0,
    authorName: '孙七',
    authorEmail: 'sunqi@example.com',
    authorUrl: '',
    authorType: 'AUTHOR_TYPE_GUEST',
    status: 'STATUS_APPROVED',
    likeCount: 2,
    dislikeCount: 0,
    replyCount: 0,
    ipAddress: '192.168.1.3',
    location: '深圳',
    parentId: 0,
    children: [],
    createdAt: generateRandomTimestamp(1, 12),
    updatedAt: generateRandomTimestamp(1, 12),
  },

  // ========== 文章 ID: 2 的评论 ==========
  {
    id: 4,
    contentType: 'CONTENT_TYPE_POST',
    objectId: 2,
    content: 'TypeScript 的类型系统真的很强大',
    authorId: 0,
    authorName: '周八',
    authorEmail: 'zhouba@example.com',
    authorUrl: '',
    authorType: 'AUTHOR_TYPE_GUEST',
    status: 'STATUS_APPROVED',
    likeCount: 8,
    dislikeCount: 0,
    replyCount: 2,
    ipAddress: '192.168.1.4',
    location: '杭州',
    parentId: 0,
    children: [
      {
        id: 10,
        contentType: 'CONTENT_TYPE_POST',
        objectId: 2,
        content: '条件类型的部分能再详细讲讲吗？',
        authorId: 0,
        authorName: '吴九',
        authorEmail: 'wujiu@example.com',
        authorUrl: '',
        authorType: 'AUTHOR_TYPE_GUEST',
        status: 'STATUS_APPROVED',
        likeCount: 1,
        dislikeCount: 0,
        replyCount: 1,
        ipAddress: '192.168.1.5',
        location: '成都',
        parentId: 4,
        replyToId: 4,
        children: [
          {
            id: 12,
            contentType: 'CONTENT_TYPE_POST',
            objectId: 2,
            content: '@吴九 好的，下期专门讲讲条件类型和分布式条件类型',
            authorId: 1,
            authorName: '博主',
            authorEmail: 'admin@example.com',
            authorUrl: 'https://example.com',
            authorType: 'AUTHOR_TYPE_ADMIN',
            status: 'STATUS_APPROVED',
            likeCount: 4,
            dislikeCount: 0,
            replyCount: 0,
            ipAddress: '192.168.1.100',
            location: '北京',
            parentId: 10,
            replyToId: 10,
            children: [],
            createdAt: generateRandomTimestamp(3, 5),
            updatedAt: generateRandomTimestamp(3, 5),
          },
        ],
        createdAt: generateRandomTimestamp(4, 8),
        updatedAt: generateRandomTimestamp(4, 8),
      },
      {
        id: 13,
        contentType: 'CONTENT_TYPE_POST',
        objectId: 2,
        content: '泛型约束也很有用',
        authorId: 0,
        authorName: '郑十',
        authorEmail: 'zhengshi@example.com',
        authorUrl: '',
        authorType: 'AUTHOR_TYPE_GUEST',
        status: 'STATUS_APPROVED',
        likeCount: 2,
        dislikeCount: 0,
        replyCount: 0,
        ipAddress: '192.168.1.10',
        location: '西安',
        parentId: 4,
        replyToId: 4,
        children: [],
        createdAt: generateRandomTimestamp(3, 10),
        updatedAt: generateRandomTimestamp(3, 10),
      },
    ],
    createdAt: generateRandomTimestamp(5, 6),
    updatedAt: generateRandomTimestamp(5, 6),
  },
  {
    id: 5,
    contentType: 'CONTENT_TYPE_POST',
    objectId: 2,
    content: '推倒类型这个概念很有意思',
    authorId: 0,
    authorName: '冯十一',
    authorEmail: 'fengshiyi@example.com',
    authorUrl: '',
    authorType: 'AUTHOR_TYPE_GUEST',
    status: 'STATUS_APPROVED',
    likeCount: 3,
    dislikeCount: 0,
    replyCount: 0,
    ipAddress: '192.168.1.11',
    location: '重庆',
    parentId: 0,
    children: [],
    createdAt: generateRandomTimestamp(4, 15),
    updatedAt: generateRandomTimestamp(4, 15),
  },

  // ========== 文章 ID: 3 的评论 ==========
  {
    id: 6,
    contentType: 'CONTENT_TYPE_POST',
    objectId: 3,
    content: 'Content Hub 是未来的趋势',
    authorId: 0,
    authorName: '郑十',
    authorEmail: 'zhengshi@example.com',
    authorUrl: '',
    authorType: 'AUTHOR_TYPE_GUEST',
    status: 'STATUS_APPROVED',
    likeCount: 6,
    dislikeCount: 0,
    replyCount: 3,
    ipAddress: '192.168.1.6',
    location: '西安',
    parentId: 0,
    children: [
      {
        id: 14,
        contentType: 'CONTENT_TYPE_POST',
        objectId: 3,
        content: '确实，前后端分离更灵活',
        authorId: 0,
        authorName: '钱十二',
        authorEmail: 'qianshier@example.com',
        authorUrl: '',
        authorType: 'AUTHOR_TYPE_GUEST',
        status: 'STATUS_APPROVED',
        likeCount: 2,
        dislikeCount: 0,
        replyCount: 0,
        ipAddress: '192.168.1.12',
        location: '苏州',
        parentId: 6,
        replyToId: 6,
        children: [],
        createdAt: generateRandomTimestamp(6, 8),
        updatedAt: generateRandomTimestamp(6, 8),
      },
      {
        id: 15,
        contentType: 'CONTENT_TYPE_POST',
        objectId: 3,
        content: 'Strapi 和 Contentful 哪个更好用？',
        authorId: 0,
        authorName: '孔十三',
        authorEmail: 'kongshisan@example.com',
        authorUrl: '',
        authorType: 'AUTHOR_TYPE_GUEST',
        status: 'STATUS_APPROVED',
        likeCount: 1,
        dislikeCount: 0,
        replyCount: 1,
        ipAddress: '192.168.1.13',
        location: '长沙',
        parentId: 6,
        replyToId: 6,
        children: [
          {
            id: 16,
            contentType: 'CONTENT_TYPE_POST',
            objectId: 3,
            content: '@孔十三 看需求，Strapi 开源免费，Contentful 功能更强大',
            authorId: 1,
            authorName: '博主',
            authorEmail: 'admin@example.com',
            authorUrl: 'https://example.com',
            authorType: 'AUTHOR_TYPE_ADMIN',
            status: 'STATUS_APPROVED',
            likeCount: 5,
            dislikeCount: 0,
            replyCount: 0,
            ipAddress: '192.168.1.100',
            location: '北京',
            parentId: 15,
            replyToId: 15,
            children: [],
            createdAt: generateRandomTimestamp(6, 5),
            updatedAt: generateRandomTimestamp(6, 5),
          },
        ],
        createdAt: generateRandomTimestamp(6, 10),
        updatedAt: generateRandomTimestamp(6, 10),
      },
      {
        id: 17,
        contentType: 'CONTENT_TYPE_POST',
        objectId: 3,
        content: 'GoWind Content Hub 也很不错',
        authorId: 0,
        authorName: '白十四',
        authorEmail: 'baisishi@example.com',
        authorUrl: '',
        authorType: 'AUTHOR_TYPE_GUEST',
        status: 'STATUS_APPROVED',
        likeCount: 4,
        dislikeCount: 0,
        replyCount: 0,
        ipAddress: '192.168.1.14',
        location: '郑州',
        parentId: 6,
        replyToId: 6,
        children: [],
        createdAt: generateRandomTimestamp(5, 20),
        updatedAt: generateRandomTimestamp(5, 20),
      },
    ],
    createdAt: generateRandomTimestamp(7, 3),
    updatedAt: generateRandomTimestamp(7, 3),
  },
]

// 获取树形评论数据 (按 objectId 筛选)
function getTreeComments(objectId: number): any[] {
  const result: any[] = []

  for (const comment of comments) {
    if (comment.objectId === objectId && comment.parentId === 0) {
      result.push({...comment})
    }
  }

  return result
}

export default defineMock([
  {
    url: '/app/v1/comments',
    method: 'GET',
    body: ({query}: any) => {
      const page = parseInt(query.page || '1')
      const pageSize = parseInt(query.pageSize || '10')

      // 获取树形评论数据
      let treeComments: any[] = []

      // 按对象 ID 和内容类型筛选
      if (query.query) {
        try {
          const queryObj = JSON.parse(query.query)
          if (queryObj.objectId) {
            treeComments = getTreeComments(parseInt(queryObj.objectId))
          }
          if (queryObj.contentType) {
            treeComments = treeComments.filter(c => c.contentType === queryObj.contentType)
          }
          if (queryObj.status) {
            treeComments = treeComments.filter(c => c.status === queryObj.status)
          }
        } catch (e) {
          console.error('Parse query error:', e)
        }
      } else {
        // 没有 query 时返回所有根评论
        treeComments = getTreeComments(0)
      }

      const start = (page - 1) * pageSize
      const end = start + pageSize

      return {
        items: treeComments.slice(start, end),
        total: treeComments.length,
      }
    },
  },
  {
    url: '/app/v1/comments',
    method: 'POST',
    body: ({body}: any) => {
      const newComment = {
        id: comments.length + 1,
        ...body.data,
        parentId: body.data.parentId || 0,
        children: [],
        createdAt: generateRandomTimestamp(0, 0),
        updatedAt: generateRandomTimestamp(0, 0),
      }

      // 如果是回复评论，添加到父评论的 children
      if (body.data.parentId && body.data.parentId > 0) {
        const parentComment = findCommentById(comments, body.data.parentId)
        if (parentComment) {
          parentComment.children.push(newComment)
          parentComment.replyCount = (parentComment.replyCount || 0) + 1
        }
      } else {
        comments.push(newComment)
      }

      return newComment
    },
  },
])

// 辅助函数：递归查找评论
function findCommentById(commentsList: any[], id: number): any | null {
  for (const comment of commentsList) {
    if (comment.id === id) {
      return comment
    }
    if (comment.children && comment.children.length > 0) {
      const found = findCommentById(comment.children, id)
      if (found) return found
    }
  }
  return null
}
