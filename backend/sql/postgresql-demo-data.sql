BEGIN;

SET LOCAL search_path = public, pg_catalog;

-- 一次性清理相关表并重置自增（包含外键依赖）
TRUNCATE TABLE public.sys_org_units,
               public.sys_positions,
               public.sys_tasks,
               public.sys_login_policies,
               public.sys_dict_types,
               public.sys_dict_entries,
               public.internal_message_categories
RESTART IDENTITY CASCADE;

-- 测试租户
INSERT INTO public.sys_tenants(id, name, code, type, audit_status, status, admin_user_id, created_at)
VALUES (1, '测试租户', 'super', 'PAID', 'APPROVED', 'ON', 2, now())
;
SELECT setval('sys_tenants_id_seq', (SELECT MAX(id) FROM sys_tenants));

-- 插入租户管理员用户
INSERT INTO public.sys_users (id, tenant_id, username, nickname, realname, email, gender, created_at)
VALUES
    -- 2. 租户管理员（TENANT_ADMIN）
    (2, 1, 'tenant_admin', '租户管理', '张管理员', 'tenant@company.com', 'MALE', now())
;
SELECT setval('sys_users_id_seq', (SELECT MAX(id) FROM sys_users));

-- 插入4个用户的凭证（密码统一为admin，哈希值与原admin一致，方便测试）
INSERT INTO public.sys_user_credentials (user_id, identity_type, identifier, credential_type, credential, status,
                                         is_primary, created_at)
VALUES
    -- 租户管理员（对应users表id=2）
    (2, 'USERNAME', 'tenant_admin', 'PASSWORD_HASH', '$2a$10$yajZDX20Y40FkG0Bu4N19eXNqRizez/S9fK63.JxGkfLq.RoNKR/a', 'ENABLED', true, now()),
    (2, 'EMAIL', 'tenant@company.com', 'PASSWORD_HASH', '$2a$10$yajZDX20Y40FkG0Bu4N19eXNqRizez/S9fK63.JxGkfLq.RoNKR/a', 'ENABLED', false, now())
;
SELECT setval('sys_user_credentials_id_seq', (SELECT MAX(id) FROM sys_user_credentials));

-- 组织架构单元
INSERT INTO public.sys_org_units (id, tenant_id, parent_id, type, name, code, description, path, sort_order, leader_id, status, created_at)
VALUES
    (1, 1, NULL, 'COMPANY', 'XX集团总部', 'HEADQUARTERS', '集团核心管理机构，统筹全集团战略规划、业务管控及资源调配', '/1', 1, 1, 'ON', now()),
    (2, 1, 1, 'DIVISION', '技术部', 'TECH', '负责集团整体技术架构规划、研发管理、系统运维及技术创新', '/1/2', 2, 5, 'ON', now()),
    (3, 1, 1, 'DIVISION', '财务部', 'FIN', '负责集团财务核算、资金管理、税务筹划、预算编制及财务风控', '/1/3', 3, 8, 'ON', now()),
    (4, 1, 1, 'DIVISION', '人事部', 'HR', '负责人力资源规划、招聘配置、薪酬绩效、员工培训及组织发展', '/1/4', 4, 9, 'ON', now()),
    (5, 1, 2, 'DEPARTMENT', '研发一部', 'DEV-1', '聚焦新能源领域产品研发、技术迭代及核心模块开发', '/1/2/5', 1, 6, 'ON', now()),
    (6, 1, 1, 'REGION', '华北大区', 'NORTH', '负责华北区域市场运营、客户维护、销售管理及本地化服务落地', '/1/6', 3, 12, 'ON', now()),
    (7, 1, 1, 'SUBSIDIARY', '广州分公司', 'GZ', '负责华南区域（广州及周边）业务拓展、客户服务及本地化运营', '/1/7', 5, 2, 'ON', now()),
    (8, 1, 1, 'SUBSIDIARY', '深圳子公司', 'SZ', '负责深圳区域市场开拓、科技创新业务落地及高端客户对接', '/1/8', 6, 4, 'ON', now()),
    (9, 1, 1, 'DIVISION', '销售部', 'SALES', '统筹集团整体销售策略制定、销售团队管理及业绩目标达成', '/1/9', 7, 16, 'ON', now()),
    (10, 1, 9, 'DEPARTMENT', '海外事业部', 'INTL', '负责海外市场拓展、国际客户合作、跨境业务管理及本地化运营', '/1/9/10', 1, 17, 'ON', now()),
    (11, 1, 10, 'TEAM', '海外销售组', 'INTL-SALES-1', '具体执行海外市场销售任务，跟进客户需求及订单落地', '/1/9/10/11', 1, 18, 'ON', now()),
    (12, 1, 5, 'PROJECT', '新能源项目组', 'NEO-PROJ', '专项负责新能源项目的研发、落地、运营及成果转化', '/1/2/5/12', 1, 6, 'ON', now()),
    (13, 1, 1, 'COMMITTEE', '审计委员会', 'AUDIT', '独立开展集团内部审计、风控检查、合规监督及问题整改跟进', '/1/13', 8, 12, 'ON', now()),
    (14, 1, 1, 'DEPARTMENT', '客服部', 'CS', '负责全集团客户咨询、投诉处理、售后服务及客户满意度提升', '/1/14', 9, 11, 'ON', now()),
    (15, 1, 14, 'TEAM', '客服一组', 'CS-1', '承接华南区域客户服务、售后问题处理及客户关系维护', '/1/14/15', 1, 20, 'ON', now())
;

-- 岗位数据
INSERT INTO public.sys_positions (id, tenant_id, type, name, code, org_unit_id, reports_to_position_id, description, job_family, job_grade, level, headcount, is_key_position, status, sort_order, created_at)
VALUES
    (1, 1, 'LEADER', '技术总监', 'TECH-DIRECTOR-001', 2, NULL, '负责公司整体技术战略规划、团队管理及核心技术决策', 'TECH', 1, 1, 1, true, 'ON', 1, now()),
    (2, 1, 'MANAGER', '技术部经理', 'TECH-MANAGER-001', 2, 1, '负责技术部日常管理、项目排期及团队协作', 'TECH', 2, 2, 1, true, 'ON', 2, now()),
    (3, 1, 'MANAGER', '前端主管', 'TECH-FE-LEADER-001', 2, 2, '负责前端团队开发管理、技术方案评审及需求落地', 'TECH', 3, 3, 3, false, 'ON', 3, now()),
    (4, 1, 'MANAGER', '后端主管', 'TECH-BE-LEADER-001', 2, 2, '负责后端服务架构设计、数据库优化及接口开发管理', 'TECH', 4, 3, 3, false, 'ON', 4, now()),
    (5, 1, 'REGULAR', '前端开发专员', 'TECH-FE-DEV-001', 2, 3, '负责Web/移动端前端页面开发、交互实现及兼容性优化', 'TECH', 5, 4, 5, false, 'ON', 5, now()),
    (6, 1, 'REGULAR', '后端开发专员', 'TECH-BE-DEV-001', 2, 4, '负责后端接口开发、业务逻辑实现及系统稳定性维护', 'TECH', 6, 4, 5, false, 'ON', 6, now()),
    (7, 1, 'REGULAR', '测试工程师', 'TECH-TEST-001', 2, 2, '负责项目功能测试、性能测试及自动化测试脚本开发', 'TECH', 3, 4, 3, false, 'ON', 7, now()),
    (8, 1, 'LEADER', '人力总监', 'HR-DIRECTOR-001', 2, NULL, '负责人力资源战略规划、组织架构设计及人才梯队建设', 'HR', 1, 1, 1, true, 'ON', 1, now()),
    (9, 1, 'MANAGER', '招聘主管', 'HR-RECRUIT-LEADER-001', 2, 8, '负责公司各部门招聘需求对接、简历筛选及面试安排', 'HR', 2, 2, 1, false, 'ON', 2, now()),
    (10, 1, 'REGULAR', '薪酬绩效专员', 'HR-C&P-001', 2, 8, '负责员工薪酬核算、绩效考核制度落地及社保公积金管理', 'HR', 3, 2, 1, false, 'ON', 3, now()),
    (11, 1, 'REGULAR', 'HRBP', 'HR-BP-001', 2, 8, '对接业务部门，提供人力资源支持（入离职、员工关系等）', 'HR', 4, 2, 1, false, 'ON', 4, now()),
    (12, 1, 'LEADER', '财务总监', 'FIN-DIRECTOR-001', 2, NULL, '负责公司财务战略、预算管理及财务风险控制', 'FIN', 1, 1, 1, true, 'ON', 1, now()),
    (13, 1, 'MANAGER', '会计主管', 'FIN-ACCOUNT-LEADER-001', 2, 12, '负责账务处理、财务报表编制及税务申报管理', 'FIN', 2, 2, 1, false, 'ON', 2, now()),
    (14, 1, 'REGULAR', '出纳专员', 'FIN-CASHIER-001', 2, 13, '负责日常资金收付、银行对账及票据管理', 'FIN', 3, 3, 1, false, 'ON', 3, now()),
    (15, 1, 'REGULAR', '成本会计', 'FIN-COST-001', 2, 13, '负责成本核算、成本分析及成本控制方案制定', 'FIN', 4, 3, 1, false, 'ON', 4, now()),
    (16, 1, 'LEADER', '市场总监', 'MKT-DIRECTOR-001', 4, NULL, '负责市场战略规划、品牌建设及营销活动策划', 'MKT', 1, 1, 1, true, 'ON', 1, now()),
    (17, 1, 'MANAGER', '新媒体运营主管', 'MKT-NEWS-LEADER-001', 4, 16, '负责新媒体平台内容运营及用户增长', 'MKT', 2, 2, 1, false, 'ON', 2, now()),
    (18, 1, 'REGULAR', '活动策划专员', 'MKT-EVENT-001', 4, 16, '负责线下活动策划、执行及效果复盘', 'MKT', 3, 3, 1, false, 'ON', 3, now()),
    (19, 1, 'REGULAR', '市场调研专员', 'MKT-RESEARCH-001', 4, 16, '负责行业动态调研、竞品分析及市场趋势报告撰写', 'MKT', 4, 3, 1, false, 'ON', 4, now()),
    (20, 1, 'REGULAR', '行政助理', 'ADMIN-ASSIST-001', 2, 8, '负责办公用品采购、会议安排等行政工作（已合并至HRBP）', 'ADMIN', 5, 5, 1, false, 'OFF', 5, now())
;
SELECT setval('sys_positions_id_seq', (SELECT COALESCE(MAX(id), 1) FROM sys_positions));

-- 用户-租户关联关系
INSERT INTO public.sys_memberships (id, tenant_id, user_id, org_unit_id, position_id, role_id, is_primary, status)
VALUES
    -- 租户管理员（TENANT_ADMIN）
    (2, 1, 2, null, null, 2, true, 'ACTIVE')
;
SELECT setval('sys_memberships_id_seq', (SELECT MAX(id) FROM sys_memberships));

-- 租户成员-角色关联关系
INSERT INTO sys_membership_roles (id, membership_id, tenant_id, role_id, is_primary, status)
VALUES
    -- 租户管理员（TENANT_ADMIN）
    (2, 2, 1, 2, true, 'ACTIVE')
;
SELECT setval('sys_membership_roles_id_seq', (SELECT MAX(id) FROM sys_membership_roles));

-- 调度任务
INSERT INTO public.sys_tasks(type, type_name, task_payload, cron_spec, enable, created_at)
VALUES
    ('PERIODIC', 'backup', '{ "name": "test"}', '0 * * * *', true, now())
;
SELECT setval('sys_tasks_id_seq', (SELECT MAX(id) FROM sys_tasks));

-- 登录策略
INSERT INTO public.sys_login_policies(id, target_id, type, method, value, reason, created_at)
VALUES
    (1, 1, 'BLACKLIST', 'IP', '127.0.0.1', '无理由', now()),
    (2, 1, 'WHITELIST', 'MAC', '00:1B:44:11:3A:B7 ', '无理由', now())
;
SELECT setval('sys_login_policies_id_seq', (SELECT MAX(id) FROM sys_login_policies));

-- 插入字典类型
INSERT INTO public.sys_dict_types (
    id, type_code, sort_order, is_enabled, created_at, updated_at
) VALUES
      (1, 'USER_STATUS', 10, true, now(), now()),
      (2, 'DEVICE_TYPE', 20, true, now(), now()),
      (3, 'ORDER_STATUS', 30, true, now(), now()),
      (4, 'GENDER', 40, true, now(), now()),
      (5, 'PAYMENT_METHOD', 50, true, now(), now())
;
SELECT setval('sys_dict_types_id_seq', (SELECT MAX(id) FROM sys_dict_types));

-- 插入字典类型国际化（zh-CN）
INSERT INTO public.sys_dict_type_i18n (
    type_id, language_code, type_name, description, tenant_id, created_at, updated_at
) VALUES
      (1, 'zh-CN', '用户状态', '系统用户的状态管理，包括正常、冻结、注销', 0, now(), now()),
      (2, 'zh-CN', '设备类型', 'IoT平台接入的设备品类，新增需同步至设备接入模块', 0, now(), now()),
      (3, 'zh-CN', '订单状态', '电商订单的全生命周期状态', 0, now(), now()),
      (4, 'zh-CN', '性别', '用户性别枚举，默认未知', 0, now(), now()),
      (5, 'zh-CN', '支付方式', '支持的支付渠道，含第三方支付和自有渠道', 0, now(), now()),

      (1, 'en-US', 'User Status', 'System user status management, including normal, frozen, and canceled', 0, now(), now()),
      (2, 'en-US', 'Device Type', 'IoT device categories connected to the platform; new types must be synchronized with the device access module', 0, now(), now()),
      (3, 'en-US', 'Order Status', 'Full lifecycle statuses for e-commerce orders', 0, now(), now()),
      (4, 'en-US', 'Gender', 'User gender enumeration, defaulting to unknown', 0, now(), now()),
      (5, 'en-US', 'Payment Method', 'Supported payment channels, including third-party and proprietary options', 0, now(), now())
;
SELECT setval('sys_dict_type_i18n_id_seq', (SELECT MAX(id) FROM sys_dict_type_i18n));

-- 插入字典条目
INSERT INTO public.sys_dict_entries (
    id, type_id, entry_value, numeric_value, sort_order, is_enabled, created_at, updated_at, tenant_id
) VALUES
      -- 用户状态
      (1, 1, 'NORMAL', 1, 1, true, now(), now(), 0),
      (2, 1, 'FROZEN', 2, 2, true, now(), now(), 0),
      (3, 1, 'CANCELED', 3, 3, true, now(), now(), 0),
      -- 设备类型
      (4, 2, 'TEMP_SENSOR', 101, 1, true, now(), now(), 0),
      (5, 2, 'CURRENT_METER', 102, 2, true, now(), now(), 0),
      (6, 2, 'GAS_DETECTOR', 103, 3, false, now(), now(), 0),
      -- 订单状态
      (7, 3, 'PENDING', 1, 1, true, now(), now(), 0),
      (8, 3, 'PAID', 2, 2, true, now(), now(), 0),
      (9, 3, 'SHIPPED', 3, 3, true, now(), now(), 0),
      (10, 3, 'COMPLETED', 4, 4, true, now(), now(), 0),
      (11, 3, 'CANCELED', 5, 5, true, now(), now(), 0),
      -- 性别
      (12, 4, 'MALE', 1, 1, true, now(), now(), 0),
      (13, 4, 'FEMALE', 2, 2, true, now(), now(), 0),
      (14, 4, 'UNKNOWN', 0, 3, true, now(), now(), 0),
      -- 支付方式
      (15, 5, 'ALIPAY', 1, 1, true, now(), now(), 0),
      (16, 5, 'WECHAT', 2, 2, true, now(), now(), 0),
      (17, 5, 'UNIONPAY', 3, 3, true, now(), now(), 0),
      (18, 5, 'CASH', 4, 4, false, now(), now(), 0)
;
SELECT setval('sys_dict_entries_id_seq', (SELECT MAX(id) FROM sys_dict_entries));

-- 插入字典条目国际化（zh-CN）
INSERT INTO public.sys_dict_entry_i18n (
    entry_id, language_code, entry_label, description, sort_order, tenant_id, created_at, updated_at
) VALUES
      -- 用户状态
      (1, 'zh-CN', '正常', '用户可正常登录和操作', 1, 0, now(), now()),
      (2, 'zh-CN', '冻结', '因违规被临时冻结，需管理员解冻', 2, 0, now(), now()),
      (3, 'zh-CN', '注销', '用户主动注销，数据保留但不可登录', 3, 0, now(), now()),
      -- 设备类型
      (4, 'zh-CN', '温湿度传感器', '支持温度（-20~80℃）和湿度（0~100%RH）采集', 1, 0, now(), now()),
      (5, 'zh-CN', '电流仪表', '交流/直流电流测量，精度0.5级', 2, 0, now(), now()),
      (6, 'zh-CN', '气体探测器', '暂不支持，待硬件适配（2025Q4计划启用）', 3, 0, now(), now()),
      -- 订单状态
      (7, 'zh-CN', '待支付', '下单后未支付，超时自动取消', 1, 0, now(), now()),
      (8, 'zh-CN', '已支付', '支付成功，等待发货', 2, 0, now(), now()),
      (9, 'zh-CN', '已发货', '商品已出库，物流配送中', 3, 0, now(), now()),
      (10, 'zh-CN', '已完成', '用户确认收货，订单结束', 4, 0, now(), now()),
      (11, 'zh-CN', '已取消', '用户或系统取消订单', 5, 0, now(), now()),
      -- 性别
      (12, 'zh-CN', '男', '', 1, 0, now(), now()),
      (13, 'zh-CN', '女', '', 2, 0, now(), now()),
      (14, 'zh-CN', '未知', '用户未填写时默认值', 3, 0, now(), now()),
      -- 支付方式
      (15, 'zh-CN', '支付宝', '支持花呗、余额宝', 1, 0, now(), now()),
      (16, 'zh-CN', '微信支付', '需绑定微信', 2, 0, now(), now()),
      (17, 'zh-CN', '银联支付', '支持信用卡、储蓄卡', 3, 0, now(), now()),
      (18, 'zh-CN', '现金支付', '线下支付，已废弃（2025-01停用）', 4, 0, now(), now()),

      -- User Status
      (1, 'en-US', 'Normal', 'User can log in and operate normally', 1, 0, now(), now()),
      (2, 'en-US', 'Frozen', 'Temporarily frozen due to violation; requires admin to unfreeze', 2, 0, now(), now()),
      (3, 'en-US', 'Canceled', 'User voluntarily canceled; data retained but login disabled', 3, 0, now(), now()),

      -- Device Type
      (4, 'en-US', 'Temperature & Humidity Sensor', 'Supports temperature (-20~80°C) and humidity (0~100% RH) measurement', 1, 0, now(), now()),
      (5, 'en-US', 'Current Meter', 'Measures AC/DC current with 0.5-class accuracy', 2, 0, now(), now()),
      (6, 'en-US', 'Gas Detector', 'Not supported yet; hardware integration planned for Q4 2025', 3, 0, now(), now()),

      -- Order Status
      (7, 'en-US', 'Pending Payment', 'Order placed but not paid; auto-canceled if timeout', 1, 0, now(), now()),
      (8, 'en-US', 'Paid', 'Payment successful; awaiting shipment', 2, 0, now(), now()),
      (9, 'en-US', 'Shipped', 'Item has left warehouse; in transit', 3, 0, now(), now()),
      (10, 'en-US', 'Completed', 'User confirmed receipt; order closed', 4, 0, now(), now()),
      (11, 'en-US', 'Canceled', 'Order canceled by user or system', 5, 0, now(), now()),

      -- Gender
      (12, 'en-US', 'Male', '', 1, 0, now(), now()),
      (13, 'en-US', 'Female', '', 2, 0, now(), now()),
      (14, 'en-US', 'Unknown', 'Default value when user does not specify', 3, 0, now(), now()),

      -- Payment Method
      (15, 'en-US', 'Alipay', 'Supports Huabei and Yu’ebao', 1, 0, now(), now()),
      (16, 'en-US', 'WeChat Pay', 'Requires WeChat account binding', 2, 0, now(), now()),
      (17, 'en-US', 'UnionPay', 'Supports credit and debit cards', 3, 0, now(), now()),
      (18, 'en-US', 'Cash', 'Offline payment; deprecated as of Jan 2025', 4, 0, now(), now())
;
SELECT setval('sys_dict_entry_i18n_id_seq', (SELECT MAX(id) FROM sys_dict_entry_i18n));

-- 站内信分类
INSERT INTO public.internal_message_categories (id, code, name, remark, sort_order, is_enabled, created_at)
VALUES
    -- 订单相关分类（原主分类+子分类平级展示）
    (1, 'order', '订单通知', '包含订单支付、发货、退款等全流程通知', 1, true, NOW()),
    (101, 'order_paid', '支付成功', '订单支付完成时触发的通知', 2, true, NOW()),
    (102, 'order_unpaid', '支付超时', '订单未在规定时间内支付的提醒', 3, true, NOW()),
    (103, 'order_shipped', '已发货', '商家发货后通知用户', 4, true, NOW()),
    (104, 'order_refunded', '已退款', '订单退款流程完成的通知', 5, true, NOW()),

    -- 系统相关分类
    (2, 'system', '系统通知', '系统公告、维护提醒、版本更新等平台级通知', 6, true, NOW()),
    (201, 'system_announcement', '系统公告', '平台规则更新、重要通知等', 7, true, NOW()),
    (202, 'system_maintenance', '维护通知', '系统计划内维护的时间提醒', 8, true, NOW()),
    (203, 'system_upgrade', '版本更新', '客户端或功能升级的提示', 9, true, NOW()),

    -- 活动相关分类
    (3, 'activity', '活动通知', '营销活动报名、开始、结束等提醒', 10, true, NOW()),
    (301, 'activity_signup', '报名成功', '用户报名活动后确认通知', 11, true, NOW()),
    (302, 'activity_start', '活动开始', '活动即将开始的倒计时提醒', 12, true, NOW()),
    (303, 'activity_end', '活动结束', '活动结束及结果公示通知', 13, true, NOW()),

    -- 用户相关分类
    (4, 'user', '用户通知', '账号安全、信息变更、权限调整等个人相关通知', 14, true, NOW()),
    (401, 'user_login_abnormal', '异地登录', '账号在陌生设备登录的安全提醒', 15, true, NOW()),
    (402, 'user_profile_updated', '资料变更', '用户手机号、邮箱等信息修改后通知', 16, true, NOW()),
    (403, 'user_permission_changed', '权限变更', '账号角色或功能权限调整通知', 17, true, NOW())
;
SELECT setval('internal_message_categories_id_seq', (SELECT MAX(id) FROM internal_message_categories));

INSERT INTO public.comments (
    created_at, updated_at, deleted_at, created_by, updated_by, deleted_by,
    content_type, object_id, content, author_id, author_name, author_email,
    author_url, author_type, status, like_count, dislike_count, reply_count,
    ip_address, location, user_agent, detected_language, is_spam, is_sticky,
    reply_to_id, parent_id
) VALUES
-- 1. 管理员发布的置顶评论（根评论）
(
    '2026-02-01 10:00:00+08', '2026-02-01 10:00:00+08', NULL, 1, 1, NULL,
    'CONTENT_TYPE_POST', 1001,
    '感谢大家的积极反馈，我们会尽快处理大家提出的问题和建议！',
    999, '管理员小明', 'admin@example.com', 'https://example.com/admin',
    'AUTHOR_TYPE_ADMIN', 'STATUS_APPROVED', 58, 2, 12,
    '192.168.1.100', '北京市',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    'zh-CN', false, true, NULL, NULL
),
-- 2. 普通用户评论（根评论）
(
    '2026-02-01 10:15:00+08', '2026-02-01 10:15:00+08', NULL, 2, 2, NULL,
    'CONTENT_TYPE_POST', 1001,
    '这篇文章内容很有价值，解决了我一直以来的疑惑，非常感谢作者的分享！',
    1001, '用户小李', 'user1001@example.com', 'https://example.com/user/1001',
    'AUTHOR_TYPE_USER', 'STATUS_APPROVED', 32, 1, 5,
    '192.168.1.101', '上海市',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15',
    'zh-CN', false, false, NULL, NULL
),
-- 3. 访客评论（根评论）
(
    '2026-02-01 10:30:00+08', '2026-02-01 10:30:00+08', NULL, 0, 0, NULL,
    'CONTENT_TYPE_POST', 1001,
    '请问这个功能在移动端是否也能正常使用？希望能补充相关说明。',
    0, '访客小张', 'guest_2026@example.com', 'https://guest-example.com',
    'AUTHOR_TYPE_GUEST', 'STATUS_PENDING', 8, 0, 3,
    '192.168.1.102', '广州市',
    'Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36',
    'zh-CN', false, false, NULL, NULL
),
-- 4. 版主评论（回复评论2）
(
    '2026-02-01 11:00:00+08', '2026-02-01 11:00:00+08', NULL, 3, 3, NULL,
    'CONTENT_TYPE_POST', 1001,
    '感谢认可！我们后续会持续更新相关内容，敬请关注。',
    888, '版主小王', 'moderator@example.com', 'https://example.com/moderator',
    'AUTHOR_TYPE_MODERATOR', 'STATUS_APPROVED', 15, 0, 1,
    '192.168.1.103', '深圳市',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/121.0',
    'zh-CN', false, false, 2, 2
),
-- 5. 垃圾评论（待审核拒绝）
(
    '2026-02-01 12:00:00+08', '2026-02-01 12:00:00+08', NULL, 0, 0, NULL,
    'CONTENT_TYPE_POST', 1001,
    '高仿奢侈品代购，全网最低价，加微信XXX-XXXX-XXXX',
    0, '匿名用户', 'spam_2026@example.com', '',
    'AUTHOR_TYPE_GUEST', 'STATUS_REJECTED', 0, 5, 0,
    '192.168.1.104', '未知位置',
    'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36',
    'zh-CN', true, false, NULL, NULL
),
-- 6. 回复访客的评论（回复评论3）
(
    '2026-02-01 13:00:00+08', '2026-02-01 13:00:00+08', NULL, 1001, 1001, NULL,
    'CONTENT_TYPE_POST', 1001,
    '您好，该功能已适配移动端，iOS和Android均可正常使用，具体操作可参考文章末尾的附录。',
    1001, '用户小李', 'user1001@example.com', 'https://example.com/user/1001',
    'AUTHOR_TYPE_USER', 'STATUS_APPROVED', 10, 0, 0,
    '192.168.1.101', '上海市',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15',
    'zh-CN', false, false, 3, 3
),
-- 7. 已删除的评论（软删除）
(
    '2026-02-01 09:00:00+08', '2026-02-01 14:00:00+08', '2026-02-01 14:00:00+08', 0, 0, 1,
    'CONTENT_TYPE_POST', 1001,
    '这个内容完全没有用，浪费时间',
    0, '匿名用户', 'deleted_user@example.com', '',
    'AUTHOR_TYPE_GUEST', 'STATUS_REJECTED', 2, 10, 0,
    '192.168.1.105', '成都市',
    'Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1',
    'zh-CN', false, false, NULL, NULL
),
-- 8. 多语言评论（英文）
(
    '2026-02-01 15:00:00+08', '2026-02-01 15:00:00+08', NULL, 4, 4, NULL,
    'CONTENT_TYPE_PAGE', 2001,
    'This article is very helpful for my project, thank you so much for sharing your experience!',
    2001, 'John Doe', 'john.doe@example.com', 'https://john-doe-example.com',
    'AUTHOR_TYPE_USER', 'STATUS_APPROVED', 7, 0, 2,
    '192.168.1.106', 'New York, USA',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0',
    'en-US', false, false, NULL, NULL
);

INSERT INTO public.media_assets (
    created_at, updated_at, deleted_at, created_by, updated_by, deleted_by,
    filename, type, mime_type, size, storage_path, url,
    width, height, duration, alt_text, title, caption,
    processing_status, processing_error, file_hash, folder_id, file_id,
    reference_count, is_private
) VALUES
-- 1. 图片资源（处理完成）
(
    '2026-02-01 09:00:00+08', '2026-02-01 09:05:00+08', NULL, 1001, 1001, NULL,
    'product_banner.jpg', 'ASSET_TYPE_IMAGE', 'image/jpeg', 2048576,
    '/uploads/2026/02/01/product_banner.jpg',
    'https://cdn.example.com/uploads/2026/02/01/product_banner.jpg',
    1920, 1080, 0,
    '产品宣传横幅', '2026春季新品横幅', '2026春季新品宣传主视觉',
    'PROCESSING_STATUS_COMPLETED', NULL,
    'md5:8a7f9b6e4d3c2a1e0f9e8d7c6b5a4f3e', 1, 10001,
    15, false
),
-- 2. 视频资源（处理中）
(
    '2026-02-01 10:30:00+08', '2026-02-01 10:45:00+08', NULL, 1002, 1002, NULL,
    'tutorial_video.mp4', 'ASSET_TYPE_VIDEO', 'video/mp4', 85937500,
    '/uploads/2026/02/01/tutorial_video.mp4',
    'https://cdn.example.com/uploads/2026/02/01/tutorial_video.mp4',
    1920, 1080, 360,
    '产品使用教程视频', '2026产品使用教程', '详细讲解产品核心功能使用方法',
    'PROCESSING_STATUS_PROCESSING', NULL,
    'md5:9b8a7c6d5e4f3a2b1c0d9e8f7a6b5c4d', 2, 10002,
    3, false
),
-- 3. 文档资源（处理完成）
(
    '2026-02-01 11:15:00+08', '2026-02-01 11:15:00+08', NULL, 1003, 1003, NULL,
    'user_manual.pdf', 'ASSET_TYPE_DOCUMENT', 'application/pdf', 1572864,
    '/uploads/2026/02/01/user_manual.pdf',
    'https://cdn.example.com/uploads/2026/02/01/user_manual.pdf',
    0, 0, 0,
    '用户手册PDF', '2026产品用户手册', '包含产品安装、使用、故障排除等内容',
    'PROCESSING_STATUS_COMPLETED', NULL,
    'md5:7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f', 3, 10003,
    28, true
),
-- 4. 音频资源（处理失败）
(
    '2026-02-01 14:00:00+08', '2026-02-01 14:10:00+08', NULL, 1004, 1004, NULL,
    'podcast_episode.mp3', 'ASSET_TYPE_AUDIO', 'audio/mpeg', 12582912,
    '/uploads/2026/02/01/podcast_episode.mp3',
    'https://cdn.example.com/uploads/2026/02/01/podcast_episode.mp3',
    0, 0, 1800,
    '播客音频文件', '2026第12期播客', '行业趋势分析专题播客',
    'PROCESSING_STATUS_FAILED', 'FFmpeg error: Invalid audio codec, could not transcode file',
    'md5:5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a', 4, 10004,
    0, false
),
-- 5. 归档文件（处理完成）
(
    '2026-02-01 15:30:00+08', '2026-02-01 15:30:00+08', NULL, 1001, 1001, NULL,
    'backup_files.zip', 'ASSET_TYPE_ARCHIVE', 'application/zip', 52428800,
    '/uploads/2026/02/01/backup_files.zip',
    'https://cdn.example.com/uploads/2026/02/01/backup_files.zip',
    0, 0, 0,
    '备份压缩包', '2026年2月数据备份', '包含用户数据和媒体文件备份',
    'PROCESSING_STATUS_COMPLETED', NULL,
    'md5:3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b', 5, 10005,
    1, true
),
-- 6. 其他类型文件（上传中）
(
    '2026-02-01 16:45:00+08', '2026-02-01 16:45:00+08', NULL, 1005, 1005, NULL,
    'data_analysis.xlsx', 'ASSET_TYPE_OTHER', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    8388608, '/uploads/2026/02/01/data_analysis.xlsx',
    'https://cdn.example.com/uploads/2026/02/01/data_analysis.xlsx',
    0, 0, 0,
    '数据分析表格', '2026用户行为分析', '月度用户行为数据分析报告',
    'PROCESSING_STATUS_UPLOADING', NULL,
    'md5:1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d', 6, 10006,
    0, false
),
-- 7. 软删除的图片资源
(
    '2026-01-20 08:00:00+08', '2026-02-01 17:00:00+08', '2026-02-01 17:00:00+08',
    1001, 1001, 1,
    'old_banner.png', 'ASSET_TYPE_IMAGE', 'image/png', 1048576,
    '/uploads/2026/01/20/old_banner.png',
    'https://cdn.example.com/uploads/2026/01/20/old_banner.png',
    1280, 720, 0,
    '旧版横幅图片', '2025冬季横幅', '2025冬季促销横幅（已下线）',
    'PROCESSING_STATUS_COMPLETED', NULL,
    'md5:8f9e8d7c6b5a4f3e2d1c0b9a8s7d6f5e', 1, 10007,
    0, false
);

INSERT INTO public.site_settings (
    created_at, updated_at, deleted_at, created_by, updated_by, deleted_by,
    site_id, locale, "group", key, value, type,
    label, description, placeholder, options, is_required, validation_regex
) VALUES
-- ===================== 全局配置（site_id=0）- 通用分组 =====================
-- 1. 站点标题（单行文本，必填）
(
    '2026-02-01 09:00:00+08', '2026-02-01 09:00:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'general', 'site_title', '我的企业官网', 'SETTING_TYPE_TEXT',
    '站点标题', '浏览器标签和页面头部显示的站点名称', '请输入站点标题',
    '{}'::jsonb, true, '^.{2,50}$'
),
-- 2. 站点描述（多行文本）
(
    '2026-02-01 09:05:00+08', '2026-02-01 09:05:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'general', 'site_description', '专注于企业数字化解决方案的专业平台，提供一站式技术服务。',
    'SETTING_TYPE_TEXTAREA',
    '站点描述', '用于SEO和社交分享的站点简介', '请输入站点描述（不超过200字）',
    '{}'::jsonb, false, '^.{0,200}$'
),
-- 3. 每页文章数（数字）
(
    '2026-02-01 09:10:00+08', '2026-02-01 09:10:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'general', 'posts_per_page', '10', 'SETTING_TYPE_NUMBER',
    '每页文章数量', '列表页默认显示的文章条数', '请输入数字（1-50）',
    '{}'::jsonb, true, '^[1-9][0-9]?$|^50$'
),
-- 4. 开启评论功能（布尔）
(
    '2026-02-01 09:15:00+08', '2026-02-01 09:15:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'general', 'enable_comments', 'true', 'SETTING_TYPE_BOOLEAN',
    '开启评论功能', '是否允许访客/用户在文章下发表评论', '',
    '{}'::jsonb, false, ''
),
-- 5. 站点LOGO（图片，关联MediaAsset ID）
(
    '2026-02-01 09:20:00+08', '2026-02-01 09:20:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'general', 'site_logo', '10001', 'SETTING_TYPE_IMAGE',
    '站点LOGO', '上传的LOGO图片ID（关联media_assets表）', '上传LOGO图片',
    '{}'::jsonb, true, ''
),

-- ===================== 全局配置（site_id=0）- SEO分组 =====================
-- 6. 网站域名（URL）
(
    '2026-02-01 09:25:00+08', '2026-02-01 09:25:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'seo', 'site_domain', 'https://www.example.com', 'SETTING_TYPE_URL',
    '网站域名', '站点的主域名（带https）', '请输入完整的网站域名',
    '{}'::jsonb, true, '^https?://.+$'
),
-- 7. 站长邮箱（邮箱）
(
    '2026-02-01 09:30:00+08', '2026-02-01 09:30:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'seo', 'webmaster_email', 'webmaster@example.com', 'SETTING_TYPE_EMAIL',
    '站长邮箱', '用于搜索引擎验证和站长工具通知', '请输入有效的邮箱地址',
    '{}'::jsonb, true, '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$'
),
-- 8. 搜索引擎索引（下拉选择）
(
    '2026-02-01 09:35:00+08', '2026-02-01 09:35:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'seo', 'robots_index', 'index', 'SETTING_TYPE_SELECT',
    '搜索引擎索引', '是否允许搜索引擎抓取站点内容', '',
    '{"index":"允许索引", "noindex":"禁止索引"}'::jsonb,
    true, ''
),

-- ===================== 全局配置（site_id=0）- 社交分组 =====================
-- 9. 社交媒体链接（JSON）
(
    '2026-02-01 09:40:00+08', '2026-02-01 09:40:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'social', 'social_links', '[{"name":"微信","icon":"wechat","url":"https://weixin.example.com"},{"name":"微博","icon":"weibo","url":"https://weibo.example.com"},{"name":"GitHub","icon":"github","url":"https://github.com/example"}]',
    'SETTING_TYPE_JSON',
    '社交媒体链接', '站点底部显示的社交平台链接（JSON格式）', '输入JSON格式的社交链接数组',
    '{}'::jsonb, false, ''
),

-- ===================== 站点1配置（site_id=1）- 英文版本 =====================
-- 10. 站点标题（英文）
(
    '2026-02-01 10:00:00+08', '2026-02-01 10:00:00+08', NULL, 1, 1, NULL,
    1, 'en-US', 'general', 'site_title', 'My Enterprise Website', 'SETTING_TYPE_TEXT',
    'Site Title', 'Site name displayed in browser tab and page header', 'Enter site title',
    '{}'::jsonb, true, '^.{2,50}$'
),
-- 11. 每页文章数（英文站点，自定义值）
(
    '2026-02-01 10:05:00+08', '2026-02-01 10:05:00+08', NULL, 1, 1, NULL,
    1, 'en-US', 'general', 'posts_per_page', '15', 'SETTING_TYPE_NUMBER',
    'Posts Per Page', 'Number of articles displayed per page', 'Enter number (1-50)',
    '{}'::jsonb, true, '^[1-9][0-9]?$|^50$'
),

-- ===================== 特殊场景配置 =====================
-- 12. 软删除的配置项
(
    '2026-01-15 08:00:00+08', '2026-02-01 11:00:00+08', '2026-02-01 11:00:00+08',
    1, 1, 1,
    0, 'zh-CN', 'general', 'old_site_footer', '© 2025 我的企业 版权所有', 'SETTING_TYPE_TEXT',
    '旧版页脚文案', '已废弃的页脚版权文案', '',
    '{}'::jsonb, false, ''
),
-- 13. 评论审核方式（下拉选择，多选项示例）
(
    '2026-02-01 11:10:00+08', '2026-02-01 11:10:00+08', NULL, 1, 1, NULL,
    0, 'zh-CN', 'comment', 'comment_moderation', 'manual', 'SETTING_TYPE_SELECT',
    '评论审核方式', '评论发布前的审核规则', '',
    '{"auto":"自动审核", "manual":"人工审核", "admin_only":"仅管理员审核"}'::jsonb,
    true, ''
);

INSERT INTO public.sites (
    created_at, updated_at, deleted_at, created_by, updated_by, deleted_by,
    tenant_id, name, slug, domain, alternate_domains, is_default,
    status, default_locale, template, theme, visit_count
) VALUES
-- 1. 主站点（默认站点，活跃状态）
(
    '2026-01-01 09:00:00+08', '2026-02-01 10:00:00+08', NULL, 1, 1, NULL,
    1, '企业官网主站', 'main-site', 'https://www.example.com',
    '["https://example.com", "https://www.example.cn"]'::jsonb, true,
    'SITE_STATUS_ACTIVE', 'zh-CN', 'enterprise', 'default-dark',
    158920
),
-- 2. 英文站点（租户1，活跃状态）
(
    '2026-01-05 10:00:00+08', '2026-02-01 11:00:00+08', NULL, 1, 1, NULL,
    1, 'Enterprise English Site', 'english-site', 'https://en.example.com',
    '["https://english.example.com"]'::jsonb, false,
    'SITE_STATUS_ACTIVE', 'en-US', 'enterprise', 'default-light',
    85640
),
-- 3. 营销活动站点（租户1，暂存状态）
(
    '2026-01-10 14:00:00+08', '2026-02-01 12:00:00+08', NULL, 1, 1, NULL,
    1, '2026春季促销活动站', 'spring-2026-promotion', 'https://promo2026.example.com',
    '[]'::jsonb, false,
    'SITE_STATUS_MAINTENANCE', 'zh-CN', 'promotion', 'spring-theme',
    0
),
-- 4. 租户2的电商站点（独立租户，活跃状态）
(
    '2026-01-15 09:30:00+08', '2026-02-01 13:00:00+08', NULL, 2, 2, NULL,
    2, '优品电商平台', 'youpin-shop', 'https://www.youpin.com',
    '["https://youpin.com", "https://m.youpin.com"]'::jsonb, true,
    'SITE_STATUS_ACTIVE', 'zh-CN', 'ecommerce', 'youpin-default',
    987540
),
-- 5. 废弃站点（租户1，已归档）
(
    '2025-06-01 08:00:00+08', '2026-01-20 15:00:00+08', NULL, 1, 1, NULL,
    1, '2025夏季活动站', 'summer-2025', 'https://summer2025.example.com',
    '[]'::jsonb, false,
    'SITE_STATUS_MAINTENANCE', 'zh-CN', 'promotion', 'summer-theme',
    45210
),
-- 6. 测试站点（租户1，禁用状态）
(
    '2026-01-20 16:00:00+08', '2026-02-01 14:00:00+08', NULL, 1, 1, NULL,
    1, '内部测试站点', 'test-site', 'https://test.example.com',
    '["https://dev.example.com"]'::jsonb, false,
    'SITE_STATUS_INACTIVE', 'zh-CN', 'test', 'test-theme',
    1250
),
-- 7. 软删除的站点（租户2，已删除）
(
    '2025-10-01 10:00:00+08', '2026-02-01 15:00:00+08', '2026-02-01 15:00:00+08',
    2, 2, 2,
    2, '旧版移动端站点', 'old-mobile-site', 'https://m.old.youpin.com',
    '[]'::jsonb, false,
    'SITE_STATUS_MAINTENANCE', 'zh-CN', 'mobile', 'old-mobile-theme',
    256800
);

-- ----------------------------
-- 插入 sites 表测试数据
-- ----------------------------
INSERT INTO public.sites (
    created_at, updated_at, tenant_id, name, slug, domain,
    alternate_domains, is_default, status, default_locale,
    template, theme, visit_count
) VALUES
-- 站点1：默认中文站点（主站）
(
    NOW(), NOW(), 0,
    'GoWind 官方博客',
    'gowind-blog',
    'blog.gowind.com',
    '["www.gowind-blog.com", "blog.gowind.cn"]'::jsonb,
    true,
    'SITE_STATUS_ACTIVE',
    'zh-CN',
    'default-blog',
    'dark-mode',
    12580
),
-- 站点2：英文子站点
(
    NOW(), NOW(), 0,
    'GoWind English Docs',
    'gowind-docs-en',
    'docs.gowind.com',
    '["en-docs.gowind.com"]'::jsonb,
    false,
    'SITE_STATUS_ACTIVE',
    'en-US',
    'default-docs',
    'light-mode',
    8920
),
-- 站点3：维护中的测试站点（用于测试状态逻辑）
(
    NOW(), NOW(), 0,
    '测试站点（维护中）',
    'test-site-maintenance',
    'test.gowind.com',
    '[]'::jsonb,
    false,
    'SITE_STATUS_MAINTENANCE',
    'zh-CN',
    'default-test',
    'default',
    1560
);

-- ----------------------------
-- 插入 site_settings 表测试数据
-- 关联上述 3 个站点，覆盖多语言/多配置组
-- ----------------------------
INSERT INTO public.site_settings (
    created_at, updated_at, site_id, locale, "group", key,
    value, type, label, description, placeholder, options,
    is_required, validation_regex
) VALUES
-- ========== 站点1（ID=1）：中文配置 ==========
-- 基础设置组
(
    NOW(), NOW(), 1, 'zh-CN', 'basic', 'site_title',
    'GoWind - 高性能Go语言CMS',
    'SETTING_TYPE_TEXT', '站点标题', '显示在浏览器标签的站点主标题', '请输入站点标题', NULL,
    true, '^.{2,50}$'
),
(
    NOW(), NOW(), 1, 'zh-CN', 'basic', 'site_description',
    '基于Go+Vue3开发的多站点/多语言CMS系统',
    'SETTING_TYPE_TEXT', '站点描述', '用于SEO的站点简介', '请输入站点描述（10-200字）', NULL,
    true, '^.{10,200}$'
),
(
    NOW(), NOW(), 1, 'zh-CN', 'basic', 'enable_comment',
    'true',
    'SETTING_TYPE_BOOLEAN', '启用评论功能', '是否允许用户在文章下评论', NULL, NULL,
    false, NULL
),
-- SEO配置组
(
    NOW(), NOW(), 1, 'zh-CN', 'seo', 'meta_keywords',
    'GoWind,CMS,Go语言,多语言,多站点',
    'SETTING_TYPE_TEXT', 'SEO关键词', '用英文逗号分隔的关键词列表', '请输入SEO关键词', NULL,
    false, NULL
),
(
    NOW(), NOW(), 1, 'zh-CN', 'seo', 'enable_baidu_verify',
    'true',
    'SETTING_TYPE_BOOLEAN', '启用百度验证', '是否添加百度站点验证代码', NULL, NULL,
    false, NULL
),
-- 主题配置组（下拉选择类型）
(
    NOW(), NOW(), 1, 'zh-CN', 'theme', 'primary_color',
    '#1890ff',
    'SETTING_TYPE_SELECT', '主题主色', '站点的核心品牌色', NULL,
    '{"#1890ff": "Sky Blue (Default)", "#e53e3e": "Red", "#389e0d": "Green", "#faad14": "Yellow"}'::jsonb,
    false, '^#[0-9a-fA-F]{6}$'
),
-- ========== 站点1（ID=1）：英文配置 ==========
(
    NOW(), NOW(), 1, 'en-US', 'basic', 'site_title',
    'GoWind - High Performance Go CMS',
    'SETTING_TYPE_TEXT', 'Site Title', 'Main title displayed in browser tab', 'Enter site title', NULL,
    true, '^.{2,50}$'
),
(
    NOW(), NOW(), 1, 'en-US', 'basic', 'site_description',
    'Multi-site & multi-language CMS built with Go + Vue3',
    'SETTING_TYPE_TEXT', 'Site Description', 'SEO-friendly site introduction', 'Enter site description (10-200 chars)', NULL,
    true, '^.{10,200}$'
),
-- ========== 站点2（ID=2）：英文配置 ==========
(
    NOW(), NOW(), 2, 'en-US', 'basic', 'site_title',
    'GoWind Documentation - English Version',
    'SETTING_TYPE_TEXT', 'Site Title', 'Main title of English docs site', 'Enter site title', NULL,
    true, '^.{2,50}$'
),
(
    NOW(), NOW(), 2, 'en-US', 'seo', 'meta_keywords',
    'GoWind,CMS,Documentation,Go,Vue3,Multi-language',
    'SETTING_TYPE_TEXT', 'SEO Keywords', 'Comma-separated keywords for SEO', 'Enter SEO keywords', NULL,
    false, NULL
),
-- ========== 站点3（ID=3）：维护中站点配置 ==========
(
    NOW(), NOW(), 3, 'zh-CN', 'maintenance', 'maintenance_title',
    '站点维护中',
    'SETTING_TYPE_TEXT', '维护页面标题', '站点维护时显示的标题', '请输入维护标题', NULL,
    true, NULL
),
(
    NOW(), NOW(), 3, 'zh-CN', 'maintenance', 'maintenance_content',
    '本站点正在维护，预计2小时后恢复访问。',
    'SETTING_TYPE_TEXT', '维护页面内容', '站点维护时显示的提示内容', '请输入维护提示', NULL,
    true, NULL
);


-- ==================== 导航组（navigations）====================
-- 中文导航组
INSERT INTO public.navigations (
    id, created_at, updated_at, name, location, locale, is_active, created_by, updated_by
) VALUES
-- Main Navigation (HEADER)
(101, NOW(), NOW(), '主导航', 'HEADER', 'zh-CN', true, 1, 1),
-- Footer Navigation (FOOTER)
(102, NOW(), NOW(), '页脚导航', 'FOOTER', 'zh-CN', true, 1, 1),
-- Sidebar Navigation (SIDEBAR)
(103, NOW(), NOW(), '侧边栏导航', 'SIDEBAR', 'zh-CN', true, 1, 1),

-- 英文导航组（对应中文组）
(201, NOW(), NOW(), 'Main Navigation', 'HEADER', 'en-US', true, 1, 1),
(202, NOW(), NOW(), 'Footer Navigation', 'FOOTER', 'en-US', true, 1, 1),
(203, NOW(), NOW(), 'Sidebar Navigation', 'SIDEBAR', 'en-US', true, 1, 1);

-- ==================== 导航项（navigation_items）====================
-- 注意：按 parentId 顺序插入（先父后子），确保外键约束有效
INSERT INTO public.navigation_items (
    id, created_at, updated_at, sort_order, link_type, navigation_id,
    title, url, object_id, icon, description, is_open_new_tab,
    is_invalid, css_class, required_permission, parent_id, created_by, updated_by
) VALUES
-- ========== 导航组 101（zh-CN 主导航） ==========
-- 首页（顶级）
(1001, NOW(), NOW(), 1, 'LINK_TYPE_CUSTOM', 101, '首页', '/', 0, 'home', '返回首页', false, false, 'nav-item header-nav', '', NULL, 1, 1),
-- 文章（顶级）
(1002, NOW(), NOW(), 2, 'LINK_TYPE_CUSTOM', 101, '文章', '/post', 0, 'document', '浏览所有文章', false, false, 'nav-item header-nav', '', NULL, 1, 1),
-- 分类（顶级，有子菜单）
(1003, NOW(), NOW(), 3, 'LINK_TYPE_CUSTOM', 101, '分类', '/category', 0, 'folder', '浏览所有分类', false, false, 'nav-item header-nav has-children', '', NULL, 1, 1),
-- 标签（顶级）
(1004, NOW(), NOW(), 4, 'LINK_TYPE_CUSTOM', 101, '标签', '/tag', 0, 'tag', '浏览所有标签', false, false, 'nav-item header-nav', '', NULL, 1, 1),
-- 关于（顶级）
(1005, NOW(), NOW(), 5, 'LINK_TYPE_PAGE', 101, '关于', '/about', 1, 'information', '关于我们', false, false, 'nav-item header-nav', '', NULL, 1, 1),

-- 分类子菜单：技术分享（parentId = 1003）
(1006, NOW(), NOW(), 1, 'LINK_TYPE_CATEGORY', 101, '技术分享', '/category/1', 1, 'code', '技术文章分类', false, false, 'nav-item child-nav', '', 1003, 1, 1),
-- 分类子菜单：生活随笔（parentId = 1003）
(1007, NOW(), NOW(), 2, 'LINK_TYPE_CATEGORY', 101, '生活随笔', '/category/2', 2, 'blog', '生活文章分类', false, false, 'nav-item child-nav', '', 1003, 1, 1),

-- ========== 导航组 102（zh-CN 页脚导航） ==========
(1008, NOW(), NOW(), 1, 'LINK_TYPE_PAGE', 102, '联系我们', '/contact', 2, 'email', '联系我们', false, false, 'nav-item footer-nav', '', NULL, 1, 1),
(1009, NOW(), NOW(), 2, 'LINK_TYPE_PAGE', 102, '隐私政策', '/privacy', 3, 'shield-checkmark', '隐私政策', false, false, 'nav-item footer-nav', '', NULL, 1, 1),
(1010, NOW(), NOW(), 3, 'LINK_TYPE_PAGE', 102, '服务条款', '/terms', 4, 'document', '服务条款', false, false, 'nav-item footer-nav', '', NULL, 1, 1),
(1011, NOW(), NOW(), 4, 'LINK_TYPE_EXTERNAL', 102, 'GitHub', 'https://github.com', 0, 'logo-github', '访问我们的GitHub', true, false, 'nav-item footer-nav', '', NULL, 1, 1),

-- ========== 导航组 103（zh-CN 侧边栏导航） ==========
(1012, NOW(), NOW(), 1, 'LINK_TYPE_CUSTOM', 103, '热门标签', '/tag', 0, 'pricetag', '浏览热门标签', false, false, 'nav-item sidebar-nav', '', NULL, 1, 1),
(1013, NOW(), NOW(), 2, 'LINK_TYPE_CUSTOM', 103, '归档', '/archive', 0, 'archive', '文章归档', false, false, 'nav-item sidebar-nav', '', NULL, 1, 1),

-- ========== 导航组 201（en-US Main Navigation） ==========
(2001, NOW(), NOW(), 1, 'LINK_TYPE_CUSTOM', 201, 'Home', '/', 0, 'home', 'Back to homepage', false, false, 'nav-item header-nav', '', NULL, 1, 1),
(2002, NOW(), NOW(), 2, 'LINK_TYPE_CUSTOM', 201, 'Posts', '/post', 0, 'document', 'Browse all posts', false, false, 'nav-item header-nav', '', NULL, 1, 1),
(2003, NOW(), NOW(), 3, 'LINK_TYPE_CUSTOM', 201, 'Categories', '/category', 0, 'folder', 'Browse all categories', false, false, 'nav-item header-nav has-children', '', NULL, 1, 1),
(2004, NOW(), NOW(), 4, 'LINK_TYPE_CUSTOM', 201, 'Tags', '/tag', 0, 'tag', 'Browse all tags', false, false, 'nav-item header-nav', '', NULL, 1, 1),
(2005, NOW(), NOW(), 5, 'LINK_TYPE_PAGE', 201, 'About', '/about', 1, 'information', 'About us', false, false, 'nav-item header-nav', '', NULL, 1, 1),
-- Categories 子菜单
(2006, NOW(), NOW(), 1, 'LINK_TYPE_CATEGORY', 201, 'Tech Sharing', '/category/1', 1, 'code', 'Tech article category', false, false, 'nav-item child-nav', '', 2003, 1, 1),
(2007, NOW(), NOW(), 2, 'LINK_TYPE_CATEGORY', 201, 'Life Notes', '/category/2', 2, 'blog', 'Life article category', false, false, 'nav-item child-nav', '', 2003, 1, 1),

-- ========== 导航组 202（en-US Footer Navigation） ==========
(2008, NOW(), NOW(), 1, 'LINK_TYPE_PAGE', 202, 'Contact', '/contact', 2, 'email', 'Contact us', false, false, 'nav-item footer-nav', '', NULL, 1, 1),
(2009, NOW(), NOW(), 2, 'LINK_TYPE_PAGE', 202, 'Privacy Policy', '/privacy', 3, 'shield-checkmark', 'Privacy policy', false, false, 'nav-item footer-nav', '', NULL, 1, 1),
(2010, NOW(), NOW(), 3, 'LINK_TYPE_PAGE', 202, 'Terms of Service', '/terms', 4, 'document', 'Terms of service', false, false, 'nav-item footer-nav', '', NULL, 1, 1),
(2011, NOW(), NOW(), 4, 'LINK_TYPE_EXTERNAL', 202, 'GitHub', 'https://github.com', 0, 'logo-github', 'Visit our GitHub', true, false, 'nav-item footer-nav', '', NULL, 1, 1),

-- ========== 导航组 203（en-US Sidebar Navigation） ==========
(2012, NOW(), NOW(), 1, 'LINK_TYPE_CUSTOM', 203, 'Popular Tags', '/tag', 0, 'pricetag', 'Browse popular tags', false, false, 'nav-item sidebar-nav', '', NULL, 1, 1),
(2013, NOW(), NOW(), 2, 'LINK_TYPE_CUSTOM', 203, 'Archive', '/archive', 0, 'archive', 'Post archive', false, false, 'nav-item sidebar-nav', '', NULL, 1, 1);


-- ----------------------------
-- 插入 comments 表测试数据
-- 包含多级评论、不同状态/作者类型、统计字段
-- ----------------------------
INSERT INTO public.comments (
    created_at, updated_at, content_type, object_id, content,
    author_id, author_name, author_email, author_url, author_type,
    status, like_count, dislike_count, reply_count, ip_address,
    location, user_agent, detected_language, is_spam, is_sticky,
    reply_to_id, parent_id
) VALUES
-- ========== 场景1：文章（object_id=101）的核心评论 ==========
-- 父评论1：管理员发布的置顶评论（已发布、置顶）
(
    NOW() - INTERVAL '3 days', NOW(),
    'CONTENT_TYPE_POST', 101,
    '感谢大家使用 GoWind CMS！如果有使用问题，欢迎在下方留言交流～',
    1, -- 管理员ID
    'GoWind 官方',
    'admin@gowind.com',
    'https://gowind.com',
    'AUTHOR_TYPE_ADMIN',
    'STATUS_APPROVED',
    89, 2, 12,
    '103.2xx.xx.1',
    '北京市 阿里云',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36',
    'zh-CN',
    false, true,
    NULL, NULL
),
-- 父评论2：登录用户的正常评论（已发布）
(
    NOW() - INTERVAL '2 days', NOW(),
    'CONTENT_TYPE_POST', 101,
    '这个CMS的多语言功能太实用了！我用它搭建了中英文双语博客，体验很棒～',
    1001, -- 普通用户ID
    '张三',
    'zhangsan@example.com',
    'https://zhangsan-blog.com',
    'AUTHOR_TYPE_USER',
    'STATUS_APPROVED',
    45, 1, 3,
    '180.1xx.xx.2',
    '上海市 电信',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/121.0.0.0 Safari/537.36',
    'zh-CN',
    false, false,
    NULL, NULL
),
-- 父评论3：访客评论（待审核）
(
    NOW() - INTERVAL '1 day', NOW(),
    'CONTENT_TYPE_POST', 101,
    '请问如何配置ES全站搜索？尝试了几次都没成功，求教程～',
    0, -- 访客无ID
    '小李',
    'xiaoli@qq.com',
    '',
    'AUTHOR_TYPE_GUEST',
    'STATUS_PENDING',
    8, 0, 1,
    '27.1xx.xx.3',
    '广东省广州市 联通',
    'Mozilla/5.0 (iPhone; CPU iPhone OS 17_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Safari/604.1',
    'zh-CN',
    false, false,
    NULL, NULL
),
-- 父评论4：垃圾评论（广告、标记为spam）
(
    NOW() - INTERVAL '12 hours', NOW(),
    'CONTENT_TYPE_POST', 101,
    '【代开发票】增值税发票/普通发票，保真可查验，联系微信：138xxxx8888',
    0,
    '诚信代开',
    'fake@spam.com',
    'http://fake-invoice.com',
    'AUTHOR_TYPE_GUEST',
    'STATUS_REJECTED',
    0, 15, 0,
    '43.2xx.xx.4',
    '海外 未知',
    'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36',
    'zh-CN',
    true, false,
    NULL, NULL
),
-- 父评论5：英文评论（海外用户）
(
    NOW() - INTERVAL '8 hours', NOW(),
    'CONTENT_TYPE_POST', 101,
    'This CMS is really lightweight and fast! I love the dark mode design, it''s very eye-friendly.',
    2001, -- 海外用户ID
    'John Doe',
    'john.doe@gmail.com',
    'https://john-doe.dev',
    'AUTHOR_TYPE_USER',
    'STATUS_APPROVED',
    23, 0, 2,
    '198.51.xx.5',
    'New York, USA',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15',
    'en-US',
    false, false,
    NULL, NULL
),

-- ========== 场景2：父评论1的子回复（多级嵌套） ==========
-- 回复1-1：回复管理员的置顶评论（子评论）
(
    NOW() - INTERVAL '2 days 18 hours', NOW(),
    'CONTENT_TYPE_POST', 101,
    '官方大大，请问多站点配置时，域名解析需要注意什么？',
    1002,
    '李四',
    'lisi@163.com',
    '',
    'AUTHOR_TYPE_USER',
    'STATUS_APPROVED',
    12, 0, 1,
    '117.1xx.xx.6',
    '浙江省杭州市 阿里云',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/122.0.0.0 Safari/537.36',
    'zh-CN',
    false, false,
    1, -- 回复的是ID=1的评论
    1  -- 父评论ID=1
),
-- 回复1-1-1：回复李四的评论（孙评论，管理员回复）
(
    NOW() - INTERVAL '2 days 12 hours', NOW(),
    'CONTENT_TYPE_POST', 101,
    '您好！多站点域名解析需确保A记录指向服务器IP，同时在CMS中配置域名白名单，避免跨域问题～',
    1,
    'GoWind 官方',
    'admin@gowind.com',
    'https://gowind.com',
    'AUTHOR_TYPE_ADMIN',
    'STATUS_APPROVED',
    8, 0, 0,
    '103.2xx.xx.1',
    '北京市 阿里云',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Chrome/122.0.0.0 Safari/537.36',
    'zh-CN',
    false, false,
    6, -- 回复的是ID=6的评论
    6  -- 父评论ID=6
),

-- ========== 场景3：父评论2的子回复 ==========
-- 回复2-1：回复张三的评论（访客）
(
    NOW() - INTERVAL '1 day 6 hours', NOW(),
    'CONTENT_TYPE_POST', 101,
    '请问中英文切换时，Tag标签的多语言是自动同步的吗？',
    0,
    '王老五',
    'wanglaowu@126.com',
    '',
    'AUTHOR_TYPE_GUEST',
    'STATUS_PENDING',
    3, 0, 0,
    '140.1xx.xx.7',
    '四川省成都市 移动',
    'Mozilla/5.0 (Android 14; Mobile; rv:123.0) Gecko/123.0 Firefox/123.0',
    'zh-CN',
    false, false,
    2, -- 回复的是ID=2的评论
    2  -- 父评论ID=2
),
-- 回复2-2：张三回复王老五（登录用户）
(
    NOW() - INTERVAL '1 day 2 hours', NOW(),
    'CONTENT_TYPE_POST', 101,
    '是的，我测试过，Tag的多语言是绑定同一个key的，切换语言会自动显示对应名称～',
    1001,
    '张三',
    'zhangsan@example.com',
    'https://zhangsan-blog.com',
    'AUTHOR_TYPE_USER',
    'STATUS_APPROVED',
    5, 0, 0,
    '180.1xx.xx.2',
    '上海市 电信',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edge/121.0.0.0 Safari/537.36',
    'zh-CN',
    false, false,
    8, -- 回复的是ID=8的评论
    8  -- 父评论ID=8
),

-- ========== 场景4：父评论3的子回复 ==========
-- 回复3-1：管理员回复小李的问题
(
    NOW() - INTERVAL '8 hours', NOW(),
    'CONTENT_TYPE_POST', 101,
    '您好！ES配置教程已更新：https://gowind.com/docs/es-setup，可参考～',
    1,
    'GoWind 官方',
    'admin@gowind.com',
    'https://gowind.com',
    'AUTHOR_TYPE_ADMIN',
    'STATUS_APPROVED',
    10, 0, 0,
    '103.2xx.xx.1',
    '北京市 阿里云',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Chrome/122.0.0.0 Safari/537.36',
    'zh-CN',
    false, false,
    3, -- 回复的是ID=3的评论
    3  -- 父评论ID=3
),

-- ========== 场景5：父评论5的子回复（英文） ==========
-- 回复5-1：回复John Doe的英文评论
(
    NOW() - INTERVAL '6 hours', NOW(),
    'CONTENT_TYPE_POST', 101,
    'I agree! The performance is amazing even with 10k+ posts. Have you tried the multi-tenant feature?',
    2002,
    'Jane Smith',
    'jane.smith@outlook.com',
    'https://jane-smith.dev',
    'AUTHOR_TYPE_USER',
    'STATUS_APPROVED',
    7, 0, 0,
    '203.0xx.xx.8',
    'London, UK',
    'Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:123.0) Gecko/123.0 Firefox/123.0',
    'en-US',
    false, false,
    5, -- 回复的是ID=5的评论
    5  -- 父评论ID=5
),

-- ========== 场景6：页面（object_id=201）的评论 ==========
-- 父评论6：页面的访客评论（已拒绝）
(
    NOW() - INTERVAL '3 days', NOW(),
    'CONTENT_TYPE_PAGE', 201,
    '这个页面的排版有问题，手机端显示不全，希望尽快修复',
    0,
    '赵六',
    'zhaoliu@foxmail.com',
    '',
    'AUTHOR_TYPE_GUEST',
    'STATUS_REJECTED',
    2, 3, 0,
    '58.2xx.xx.9',
    '江苏省南京市 电信',
    'Mozilla/5.0 (iPad; CPU OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/121.0.6167.140 Mobile/15E148 Safari/604.1',
    'zh-CN',
    false, false,
    NULL, NULL
),
-- 父评论7：页面的垃圾评论（广告）
(
    NOW() - INTERVAL '1 day', NOW(),
    'CONTENT_TYPE_PAGE', 201,
    '低价出售阿里云服务器/域名，一年仅需99元，联系QQ：123456789',
    0,
    '优惠促销',
    'sale@fake.com',
    'http://fake-sale.com',
    'AUTHOR_TYPE_GUEST',
    'STATUS_REJECTED',
    0, 8, 0,
    '124.1xx.xx.10',
    '海外 未知',
    'Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)',
    'zh-CN',
    true, false,
    NULL, NULL
),
-- 父评论8：页面的正常评论（已发布）
(
    NOW() - INTERVAL '5 hours', NOW(),
    'CONTENT_TYPE_PAGE', 201,
    '页面的暗黑模式切换很丝滑，希望能增加自定义主题色的功能～',
    1003,
    '小七',
    'xiaoqi@139.com',
    'https://xiaoqi.tech',
    'AUTHOR_TYPE_USER',
    'STATUS_APPROVED',
    15, 0, 0,
    '220.1xx.xx.11',
    '湖北省武汉市 移动',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36',
    'zh-CN',
    false, false,
    NULL, NULL
),
-- 父评论9：页面的审核中评论（英文）
(
    NOW() - INTERVAL '2 hours', NOW(),
    'CONTENT_TYPE_PAGE', 201,
    'The SEO settings for this page are very comprehensive, it helped me improve my Google ranking!',
    2003,
    'Mike Wilson',
    'mike.wilson@yahoo.com',
    'https://mike-wilson.com',
    'AUTHOR_TYPE_USER',
    'STATUS_PENDING',
    4, 0, 0,
    '172.1xx.xx.12',
    'California, USA',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36',
    'en-US',
    false, false,
    NULL, NULL
);


-- ----------------------------
-- 插入 pages 表（页面主表）测试数据
-- ----------------------------
INSERT INTO public.pages (
    created_at, updated_at, sort_order, path, editor_type,
    status, type, slug, author_id, author_name,
    disallow_comment, redirect_url, show_in_navigation,
    template, is_custom_template, visits, custom_fields,
    custom_head, custom_foot, depth, parent_id
) VALUES
-- 页面1：首页（PAGE_TYPE_HOME=2）
(
    NOW() - INTERVAL '30 days', NOW(),
    1, '/', 'EDITOR_TYPE_MARKDOWN',
    'PAGE_STATUS_PUBLISHED', 'PAGE_TYPE_HOME', 'home',
    1, 'GoWind 官方',
    false, '', true,
    'default-home', false, 15890,
    '{"banner_show": "true", "banner_delay": "3000", "show_hot_articles": "true"}'::jsonb,
    '<meta name="viewport" content="width=device-width, initial-scale=1.0">',
    '<script src="/js/home-footer.js"></script>',
    0, NULL
),
-- 页面2：关于我们（普通页面 PAGE_TYPE_DEFAULT=1）
(
    NOW() - INTERVAL '25 days', NOW(),
    2, '/about', 'EDITOR_TYPE_RICH_TEXT',
    'PAGE_STATUS_PUBLISHED', 'PAGE_TYPE_DEFAULT', 'about',
    1, 'GoWind 官方',
    true, '', true,
    'default-static', false, 8760,
    '{"show_team_avatar": "true", "team_size": "15", "founded_year": "2024"}'::jsonb,
    '', '',
    0, NULL
),
-- 页面3：文档中心（普通页面 PAGE_TYPE_DEFAULT=1）
(
    NOW() - INTERVAL '20 days', NOW(),
    3, '/docs', 'EDITOR_TYPE_MARKDOWN',
    'PAGE_STATUS_PUBLISHED', 'PAGE_TYPE_DEFAULT', 'docs',
    1, 'GoWind 官方',
    false, '', true,
    'default-docs', false, 12540,
    '{"sidebar_collapse": "false", "edit_on_github": "true", "github_repo": "gowind/cms-docs"}'::jsonb,
    '<link rel="stylesheet" href="/css/docs.css">',
    '<script src="/js/docs-sidebar.js"></script>',
    0, NULL
),
-- 页面4：快速开始（普通页面 PAGE_TYPE_DEFAULT=1，二级）
(
    NOW() - INTERVAL '18 days', NOW(),
    1, '/docs/quick-start', 'EDITOR_TYPE_MARKDOWN',
    'PAGE_STATUS_PUBLISHED', 'PAGE_TYPE_DEFAULT', 'quick-start',
    1, 'GoWind 官方',
    false, '', true,
    'default-docs', false, 9870,
    '{"difficulty": "beginner", "estimated_time": "5分钟"}'::jsonb,
    '', '',
    1, 3
),
-- 页面5：404错误页（PAGE_TYPE_ERROR_404=3，特殊页面）
(
    NOW() - INTERVAL '15 days', NOW(),
    0, '/404', 'EDITOR_TYPE_MARKDOWN',
    'PAGE_STATUS_PUBLISHED', 'PAGE_TYPE_ERROR_404', '404',
    1, 'GoWind 官方',
    true, '', false,
    'default-error', false, 5680,
    '{"show_search": "true", "show_home_button": "true", "custom_message": "您访问的页面不存在～"}'::jsonb,
    '', '<script src="/js/404-tracking.js"></script>',
    0, NULL
),
-- 页面6：500错误页（PAGE_TYPE_ERROR_500=4，特殊页面）
(
    NOW() - INTERVAL '15 days', NOW(),
    0, '/500', 'EDITOR_TYPE_MARKDOWN',
    'PAGE_STATUS_PUBLISHED', 'PAGE_TYPE_ERROR_500', '500',
    1, 'GoWind 官方',
    true, '', false,
    'default-error', false, 1230,
    '{"show_contact_button": "true", "maintenance_phone": "400-123-4567"}'::jsonb,
    '', '',
    0, NULL
),
-- 页面7：登录页（自定义逻辑页 PAGE_TYPE_CUSTOM=5）
(
    NOW() - INTERVAL '12 days', NOW(),
    0, '/login', 'EDITOR_TYPE_RICH_TEXT', -- 自定义编辑器（无内容编辑）
    'PAGE_STATUS_PUBLISHED', 'PAGE_TYPE_CUSTOM', 'login',
    1, 'GoWind 官方',
    true, '', false,
    'custom-login', true, 8920,
    '{"show_captcha": "true", "remember_me_days": "7", "oauth_github": "true", "oauth_google": "false"}'::jsonb,
    '<link rel="stylesheet" href="/css/login.css">',
    '<script src="/js/login-validation.js"></script>',
    0, NULL
),
-- 页面8：隐私政策（普通页面 PAGE_TYPE_DEFAULT=1）
(
    NOW() - INTERVAL '12 days', NOW(),
    4, '/privacy', 'EDITOR_TYPE_RICH_TEXT',
    'PAGE_STATUS_PUBLISHED', 'PAGE_TYPE_DEFAULT', 'privacy-policy',
    1, 'GoWind 官方',
    true, '', true,
    'default-static', false, 4320,
    '{"last_updated": "2024-03-01", "version": "1.0"}'::jsonb,
    '', '',
    0, NULL
),
-- 页面9：注册页（自定义逻辑页 PAGE_TYPE_CUSTOM=5）
(
    NOW() - INTERVAL '10 days', NOW(),
    0, '/register', 'EDITOR_TYPE_RICH_TEXT',
    'PAGE_STATUS_PUBLISHED', 'PAGE_TYPE_CUSTOM', 'register',
    1, 'GoWind 官方',
    true, '', false,
    'custom-register', true, 5430,
    '{"need_email_verify": "true", "default_role": "user", "invite_code_required": "false"}'::jsonb,
    '<link rel="stylesheet" href="/css/register.css">',
    '<script src="/js/register-verify.js"></script>',
    0, NULL
);

-- ----------------------------
-- 插入 page_translations 表（页面多语言翻译）测试数据
-- ----------------------------
INSERT INTO public.page_translations (
    created_at, updated_at, page_id, language_code, title,
    slug, summary, content, original_content, thumbnail,
    cover_image, full_path, word_count, meta_keywords,
    meta_description, seo_title
) VALUES
-- ========== 页面1（首页） - 中文翻译（真实换行） ==========
(
    NOW() - INTERVAL '30 days', NOW(), 1, 'zh-CN',
    'GoWind CMS - 高性能Go语言多站点多语言CMS系统',
    'home',
    'GoWind CMS 是基于Go+Vue3开发的高性能、轻量级CMS系统，支持多站点、多语言、多租户，开箱即用。',
    '# GoWind CMS

## 核心特性
- 多站点管理：一个系统管理多个独立站点
- 多语言支持：内置多语言翻译体系
- 高性能：基于Go语言开发，QPS可达10万+
- 易扩展：插件化架构，支持自定义模板/主题',
    '# GoWind CMS

## 核心特性
- 多站点管理：一个系统管理多个独立站点
- 多语言支持：内置多语言翻译体系
- 高性能：基于Go语言开发，QPS可达10万+
- 易扩展：插件化架构，支持自定义模板/主题',
    '/images/thumbnails/home-zh.jpg',
    '/images/covers/home-zh.jpg',
    '/',
    520,
    'GoWind,CMS,Go语言,多站点,多语言,高性能',
    'GoWind CMS 是基于Go+Vue3开发的高性能CMS系统，支持多站点、多语言、多租户，开箱即用。',
    'GoWind CMS - 高性能Go语言多站点多语言CMS系统'
),
-- ========== 页面1（首页） - 英文翻译（真实换行） ==========
(
    NOW() - INTERVAL '30 days', NOW(), 1, 'en-US',
    'GoWind CMS - High Performance Go CMS for Multi-site & Multi-language',
    'home',
    'GoWind CMS is a high-performance, lightweight CMS built with Go+Vue3, supporting multi-site, multi-language, and multi-tenant features, ready to use out of the box.',
    '# GoWind CMS

## Core Features
- Multi-site Management: Manage multiple independent sites with one system
- Multi-language Support: Built-in multi-language translation system
- High Performance: Developed with Go, QPS up to 100,000+
- Extensible: Plugin-based architecture, support custom templates/themes',
    '# GoWind CMS

## Core Features
- Multi-site Management: Manage multiple independent sites with one system
- Multi-language Support: Built-in multi-language translation system
- High Performance: Developed with Go, QPS up to 100,000+
- Extensible: Plugin-based architecture, support custom templates/themes',
    '/images/thumbnails/home-en.jpg',
    '/images/covers/home-en.jpg',
    '/en',
    480,
    'GoWind,CMS,Go,Multi-site,Multi-language,High Performance',
    'GoWind CMS is a high-performance CMS built with Go+Vue3, supporting multi-site, multi-language, and multi-tenant features, ready to use out of the box.',
    'GoWind CMS - High Performance Go CMS for Multi-site & Multi-language'
),
-- ========== 页面5（404错误页） - 中文翻译（真实换行） ==========
(
    NOW() - INTERVAL '15 days', NOW(), 5, 'zh-CN',
    '404 - 页面不存在',
    '404',
    '您访问的页面不存在或已被删除，可返回首页或使用搜索功能查找相关内容。',
    '# 404 - 页面不存在

抱歉，您访问的页面不存在或已被删除。

### 您可以：
- 点击 <a href="/">首页</a> 返回网站首页
- 使用下方搜索框查找相关内容
- 联系客服：400-123-4567',
    '# 404 - 页面不存在

抱歉，您访问的页面不存在或已被删除。

### 您可以：
- 点击 <a href="/">首页</a> 返回网站首页
- 使用下方搜索框查找相关内容
- 联系客服：400-123-4567',
    '/images/thumbnails/404-zh.jpg',
    '/images/covers/404-zh.jpg',
    '/404',
    280,
    '404,页面不存在,返回首页',
    '您访问的页面不存在或已被删除，可返回首页或使用搜索功能查找相关内容。',
    '404 - 页面不存在 | GoWind CMS'
),
-- ========== 页面5（404错误页） - 英文翻译（真实换行） ==========
(
    NOW() - INTERVAL '15 days', NOW(), 5, 'en-US',
    '404 - Page Not Found',
    '404',
    'The page you are looking for does not exist or has been deleted. You can return to the homepage or use the search function to find relevant content.',
    '# 404 - Page Not Found

Sorry, the page you are looking for does not exist or has been deleted.

### You can:
- Click <a href="/en">Home</a> to return to the homepage
- Use the search box below to find relevant content
- Contact support: 400-123-4567',
    '# 404 - Page Not Found

Sorry, the page you are looking for does not exist or has been deleted.

### You can:
- Click <a href="/en">Home</a> to return to the homepage
- Use the search box below to find relevant content
- Contact support: 400-123-4567',
    '/images/thumbnails/404-en.jpg',
    '/images/covers/404-en.jpg',
    '/en/404',
    260,
    '404,Page Not Found,Return Home',
    'The page you are looking for does not exist or has been deleted. You can return to the homepage or use the search function to find relevant content.',
    '404 - Page Not Found | GoWind CMS'
),
-- ========== 页面7（登录页） - 中文翻译 ==========
(
    NOW() - INTERVAL '12 days', NOW(), 7, 'zh-CN',
    '登录 - GoWind CMS',
    'login',
    '登录GoWind CMS系统，管理您的站点内容、配置和用户。',
    '', -- 自定义逻辑页无Markdown/富文本内容
    '',
    '/images/thumbnails/login-zh.jpg',
    '/images/covers/login-zh.jpg',
    '/login',
    0,
    '登录,GoWind CMS,账号登录,用户中心',
    '登录GoWind CMS系统，管理您的站点内容、配置和用户。',
    '登录 - GoWind CMS'
),
-- ========== 页面2（关于我们） - 中文翻译 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 2, 'zh-CN',
    '关于我们 - GoWind CMS',
    'about',
    'GoWind 团队成立于2024年，专注于高性能开源CMS系统的研发，致力于为开发者提供简单、高效的内容管理解决方案。',
    '<h1>关于 GoWind CMS</h1>
<p>GoWind 团队成立于2024年，核心成员均有10年以上互联网研发经验，专注于高性能开源CMS系统的研发。</p>
<h2>我们的使命</h2>
<p>让内容管理更简单，让开发者更高效。</p>
<h2>团队规模</h2>
<p>目前团队共15人，涵盖后端、前端、产品、设计等多个领域。</p>',
    '<h1>关于 GoWind CMS</h1>
<p>GoWind 团队成立于2024年，核心成员均有10年以上互联网研发经验，专注于高性能开源CMS系统的研发。</p>
<h2>我们的使命</h2>
<p>让内容管理更简单，让开发者更高效。</p>
<h2>团队规模</h2>
<p>目前团队共15人，涵盖后端、前端、产品、设计等多个领域。</p>',
    '/images/thumbnails/about-zh.jpg',
    '/images/covers/about-zh.jpg',
    '/about',
    850,
    'GoWind,关于我们,团队介绍,CMS研发',
    'GoWind 团队成立于2024年，专注于高性能开源CMS系统的研发，致力于为开发者提供简单、高效的内容管理解决方案。',
    '关于我们 - GoWind CMS'
);


-- ----------------------------
-- 插入 categories 表（分类主表）测试数据
-- ----------------------------
INSERT INTO public.categories (
    created_at, updated_at, sort_order, path, status,
    depth, is_nav, icon, post_count, direct_post_count,
    custom_fields, parent_id
) VALUES
-- 分类1：博客（一级、显示在导航、已发布）
(
    NOW() - INTERVAL '30 days', NOW(),
    1, '/blog', 'CATEGORY_STATUS_ACTIVE',
    0, true, 'icon-blog', 128, 128,
    '{"show_banner": "true", "banner_image": "/images/category/blog-banner.jpg", "layout": "list"}'::jsonb,
    NULL
),
-- 分类2：技术文档（一级、显示在导航、已发布）
(
    NOW() - INTERVAL '28 days', NOW(),
    2, '/docs', 'CATEGORY_STATUS_ACTIVE',
    0, true, 'icon-docs', 89, 0,
    '{"sidebar_type": "collapsible", "edit_link": "https://github.com/gowind/docs/edit/main"}'::jsonb,
    NULL
),
-- 分类3：产品中心（一级、显示在导航、已发布）
(
    NOW() - INTERVAL '25 days', NOW(),
    3, '/products', 'CATEGORY_STATUS_ACTIVE',
    0, true, 'icon-product', 45, 45,
    '{"show_filter": "true", "filter_fields": "price, type", "default_sort": "newest"}'::jsonb,
    NULL
),
-- 分类4：安装部署（二级、不显示在导航、已发布，父ID=2）
(
    NOW() - INTERVAL '22 days', NOW(),
    1, '/docs/deployment', 'CATEGORY_STATUS_ACTIVE',
    1, false, 'icon-deploy', 32, 32,
    '{"difficulty": "beginner", "estimated_time": "10分钟"}'::jsonb,
    2
),
-- 分类5：环境配置（三级、不显示在导航、已发布，父ID=4）
(
    NOW() - INTERVAL '20 days', NOW(),
    1, '/docs/deployment/env', 'CATEGORY_STATUS_ACTIVE',
    2, false, 'icon-env', 18, 18,
    '{"os_support": "Linux, Windows, macOS", "min_requirements": "2GB RAM, 1 CPU"}'::jsonb,
    4
),
-- 分类6：功能教程（二级、不显示在导航、已发布，父ID=2）
(
    NOW() - INTERVAL '18 days', NOW(),
    2, '/docs/tutorials', 'CATEGORY_STATUS_ACTIVE',
    1, false, 'icon-tutorial', 35, 35,
    '{"video_support": "true", "tutorial_level": "beginner, intermediate"}'::jsonb,
    2
),
-- 分类7：常见问题（二级、不显示在导航、草稿，父ID=2）
(
    NOW() - INTERVAL '15 days', NOW(),
    3, '/docs/faq', 'CATEGORY_STATUS_HIDDEN',
    1, false, 'icon-faq', 0, 0,
    '{"show_search": "true", "hot_tags": "installation configuration performance"}'::jsonb,
    2
),
-- 分类8：废弃分类（一级、不显示在导航、已归档）
(
    NOW() - INTERVAL '10 days', NOW(),
    4, '/old-category', 'CATEGORY_STATUS_ARCHIVED',
    0, false, 'icon-archive', 0, 0,
    '{"archive_reason": "内容合并至技术文档", "archive_time": "2024-03-01"}'::jsonb,
    NULL
),
-- 分类9：行业资讯（一级、显示在导航、已发布）
(
    NOW() - INTERVAL '8 days', NOW(),
    5, '/news', 'CATEGORY_STATUS_ACTIVE',
    0, true, 'icon-news', 67, 67,
    '{"show_date": "true", "author_display": "true", "comment_support": "true"}'::jsonb,
    NULL
);

-- ----------------------------
-- 插入 category_translations 表（分类多语言翻译）测试数据
-- 关联上述分类，覆盖 zh-CN/en-US 双语言，真实换行符
-- ----------------------------
INSERT INTO public.category_translations (
    created_at, updated_at, category_id, language_code,
    name, slug, description, thumbnail, cover_image,
    template, full_path, meta_keywords, meta_description, seo_title
) VALUES
-- ========== 分类1（博客） - 中文翻译 ==========
(
    NOW() - INTERVAL '30 days', NOW(), 1, 'zh-CN',
    '博客', 'blog',
    'GoWind CMS 官方博客，分享产品更新、技术实践、行业动态等内容。

内容涵盖：
- Go 语言开发技巧
- CMS 产品使用经验
- 开源项目维护心得',
    '/images/thumbnails/category-blog-zh.jpg',
    '/images/covers/category-blog-zh.jpg',
    'category-blog', '/blog',
    'GoWind,博客,技术实践,产品更新',
    'GoWind CMS 官方博客，分享产品更新、技术实践、行业动态等内容。',
    '博客 | GoWind CMS'
),
-- ========== 分类1（博客） - 英文翻译 ==========
(
    NOW() - INTERVAL '30 days', NOW(), 1, 'en-US',
    'Blog', 'blog',
    'GoWind CMS official blog, sharing product updates, technical practices, industry trends and other content.

Content includes:
- Go language development skills
- CMS product usage experience
- Open source project maintenance experience',
    '/images/thumbnails/category-blog-en.jpg',
    '/images/covers/category-blog-en.jpg',
    'category-blog', '/en/blog',
    'GoWind,Blog,Technical Practices,Product Updates',
    'GoWind CMS official blog, sharing product updates, technical practices, industry trends and other content.',
    'Blog | GoWind CMS'
),
-- ========== 分类2（技术文档） - 中文翻译 ==========
(
    NOW() - INTERVAL '28 days', NOW(), 2, 'zh-CN',
    '技术文档', 'docs',
    'GoWind CMS 完整技术文档，包含安装部署、功能教程、API 参考等内容，帮助开发者快速上手。

核心板块：
- 安装部署：环境配置、快速开始
- 功能教程：核心功能、高级用法
- API 参考：接口文档、调用示例',
    '/images/thumbnails/category-docs-zh.jpg',
    '/images/covers/category-docs-zh.jpg',
    'category-docs', '/docs',
    'GoWind,技术文档,安装部署,API参考',
    'GoWind CMS 完整技术文档，包含安装部署、功能教程、API 参考等内容，帮助开发者快速上手。',
    '技术文档 | GoWind CMS'
),
-- ========== 分类2（技术文档） - 英文翻译 ==========
(
    NOW() - INTERVAL '28 days', NOW(), 2, 'en-US',
    'Documentation', 'docs',
    'Complete technical documentation for GoWind CMS, including installation and deployment, feature tutorials, API references, etc., to help developers get started quickly.

Core sections:
- Installation & Deployment: Environment configuration, quick start
- Feature Tutorials: Core features, advanced usage
- API Reference: Interface documentation, call examples',
    '/images/thumbnails/category-docs-en.jpg',
    '/images/covers/category-docs-en.jpg',
    'category-docs', '/en/docs',
    'GoWind,Documentation,Installation,API Reference',
    'Complete technical documentation for GoWind CMS, including installation and deployment, feature tutorials, API references, etc., to help developers get started quickly.',
    'Documentation | GoWind CMS'
),
-- ========== 分类4（安装部署） - 中文翻译 ==========
(
    NOW() - INTERVAL '22 days', NOW(), 4, 'zh-CN',
    '安装部署', 'deployment',
    'GoWind CMS 安装部署相关文档，包含环境配置、快速安装、Docker 部署等内容，适配不同操作系统。',
    '/images/thumbnails/category-deploy-zh.jpg',
    '/images/covers/category-deploy-zh.jpg',
    'category-docs-child', '/docs/deployment',
    'GoWind,安装部署,环境配置,Docker部署',
    'GoWind CMS 安装部署相关文档，包含环境配置、快速安装、Docker 部署等内容，适配不同操作系统。',
    '安装部署 | GoWind CMS 技术文档'
),
-- ========== 分类5（环境配置） - 中文翻译 ==========
(
    NOW() - INTERVAL '20 days', NOW(), 5, 'zh-CN',
    '环境配置', 'env',
    'GoWind CMS 环境配置详细教程，包含 Linux/Windows/macOS 系统的环境搭建、依赖安装、端口配置等。',
    '/images/thumbnails/category-env-zh.jpg',
    '/images/covers/category-env-zh.jpg',
    'category-docs-grandchild', '/docs/deployment/env',
    'GoWind,环境配置,依赖安装,端口配置',
    'GoWind CMS 环境配置详细教程，包含 Linux/Windows/macOS 系统的环境搭建、依赖安装、端口配置等。',
    '环境配置 | GoWind CMS 安装部署'
),
-- ========== 分类3（产品中心） - 中文翻译 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 3, 'zh-CN',
    '产品中心', 'products',
    'GoWind 旗下全系列产品，包含开源版、企业版、云服务等，满足不同规模用户的内容管理需求。',
    '/images/thumbnails/category-product-zh.jpg',
    '/images/covers/category-product-zh.jpg',
    'category-product', '/products',
    'GoWind,产品中心,开源版,企业版,云服务',
    'GoWind 旗下全系列产品，包含开源版、企业版、云服务等，满足不同规模用户的内容管理需求。',
    '产品中心 | GoWind CMS'
),
-- ========== 分类9（行业资讯） - 中文翻译 ==========
(
    NOW() - INTERVAL '8 days', NOW(), 9, 'zh-CN',
    '行业资讯', 'news',
    'CMS 行业最新资讯、政策解读、市场分析，帮助您把握行业发展趋势。',
    '/images/thumbnails/category-news-zh.jpg',
    '/images/covers/category-news-zh.jpg',
    'category-news', '/news',
    'GoWind,行业资讯,CMS行业,市场分析',
    'CMS 行业最新资讯、政策解读、市场分析，帮助您把握行业发展趋势。',
    '行业资讯 | GoWind CMS'
);


-- ----------------------------
-- 插入 posts 表
-- ----------------------------
INSERT INTO public.posts (
    created_at, updated_at, sort_order, editor_type,
    status, slug, disallow_comment, in_progress,
    auto_summary, is_featured, visits, likes,
    comment_count, author_id, author_name, password_hash,
    custom_fields, category_ids, tag_ids
) VALUES
-- 文章1：GoWind CMS 快速上手（已发布、精选）
(
    NOW() - INTERVAL '30 days', NOW(),
    1, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'gowind-cms-quick-start',
    false, false, true, true,
    15890, 892, 156,
    1, 'GoWind 官方', '',
    '{"show_toc": "true", "toc_depth": "3", "allow_copy": "true", "copyright_notice": "GoWind 官方原创"}'::jsonb,
    '[2,4]'::jsonb,
    '[1,2,3]'::jsonb
),
-- 文章2：GoWind v2.0 版本发布公告（已发布、精选）
(
    NOW() - INTERVAL '25 days', NOW(),
    2, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'gowind-v2-0-release',
    false, false, true, true,
    9870, 654, 89,
    1, 'GoWind 官方', '',
    '{"show_changelog": "true", "release_date": "2024-03-01", "upgrade_guide_url": "/docs/upgrade/v2.0"}'::jsonb,
    '[1,3]'::jsonb,
    '[4,5,6]'::jsonb
),
-- 文章3：Linux 环境下部署 GoWind CMS（已发布）
(
    NOW() - INTERVAL '22 days', NOW(),
    3, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'deploy-gowind-on-linux',
    false, false, true, false,
    7650, 432, 78,
    1001, '张三', '',
    '{"os_type": "Linux", "distro": "Ubuntu, CentOS", "tested_version": "v1.9.0"}'::jsonb,
    '[2,4,5]'::jsonb,
    '[1,7,8]'::jsonb
),
-- 文章4：2024 CMS 行业发展趋势分析（已发布）
(
    NOW() - INTERVAL '20 days', NOW(),
    4, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', '2024-cms-industry-trends',
    false, false, true, false,
    6540, 389, 67,
    1002, '李四', '',
    '{"data_source": "IDC 2024 行业报告", "chart_support": "true", "downloadable": "true"}'::jsonb,
    '[1,9]'::jsonb,
    '[9,10,11]'::jsonb
),
-- 文章5：GoWind CMS 自定义模板开发（草稿、未完成）
(
    NOW() - INTERVAL '15 days', NOW(),
    5, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_DRAFT', 'gowind-custom-template-dev',
    true, true, false, false,
    1230, 0, 0,
    1001, '张三', '',
    '{"dev_status": "50%", "expected_release": "2024-04-01", "required_skills": "Go, Vue3, HTML/CSS"}'::jsonb,
    '[2,6]'::jsonb,
    '[1,12,13]'::jsonb
),
-- 文章6：GoWind 企业版功能详解（加密、已发布）
(
    NOW() - INTERVAL '12 days', NOW(),
    6, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'gowind-enterprise-features',
    true, false, true, false,
    4320, 256, 0,
    1, 'GoWind 官方', '$2a$10$89jZk54G89sdkf89sdf89sd89sdf89sdf89sdf',
    '{"is_enterprise": "true", "price_range": "¥9999-¥19999", "trial_available": "true"}'::jsonb,
    '[3]'::jsonb,
    '[5,14,15]'::jsonb
),
-- 文章7：常见问题解答（草稿）
(
    NOW() - INTERVAL '10 days', NOW(),
    7, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_DRAFT', 'gowind-faq',
    true, false, true, false,
    890, 0, 0,
    1, 'GoWind 官方', '',
    '{"faq_category": "installation, configuration, performance", "update_frequency": "monthly"}'::jsonb,
    '[2,7]'::jsonb,
    '[16,17,18]'::jsonb
),
-- 文章8：GoWind CMS 性能优化指南（已发布、精选）
(
    NOW() - INTERVAL '8 days', NOW(),
    8, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'gowind-cms-performance-optimization',
    false, false, true, true,
    5680, 789, 123,
    1003, '王五', '',
    '{"benchmark_data": "true", "qps_before": "50000", "qps_after": "100000", "optimization_points": "DB, Cache, Code"}'::jsonb,
    '[1,2]'::jsonb,
    '[1,19,20]'::jsonb
);

-- ----------------------------
-- 插入 post_translations 表
-- ----------------------------
INSERT INTO public.post_translations (
    created_at, updated_at, post_id, language_code, title,
    slug, summary, content, original_content, thumbnail,
    template, full_path, word_count, meta_keywords,
    meta_description, seo_title
) VALUES
-- ========== 文章1：GoWind CMS 快速上手 - 中文 ==========
(
    NOW() - INTERVAL '30 days', NOW(), 1, 'zh-CN',
    'GoWind CMS 快速上手：5分钟搭建你的第一个CMS站点',
    'gowind-cms-quick-start',
    '本文带你5分钟快速搭建GoWind CMS站点，涵盖环境准备、代码克隆、配置启动、初始登录全流程。',
    '# GoWind CMS 快速上手

## 环境准备
### 前置依赖
- Go 1.21+（推荐1.22最新版）
- PostgreSQL 14+
- Git（可选）
- Docker（可选）

## 安装步骤
1. 克隆代码仓库：git clone https://github.com/gowind/cms.git
2. 配置环境变量：复制.env.example为.env并修改数据库配置
3. 启动服务：go run main.go 或 docker-compose up -d
4. 初始登录：http://localhost:8080，默认账号admin/admin

> 首次登录请立即修改密码！',
    '# GoWind CMS 快速上手

## 环境准备
### 前置依赖
- Go 1.21+（推荐1.22最新版）
- PostgreSQL 14+
- Git（可选）
- Docker（可选）

## 安装步骤
1. 克隆代码仓库：git clone https://github.com/gowind/cms.git
2. 配置环境变量：复制.env.example为.env并修改数据库配置
3. 启动服务：go run main.go 或 docker-compose up -d
4. 初始登录：http://localhost:8080，默认账号admin/admin

> 首次登录请立即修改密码！',
    '/images/thumbnails/post-quick-start-zh.jpg',
    'post-default', '/blog/gowind-cms-quick-start',
    2580,
    'GoWind,CMS,快速上手,安装部署,Go语言,PostgreSQL',
    '本文带你5分钟快速搭建GoWind CMS站点，涵盖环境准备、代码克隆、配置启动、初始登录全流程。',
    'GoWind CMS 快速上手：5分钟搭建你的第一个CMS站点 | GoWind 官方文档'
),
-- ========== 文章2：GoWind v2.0 版本发布公告 - 中文 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 2, 'zh-CN',
    'GoWind CMS v2.0 正式发布：新增多租户、性能提升100%',
    'gowind-v2-0-release',
    'GoWind CMS v2.0版本发布，核心更新多租户支持、性能优化、UI重构，QPS突破10万。',
    '# GoWind CMS v2.0 正式发布

## 发布说明
GoWind CMS v2.0于2024年3月1日发布，是开源以来的重大版本更新！

## 核心新功能
1. 多租户支持：单实例部署多套独立站点，数据隔离
2. 性能优化：QPS从5万提升至10万，响应时间降低60%
3. UI重构：基于Vue3+Element Plus重构后台，移动端适配

## 升级指南
- 从v1.9升级：备份数据库后执行go run scripts/upgrade/v2.0.go
- 全新安装：直接克隆v2.0分支代码部署',
    '# GoWind CMS v2.0 正式发布

## 发布说明
GoWind CMS v2.0于2024年3月1日发布，是开源以来的重大版本更新！

## 核心新功能
1. 多租户支持：单实例部署多套独立站点，数据隔离
2. 性能优化：QPS从5万提升至10万，响应时间降低60%
3. UI重构：基于Vue3+Element Plus重构后台，移动端适配

## 升级指南
- 从v1.9升级：备份数据库后执行go run scripts/upgrade/v2.0.go
- 全新安装：直接克隆v2.0分支代码部署',
    '/images/thumbnails/post-v2-release-zh.jpg',
    'post-announcement', '/blog/gowind-v2-0-release',
    3200,
    'GoWind,CMS,v2.0,版本发布,多租户,性能优化,UI重构',
    'GoWind CMS v2.0版本发布，核心更新多租户支持、性能优化、UI重构，QPS突破10万。',
    'GoWind CMS v2.0 正式发布：新增多租户、性能提升100% | GoWind 官方公告'
),
-- ========== 文章3：Linux 环境下部署 GoWind CMS - 中文 ==========
(
    NOW() - INTERVAL '22 days', NOW(), 3, 'zh-CN',
    'Linux 环境下部署 GoWind CMS（Ubuntu/CentOS通用）',
    'deploy-gowind-on-linux',
    '详解Linux环境（Ubuntu/CentOS）下GoWind CMS的部署步骤，含依赖安装、端口配置、开机自启。',
    '# Linux 环境下部署 GoWind CMS

## 适用系统
- Ubuntu 20.04/22.04
- CentOS 7/8

## 依赖安装
### Ubuntu
apt update && apt install -y golang postgresql git
### CentOS
yum install -y golang postgresql git

## 部署步骤
1. 创建数据库：createdb gowind
2. 克隆代码：git clone https://github.com/gowind/cms.git
3. 配置数据库连接：修改.env文件
4. 启动服务：nohup go run main.go > app.log 2>&1 &

## 开机自启
创建systemd服务文件：/etc/systemd/system/gowind.service',
    '# Linux 环境下部署 GoWind CMS

## 适用系统
- Ubuntu 20.04/22.04
- CentOS 7/8

## 依赖安装
### Ubuntu
apt update && apt install -y golang postgresql git
### CentOS
yum install -y golang postgresql git

## 部署步骤
1. 创建数据库：createdb gowind
2. 克隆代码：git clone https://github.com/gowind/cms.git
3. 配置数据库连接：修改.env文件
4. 启动服务：nohup go run main.go > app.log 2>&1 &

## 开机自启
创建systemd服务文件：/etc/systemd/system/gowind.service',
    '/images/thumbnails/post-deploy-linux-zh.jpg',
    'post-tech', '/blog/deploy-gowind-on-linux',
    2800,
    'GoWind,CMS,Linux部署,Ubuntu,CentOS,Go语言',
    '详解Linux环境（Ubuntu/CentOS）下GoWind CMS的部署步骤，含依赖安装、端口配置、开机自启。',
    'Linux 环境下部署 GoWind CMS（Ubuntu/CentOS通用） | GoWind 技术博客'
),
-- ========== 文章4：2024 CMS 行业发展趋势分析 - 中文 ==========
(
    NOW() - INTERVAL '20 days', NOW(), 4, 'zh-CN',
    '2024 CMS 行业发展趋势：轻量化、私有化、AI赋能',
    '2024-cms-industry-trends',
    '基于IDC 2024行业报告，分析CMS行业三大趋势：轻量化、私有化部署、AI智能赋能。',
    '# 2024 CMS 行业发展趋势分析

## 数据来源
IDC 2024年全球CMS市场研究报告

## 核心趋势
1. 轻量化：轻量级CMS占比提升至65%，替代重型系统
2. 私有化：企业级用户私有化部署需求增长40%
3. AI赋能：AI生成内容、智能排版成为标配功能

## 市场规模
2024年全球CMS市场规模预计达89亿美元，年增长率18%。

## 国内趋势
国产化替代加速，Go/Java语言开发的CMS占比提升。',
    '# 2024 CMS 行业发展趋势分析

## 数据来源
IDC 2024年全球CMS市场研究报告

## 核心趋势
1. 轻量化：轻量级CMS占比提升至65%，替代重型系统
2. 私有化：企业级用户私有化部署需求增长40%
3. AI赋能：AI生成内容、智能排版成为标配功能

## 市场规模
2024年全球CMS市场规模预计达89亿美元，年增长率18%。

## 国内趋势
国产化替代加速，Go/Java语言开发的CMS占比提升。',
    '/images/thumbnails/post-cms-trends-zh.jpg',
    'post-analysis', '/blog/2024-cms-industry-trends',
    2600,
    'CMS,2024趋势,行业分析,轻量化,私有化,AI赋能',
    '基于IDC 2024行业报告，分析CMS行业三大趋势：轻量化、私有化部署、AI智能赋能。',
    '2024 CMS 行业发展趋势：轻量化、私有化、AI赋能 | GoWind 行业洞察'
),
-- ========== 文章5：GoWind CMS 自定义模板开发 - 中文（草稿） ==========
(
    NOW() - INTERVAL '15 days', NOW(), 5, 'zh-CN',
    'GoWind CMS 自定义模板开发教程（草稿）',
    'gowind-custom-template-dev',
    'GoWind CMS自定义模板开发教程，涵盖模板语法、数据调用、样式定制，当前开发进度50%。',
    '# GoWind CMS 自定义模板开发

> 本文正在编写中，开发进度50%，预计2024年4月1日完成。

## 开发准备
### 所需技能
- Go语言基础
- Vue3 + Element Plus
- HTML/CSS/JS

## 模板目录结构
/templates/custom/
  - index.tpl # 首页模板
  - post.tpl # 文章模板
  - style.css # 自定义样式

## 待编写内容
1. 模板语法详解
2. 数据调用示例
3. 自定义组件开发',
    '# GoWind CMS 自定义模板开发

> 本文正在编写中，开发进度50%，预计2024年4月1日完成。

## 开发准备
### 所需技能
- Go语言基础
- Vue3 + Element Plus
- HTML/CSS/JS

## 模板目录结构
/templates/custom/
  - index.tpl # 首页模板
  - post.tpl # 文章模板
  - style.css # 自定义样式

## 待编写内容
1. 模板语法详解
2. 数据调用示例
3. 自定义组件开发',
    '/images/thumbnails/post-template-dev-zh.jpg',
    'post-draft', '/blog/gowind-custom-template-dev',
    1800,
    'GoWind,CMS,自定义模板,开发教程,Go,Vue3',
    'GoWind CMS自定义模板开发教程，涵盖模板语法、数据调用、样式定制，当前开发进度50%。',
    'GoWind CMS 自定义模板开发教程（草稿） | GoWind 开发文档'
),
-- ========== 文章6：GoWind 企业版功能详解 - 中文（加密） ==========
(
    NOW() - INTERVAL '12 days', NOW(), 6, 'zh-CN',
    'GoWind CMS 企业版功能详解（付费专属）',
    'gowind-enterprise-features',
    'GoWind企业版专属功能：多租户管理、高级权限、数据备份、专属客服，价格9999-19999元/年。',
    '# GoWind CMS 企业版功能详解

## 专属功能
1. 多租户管理：单系统管理多站点，数据完全隔离
2. 高级权限：按角色/部门精细化权限控制
3. 数据备份：自动定时备份，支持异地容灾
4. 专属客服：7*24小时技术支持
5. 定制开发：按需定制功能模块

## 价格方案
- 基础版：¥9999/年（10租户以内）
- 标准版：¥14999/年（50租户以内）
- 旗舰版：¥19999/年（不限租户）

## 试用申请
联系客服：400-123-4567，可申请15天免费试用。',
    '# GoWind CMS 企业版功能详解

## 专属功能
1. 多租户管理：单系统管理多站点，数据完全隔离
2. 高级权限：按角色/部门精细化权限控制
3. 数据备份：自动定时备份，支持异地容灾
4. 专属客服：7*24小时技术支持
5. 定制开发：按需定制功能模块

## 价格方案
- 基础版：¥9999/年（10租户以内）
- 标准版：¥14999/年（50租户以内）
- 旗舰版：¥19999/年（不限租户）

## 试用申请
联系客服：400-123-4567，可申请15天免费试用。',
    '/images/thumbnails/post-enterprise-zh.jpg',
    'post-enterprise', '/blog/gowind-enterprise-features',
    2200,
    'GoWind,CMS,企业版,付费功能,多租户,高级权限',
    'GoWind企业版专属功能：多租户管理、高级权限、数据备份、专属客服，价格9999-19999元/年。',
    'GoWind CMS 企业版功能详解（付费专属） | GoWind 企业服务'
),
-- ========== 文章7：常见问题解答 - 中文（草稿） ==========
(
    NOW() - INTERVAL '10 days', NOW(), 7, 'zh-CN',
    'GoWind CMS 常见问题解答（FAQ）',
    'gowind-faq',
    'GoWind CMS常见问题汇总，涵盖安装、配置、性能、升级等方向，每月更新。',
    '# GoWind CMS 常见问题解答

## 安装相关
Q1：安装时提示数据库连接失败？
A1：检查.env文件中的数据库地址、端口、账号密码是否正确。

Q2：启动服务后访问不了？
A2：检查端口是否被占用，防火墙是否开放8080端口。

## 配置相关
Q3：如何开启多语言支持？
A3：在后台设置-多语言中启用，上传翻译文件。

## 待补充
- 性能优化相关问题
- 升级相关问题',
    '# GoWind CMS 常见问题解答

## 安装相关
Q1：安装时提示数据库连接失败？
A1：检查.env文件中的数据库地址、端口、账号密码是否正确。

Q2：启动服务后访问不了？
A2：检查端口是否被占用，防火墙是否开放8080端口。

## 配置相关
Q3：如何开启多语言支持？
A3：在后台设置-多语言中启用，上传翻译文件。

## 待补充
- 性能优化相关问题
- 升级相关问题',
    '/images/thumbnails/post-faq-zh.jpg',
    'post-faq', '/blog/gowind-faq',
    1500,
    'GoWind,CMS,FAQ,常见问题,安装,配置,性能',
    'GoWind CMS常见问题汇总，涵盖安装、配置、性能、升级等方向，每月更新。',
    'GoWind CMS 常见问题解答（FAQ） | GoWind 帮助中心'
),
-- ========== 文章8：GoWind CMS 性能优化指南 - 中文 ==========
(
    NOW() - INTERVAL '8 days', NOW(), 8, 'zh-CN',
    'GoWind CMS 性能优化指南：从5万QPS到10万的实战经验',
    'gowind-cms-performance-optimization',
    '分享GoWind CMS性能优化实战经验，涵盖数据库优化、缓存策略、代码层面，QPS提升100%。',
    '# GoWind CMS 性能优化指南

## 优化背景
v1.9版本QPS仅5万，响应时间200ms，无法满足高并发需求。

## 优化指标
| 指标 | 优化前 | 优化后 | 提升 |
|------|--------|--------|------|
| QPS | 50000 | 100000 | 100% |
| 响应时间 | 200ms | 80ms | 60% |
| 数据库负载 | 80% | 30% | 62.5% |

## 核心优化点
1. 数据库：新增索引、慢查询优化、读写分离
2. 缓存：Redis缓存分类/文章，动态过期策略
3. 代码：优化Goroutine、JSON序列化、静态资源压缩',
    '# GoWind CMS 性能优化指南

## 优化背景
v1.9版本QPS仅5万，响应时间200ms，无法满足高并发需求。

## 优化指标
| 指标 | 优化前 | 优化后 | 提升 |
|------|--------|--------|------|
| QPS | 50000 | 100000 | 100% |
| 响应时间 | 200ms | 80ms | 60% |
| 数据库负载 | 80% | 30% | 62.5% |

## 核心优化点
1. 数据库：新增索引、慢查询优化、读写分离
2. 缓存：Redis缓存分类/文章，动态过期策略
3. 代码：优化Goroutine、JSON序列化、静态资源压缩',
    '/images/thumbnails/post-performance-zh.jpg',
    'post-tech', '/blog/gowind-cms-performance-optimization',
    3000,
    'GoWind,CMS,性能优化,QPS,数据库优化,缓存策略,Go语言',
    '分享GoWind CMS性能优化实战经验，涵盖数据库优化、缓存策略、代码层面，QPS提升100%。',
    'GoWind CMS 性能优化指南：从5万QPS到10万的实战经验 | GoWind 技术博客'
);


-- ----------------------------
-- 插入 tags 表（标签主表）测试数据（ID 1-20，匹配前文文章tag_ids）
-- ----------------------------
INSERT INTO public.tags (
    created_at, updated_at, sort_order, status,
    color, icon, "group", is_featured, post_count
) VALUES
-- 1. Go语言（技术分组、精选、启用）
(
    NOW() - INTERVAL '60 days', NOW(),
    1, 'TAG_STATUS_ACTIVE',
    '#00ADD8', 'icon-go', 'TECH', true, 15
),
-- 2. CMS（产品分组、精选、启用）
(
    NOW() - INTERVAL '60 days', NOW(),
    2, 'TAG_STATUS_ACTIVE',
    '#4CAF50', 'icon-cms', 'PRODUCT', true, 22
),
-- 3. 快速上手（教程分组、启用）
(
    NOW() - INTERVAL '55 days', NOW(),
    3, 'TAG_STATUS_ACTIVE',
    '#FF9800', 'icon-quickstart', 'TUTORIAL', false, 8
),
-- 4. 版本更新（产品分组、启用）
(
    NOW() - INTERVAL '50 days', NOW(),
    4, 'TAG_STATUS_ACTIVE',
    '#9C27B0', 'icon-update', 'PRODUCT', false, 6
),
-- 5. 新功能（产品分组、精选、启用）
(
    NOW() - INTERVAL '50 days', NOW(),
    5, 'TAG_STATUS_ACTIVE',
    '#E91E63', 'icon-feature', 'PRODUCT', true, 10
),
-- 6. 升级指南（教程分组、启用）
(
    NOW() - INTERVAL '45 days', NOW(),
    6, 'TAG_STATUS_ACTIVE',
    '#607D8B', 'icon-upgrade', 'TUTORIAL', false, 4
),
-- 7. Linux（技术分组、启用）
(
    NOW() - INTERVAL '45 days', NOW(),
    7, 'TAG_STATUS_ACTIVE',
    '#FCC624', 'icon-linux', 'TECH', false, 7
),
-- 8. 部署教程（教程分组、启用）
(
    NOW() - INTERVAL '40 days', NOW(),
    8, 'TAG_STATUS_ACTIVE',
    '#795548', 'icon-deploy', 'TUTORIAL', false, 5
),
-- 9. 行业分析（行业分组、精选、启用）
(
    NOW() - INTERVAL '40 days', NOW(),
    9, 'TAG_STATUS_ACTIVE',
    '#673AB7', 'icon-analysis', 'INDUSTRY', true, 9
),
-- 10. 2024趋势（行业分组、启用）
(
    NOW() - INTERVAL '35 days', NOW(),
    10, 'TAG_STATUS_ACTIVE',
    '#3F51B5', 'icon-trend', 'INDUSTRY', false, 3
),
-- 11. 轻量化（行业分组、启用）
(
    NOW() - INTERVAL '35 days', NOW(),
    11, 'TAG_STATUS_ACTIVE',
    '#009688', 'icon-light', 'INDUSTRY', false, 2
),
-- 12. 自定义模板（技术分组、启用）
(
    NOW() - INTERVAL '30 days', NOW(),
    12, 'TAG_STATUS_ACTIVE',
    '#8BC34A', 'icon-template', 'TECH', false, 4
),
-- 13. 开发教程（教程分组、启用）
(
    NOW() - INTERVAL '30 days', NOW(),
    13, 'TAG_STATUS_ACTIVE',
    '#CDDC39', 'icon-dev', 'TUTORIAL', false, 6
),
-- 14. 企业版（产品分组、启用）
(
    NOW() - INTERVAL '25 days', NOW(),
    14, 'TAG_STATUS_ACTIVE',
    '#FF5722', 'icon-enterprise', 'PRODUCT', false, 3
),
-- 15. 付费功能（产品分组、启用）
(
    NOW() - INTERVAL '25 days', NOW(),
    15, 'TAG_STATUS_ACTIVE',
    '#795548', 'icon-paid', 'PRODUCT', false, 2
),
-- 16. FAQ（帮助分组、启用）
(
    NOW() - INTERVAL '20 days', NOW(),
    16, 'TAG_STATUS_ACTIVE',
    '#9E9E9E', 'icon-faq', 'HELP', false, 5
),
-- 17. 问题解答（帮助分组、启用）
(
    NOW() - INTERVAL '20 days', NOW(),
    17, 'TAG_STATUS_ACTIVE',
    '#607D8B', 'icon-answer', 'HELP', false, 4
),
-- 18. 故障排除（帮助分组、启用）
(
    NOW() - INTERVAL '15 days', NOW(),
    18, 'TAG_STATUS_ACTIVE',
    '#F44336', 'icon-troubleshoot', 'HELP', false, 3
),
-- 19. 性能优化（技术分组、精选、启用）
(
    NOW() - INTERVAL '15 days', NOW(),
    19, 'TAG_STATUS_ACTIVE',
    '#2196F3', 'icon-performance', 'TECH', true, 8
),
-- 20. 高并发（技术分组、启用）
(
    NOW() - INTERVAL '10 days', NOW(),
    20, 'TAG_STATUS_ACTIVE',
    '#4CAF50', 'icon-concurrency', 'TECH', false, 2
);

-- ----------------------------
-- 插入 tag_translations 表（标签多语言翻译）测试数据
-- 覆盖zh-CN/en-US，内容为单行字符串+ \n 换行，可直接拷贝
-- ----------------------------
INSERT INTO public.tag_translations (
    created_at, updated_at, tag_id, language_code,
    name, slug, description, cover_image,
    template, full_path, canonical_url,
    meta_keywords, meta_description, seo_title
) VALUES
-- ========== 标签1：Go语言 - 中文 ==========
(
    NOW() - INTERVAL '60 days', NOW(), 1, 'zh-CN',
    'Go语言', 'go',
    'Go（Golang）是Google开发的开源编程语言，简洁、高效、高并发，GoWind CMS核心开发语言。',
    '/images/tags/go-zh.jpg',
    'tag-default', '/tags/go', 'https://gowind.com/tags/go',
    'Go语言,Golang,GoWind,高并发,编程语言',
    'Go（Golang）是Google开发的开源编程语言，简洁、高效、高并发，GoWind CMS核心开发语言。',
    'Go语言 | GoWind 标签'
),
-- ========== 标签1：Go语言 - 英文 ==========
(
    NOW() - INTERVAL '60 days', NOW(), 1, 'en-US',
    'Go Language', 'go',
    'Go (Golang) is an open-source programming language developed by Google, concise, efficient, and high-concurrency, the core development language of GoWind CMS.',
    '/images/tags/go-en.jpg',
    'tag-default', '/en/tags/go', 'https://gowind.com/en/tags/go',
    'Go Language,Golang,GoWind,High Concurrency,Programming Language',
    'Go (Golang) is an open-source programming language developed by Google, concise, efficient, and high-concurrency, the core development language of GoWind CMS.',
    'Go Language | GoWind Tags'
),
-- ========== 标签2：CMS - 中文 ==========
(
    NOW() - INTERVAL '60 days', NOW(), 2, 'zh-CN',
    'CMS', 'cms',
    '内容管理系统（CMS）是用于创建和管理数字内容的软件应用，GoWind CMS是轻量级高性能CMS系统。',
    '/images/tags/cms-zh.jpg',
    'tag-default', '/tags/cms', 'https://gowind.com/tags/cms',
    'CMS,内容管理系统,GoWind,轻量级,高性能',
    '内容管理系统（CMS）是用于创建和管理数字内容的软件应用，GoWind CMS是轻量级高性能CMS系统。',
    'CMS | GoWind 标签'
),
-- ========== 标签2：CMS - 英文 ==========
(
    NOW() - INTERVAL '60 days', NOW(), 2, 'en-US',
    'CMS', 'cms',
    'Content Management System (CMS) is a software application for creating and managing digital content, GoWind CMS is a lightweight and high-performance CMS system.',
    '/images/tags/cms-en.jpg',
    'tag-default', '/en/tags/cms', 'https://gowind.com/en/tags/cms',
    'CMS,Content Management System,GoWind,Lightweight,High Performance',
    'Content Management System (CMS) is a software application for creating and managing digital content, GoWind CMS is a lightweight and high-performance CMS system.',
    'CMS | GoWind Tags'
),
-- ========== 标签5：新功能 - 中文 ==========
(
    NOW() - INTERVAL '50 days', NOW(), 5, 'zh-CN',
    '新功能', 'new-features',
    'GoWind CMS 新增功能模块，包含多租户、AI内容生成、性能优化等核心新特性。',
    '/images/tags/new-features-zh.jpg',
    'tag-default', '/tags/new-features', 'https://gowind.com/tags/new-features',
    '新功能,GoWind,CMS,多租户,AI内容生成',
    'GoWind CMS 新增功能模块，包含多租户、AI内容生成、性能优化等核心新特性。',
    '新功能 | GoWind 标签'
),
-- ========== 标签9：行业分析 - 中文 ==========
(
    NOW() - INTERVAL '40 days', NOW(), 9, 'zh-CN',
    '行业分析', 'industry-analysis',
    'CMS行业发展趋势、市场规模、技术方向分析，基于IDC、艾瑞等权威机构数据。',
    '/images/tags/industry-analysis-zh.jpg',
    'tag-default', '/tags/industry-analysis', 'https://gowind.com/tags/industry-analysis',
    '行业分析,CMS,2024趋势,市场规模,技术方向',
    'CMS行业发展趋势、市场规模、技术方向分析，基于IDC、艾瑞等权威机构数据。',
    '行业分析 | GoWind 标签'
),
-- ========== 标签19：性能优化 - 中文 ==========
(
    NOW() - INTERVAL '15 days', NOW(), 19, 'zh-CN',
    '性能优化', 'performance-optimization',
    'GoWind CMS 性能优化相关内容，涵盖数据库优化、缓存策略、代码层面优化，提升QPS和响应速度。',
    '/images/tags/performance-zh.jpg',
    'tag-default', '/tags/performance-optimization', 'https://gowind.com/tags/performance-optimization',
    '性能优化,QPS,数据库优化,缓存策略,高并发',
    'GoWind CMS 性能优化相关内容，涵盖数据库优化、缓存策略、代码层面优化，提升QPS和响应速度。',
    '性能优化 | GoWind 标签'
),
-- ========== 标签16：FAQ - 中文 ==========
(
    NOW() - INTERVAL '20 days', NOW(), 16, 'zh-CN',
    'FAQ', 'faq',
    'GoWind CMS 常见问题解答，涵盖安装部署、配置使用、性能优化、升级迁移等方向。',
    '/images/tags/faq-zh.jpg',
    'tag-default', '/tags/faq', 'https://gowind.com/tags/faq',
    'FAQ,常见问题,GoWind,CMS,问题解答',
    'GoWind CMS 常见问题解答，涵盖安装部署、配置使用、性能优化、升级迁移等方向。',
    'FAQ | GoWind 标签'
),
-- ========== 标签14：企业版 - 中文 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 14, 'zh-CN',
    '企业版', 'enterprise',
    'GoWind CMS 企业版，包含多租户、高级权限、数据备份、专属客服等付费专属功能。',
    '/images/tags/enterprise-zh.jpg',
    'tag-default', '/tags/enterprise', 'https://gowind.com/tags/enterprise',
    '企业版,GoWind,CMS,付费功能,多租户',
    'GoWind CMS 企业版，包含多租户、高级权限、数据备份、专属客服等付费专属功能。',
    '企业版 | GoWind 标签'
);


COMMIT;
