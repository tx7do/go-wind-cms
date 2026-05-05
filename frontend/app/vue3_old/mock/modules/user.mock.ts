import { defineMock } from 'vite-plugin-mock-dev-server'

export default defineMock([
  {
    url: '/app/v1/me',
    method: 'GET',
    body: () => ({
      id: '1',
      username: 'admin',
      nickname: '系统管理员',
      realname: '张三',
      email: 'admin@example.com',
      mobile: '13800138000',
      avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=admin',
      description: '系统超级管理员，拥有所有权限。热爱技术，致力于打造优秀的内容管理平台。',
      gender: 'MALE',
      region: '北京市朝阳区',
      tenantName: '总公司',
      tenantId: '1',
      status: 'NORMAL',
      roleNames: ['超级管理员', '系统管理员'],
      positionNames: ['系统管理员', '技术总监', 'CTO'],
      followers: 196,
      following: 1,
      postCount: 2,
      likeCount: 947,
      createdAt: new Date(Date.now() - 86400000 * 365),
      updatedAt: new Date(Date.now() - 3600000),
      lastLoginAt: new Date(Date.now() - 1800000),
    }),
  },
])

