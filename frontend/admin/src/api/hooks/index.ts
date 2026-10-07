/**
 * API Hooks 索引文件
 * 导出所有业务模块的 hooks
 */

// 认证相关
export * from './auth';

// 管理门户相关
export * from './admin-portal';

// 用户相关
export * from './user';

// 用户个人资料
export * from './user-profile';

// 租户管理
export * from './tenant';

// 组织人员管理 (OPM)
export * from './org-unit';
export * from './position';

// 权限管理
export * from './permission';
export * from './permission-group';
export * from './role';

// 系统管理
export * from './menu';
export * from './api';
export * from './dict';
export * from './file';
export * from './file-transfer';
export * from './task';
export * from './login-policy';
export * from './language';

// 站点配置（CMS 域）
export * from './site';
export * from './site-setting';
export * from './navigation';
export * from './navigation-item';

// 媒体中心（CMS 域）
export * from './media-asset';

// 互动运营（CMS 域）
export * from './comment';

// 内容管理（CMS 域）
export * from './tag';
export * from './category';
export * from './content-model';
export * from './page';
export * from './post';

// 日志审计
export * from './login-audit-log';
export * from './api-audit-log';
export * from './operation-audit-log';
export * from './data-access-audit-log';
export * from './permission-audit-log';
export * from './policy-evaluation-log';

// 缓存监控

// 在线会话

// 服务监控

// 通知渠道

// 首页分析概览
export * from './dashboard';

// 内部消息
export * from './internal-message';

// 脚本管理

// 脚本执行日志
