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

-- ----------------------------
-- 插入 sys_tenants 租户
-- ----------------------------
INSERT INTO public.sys_tenants(id, name, code, type, audit_status, status, admin_user_id, created_at)
VALUES (1, '测试租户', 'super', 'PAID', 'APPROVED', 'ON', 2, now())
;
SELECT setval('sys_tenants_id_seq', (SELECT MAX(id) FROM sys_tenants));

-- ----------------------------
-- 插入 sys_users 租户管理员用户
-- ----------------------------
INSERT INTO public.sys_users (id, tenant_id, username, nickname, realname, email, gender, created_at)
VALUES
    -- 2. 租户管理员（TENANT_ADMIN）
    (2, 1, 'tenant_admin', '租户管理', '张管理员', 'tenant@company.com', 'MALE', now())
;
SELECT setval('sys_users_id_seq', (SELECT MAX(id) FROM sys_users));

-- ----------------------------
-- 插入 sys_user_credentials 用户凭证（密码统一为admin，哈希值与原admin一致，方便测试）
-- ----------------------------
INSERT INTO public.sys_user_credentials (user_id, identity_type, identifier, credential_type, credential, status,
                                         is_primary, created_at)
VALUES
    -- 租户管理员（对应users表id=2）
    (2, 'USERNAME', 'tenant_admin', 'PASSWORD_HASH', '$2a$10$yajZDX20Y40FkG0Bu4N19eXNqRizez/S9fK63.JxGkfLq.RoNKR/a', 'ENABLED', true, now()),
    (2, 'EMAIL', 'tenant@company.com', 'PASSWORD_HASH', '$2a$10$yajZDX20Y40FkG0Bu4N19eXNqRizez/S9fK63.JxGkfLq.RoNKR/a', 'ENABLED', false, now())
;
SELECT setval('sys_user_credentials_id_seq', (SELECT MAX(id) FROM sys_user_credentials));

-- ----------------------------
-- 插入 sys_org_units 组织架构单元
-- ----------------------------
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
SELECT setval('sys_org_units_id_seq', (SELECT MAX(id) FROM sys_org_units));

-- ----------------------------
-- 插入 sys_positions 岗位数据
-- ----------------------------
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

-- ----------------------------
-- 插入 sys_tasks 调度任务
-- ----------------------------
INSERT INTO public.sys_tasks(type, type_name, task_payload, cron_spec, enable, created_at)
VALUES
    ('PERIODIC', 'backup', '{ "name": "test"}', '0 * * * *', true, now())
;
SELECT setval('sys_tasks_id_seq', (SELECT MAX(id) FROM sys_tasks));

-- ----------------------------
-- 插入 sys_login_policies 登录策略
-- ----------------------------
INSERT INTO public.sys_login_policies(id, target_id, type, method, value, reason, created_at)
VALUES
    (1, 1, 'BLACKLIST', 'IP', '127.0.0.1', '无理由', now()),
    (2, 1, 'WHITELIST', 'MAC', '00:1B:44:11:3A:B7 ', '无理由', now())
;
SELECT setval('sys_login_policies_id_seq', (SELECT MAX(id) FROM sys_login_policies));

-- ----------------------------
-- 插入 sys_dict_types 字典类型
-- ----------------------------
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

-- ----------------------------
-- 插入 sys_dict_type_i18n 字典类型国际化
-- ----------------------------
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

-- ----------------------------
-- 插入 sys_dict_entries 字典条目
-- ----------------------------
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

-- ----------------------------
-- 插入 sys_dict_entry_i18n 字典条目国际化
-- ----------------------------
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

-- ----------------------------
-- 插入 internal_message_categories 站内信分类
-- ----------------------------
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

-- ----------------------------
-- 插入 comments 表的测试数据，（评论数据）
-- 注意：按父子顺序插入（先顶级评论，再子评论），避免外键约束失败
-- ----------------------------
INSERT INTO public.comments (
    id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by,
    content_type, object_id, content, author_id, author_name, author_email,
    author_url, author_type, status, like_count, dislike_count, reply_count,
    ip_address, location, user_agent, detected_language, is_spam, is_sticky,
    reply_to_id, parent_id
) VALUES
-- ========== 文章 ID=1 的顶级评论 ==========
-- 评论1：王五（顶级）
(1,
 NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 1,
 '非常好的文章，学到了很多！',
 0, '王五', 'wangwu@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 5, 0, 2,
 '192.168.1.1', '北京', '', 'zh-CN', false, false,
 NULL, NULL),
-- 评论2：赵六（顶级）
(2,
 NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 1,
 'Composition API 确实很强大，感谢分享！',
 0, '赵六', 'zhaoliu@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 3, 0, 1,
 '192.168.1.2', '上海', '', 'zh-CN', false, false,
 NULL, NULL),
-- 评论3：孙七（顶级）
(3,
 NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 1,
 '期待更多关于 Vue 3 的文章',
 0, '孙七', 'sunqi@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 2, 0, 0,
 '192.168.1.3', '深圳', '', 'zh-CN', false, false,
 NULL, NULL),
-- ========== 文章 ID=2 的顶级评论 ==========
-- 评论4：周八（顶级）
(4,
 NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 2,
 'TypeScript 的类型系统真的很强大',
 0, '周八', 'zhouba@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 8, 0, 2,
 '192.168.1.4', '杭州', '', 'zh-CN', false, false,
 NULL, NULL),
-- 评论5：冯十一（顶级）
(5,
 NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 2,
 '推倒类型这个概念很有意思',
 0, '冯十一', 'fengshiyi@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 3, 0, 0,
 '192.168.1.11', '重庆', '', 'zh-CN', false, false,
 NULL, NULL),
-- ========== 文章 ID=3 的顶级评论 ==========
-- 评论6：郑十（顶级）
(6,
 NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 3,
 '内容中台是未来的趋势',
 0, '郑十', 'zhengshi@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 6, 0, 3,
 '192.168.1.6', '西安', '', 'zh-CN', false, false,
 NULL, NULL),

-- ========== 文章 ID=1 的二级评论（父评论1） ==========
-- 评论7：张三 回复 评论1
(7,
 NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 1,
 '确实很不错，尤其是关于 Composition API 的部分',
 0, '张三', 'zhangsan@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 2, 0, 0,
 '192.168.1.7', '广州', '', 'zh-CN', false, false,
 1, 1),  -- reply_to_id=1, parent_id=1

-- 评论8：李四 回复 评论1
(8,
 NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 1,
 '@王五 同意！博主写得很用心',
 0, '李四', 'lisi@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 1, 0, 1,
 '192.168.1.8', '武汉', '', 'zh-CN', false, false,
 1, 1),  -- reply_to_id=1, parent_id=1

-- ========== 文章 ID=1 的三级评论（父评论8） ==========
-- 评论11：博主 回复 评论8
(11,
 NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', NULL,
 1, 1, NULL,
 'CONTENT_TYPE_POST', 1,
 '@李四 谢谢支持！会继续努力的',
 1, '博主', 'admin@example.com',
 'https://example.com', 'AUTHOR_TYPE_ADMIN', 'STATUS_APPROVED',
 3, 0, 0,
 '192.168.1.100', '北京', '', 'zh-CN', false, false,
 8, 8),  -- reply_to_id=8, parent_id=8

-- ========== 文章 ID=1 的二级评论（父评论2） ==========
-- 评论9：孙七 回复 评论2
(9,
 NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 1,
 '@赵六 对，比 Options API 灵活多了',
 0, '孙七', 'sunqi@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 1, 0, 0,
 '192.168.1.9', '南京', '', 'zh-CN', false, false,
 2, 2),  -- reply_to_id=2, parent_id=2

-- ========== 文章 ID=2 的二级评论（父评论4） ==========
-- 评论10：吴九 回复 评论4
(10,
 NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 2,
 '条件类型的部分能再详细讲讲吗？',
 0, '吴九', 'wujiu@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 1, 0, 1,
 '192.168.1.5', '成都', '', 'zh-CN', false, false,
 4, 4),  -- reply_to_id=4, parent_id=4

-- 评论13：郑十 回复 评论4
(13,
 NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 2,
 '泛型约束也很有用',
 0, '郑十', 'zhengshi@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 2, 0, 0,
 '192.168.1.10', '西安', '', 'zh-CN', false, false,
 4, 4),  -- reply_to_id=4, parent_id=4

-- ========== 文章 ID=2 的三级评论（父评论10） ==========
-- 评论12：博主 回复 评论10
(12,
 NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', NULL,
 1, 1, NULL,
 'CONTENT_TYPE_POST', 2,
 '@吴九 好的，下期专门讲讲条件类型和分布式条件类型',
 1, '博主', 'admin@example.com',
 'https://example.com', 'AUTHOR_TYPE_ADMIN', 'STATUS_APPROVED',
 4, 0, 0,
 '192.168.1.100', '北京', '', 'zh-CN', false, false,
 10, 10),  -- reply_to_id=10, parent_id=10

-- ========== 文章 ID=3 的二级评论（父评论6） ==========
-- 评论14：钱十二 回复 评论6
(14,
 NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 3,
 '确实，前后端分离更灵活',
 0, '钱十二', 'qianshier@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 2, 0, 0,
 '192.168.1.12', '苏州', '', 'zh-CN', false, false,
 6, 6),  -- reply_to_id=6, parent_id=6

-- 评论15：孔十三 回复 评论6
(15,
 NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 3,
 'Strapi 和 Contentful 哪个更好用？',
 0, '孔十三', 'kongshisan@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 1, 0, 1,
 '192.168.1.13', '长沙', '', 'zh-CN', false, false,
 6, 6),  -- reply_to_id=6, parent_id=6

-- 评论17：白十四 回复 评论6
(17,
 NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', NULL,
 0, 0, NULL,
 'CONTENT_TYPE_POST', 3,
 'GoWind Content Hub 也很不错',
 0, '白十四', 'baisishi@example.com',
 '', 'AUTHOR_TYPE_GUEST', 'STATUS_APPROVED',
 4, 0, 0,
 '192.168.1.14', '郑州', '', 'zh-CN', false, false,
 6, 6),  -- reply_to_id=6, parent_id=6

-- ========== 文章 ID=3 的三级评论（父评论15） ==========
-- 评论16：博主 回复 评论15
(16,
 NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days', NULL,
 1, 1, NULL,
 'CONTENT_TYPE_POST', 3,
 '@孔十三 看需求，Strapi 开源免费，Contentful 功能更强大',
 1, '博主', 'admin@example.com',
 'https://example.com', 'AUTHOR_TYPE_ADMIN', 'STATUS_APPROVED',
 5, 0, 0,
 '192.168.1.100', '北京', '', 'zh-CN', false, false,
 15, 15)
;
SELECT setval('comments_id_seq', (SELECT MAX(id) FROM comments));

-- ----------------------------
-- 插入 media_assets 媒体资源数据
-- ----------------------------
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
SELECT setval('media_assets_id_seq', (SELECT MAX(id) FROM media_assets));

-- ----------------------------
-- 插入 site_settings 站点设置数据
-- ----------------------------
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
SELECT setval('site_settings_id_seq', (SELECT MAX(id) FROM site_settings));

-- ----------------------------
-- 插入 sites 表测试数据
-- ----------------------------
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
SELECT setval('sites_id_seq', (SELECT MAX(id) FROM sites));

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
    'GoWind - 高性能Go语言内容中台',
    'SETTING_TYPE_TEXT', '站点标题', '显示在浏览器标签的站点主标题', '请输入站点标题', NULL,
    true, '^.{2,50}$'
),
(
    NOW(), NOW(), 1, 'zh-CN', 'basic', 'site_description',
    '基于Go+Vue3开发的多站点/多语言内容中台',
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
    'GoWind,风行内容中台,Go语言,多语言,多站点',
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
    'GoWind - High Performance Go Content Hub',
    'SETTING_TYPE_TEXT', 'Site Title', 'Main title displayed in browser tab', 'Enter site title', NULL,
    true, '^.{2,50}$'
),
(
    NOW(), NOW(), 1, 'en-US', 'basic', 'site_description',
    'Multi-site & multi-language Content Hub built with Go + Vue3',
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
    'GoWind,Content Hub,Documentation,Go,Vue3,Multi-language',
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
SELECT setval('site_settings_id_seq', (SELECT MAX(id) FROM site_settings));

-- ----------------------------
-- 插入 navigations 表（导航组）
-- ----------------------------
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
-- 手机端底部导航
(104, NOW(), NOW(), '手机底部导航', 'MOBILE', 'zh-CN', true, 1, 1),

-- 英文导航组（对应中文组）
(201, NOW(), NOW(), 'Main Navigation', 'HEADER', 'en-US', true, 1, 1),
(202, NOW(), NOW(), 'Footer Navigation', 'FOOTER', 'en-US', true, 1, 1),
(203, NOW(), NOW(), 'Sidebar Navigation', 'SIDEBAR', 'en-US', true, 1, 1),
(204, NOW(), NOW(), 'Mobile Bottom Navigation', 'MOBILE', 'en-US', true, 1, 1);
SELECT setval('navigations_id_seq', (SELECT MAX(id) FROM navigations));

-- ----------------------------
-- 插入 navigation_items 表（导航项）
-- ----------------------------
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
(1011, NOW(), NOW(), 4, 'LINK_TYPE_EXTERNAL', 102, 'GitHub', 'https://github.com/tx7do/go-wind-cms', 0, 'logo-github', '访问我们的GitHub', true, false, 'nav-item footer-nav', '', NULL, 1, 1),
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
(2011, NOW(), NOW(), 4, 'LINK_TYPE_EXTERNAL', 202, 'GitHub', 'https://github.com/tx7do/go-wind-cms', 0, 'logo-github', 'Visit our GitHub', true, false, 'nav-item footer-nav', '', NULL, 1, 1),
-- ========== 导航组 203（en-US Sidebar Navigation） ==========
(2012, NOW(), NOW(), 1, 'LINK_TYPE_CUSTOM', 203, 'Popular Tags', '/tag', 0, 'pricetag', 'Browse popular tags', false, false, 'nav-item sidebar-nav', '', NULL, 1, 1),
(2013, NOW(), NOW(), 2, 'LINK_TYPE_CUSTOM', 203, 'Archive', '/archive', 0, 'archive', 'Post archive', false, false, 'nav-item sidebar-nav', '', NULL, 1, 1),
-- ========== 导航组 104（zh-CN 手机底部导航） ==========
(1014, NOW(), NOW(), 1, 'LINK_TYPE_CUSTOM', 104, '首页', '/', 0, 'home', '返回首页', false, false, 'nav-item mobile-nav', '', NULL, 1, 1),
(1015, NOW(), NOW(), 2, 'LINK_TYPE_CUSTOM', 104, '发现', '/explore', 0, 'explore', '发现内容', false, false, 'nav-item mobile-nav', '', NULL, 1, 1),
(1016, NOW(), NOW(), 3, 'LINK_TYPE_CUSTOM', 104, '收藏', '/bookmarks', 0, 'bookmark', '我的收藏', false, false, 'nav-item mobile-nav', '', NULL, 1, 1),
(1017, NOW(), NOW(), 4, 'LINK_TYPE_CUSTOM', 104, '我的', '/profile', 0, 'person', '个人中心', false, false, 'nav-item mobile-nav', '', NULL, 1, 1),
-- ========== 导航组 204（en-US Mobile Bottom Navigation） ==========
(2014, NOW(), NOW(), 1, 'LINK_TYPE_CUSTOM', 204, 'Home', '/', 0, 'home', 'Back to homepage', false, false, 'nav-item mobile-nav', '', NULL, 1, 1),
(2015, NOW(), NOW(), 2, 'LINK_TYPE_CUSTOM', 204, 'Discover', '/explore', 0, 'explore', 'Discover content', false, false, 'nav-item mobile-nav', '', NULL, 1, 1),
(2016, NOW(), NOW(), 3, 'LINK_TYPE_CUSTOM', 204, 'Bookmarks', '/bookmarks', 0, 'bookmark', 'My bookmarks', false, false, 'nav-item mobile-nav', '', NULL, 1, 1),
(2017, NOW(), NOW(), 4, 'LINK_TYPE_CUSTOM', 204, 'Me', '/profile', 0, 'person', 'Personal center', false, false, 'nav-item mobile-nav', '', NULL, 1, 1);
SELECT setval('navigation_items_id_seq', (SELECT MAX(id) FROM navigation_items));

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
SELECT setval('pages_id_seq', (SELECT MAX(id) FROM pages));

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
    '风行内容中台 - 高性能Go语言多站点多语言Content Hub系统',
    'home',
    '风行内容中台 是基于Go+Vue3开发的高性能、轻量级Content Hub系统，支持多站点、多语言、多租户，开箱即用。',
    '# 风行内容中台

## 核心特性
- 多站点管理：一个系统管理多个独立站点
- 多语言支持：内置多语言翻译体系
- 高性能：基于Go语言开发，QPS可达10万+
- 易扩展：插件化架构，支持自定义模板/主题',
    '# 风行内容中台

## 核心特性
- 多站点管理：一个系统管理多个独立站点
- 多语言支持：内置多语言翻译体系
- 高性能：基于Go语言开发，QPS可达10万+
- 易扩展：插件化架构，支持自定义模板/主题',
    '/images/thumbnails/home-zh.jpg',
    '/images/covers/home-zh.jpg',
    '/',
    520,
    'GoWind,Content Hub,Go语言,多站点,多语言,高性能',
    '风行内容中台 是基于Go+Vue3开发的高性能Content Hub系统，支持多站点、多语言、多租户，开箱即用。',
    '风行内容中台 - 高性能Go语言多站点多语言Content Hub系统'
),
-- ========== 页面1（首页） - 英文翻译（真实换行） ==========
(
    NOW() - INTERVAL '30 days', NOW(), 1, 'en-US',
    'GoWind Content Hub - High Performance Go Content Hub for Multi-site & Multi-language',
    'home',
    'GoWind Content Hub is a high-performance, lightweight Content Hub built with Go+Vue3, supporting multi-site, multi-language, and multi-tenant features, ready to use out of the box.',
    '# GoWind Content Hub

## Core Features
- Multi-site Management: Manage multiple independent sites with one system
- Multi-language Support: Built-in multi-language translation system
- High Performance: Developed with Go, QPS up to 100,000+
- Extensible: Plugin-based architecture, support custom templates/themes',
    '# GoWind Content Hub

## Core Features
- Multi-site Management: Manage multiple independent sites with one system
- Multi-language Support: Built-in multi-language translation system
- High Performance: Developed with Go, QPS up to 100,000+
- Extensible: Plugin-based architecture, support custom templates/themes',
    '/images/thumbnails/home-en.jpg',
    '/images/covers/home-en.jpg',
    '/en',
    480,
    'GoWind,Content Hub,Go,Multi-site,Multi-language,High Performance',
    'GoWind Content Hub is a high-performance Content Hub built with Go+Vue3, supporting multi-site, multi-language, and multi-tenant features, ready to use out of the box.',
    'GoWind Content Hub - High Performance Go Content Hub for Multi-site & Multi-language'
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
    '404 - 页面不存在 | 风行内容中台'
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
    '404 - Page Not Found | GoWind Content Hub'
),
-- ========== 页面7（登录页） - 中文翻译 ==========
(
    NOW() - INTERVAL '12 days', NOW(), 7, 'zh-CN',
    '登录 - 风行内容中台',
    'login',
    '登录风行内容中台，管理您的站点内容、配置和用户。',
    '', -- 自定义逻辑页无Markdown/富文本内容
    '',
    '/images/thumbnails/login-zh.jpg',
    '/images/covers/login-zh.jpg',
    '/login',
    0,
    '登录,风行内容中台,账号登录,用户中心',
    '登录风行内容中台，管理您的站点内容、配置和用户。',
    '登录 - 风行内容中台'
),
-- ========== 页面2（关于我们） - 中文翻译 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 2, 'zh-CN',
    '关于我们 - 风行内容中台',
    'about',
    'GoWind 团队成立于2024年，专注于高性能开源Content Hub系统的研发，致力于为开发者提供简单、高效的内容管理解决方案。',
    '<h1>关于 风行内容中台</h1>
<p>GoWind 团队成立于2024年，核心成员均有10年以上互联网研发经验，专注于高性能开源Content Hub系统的研发。</p>
<h2>我们的使命</h2>
<p>让内容管理更简单，让开发者更高效。</p>
<h2>团队规模</h2>
<p>目前团队共15人，涵盖后端、前端、产品、设计等多个领域。</p>',
    '<h1>关于 风行内容中台</h1>
<p>GoWind 团队成立于2024年，核心成员均有10年以上互联网研发经验，专注于高性能开源Content Hub系统的研发。</p>
<h2>我们的使命</h2>
<p>让内容管理更简单，让开发者更高效。</p>
<h2>团队规模</h2>
<p>目前团队共15人，涵盖后端、前端、产品、设计等多个领域。</p>',
    '/images/thumbnails/about-zh.jpg',
    '/images/covers/about-zh.jpg',
    '/about',
    850,
    'GoWind,关于我们,团队介绍,Content Hub研发',
    'GoWind 团队成立于2024年，专注于高性能开源Content Hub系统的研发，致力于为开发者提供简单、高效的内容管理解决方案。',
    '关于我们 - 风行内容中台'
);
SELECT setval('page_translations_id_seq', (SELECT MAX(id) FROM page_translations));

-- ----------------------------
-- 插入 categories 表（分类主表）测试数据
-- ----------------------------
-- 注意：按 parentId 顺序插入（先父后子），确保外键约束有效
INSERT INTO public.categories (
    id, created_at, updated_at, sort_order, path, status,
    depth, is_nav, icon, post_count, direct_post_count,
    custom_fields, parent_id, created_by, updated_by
) VALUES
-- ========== 一级分类 ==========
-- 技术分享
(1, NOW() - INTERVAL '30 days', NOW() - INTERVAL '15 days', 1, '/tech', 'CATEGORY_STATUS_ACTIVE', 0, true, 'carbon:document', 45, 15, '{}'::jsonb, NULL, 1, 1),
-- 生活随笔
(2, NOW() - INTERVAL '25 days', NOW() - INTERVAL '12 days', 2, '/life', 'CATEGORY_STATUS_ACTIVE', 0, true, 'carbon:blog', 30, 12, '{}'::jsonb, NULL, 1, 1),
-- 产品设计
(3, NOW() - INTERVAL '20 days', NOW() - INTERVAL '10 days', 3, '/design', 'CATEGORY_STATUS_ACTIVE', 0, true, 'carbon:chart-line', 25, 8, '{}'::jsonb, NULL, 1, 1),
-- 创业思考
(4, NOW() - INTERVAL '15 days', NOW() - INTERVAL '8 days', 4, '/startup', 'CATEGORY_STATUS_ACTIVE', 0, true, 'carbon:idea', 18, 10, '{}'::jsonb, NULL, 1, 1),
-- ========== 二级分类（父ID=1：技术分享） ==========
-- 前端开发
(11, NOW() - INTERVAL '25 days', NOW() - INTERVAL '10 days', 1, '/tech/frontend', 'CATEGORY_STATUS_ACTIVE', 1, false, 'carbon:code', 20, 20, '{}'::jsonb, 1, 1, 1),
-- 后端开发
(12, NOW() - INTERVAL '24 days', NOW() - INTERVAL '9 days', 2, '/tech/backend', 'CATEGORY_STATUS_ACTIVE', 1, false, 'carbon:cloud', 15, 15, '{}'::jsonb, 1, 1, 1),
-- 移动开发
(13, NOW() - INTERVAL '23 days', NOW() - INTERVAL '8 days', 3, '/tech/mobile', 'CATEGORY_STATUS_ACTIVE', 1, false, 'carbon:mobile', 10, 10, '{}'::jsonb, 1, 1, 1),
-- ========== 二级分类（父ID=2：生活随笔） ==========
-- 旅行游记
(21, NOW() - INTERVAL '20 days', NOW() - INTERVAL '7 days', 1, '/life/travel', 'CATEGORY_STATUS_ACTIVE', 1, false, 'carbon:map', 10, 10, '{}'::jsonb, 2, 1, 1),
-- 美食探店
(22, NOW() - INTERVAL '19 days', NOW() - INTERVAL '6 days', 2, '/life/food', 'CATEGORY_STATUS_ACTIVE', 1, false, 'carbon:favorite', 8, 8, '{}'::jsonb, 2, 1, 1),
-- ========== 二级分类（父ID=3：产品设计） ==========
-- UI 设计
(31, NOW() - INTERVAL '18 days', NOW() - INTERVAL '5 days', 1, '/design/ui-design', 'CATEGORY_STATUS_ACTIVE', 1, false, 'carbon:color-switch', 10, 10, '{}'::jsonb, 3, 1, 1),
-- UX 设计
(32, NOW() - INTERVAL '17 days', NOW() - INTERVAL '4 days', 2, '/design/ux-design', 'CATEGORY_STATUS_ACTIVE', 1, false, 'carbon:user-profile', 7, 7, '{}'::jsonb, 3, 1, 1),
-- ========== 二级分类（父ID=4：创业思考） ==========
-- 团队管理
(41, NOW() - INTERVAL '14 days', NOW() - INTERVAL '3 days', 1, '/startup/team-management', 'CATEGORY_STATUS_ACTIVE', 1, false, 'carbon:group', 5, 5, '{}'::jsonb, 4, 1, 1),
-- 产品思考
(42, NOW() - INTERVAL '13 days', NOW() - INTERVAL '2 days', 2, '/startup/product-thinking', 'CATEGORY_STATUS_ACTIVE', 1, false, 'carbon:product', 3, 3, '{}'::jsonb, 4, 1, 1);
SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));

-- ----------------------------
-- 插入 category_translations 表（分类翻译）
-- ----------------------------
INSERT INTO public.category_translations (
    id, created_at, updated_at, category_id, language_code,
    name, slug, description, thumbnail, cover_image,
    template, full_path, meta_keywords, meta_description, seo_title,
    created_by, updated_by
) VALUES
-- ========== 分类1（技术分享）- 中文 ==========
(1, NOW() - INTERVAL '30 days', NOW() - INTERVAL '15 days', 1, 'zh-CN', '技术分享', 'tech', '分享最新的技术文章和教程', 'https://picsum.photos/400/300?random=1', 'https://picsum.photos/1200/400?random=1', NULL, '/tech', '技术分享,教程,开发', '分享最新的技术文章和教程', '技术分享 | 风行内容中台', 1, 1),
-- ========== 分类1（技术分享）- 英文 ==========
(101, NOW() - INTERVAL '30 days', NOW() - INTERVAL '15 days', 1, 'en-US', 'Tech Sharing', 'tech', 'Share the latest technical articles and tutorials', 'https://picsum.photos/400/300?random=1', 'https://picsum.photos/1200/400?random=1', NULL, '/en/tech', 'Tech Sharing,Tutorials,Development', 'Share the latest technical articles and tutorials', 'Tech Sharing | GoWind Content Hub', 1, 1),
-- ========== 分类2（生活随笔）- 中文 ==========
(2, NOW() - INTERVAL '25 days', NOW() - INTERVAL '12 days', 2, 'zh-CN', '生活随笔', 'life', '记录生活中的点点滴滴', 'https://picsum.photos/400/300?random=2', 'https://picsum.photos/1200/400?random=2', NULL, '/life', '生活随笔,记录,日常', '记录生活中的点点滴滴', '生活随笔 | 风行内容中台', 1, 1),
-- ========== 分类2（生活随笔）- 英文 ==========
(102, NOW() - INTERVAL '25 days', NOW() - INTERVAL '12 days', 2, 'en-US', 'Life Notes', 'life', 'Record moments and thoughts from daily life', 'https://picsum.photos/400/300?random=2', 'https://picsum.photos/1200/400?random=2', NULL, '/en/life', 'Life Notes,Records,Daily', 'Record moments and thoughts from daily life', 'Life Notes | GoWind Content Hub', 1, 1),
-- ========== 分类3（产品设计）- 中文 ==========
(3, NOW() - INTERVAL '20 days', NOW() - INTERVAL '10 days', 3, 'zh-CN', '产品设计', 'design', '产品设计理念与实践', 'https://picsum.photos/400/300?random=3', 'https://picsum.photos/1200/400?random=3', NULL, '/design', '产品设计,理念,实践', '产品设计理念与实践', '产品设计 | 风行内容中台', 1, 1),
-- ========== 分类3（产品设计）- 英文 ==========
(103, NOW() - INTERVAL '20 days', NOW() - INTERVAL '10 days', 3, 'en-US', 'Product Design', 'design', 'Product design concepts and practices', 'https://picsum.photos/400/300?random=3', 'https://picsum.photos/1200/400?random=3', NULL, '/en/design', 'Product Design,Concepts,Practices', 'Product design concepts and practices', 'Product Design | GoWind Content Hub', 1, 1),
-- ========== 分类4（创业思考）- 中文 ==========
(4, NOW() - INTERVAL '15 days', NOW() - INTERVAL '8 days', 4, 'zh-CN', '创业思考', 'startup', '创业路上的思考与总结', 'https://picsum.photos/400/300?random=4', 'https://picsum.photos/1200/400?random=4', NULL, '/startup', '创业思考,总结,经验', '创业路上的思考与总结', '创业思考 | 风行内容中台', 1, 1),
-- ========== 分类4（创业思考）- 英文 ==========
(104, NOW() - INTERVAL '15 days', NOW() - INTERVAL '8 days', 4, 'en-US', 'Startup Insights', 'startup', 'Reflections and summaries from startup journey', 'https://picsum.photos/400/300?random=4', 'https://picsum.photos/1200/400?random=4', NULL, '/en/startup', 'Startup Insights,Reflections,Journey', 'Reflections and summaries from startup journey', 'Startup Insights | GoWind Content Hub', 1, 1),
-- ========== 分类11（前端开发）- 中文 ==========
(11, NOW() - INTERVAL '25 days', NOW() - INTERVAL '10 days', 11, 'zh-CN', '前端开发', 'frontend', '前端开发技术和框架', 'https://picsum.photos/400/300?random=11', 'https://picsum.photos/1200/400?random=11', NULL, '/tech/frontend', '前端开发,框架,技术', '前端开发技术和框架', '前端开发 | 技术分享', 1, 1),
-- ========== 分类11（前端开发）- 英文 ==========
(111, NOW() - INTERVAL '25 days', NOW() - INTERVAL '10 days', 11, 'en-US', 'Frontend Development', 'frontend', 'Frontend development technologies and frameworks', 'https://picsum.photos/400/300?random=11', 'https://picsum.photos/1200/400?random=11', NULL, '/en/tech/frontend', 'Frontend Development,Frameworks,Technologies', 'Frontend development technologies and frameworks', 'Frontend Development | Tech Sharing', 1, 1),
-- ========== 分类12（后端开发）- 中文 ==========
(12, NOW() - INTERVAL '24 days', NOW() - INTERVAL '9 days', 12, 'zh-CN', '后端开发', 'backend', '后端开发技术和架构', 'https://picsum.photos/400/300?random=12', 'https://picsum.photos/1200/400?random=12', NULL, '/tech/backend', '后端开发,架构,技术', '后端开发技术和架构', '后端开发 | 技术分享', 1, 1),
-- ========== 分类12（后端开发）- 英文 ==========
(112, NOW() - INTERVAL '24 days', NOW() - INTERVAL '9 days', 12, 'en-US', 'Backend Development', 'backend', 'Backend development technologies and architecture', 'https://picsum.photos/400/300?random=12', 'https://picsum.photos/1200/400?random=12', NULL, '/en/tech/backend', 'Backend Development,Architecture,Technologies', 'Backend development technologies and architecture', 'Backend Development | Tech Sharing', 1, 1),
-- ========== 分类13（移动开发）- 中文 ==========
(13, NOW() - INTERVAL '23 days', NOW() - INTERVAL '8 days', 13, 'zh-CN', '移动开发', 'mobile', '移动端开发技术', 'https://picsum.photos/400/300?random=13', 'https://picsum.photos/1200/400?random=13', NULL, '/tech/mobile', '移动开发,技术,移动端', '移动端开发技术', '移动开发 | 技术分享', 1, 1),
-- ========== 分类13（移动开发）- 英文 ==========
(113, NOW() - INTERVAL '23 days', NOW() - INTERVAL '8 days', 13, 'en-US', 'Mobile Development', 'mobile', 'Mobile development technologies', 'https://picsum.photos/400/300?random=13', 'https://picsum.photos/1200/400?random=13', NULL, '/en/tech/mobile', 'Mobile Development,Technologies', 'Mobile development technologies', 'Mobile Development | Tech Sharing', 1, 1),
-- ========== 分类21（旅行游记）- 中文 ==========
(21, NOW() - INTERVAL '20 days', NOW() - INTERVAL '7 days', 21, 'zh-CN', '旅行游记', 'travel', '旅行见闻和游记', 'https://picsum.photos/400/300?random=21', 'https://picsum.photos/1200/400?random=21', NULL, '/life/travel', '旅行游记,见闻,游记', '旅行见闻和游记', '旅行游记 | 生活随笔', 1, 1),
-- ========== 分类21（旅行游记）- 英文 ==========
(121, NOW() - INTERVAL '20 days', NOW() - INTERVAL '7 days', 21, 'en-US', 'Travel', 'travel', 'Travel experiences and journals', 'https://picsum.photos/400/300?random=21', 'https://picsum.photos/1200/400?random=21', NULL, '/en/life/travel', 'Travel,Experiences,Journals', 'Travel experiences and journals', 'Travel | Life Notes', 1, 1),
-- ========== 分类22（美食探店）- 中文 ==========
(22, NOW() - INTERVAL '19 days', NOW() - INTERVAL '6 days', 22, 'zh-CN', '美食探店', 'food', '探索城市美食', 'https://picsum.photos/400/300?random=22', 'https://picsum.photos/1200/400?random=22', NULL, '/life/food', '美食探店,城市美食,探索', '探索城市美食', '美食探店 | 生活随笔', 1, 1),
-- ========== 分类22（美食探店）- 英文 ==========
(122, NOW() - INTERVAL '19 days', NOW() - INTERVAL '6 days', 22, 'en-US', 'Food Exploration', 'food', 'Explore city delicacies', 'https://picsum.photos/400/300?random=22', 'https://picsum.photos/1200/400?random=22', NULL, '/en/life/food', 'Food Exploration,Delicacies,City', 'Explore city delicacies', 'Food Exploration | Life Notes', 1, 1),
-- ========== 分类31（UI 设计）- 中文 ==========
(31, NOW() - INTERVAL '18 days', NOW() - INTERVAL '5 days', 31, 'zh-CN', 'UI 设计', 'ui-design', '用户界面设计', 'https://picsum.photos/400/300?random=31', 'https://picsum.photos/1200/400?random=31', NULL, '/design/ui-design', 'UI设计,用户界面,设计', '用户界面设计', 'UI 设计 | 产品设计', 1, 1),
-- ========== 分类31（UI 设计）- 英文 ==========
(131, NOW() - INTERVAL '18 days', NOW() - INTERVAL '5 days', 31, 'en-US', 'UI Design', 'ui-design', 'User Interface Design', 'https://picsum.photos/400/300?random=31', 'https://picsum.photos/1200/400?random=31', NULL, '/en/design/ui-design', 'UI Design,User Interface', 'User Interface Design', 'UI Design | Product Design', 1, 1),
-- ========== 分类32（UX 设计）- 中文 ==========
(32, NOW() - INTERVAL '17 days', NOW() - INTERVAL '4 days', 32, 'zh-CN', 'UX 设计', 'ux-design', '用户体验设计', 'https://picsum.photos/400/300?random=32', 'https://picsum.photos/1200/400?random=32', NULL, '/design/ux-design', 'UX设计,用户体验,设计', '用户体验设计', 'UX 设计 | 产品设计', 1, 1),
-- ========== 分类32（UX 设计）- 英文 ==========
(132, NOW() - INTERVAL '17 days', NOW() - INTERVAL '4 days', 32, 'en-US', 'UX Design', 'ux-design', 'User Experience Design', 'https://picsum.photos/400/300?random=32', 'https://picsum.photos/1200/400?random=32', NULL, '/en/design/ux-design', 'UX Design,User Experience', 'User Experience Design', 'UX Design | Product Design', 1, 1),
-- ========== 分类41（团队管理）- 中文 ==========
(41, NOW() - INTERVAL '14 days', NOW() - INTERVAL '3 days', 41, 'zh-CN', '团队管理', 'team-management', '团队建设和管理经验', 'https://picsum.photos/400/300?random=41', 'https://picsum.photos/1200/400?random=41', NULL, '/startup/team-management', '团队管理,建设,经验', '团队建设和管理经验', '团队管理 | 创业思考', 1, 1),
-- ========== 分类41（团队管理）- 英文 ==========
(141, NOW() - INTERVAL '14 days', NOW() - INTERVAL '3 days', 41, 'en-US', 'Team Management', 'team-management', 'Team building and management experience', 'https://picsum.photos/400/300?random=41', 'https://picsum.photos/1200/400?random=41', NULL, '/en/startup/team-management', 'Team Management,Building,Experience', 'Team building and management experience', 'Team Management | Startup Insights', 1, 1),
-- ========== 分类42（产品思考）- 中文 ==========
(42, NOW() - INTERVAL '13 days', NOW() - INTERVAL '2 days', 42, 'zh-CN', '产品思考', 'product-thinking', '产品规划和思考', 'https://picsum.photos/400/300?random=42', 'https://picsum.photos/1200/400?random=42', NULL, '/startup/product-thinking', '产品思考,规划,产品', '产品规划和思考', '产品思考 | 创业思考', 1, 1),
-- ========== 分类42（产品思考）- 英文 ==========
(142, NOW() - INTERVAL '13 days', NOW() - INTERVAL '2 days', 42, 'en-US', 'Product Thinking', 'product-thinking', 'Product planning and thinking', 'https://picsum.photos/400/300?random=42', 'https://picsum.photos/1200/400?random=42', NULL, '/en/startup/product-thinking', 'Product Thinking,Planning', 'Product planning and thinking', 'Product Thinking | Startup Insights', 1, 1);
SELECT setval('category_translations_id_seq', (SELECT MAX(id) FROM category_translations));

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
-- 2. Content Hub（产品分组、精选、启用）
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
SELECT setval('tags_id_seq', (SELECT MAX(id) FROM tags));

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
    'Go（Golang）是Google开发的开源编程语言，简洁、高效、高并发，风行内容中台核心开发语言。',
    '/images/tags/go-zh.jpg',
    'tag-default', '/tags/go', 'https://gowind.com/tags/go',
    'Go语言,Golang,GoWind,高并发,编程语言',
    'Go（Golang）是Google开发的开源编程语言，简洁、高效、高并发，风行内容中台核心开发语言。',
    'Go语言 | GoWind 标签'
),
-- ========== 标签1：Go语言 - 英文 ==========
(
    NOW() - INTERVAL '60 days', NOW(), 1, 'en-US',
    'Go Language', 'go',
    'Go (Golang) is an open-source programming language developed by Google, concise, efficient, and high-concurrency, the core development language of GoWind Content Hub.',
    '/images/tags/go-en.jpg',
    'tag-default', '/en/tags/go', 'https://gowind.com/en/tags/go',
    'Go Language,Golang,GoWind,High Concurrency,Programming Language',
    'Go (Golang) is an open-source programming language developed by Google, concise, efficient, and high-concurrency, the core development language of GoWind Content Hub.',
    'Go Language | GoWind Tags'
),
-- ========== 标签2：Content Hub - 中文 ==========
(
    NOW() - INTERVAL '60 days', NOW(), 2, 'zh-CN',
    'Content Hub', 'cms',
    '内容管理系统（Content Hub）是用于创建和管理数字内容的软件应用，风行内容中台是轻量级高性能Content Hub系统。',
    '/images/tags/cms-zh.jpg',
    'tag-default', '/tags/cms', 'https://gowind.com/tags/cms',
    'Content Hub,内容管理系统,GoWind,轻量级,高性能',
    '内容管理系统（Content Hub）是用于创建和管理数字内容的软件应用，风行内容中台是轻量级高性能Content Hub系统。',
    'Content Hub | GoWind 标签'
),
-- ========== 标签2：Content Hub - 英文 ==========
(
    NOW() - INTERVAL '60 days', NOW(), 2, 'en-US',
    'Content Hub', 'cms',
    'Content Management System (Content Hub) is a software application for creating and managing digital content, GoWind Content Hub is a lightweight and high-performance Content Hub system.',
    '/images/tags/cms-en.jpg',
    'tag-default', '/en/tags/cms', 'https://gowind.com/en/tags/cms',
    'Content Hub,Content Management System,GoWind,Lightweight,High Performance',
    'Content Management System (Content Hub) is a software application for creating and managing digital content, GoWind Content Hub is a lightweight and high-performance Content Hub system.',
    'Content Hub | GoWind Tags'
),
-- ========== 标签3：快速上手 - 中文 ==========
(
    NOW() - INTERVAL '55 days', NOW(), 3, 'zh-CN',
    '快速上手', 'quick-start',
    '风行内容中台 快速上手指南，帮助用户在 5 分钟内完成环境搭建、代码克隆、配置启动和初始登录。',
    '/images/tags/quick-start-zh.jpg',
    'tag-default', '/tags/quick-start', 'https://gowind.com/tags/quick-start',
    '快速上手,GoWind,风行内容中台,安装部署,入门指南',
    '风行内容中台 快速上手指南，帮助用户在 5 分钟内完成环境搭建、代码克隆、配置启动和初始登录。',
    '快速上手 | GoWind 标签'
),
-- ========== 标签3：快速上手 - 英文 ==========
(
    NOW() - INTERVAL '55 days', NOW(), 3, 'en-US',
    'Quick Start', 'quick-start',
    'Quick start guide for GoWind Content Hub, helping users complete environment setup, code cloning, configuration startup and initial login within 5 minutes.',
    '/images/tags/quick-start-en.jpg',
    'tag-default', '/en/tags/quick-start', 'https://gowind.com/en/tags/quick-start',
    'Quick Start,GoWind,Content Hub,Installation,Setup Guide',
    'Quick start guide for GoWind Content Hub, helping users complete environment setup, code cloning, configuration startup and initial login within 5 minutes.',
    'Quick Start | GoWind Tags'
),
-- ========== 标签4：版本更新 - 中文 ==========
(
    NOW() - INTERVAL '50 days', NOW(), 4, 'zh-CN',
    '版本更新', 'version-update',
    '风行内容中台 版本更新记录与变更日志，包含功能新增、问题修复及兼容性改进。',
    '/images/tags/version-update-zh.jpg',
    'tag-default', '/tags/version-update', 'https://gowind.com/tags/version-update',
    '版本更新,GoWind,风行内容中台,变更日志,发布说明',
    '风行内容中台 版本更新记录与变更日志，包含功能新增、问题修复及兼容性改进。',
    '版本更新 | GoWind 标签'
),
-- ========== 标签4：版本更新 - 英文 ==========
(
    NOW() - INTERVAL '50 days', NOW(), 4, 'en-US',
    'Version Update', 'version-update',
    'Version update records and changelogs for GoWind Content Hub, including feature additions, bug fixes, and compatibility improvements.',
    '/images/tags/version-update-en.jpg',
    'tag-default', '/en/tags/version-update', 'https://gowind.com/en/tags/version-update',
    'Version Update,GoWind,Content Hub,Changelog,Release Notes',
    'Version update records and changelogs for GoWind Content Hub, including feature additions, bug fixes, and compatibility improvements.',
    'Version Update | GoWind Tags'
),
-- ========== 标签5：新功能 - 中文 ==========
(
    NOW() - INTERVAL '50 days', NOW(), 5, 'zh-CN',
    '新功能', 'new-features',
    '风行内容中台 新增功能模块，包含多租户、AI内容生成、性能优化等核心新特性。',
    '/images/tags/new-features-zh.jpg',
    'tag-default', '/tags/new-features', 'https://gowind.com/tags/new-features',
    '新功能,GoWind,Content Hub,多租户,AI内容生成',
    '风行内容中台 新增功能模块，包含多租户、AI内容生成、性能优化等核心新特性。',
    '新功能 | GoWind 标签'
),
-- ========== 标签5：新功能 - 英文 ==========
(
    NOW() - INTERVAL '50 days', NOW(), 5, 'en-US',
    'New Features', 'new-features',
    'New feature modules of GoWind Content Hub, including multi-tenancy, AI content generation, performance optimization and other core new features.',
    '/images/tags/new-features-en.jpg',
    'tag-default', '/en/tags/new-features', 'https://gowind.com/en/tags/new-features',
    'New Features,GoWind,Content Hub,Multi-tenancy,AI Content Generation',
    'New feature modules of GoWind Content Hub, including multi-tenancy, AI content generation, performance optimization and other core new features.',
    'New Features | GoWind Tags'
),
-- ========== 标签6：升级指南 - 中文 ==========
(
    NOW() - INTERVAL '45 days', NOW(), 6, 'zh-CN',
    '升级指南', 'upgrade-guide',
    '风行内容中台 逐步升级指南，涵盖版本兼容性检查、数据库迁移、配置更新及回滚流程。',
    '/images/tags/upgrade-guide-zh.jpg',
    'tag-default', '/tags/upgrade-guide', 'https://gowind.com/tags/upgrade-guide',
    '升级指南,GoWind,风行内容中台,迁移,版本兼容',
    '风行内容中台 逐步升级指南，涵盖版本兼容性检查、数据库迁移、配置更新及回滚流程。',
    '升级指南 | GoWind 标签'
),
-- ========== 标签6：升级指南 - 英文 ==========
(
    NOW() - INTERVAL '45 days', NOW(), 6, 'en-US',
    'Upgrade Guide', 'upgrade-guide',
    'Step-by-step upgrade guide for GoWind Content Hub, covering version compatibility checks, database migration, configuration updates and rollback procedures.',
    '/images/tags/upgrade-guide-en.jpg',
    'tag-default', '/en/tags/upgrade-guide', 'https://gowind.com/en/tags/upgrade-guide',
    'Upgrade Guide,GoWind,Content Hub,Migration,Version Compatibility',
    'Step-by-step upgrade guide for GoWind Content Hub, covering version compatibility checks, database migration, configuration updates and rollback procedures.',
    'Upgrade Guide | GoWind Tags'
),
-- ========== 标签7：Linux - 中文 ==========
(
    NOW() - INTERVAL '45 days', NOW(), 7, 'zh-CN',
    'Linux', 'linux',
    '风行内容中台 在 Linux 系统下的部署与运维内容，包含 Ubuntu/CentOS 系统配置、服务管理及性能调优。',
    '/images/tags/linux-zh.jpg',
    'tag-default', '/tags/linux', 'https://gowind.com/tags/linux',
    'Linux,GoWind,风行内容中台,部署,Ubuntu,CentOS',
    '风行内容中台 在 Linux 系统下的部署与运维内容，包含 Ubuntu/CentOS 系统配置、服务管理及性能调优。',
    'Linux | GoWind 标签'
),
-- ========== 标签7：Linux - 英文 ==========
(
    NOW() - INTERVAL '45 days', NOW(), 7, 'en-US',
    'Linux', 'linux',
    'Linux-related content for GoWind Content Hub deployment and operation, including Ubuntu/CentOS system configuration, service management and performance tuning.',
    '/images/tags/linux-en.jpg',
    'tag-default', '/en/tags/linux', 'https://gowind.com/en/tags/linux',
    'Linux,GoWind,Content Hub,Deployment,Ubuntu,CentOS',
    'Linux-related content for GoWind Content Hub deployment and operation, including Ubuntu/CentOS system configuration, service management and performance tuning.',
    'Linux | GoWind Tags'
),
-- ========== 标签8：部署教程 - 中文 ==========
(
    NOW() - INTERVAL '40 days', NOW(), 8, 'zh-CN',
    '部署教程', 'deployment-tutorial',
    '风行内容中台 全平台部署教程，涵盖云服务器、Docker 容器及 Kubernetes 集群等多种部署方式。',
    '/images/tags/deployment-tutorial-zh.jpg',
    'tag-default', '/tags/deployment-tutorial', 'https://gowind.com/tags/deployment-tutorial',
    '部署教程,GoWind,风行内容中台,Docker,Kubernetes,云服务器',
    '风行内容中台 全平台部署教程，涵盖云服务器、Docker 容器及 Kubernetes 集群等多种部署方式。',
    '部署教程 | GoWind 标签'
),
-- ========== 标签8：部署教程 - 英文 ==========
(
    NOW() - INTERVAL '40 days', NOW(), 8, 'en-US',
    'Deployment Tutorial', 'deployment-tutorial',
    'Comprehensive deployment tutorials for GoWind Content Hub on various platforms, including cloud servers, Docker containers and Kubernetes clusters.',
    '/images/tags/deployment-tutorial-en.jpg',
    'tag-default', '/en/tags/deployment-tutorial', 'https://gowind.com/en/tags/deployment-tutorial',
    'Deployment Tutorial,GoWind,Content Hub,Docker,Kubernetes,Cloud',
    'Comprehensive deployment tutorials for GoWind Content Hub on various platforms, including cloud servers, Docker containers and Kubernetes clusters.',
    'Deployment Tutorial | GoWind Tags'
),
-- ========== 标签9：行业分析 - 中文 ==========
(
    NOW() - INTERVAL '40 days', NOW(), 9, 'zh-CN',
    '行业分析', 'industry-analysis',
    'Content Hub行业发展趋势、市场规模、技术方向分析，基于IDC、艾瑞等权威机构数据。',
    '/images/tags/industry-analysis-zh.jpg',
    'tag-default', '/tags/industry-analysis', 'https://gowind.com/tags/industry-analysis',
    '行业分析,Content Hub,2024趋势,市场规模,技术方向',
    'Content Hub行业发展趋势、市场规模、技术方向分析，基于IDC、艾瑞等权威机构数据。',
    '行业分析 | GoWind 标签'
),
-- ========== 标签9：行业分析 - 英文 ==========
(
    NOW() - INTERVAL '40 days', NOW(), 9, 'en-US',
    'Industry Analysis', 'industry-analysis',
    'Content Hub industry development trends, market size, and technology direction analysis based on authoritative data from IDC, iResearch and other institutions.',
    '/images/tags/industry-analysis-en.jpg',
    'tag-default', '/en/tags/industry-analysis', 'https://gowind.com/en/tags/industry-analysis',
    'Industry Analysis,Content Hub,2024 Trends,Market Size,Technology Direction',
    'Content Hub industry development trends, market size, and technology direction analysis based on authoritative data from IDC, iResearch and other institutions.',
    'Industry Analysis | GoWind Tags'
),
-- ========== 标签10：2024趋势 - 中文 ==========
(
    NOW() - INTERVAL '35 days', NOW(), 10, 'zh-CN',
    '2024趋势', '2024-trends',
    '2024 年 Content Hub 行业趋势分析，涵盖轻量化架构、私有化部署、AI 赋能及市场增长预测。',
    '/images/tags/2024-trends-zh.jpg',
    'tag-default', '/tags/2024-trends', 'https://gowind.com/tags/2024-trends',
    '2024趋势,Content Hub,行业预测,轻量化,AI赋能',
    '2024 年 Content Hub 行业趋势分析，涵盖轻量化架构、私有化部署、AI 赋能及市场增长预测。',
    '2024趋势 | GoWind 标签'
),
-- ========== 标签10：2024趋势 - 英文 ==========
(
    NOW() - INTERVAL '35 days', NOW(), 10, 'en-US',
    '2024 Trends', '2024-trends',
    'Analysis of Content Hub industry trends in 2024, including lightweight architecture, private deployment, AI empowerment and market growth forecasts.',
    '/images/tags/2024-trends-en.jpg',
    'tag-default', '/en/tags/2024-trends', 'https://gowind.com/en/tags/2024-trends',
    '2024 Trends,Content Hub,Industry Forecast,Lightweight,AI',
    'Analysis of Content Hub industry trends in 2024, including lightweight architecture, private deployment, AI empowerment and market growth forecasts.',
    '2024 Trends | GoWind Tags'
),
-- ========== 标签11：轻量化 - 中文 ==========
(
    NOW() - INTERVAL '35 days', NOW(), 11, 'zh-CN',
    '轻量化', 'lightweight',
    '风行内容中台 轻量化架构设计理念，聚焦资源消耗最小化、快速启动及高效运行。',
    '/images/tags/lightweight-zh.jpg',
    'tag-default', '/tags/lightweight', 'https://gowind.com/tags/lightweight',
    '轻量化,GoWind,风行内容中台,架构设计,资源优化',
    '风行内容中台 轻量化架构设计理念，聚焦资源消耗最小化、快速启动及高效运行。',
    '轻量化 | GoWind 标签'
),
-- ========== 标签11：轻量化 - 英文 ==========
(
    NOW() - INTERVAL '35 days', NOW(), 11, 'en-US',
    'Lightweight', 'lightweight',
    'Lightweight architecture design philosophy of GoWind Content Hub, focusing on minimal resource consumption, fast startup and efficient operation.',
    '/images/tags/lightweight-en.jpg',
    'tag-default', '/en/tags/lightweight', 'https://gowind.com/en/tags/lightweight',
    'Lightweight,GoWind,Content Hub,Architecture,Minimal Resource',
    'Lightweight architecture design philosophy of GoWind Content Hub, focusing on minimal resource consumption, fast startup and efficient operation.',
    'Lightweight | GoWind Tags'
),
-- ========== 标签12：自定义模板 - 中文 ==========
(
    NOW() - INTERVAL '30 days', NOW(), 12, 'zh-CN',
    '自定义模板', 'custom-template',
    '风行内容中台 自定义模板开发指南，涵盖模板语法、数据绑定、组件开发及样式定制。',
    '/images/tags/custom-template-zh.jpg',
    'tag-default', '/tags/custom-template', 'https://gowind.com/tags/custom-template',
    '自定义模板,GoWind,Content Hub,模板引擎,前端开发',
    '风行内容中台 自定义模板开发指南，涵盖模板语法、数据绑定、组件开发及样式定制。',
    '自定义模板 | GoWind 标签'
),
-- ========== 标签12：自定义模板 - 英文 ==========
(
    NOW() - INTERVAL '30 days', NOW(), 12, 'en-US',
    'Custom Template', 'custom-template',
    'Custom template development guide for GoWind Content Hub, covering template syntax, data binding, component development and style customization.',
    '/images/tags/custom-template-en.jpg',
    'tag-default', '/en/tags/custom-template', 'https://gowind.com/en/tags/custom-template',
    'Custom Template,GoWind,Content Hub,Template Engine,Frontend Development',
    'Custom template development guide for GoWind Content Hub, covering template syntax, data binding, component development and style customization.',
    'Custom Template | GoWind Tags'
),
-- ========== 标签13：开发教程 - 中文 ==========
(
    NOW() - INTERVAL '30 days', NOW(), 13, 'zh-CN',
    '开发教程', 'development-tutorial',
    '风行内容中台 扩展与定制开发教程，包含插件开发、API 集成及主题定制。',
    '/images/tags/development-tutorial-zh.jpg',
    'tag-default', '/tags/development-tutorial', 'https://gowind.com/tags/development-tutorial',
    '开发教程,GoWind,Content Hub,插件开发,API集成',
    '风行内容中台 扩展与定制开发教程，包含插件开发、API 集成及主题定制。',
    '开发教程 | GoWind 标签'
),
-- ========== 标签13：开发教程 - 英文 ==========
(
    NOW() - INTERVAL '30 days', NOW(), 13, 'en-US',
    'Development Tutorial', 'development-tutorial',
    'Development tutorials for GoWind Content Hub extension and customization, including plugin development, API integration and theme customization.',
    '/images/tags/development-tutorial-en.jpg',
    'tag-default', '/en/tags/development-tutorial', 'https://gowind.com/en/tags/development-tutorial',
    'Development Tutorial,GoWind,Content Hub,Plugin Development,API Integration',
    'Development tutorials for GoWind Content Hub extension and customization, including plugin development, API integration and theme customization.',
    'Development Tutorial | GoWind Tags'
),
-- ========== 标签14：企业版 - 中文 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 14, 'zh-CN',
    '企业版', 'enterprise',
    '风行内容中台 企业版，包含多租户、高级权限、数据备份、专属客服等付费专属功能。',
    '/images/tags/enterprise-zh.jpg',
    'tag-default', '/tags/enterprise', 'https://gowind.com/tags/enterprise',
    '企业版,GoWind,风行内容中台,付费功能,多租户',
    '风行内容中台 企业版，包含多租户、高级权限、数据备份、专属客服等付费专属功能。',
    '企业版 | GoWind 标签'
),
-- ========== 标签14：企业版 - 英文 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 14, 'en-US',
    'Enterprise Edition', 'enterprise',
    'GoWind Content Hub Enterprise Edition, including paid exclusive features such as multi-tenancy, advanced permissions, data backup, and dedicated customer service.',
    '/images/tags/enterprise-en.jpg',
    'tag-default', '/en/tags/enterprise', 'https://gowind.com/en/tags/enterprise',
    'Enterprise Edition,GoWind,Content Hub,Paid Features,Multi-tenancy',
    'GoWind Content Hub Enterprise Edition, including paid exclusive features such as multi-tenancy, advanced permissions, data backup, and dedicated customer service.',
    'Enterprise Edition | GoWind Tags'
),
-- ========== 标签15：付费功能 - 中文 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 15, 'zh-CN',
    '付费功能', 'paid-features',
    '风行内容中台 企业版专属付费功能，包含高级分析、自定义域名、优先支持及 SLA 保障。',
    '/images/tags/paid-features-zh.jpg',
    'tag-default', '/tags/paid-features', 'https://gowind.com/tags/paid-features',
    '付费功能,GoWind,风行内容中台,企业版,商业化',
    '风行内容中台 企业版专属付费功能，包含高级分析、自定义域名、优先支持及 SLA 保障。',
    '付费功能 | GoWind 标签'
),
-- ========== 标签15：付费功能 - 英文 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 15, 'en-US',
    'Paid Features', 'paid-features',
    'Exclusive paid features of GoWind Content Hub Enterprise Edition, including advanced analytics, custom domains, priority support and SLA guarantees.',
    '/images/tags/paid-features-en.jpg',
    'tag-default', '/en/tags/paid-features', 'https://gowind.com/en/tags/paid-features',
    'Paid Features,GoWind,Content Hub,Enterprise,Monetization',
    'Exclusive paid features of GoWind Content Hub Enterprise Edition, including advanced analytics, custom domains, priority support and SLA guarantees.',
    'Paid Features | GoWind Tags'
),
-- ========== 标签16：FAQ - 中文 ==========
(
    NOW() - INTERVAL '20 days', NOW(), 16, 'zh-CN',
    'FAQ', 'faq',
    '风行内容中台 常见问题解答，涵盖安装部署、配置使用、性能优化、升级迁移等方向。',
    '/images/tags/faq-zh.jpg',
    'tag-default', '/tags/faq', 'https://gowind.com/tags/faq',
    'FAQ,常见问题,GoWind,风行内容中台,问题解答',
    '风行内容中台 常见问题解答，涵盖安装部署、配置使用、性能优化、升级迁移等方向。',
    'FAQ | GoWind 标签'
),
-- ========== 标签16：FAQ - 英文 ==========
(
    NOW() - INTERVAL '20 days', NOW(), 16, 'en-US',
    'FAQ', 'faq',
    'Frequently Asked Questions about GoWind Content Hub, covering installation and deployment, configuration usage, performance optimization, upgrade and migration.',
    '/images/tags/faq-en.jpg',
    'tag-default', '/en/tags/faq', 'https://gowind.com/en/tags/faq',
    'FAQ,Frequently Asked Questions,GoWind,Content Hub,Problem Solving',
    'Frequently Asked Questions about GoWind Content Hub, covering installation and deployment, configuration usage, performance optimization, upgrade and migration.',
    'FAQ | GoWind Tags'
),
-- ========== 标签17：问题解答 - 中文 ==========
(
    NOW() - INTERVAL '20 days', NOW(), 17, 'zh-CN',
    '问题解答', 'problem-solving',
    '风行内容中台 使用过程中常见问题的解决方案，提供分步排查指南及配置示例。',
    '/images/tags/problem-solving-zh.jpg',
    'tag-default', '/tags/problem-solving', 'https://gowind.com/tags/problem-solving',
    '问题解答,GoWind,风行内容中台,故障排查,配置',
    '风行内容中台 使用过程中常见问题的解决方案，提供分步排查指南及配置示例。',
    '问题解答 | GoWind 标签'
),
-- ========== 标签17：问题解答 - 英文 ==========
(
    NOW() - INTERVAL '20 days', NOW(), 17, 'en-US',
    'Problem Solving', 'problem-solving',
    'Solutions to common problems encountered when using GoWind Content Hub, with step-by-step troubleshooting guides and configuration examples.',
    '/images/tags/problem-solving-en.jpg',
    'tag-default', '/en/tags/problem-solving', 'https://gowind.com/en/tags/problem-solving',
    'Problem Solving,GoWind,Content Hub,Troubleshooting,Configuration',
    'Solutions to common problems encountered when using GoWind Content Hub, with step-by-step troubleshooting guides and configuration examples.',
    'Problem Solving | GoWind Tags'
),
-- ========== 标签18：故障排除 - 中文 ==========
(
    NOW() - INTERVAL '15 days', NOW(), 18, 'zh-CN',
    '故障排除', 'troubleshooting',
    '风行内容中台 高级故障排除技术，包含日志分析、性能剖析、数据库调试及网络诊断。',
    '/images/tags/troubleshooting-zh.jpg',
    'tag-default', '/tags/troubleshooting', 'https://gowind.com/tags/troubleshooting',
    '故障排除,GoWind,风行内容中台,调试,日志分析',
    '风行内容中台 高级故障排除技术，包含日志分析、性能剖析、数据库调试及网络诊断。',
    '故障排除 | GoWind 标签'
),
-- ========== 标签18：故障排除 - 英文 ==========
(
    NOW() - INTERVAL '15 days', NOW(), 18, 'en-US',
    'Troubleshooting', 'troubleshooting',
    'Advanced troubleshooting techniques for GoWind Content Hub, including log analysis, performance profiling, database debugging and network diagnostics.',
    '/images/tags/troubleshooting-en.jpg',
    'tag-default', '/en/tags/troubleshooting', 'https://gowind.com/en/tags/troubleshooting',
    'Troubleshooting,GoWind,Content Hub,Debugging,Log Analysis',
    'Advanced troubleshooting techniques for GoWind Content Hub, including log analysis, performance profiling, database debugging and network diagnostics.',
    'Troubleshooting | GoWind Tags'
),
-- ========== 标签19：性能优化 - 中文 ==========
(
    NOW() - INTERVAL '15 days', NOW(), 19, 'zh-CN',
    '性能优化', 'performance-optimization',
    '风行内容中台 性能优化相关内容，涵盖数据库优化、缓存策略、代码层面优化，提升QPS和响应速度。',
    '/images/tags/performance-zh.jpg',
    'tag-default', '/tags/performance-optimization', 'https://gowind.com/tags/performance-optimization',
    '性能优化,QPS,数据库优化,缓存策略,高并发',
    '风行内容中台 性能优化相关内容，涵盖数据库优化、缓存策略、代码层面优化，提升QPS和响应速度。',
    '性能优化 | GoWind 标签'
),
-- ========== 标签19：性能优化 - 英文 ==========
(
    NOW() - INTERVAL '15 days', NOW(), 19, 'en-US',
    'Performance Optimization', 'performance-optimization',
    'Performance optimization content for GoWind Content Hub, covering database optimization, caching strategies, code-level optimization to improve QPS and response speed.',
    '/images/tags/performance-en.jpg',
    'tag-default', '/en/tags/performance-optimization', 'https://gowind.com/en/tags/performance-optimization',
    'Performance Optimization,QPS,Database Optimization,Caching Strategy,High Concurrency',
    'Performance optimization content for GoWind Content Hub, covering database optimization, caching strategies, code-level optimization to improve QPS and response speed.',
    'Performance Optimization | GoWind Tags'
),
-- ========== 标签20：高并发 - 中文 ==========
(
    NOW() - INTERVAL '10 days', NOW(), 20, 'zh-CN',
    '高并发', 'high-concurrency',
    '风行内容中台 高并发处理能力，包含负载均衡、连接池、请求队列及水平扩展策略。',
    '/images/tags/high-concurrency-zh.jpg',
    'tag-default', '/tags/high-concurrency', 'https://gowind.com/tags/high-concurrency',
    '高并发,GoWind,风行内容中台,负载均衡,扩展性',
    '风行内容中台 高并发处理能力，包含负载均衡、连接池、请求队列及水平扩展策略。',
    '高并发 | GoWind 标签'
),
-- ========== 标签20：高并发 - 英文 ==========
(
    NOW() - INTERVAL '10 days', NOW(), 20, 'en-US',
    'High Concurrency', 'high-concurrency',
    'High concurrency handling capabilities of GoWind Content Hub, including load balancing, connection pooling, request queuing and horizontal scaling strategies.',
    '/images/tags/high-concurrency-en.jpg',
    'tag-default', '/en/tags/high-concurrency', 'https://gowind.com/en/tags/high-concurrency',
    'High Concurrency,GoWind,Content Hub,Load Balancing,Scaling',
    'High concurrency handling capabilities of GoWind Content Hub, including load balancing, connection pooling, request queuing and horizontal scaling strategies.',
    'High Concurrency | GoWind Tags'
)
;
SELECT setval('tag_translations_id_seq', (SELECT MAX(id) FROM tag_translations));

-- ----------------------------
-- 插入 posts 表
-- ----------------------------
INSERT INTO public.posts (
    created_at, updated_at, sort_order, editor_type,
    status, code, disallow_comment, in_progress,
    auto_summary, is_featured, visits, likes,
    comment_count, author_id, author_name, password_hash,
    custom_fields
) VALUES
-- 文章1：风行内容中台 快速上手（已发布、精选）
(
    NOW() - INTERVAL '30 days', NOW(),
    1, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'gowind-cms-quick-start',
    false, false, true, true,
    15890, 892, 156,
    1, 'GoWind 官方', '',
    '{"show_toc": "true", "toc_depth": "3", "allow_copy": "true", "copyright_notice": "GoWind 官方原创"}'::jsonb
),
-- 文章2：GoWind v2.0 版本发布公告（已发布、精选）
(
    NOW() - INTERVAL '25 days', NOW(),
    2, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'gowind-v2-0-release',
    false, false, true, true,
    9870, 654, 89,
    1, 'GoWind 官方', '',
    '{"show_changelog": "true", "release_date": "2024-03-01", "upgrade_guide_url": "/docs/upgrade/v2.0"}'::jsonb
),
-- 文章3：Linux 环境下部署 风行内容中台（已发布）
(
    NOW() - INTERVAL '22 days', NOW(),
    3, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'deploy-gowind-on-linux',
    false, false, true, false,
    7650, 432, 78,
    1001, '张三', '',
    '{"os_type": "Linux", "distro": "Ubuntu, CentOS", "tested_version": "v1.9.0"}'::jsonb
),
-- 文章4：2024 Content Hub 行业发展趋势分析（已发布）
(
    NOW() - INTERVAL '20 days', NOW(),
    4, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', '2024-cms-industry-trends',
    false, false, true, false,
    6540, 389, 67,
    1002, '李四', '',
    '{"data_source": "IDC 2024 行业报告", "chart_support": "true", "downloadable": "true"}'::jsonb
),
-- 文章5：风行内容中台 自定义模板开发（草稿、未完成）
(
    NOW() - INTERVAL '15 days', NOW(),
    5, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_DRAFT', 'gowind-custom-template-dev',
    true, true, false, false,
    1230, 0, 0,
    1001, '张三', '',
    '{"dev_status": "50%", "expected_release": "2024-04-01", "required_skills": "Go, Vue3, HTML/CSS"}'::jsonb
),
-- 文章6：GoWind 企业版功能详解（加密、已发布）
(
    NOW() - INTERVAL '12 days', NOW(),
    6, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'gowind-enterprise-features',
    true, false, true, false,
    4320, 256, 0,
    1, 'GoWind 官方', '$2a$10$89jZk54G89sdkf89sdf89sd89sdf89sdf89sdf',
    '{"is_enterprise": "true", "price_range": "¥9999-¥19999", "trial_available": "true"}'::jsonb
),
-- 文章7：常见问题解答（草稿）
(
    NOW() - INTERVAL '10 days', NOW(),
    7, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_DRAFT', 'gowind-faq',
    true, false, true, false,
    890, 0, 0,
    1, 'GoWind 官方', '',
    '{"faq_category": "installation, configuration, performance", "update_frequency": "monthly"}'::jsonb
),
-- 文章8：风行内容中台 性能优化指南（已发布、精选）
(
    NOW() - INTERVAL '8 days', NOW(),
    8, 'EDITOR_TYPE_MARKDOWN',
    'POST_STATUS_PUBLISHED', 'gowind-cms-performance-optimization',
    false, false, true, true,
    5680, 789, 123,
    1003, '王五', '',
    '{"benchmark_data": "true", "qps_before": "50000", "qps_after": "100000", "optimization_points": "DB, Cache, Code"}'::jsonb
);
SELECT setval('posts_id_seq', (SELECT MAX(id) FROM posts));

-- ----------------------------
-- 插入 post_translations 表
-- ----------------------------
INSERT INTO public.post_translations (
    created_at, updated_at, post_id, language_code, title,
    slug, summary, content, original_content, thumbnail,
    template, full_path, word_count, meta_keywords,
    meta_description, seo_title
) VALUES
-- ========== 文章1：风行内容中台 快速上手 - 中文 ==========
(
    NOW() - INTERVAL '30 days', NOW(), 1, 'zh-CN',
    '风行内容中台 快速上手：5分钟搭建你的第一个内容中台',
    'gowind-content-hub-quick-start',
    '本文带你5分钟快速搭建风行内容中台，涵盖环境准备、代码克隆、配置启动、初始登录全流程。',
    '# 风行内容中台 快速上手

## 环境准备
### 前置依赖
- Go 1.21+（推荐1.22最新版）
- PostgreSQL 14+
- Git（可选）
- Docker（可选）

## 安装步骤
1. 克隆代码仓库：git clone https://github.com/tx7do/go-wind-cms.git
2. 配置环境变量：复制.env.example为.env并修改数据库配置
3. 启动服务：go run main.go 或 docker-compose up -d
4. 初始登录：http://localhost:8080，默认账号admin/admin

> 首次登录请立即修改密码！',
    '# GoWind Content Hub 快速上手

## 环境准备
### 前置依赖
- Go 1.21+（推荐1.22最新版）
- PostgreSQL 14+
- Git（可选）
- Docker（可选）

## 安装步骤
1. 克隆代码仓库：git clone https://github.com/tx7do/go-wind-cms.git
2. 配置环境变量：复制.env.example为.env并修改数据库配置
3. 启动服务：go run main.go 或 docker-compose up -d
4. 初始登录：http://localhost:8080，默认账号admin/admin

> 首次登录请立即修改密码！',
    'https://picsum.photos/800/450?random=1',
    'post-default', '/blog/gowind-content-hub-quick-start',
    2580,
    'GoWind,风行内容中台,快速上手,安装部署,Go语言,PostgreSQL',
    '本文带你5分钟快速搭建风行内容中台，涵盖环境准备、代码克隆、配置启动、初始登录全流程。',
    '风行内容中台 快速上手：5分钟搭建你的第一个内容中台 | GoWind 官方文档'
),
-- ========== 文章2：GoWind v2.0 版本发布公告 - 中文 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 2, 'zh-CN',
    '风行内容中台 v2.0 正式发布：新增多租户、性能提升100%',
    'gowind-v2-0-release',
    'GoWind Content Hub v2.0版本发布，核心更新多租户支持、性能优化、UI重构，QPS突破10万。',
    '# GoWind Content Hub v2.0 正式发布

## 发布说明
GoWind Content Hub v2.0于2024年3月1日发布，是开源以来的重大版本更新！

## 核心新功能
1. 多租户支持：单实例部署多套独立站点，数据隔离
2. 性能优化：QPS从5万提升至10万，响应时间降低60%
3. UI重构：基于Vue3+Element Plus重构后台，移动端适配

## 升级指南
- 从v1.9升级：备份数据库后执行go run scripts/upgrade/v2.0.go
- 全新安装：直接克隆v2.0分支代码部署',
    '# GoWind Content Hub v2.0 正式发布

## 发布说明
GoWind Content Hub v2.0于2024年3月1日发布，是开源以来的重大版本更新！

## 核心新功能
1. 多租户支持：单实例部署多套独立站点，数据隔离
2. 性能优化：QPS从5万提升至10万，响应时间降低60%
3. UI重构：基于Vue3+Element Plus重构后台，移动端适配

## 升级指南
- 从v1.9升级：备份数据库后执行go run scripts/upgrade/v2.0.go
- 全新安装：直接克隆v2.0分支代码部署',
    'https://picsum.photos/800/450?random=2',
    'post-announcement', '/blog/gowind-v2-0-release',
    3200,
    'GoWind,Content Hub,v2.0,版本发布,多租户,性能优化,UI重构',
    'GoWind Content Hub v2.0版本发布，核心更新多租户支持、性能优化、UI重构，QPS突破10万。',
    'GoWind Content Hub v2.0 正式发布：新增多租户、性能提升100% | GoWind 官方公告'
),
-- ========== 文章3：Linux 环境下部署 GoWind Content Hub - 中文 ==========
(
    NOW() - INTERVAL '22 days', NOW(), 3, 'zh-CN',
    'Linux 环境下部署 GoWind Content Hub（Ubuntu/CentOS通用）',
    'deploy-gowind-on-linux',
    '详解Linux环境（Ubuntu/CentOS）下GoWind Content Hub的部署步骤，含依赖安装、端口配置、开机自启。',
    '# Linux 环境下部署 GoWind Content Hub

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
2. 克隆代码：git clone https://github.com/tx7do/go-wind-cms.git
3. 配置数据库连接：修改.env文件
4. 启动服务：nohup go run main.go > app.log 2>&1 &

## 开机自启
创建systemd服务文件：/etc/systemd/system/gowind.service',
    '# Linux 环境下部署 GoWind Content Hub

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
2. 克隆代码：git clone https://github.com/tx7do/go-wind-cms.git
3. 配置数据库连接：修改.env文件
4. 启动服务：nohup go run main.go > app.log 2>&1 &

## 开机自启
创建systemd服务文件：/etc/systemd/system/gowind.service',
    'https://picsum.photos/800/450?random=3',
    'post-tech', '/blog/deploy-gowind-on-linux',
    2800,
    'GoWind,Content Hub,Linux部署,Ubuntu,CentOS,Go语言',
    '详解Linux环境（Ubuntu/CentOS）下GoWind Content Hub的部署步骤，含依赖安装、端口配置、开机自启。',
    'Linux 环境下部署 GoWind Content Hub（Ubuntu/CentOS通用） | GoWind 技术博客'
),
-- ========== 文章4：2024 Content Hub 行业发展趋势分析 - 中文 ==========
(
    NOW() - INTERVAL '20 days', NOW(), 4, 'zh-CN',
    '2024 Content Hub 行业发展趋势：轻量化、私有化、AI赋能',
    '2024-cms-industry-trends',
    '基于IDC 2024行业报告，分析Content Hub行业三大趋势：轻量化、私有化部署、AI智能赋能。',
    '# 2024 Content Hub 行业发展趋势分析

## 数据来源
IDC 2024年全球Content Hub市场研究报告

## 核心趋势
1. 轻量化：轻量级Content Hub占比提升至65%，替代重型系统
2. 私有化：企业级用户私有化部署需求增长40%
3. AI赋能：AI生成内容、智能排版成为标配功能

## 市场规模
2024年全球Content Hub市场规模预计达89亿美元，年增长率18%。

## 国内趋势
国产化替代加速，Go/Java语言开发的Content Hub占比提升。',
    '# 2024 Content Hub 行业发展趋势分析

## 数据来源
IDC 2024年全球Content Hub市场研究报告

## 核心趋势
1. 轻量化：轻量级Content Hub占比提升至65%，替代重型系统
2. 私有化：企业级用户私有化部署需求增长40%
3. AI赋能：AI生成内容、智能排版成为标配功能

## 市场规模
2024年全球Content Hub市场规模预计达89亿美元，年增长率18%。

## 国内趋势
国产化替代加速，Go/Java语言开发的Content Hub占比提升。',
    'https://picsum.photos/800/450?random=4',
    'post-analysis', '/blog/2024-cms-industry-trends',
    2600,
    'Content Hub,2024趋势,行业分析,轻量化,私有化,AI赋能',
    '基于IDC 2024行业报告，分析Content Hub行业三大趋势：轻量化、私有化部署、AI智能赋能。',
    '2024 Content Hub 行业发展趋势：轻量化、私有化、AI赋能 | GoWind 行业洞察'
),
-- ========== 文章5：GoWind Content Hub 自定义模板开发 - 中文（草稿） ==========
(
    NOW() - INTERVAL '15 days', NOW(), 5, 'zh-CN',
    'GoWind Content Hub 自定义模板开发教程（草稿）',
    'gowind-custom-template-dev',
    'GoWind Content Hub自定义模板开发教程，涵盖模板语法、数据调用、样式定制，当前开发进度50%。',
    '# GoWind Content Hub 自定义模板开发

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
    '# GoWind Content Hub 自定义模板开发

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
    'https://picsum.photos/800/450?random=5',
    'post-draft', '/blog/gowind-custom-template-dev',
    1800,
    'GoWind,Content Hub,自定义模板,开发教程,Go,Vue3',
    'GoWind Content Hub自定义模板开发教程，涵盖模板语法、数据调用、样式定制，当前开发进度50%。',
    'GoWind Content Hub 自定义模板开发教程（草稿） | GoWind 开发文档'
),
-- ========== 文章6：GoWind 企业版功能详解 - 中文（加密） ==========
(
    NOW() - INTERVAL '12 days', NOW(), 6, 'zh-CN',
    'GoWind Content Hub 企业版功能详解（付费专属）',
    'gowind-enterprise-features',
    'GoWind企业版专属功能：多租户管理、高级权限、数据备份、专属客服，价格9999-19999元/年。',
    '# GoWind Content Hub 企业版功能详解

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
    '# GoWind Content Hub 企业版功能详解

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
    'https://picsum.photos/800/450?random=6',
    'post-enterprise', '/blog/gowind-enterprise-features',
    2200,
    'GoWind,Content Hub,企业版,付费功能,多租户,高级权限',
    'GoWind企业版专属功能：多租户管理、高级权限、数据备份、专属客服，价格9999-19999元/年。',
    'GoWind Content Hub 企业版功能详解（付费专属） | GoWind 企业服务'
),
-- ========== 文章7：常见问题解答 - 中文（草稿） ==========
(
    NOW() - INTERVAL '10 days', NOW(), 7, 'zh-CN',
    'GoWind Content Hub 常见问题解答（FAQ）',
    'gowind-faq',
    'GoWind Content Hub常见问题汇总，涵盖安装、配置、性能、升级等方向，每月更新。',
    '# GoWind Content Hub 常见问题解答

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
    '# GoWind Content Hub 常见问题解答

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
    'https://picsum.photos/800/450?random=7',
    'post-faq', '/blog/gowind-faq',
    1500,
    'GoWind,Content Hub,FAQ,常见问题,安装,配置,性能',
    'GoWind Content Hub常见问题汇总，涵盖安装、配置、性能、升级等方向，每月更新。',
    'GoWind Content Hub 常见问题解答（FAQ） | GoWind 帮助中心'
),
-- ========== 文章8：GoWind Content Hub 性能优化指南 - 中文 ==========
(
    NOW() - INTERVAL '8 days', NOW(), 8, 'zh-CN',
    'GoWind Content Hub 性能优化指南：从5万QPS到10万的实战经验',
    'gowind-cms-performance-optimization',
    '分享GoWind Content Hub性能优化实战经验，涵盖数据库优化、缓存策略、代码层面，QPS提升100%。',
    '# GoWind Content Hub 性能优化指南

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
    '# GoWind Content Hub 性能优化指南

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
    'https://picsum.photos/800/450?random=8',
    'post-tech', '/blog/gowind-cms-performance-optimization',
    3000,
    'GoWind,Content Hub,性能优化,QPS,数据库优化,缓存策略,Go语言',
    '分享GoWind Content Hub性能优化实战经验，涵盖数据库优化、缓存策略、代码层面，QPS提升100%。',
    'GoWind Content Hub 性能优化指南：从5万QPS到10万的实战经验 | GoWind 技术博客'
),
-- ========== 文章1：GoWind Content Hub 快速上手 - 英文 ==========
(
    NOW() - INTERVAL '30 days', NOW(), 1, 'en-US',
    'GoWind Content Hub Quick Start: Build Your First Content Hub Site in 5 Minutes',
    'gowind-cms-quick-start',
    'This guide helps you set up a GoWind Content Hub site in 5 minutes, covering environment preparation, code cloning, configuration startup, and initial login.',
    '# GoWind Content Hub Quick Start

## Environment Preparation
### Prerequisites
- Go 1.21+ (recommended 1.22 latest version)
- PostgreSQL 14+
- Git (optional)
- Docker (optional)

## Installation Steps
1. Clone repository: git clone https://github.com/tx7do/go-wind-cms.git
2. Configure environment variables: copy .env.example to .env and modify database settings
3. Start service: go run main.go or docker-compose up -d
4. Initial login: http://localhost:8080, default account admin/admin

> Please change password immediately after first login!',
    '# GoWind Content Hub Quick Start

## Environment Preparation
### Prerequisites
- Go 1.21+ (recommended 1.22 latest version)
- PostgreSQL 14+
- Git (optional)
- Docker (optional)

## Installation Steps
1. Clone repository: git clone https://github.com/tx7do/go-wind-cms.git
2. Configure environment variables: copy .env.example to .env and modify database settings
3. Start service: go run main.go or docker-compose up -d
4. Initial login: http://localhost:8080, default account admin/admin

> Please change password immediately after first login!',
    'https://picsum.photos/800/450?random=9',
    'post-default', '/en/blog/gowind-cms-quick-start',
    2580,
    'GoWind,Content Hub,Quick Start,Installation,Go Language,PostgreSQL',
    'This guide helps you set up a GoWind Content Hub site in 5 minutes, covering environment preparation, code cloning, configuration startup, and initial login.',
    'GoWind Content Hub Quick Start: Build Your First Content Hub Site in 5 Minutes | GoWind Official Documentation'
),
-- ========== 文章2：GoWind v2.0 版本发布公告 - 英文 ==========
(
    NOW() - INTERVAL '25 days', NOW(), 2, 'en-US',
    'GoWind Content Hub v2.0 Official Release: Multi-tenancy Support, 100% Performance Boost',
    'gowind-v2-0-release',
    'GoWind Content Hub v2.0 release introduces core updates including multi-tenancy support, performance optimization, UI refactoring, with QPS exceeding 100,000.',
    '# GoWind Content Hub v2.0 Official Release

## Release Notes
GoWind Content Hub v2.0 was released on March 1, 2024, marking the most significant update since open-sourcing!

## Core New Features
1. Multi-tenancy Support: Deploy multiple independent sites on a single instance with complete data isolation
2. Performance Optimization: QPS increased from 50,000 to 100,000, response time reduced by 60%
3. UI Refactoring: Backend rebuilt with Vue3 + Element Plus, mobile-responsive design

## Upgrade Guide
- Upgrade from v1.9: Backup database and execute go run scripts/upgrade/v2.0.go
- Fresh installation: Clone v2.0 branch code directly for deployment',
    '# GoWind Content Hub v2.0 Official Release

## Release Notes
GoWind Content Hub v2.0 was released on March 1, 2024, marking the most significant update since open-sourcing!

## Core New Features
1. Multi-tenancy Support: Deploy multiple independent sites on a single instance with complete data isolation
2. Performance Optimization: QPS increased from 50,000 to 100,000, response time reduced by 60%
3. UI Refactoring: Backend rebuilt with Vue3 + Element Plus, mobile-responsive design

## Upgrade Guide
- Upgrade from v1.9: Backup database and execute go run scripts/upgrade/v2.0.go
- Fresh installation: Clone v2.0 branch code directly for deployment',
    'https://picsum.photos/800/450?random=10',
    'post-announcement', '/en/blog/gowind-v2-0-release',
    3200,
    'GoWind,Content Hub,v2.0,Release Notes,Multi-tenancy,Performance Optimization,UI Refactoring',
    'GoWind Content Hub v2.0 release introduces core updates including multi-tenancy support, performance optimization, UI refactoring, with QPS exceeding 100,000.',
    'GoWind Content Hub v2.0 Official Release: Multi-tenancy Support, 100% Performance Boost | GoWind Official Announcement'
),
-- ========== 文章3：Linux 环境下部署 GoWind Content Hub - 英文 ==========
(
    NOW() - INTERVAL '22 days', NOW(), 3, 'en-US',
    'Deploy GoWind Content Hub on Linux (Ubuntu/CentOS Universal Guide)',
    'deploy-gowind-on-linux',
    'Detailed deployment guide for GoWind Content Hub on Linux environments (Ubuntu/CentOS), including dependency installation, port configuration, and auto-start setup.',
    '# Deploy GoWind Content Hub on Linux

## Supported Systems
- Ubuntu 20.04/22.04
- CentOS 7/8

## Dependency Installation
### Ubuntu
apt update && apt install -y golang postgresql git
### CentOS
yum install -y golang postgresql git

## Deployment Steps
1. Create database: createdb gowind
2. Clone code: git clone https://github.com/tx7do/go-wind-cms.git
3. Configure database connection: modify .env file
4. Start service: nohup go run main.go > app.log 2>&1 &

## Auto-start on Boot
Create systemd service file: /etc/systemd/system/gowind.service',
    '# Deploy GoWind Content Hub on Linux

## Supported Systems
- Ubuntu 20.04/22.04
- CentOS 7/8

## Dependency Installation
### Ubuntu
apt update && apt install -y golang postgresql git
### CentOS
yum install -y golang postgresql git

## Deployment Steps
1. Create database: createdb gowind
2. Clone code: git clone https://github.com/tx7do/go-wind-cms.git
3. Configure database connection: modify .env file
4. Start service: nohup go run main.go > app.log 2>&1 &

## Auto-start on Boot
Create systemd service file: /etc/systemd/system/gowind.service',
    'https://picsum.photos/800/450?random=11',
    'post-tech', '/en/blog/deploy-gowind-on-linux',
    2800,
    'GoWind,Content Hub,Linux Deployment,Ubuntu,CentOS,Go Language',
    'Detailed deployment guide for GoWind Content Hub on Linux environments (Ubuntu/CentOS), including dependency installation, port configuration, and auto-start setup.',
    'Deploy GoWind Content Hub on Linux (Ubuntu/CentOS Universal Guide) | GoWind Technical Blog'
),
-- ========== 文章4：2024 Content Hub 行业发展趋势分析 - 英文 ==========
(
    NOW() - INTERVAL '20 days', NOW(), 4, 'en-US',
    '2024 Content Hub Industry Trend Analysis: Lightweight, Private Deployment, AI Empowerment',
    '2024-cms-industry-trends',
    'Based on IDC 2024 industry report, analysis of three major Content Hub industry trends: lightweight architecture, private deployment, and AI-powered features.',
    '# 2024 Content Hub Industry Trend Analysis

## Data Source
IDC 2024 Global Content Hub Market Research Report

## Core Trends
1. Lightweight: Lightweight Content Hub market share increased to 65%, replacing heavyweight systems
2. Private Deployment: Enterprise private deployment demand grew by 40%
3. AI Empowerment: AI-generated content and intelligent layout becoming standard features

## Market Size
Global Content Hub market size expected to reach $8.9 billion in 2024, with 18% annual growth rate.

## Domestic Trends
Accelerated domestic substitution, Content Hub developed with Go/Java languages gaining market share.',
    '# 2024 Content Hub Industry Trend Analysis

## Data Source
IDC 2024 Global Content Hub Market Research Report

## Core Trends
1. Lightweight: Lightweight Content Hub market share increased to 65%, replacing heavyweight systems
2. Private Deployment: Enterprise private deployment demand grew by 40%
3. AI Empowerment: AI-generated content and intelligent layout becoming standard features

## Market Size
Global Content Hub market size expected to reach $8.9 billion in 2024, with 18% annual growth rate.

## Domestic Trends
Accelerated domestic substitution, Content Hub developed with Go/Java languages gaining market share.',
    'https://picsum.photos/800/450?random=12',
    'post-analysis', '/en/blog/2024-cms-industry-trends',
    2600,
    'Content Hub,2024 Trends,Industry Analysis,Lightweight,Private Deployment,AI Empowerment',
    'Based on IDC 2024 industry report, analysis of three major Content Hub industry trends: lightweight architecture, private deployment, and AI-powered features.',
    '2024 Content Hub Industry Trend Analysis: Lightweight, Private Deployment, AI Empowerment | GoWind Industry Insights'
),
-- ========== 文章5：GoWind Content Hub 自定义模板开发 - 英文（草稿） ==========
(
    NOW() - INTERVAL '15 days', NOW(), 5, 'en-US',
    'GoWind Content Hub Custom Template Development Tutorial (Draft)',
    'gowind-custom-template-dev',
    'GoWind Content Hub custom template development tutorial covering template syntax, data binding, style customization, current progress 50%.',
    '# GoWind Content Hub Custom Template Development

> This article is under development, current progress 50%, expected completion date April 1, 2024.

## Development Prerequisites
### Required Skills
- Go language fundamentals
- Vue3 + Element Plus
- HTML/CSS/JS

## Template Directory Structure
/templates/custom/
  - index.tpl # Homepage template
  - post.tpl # Article template
  - style.css # Custom styles

## To Be Written
1. Template syntax detailed explanation
2. Data binding examples
3. Custom component development',
    '# GoWind Content Hub Custom Template Development

> This article is under development, current progress 50%, expected completion date April 1, 2024.

## Development Prerequisites
### Required Skills
- Go language fundamentals
- Vue3 + Element Plus
- HTML/CSS/JS

## Template Directory Structure
/templates/custom/
  - index.tpl # Homepage template
  - post.tpl # Article template
  - style.css # Custom styles

## To Be Written
1. Template syntax detailed explanation
2. Data binding examples
3. Custom component development',
    'https://picsum.photos/800/450?random=13',
    'post-draft', '/en/blog/gowind-custom-template-dev',
    1800,
    'GoWind,Content Hub,Custom Template,Development Tutorial,Go,Vue3',
    'GoWind Content Hub custom template development tutorial covering template syntax, data binding, style customization, current progress 50%.',
    'GoWind Content Hub Custom Template Development Tutorial (Draft) | GoWind Development Documentation'
),
-- ========== 文章6：GoWind 企业版功能详解 - 英文（加密） ==========
(
    NOW() - INTERVAL '12 days', NOW(), 6, 'en-US',
    'GoWind Content Hub Enterprise Edition Feature Details (Paid Exclusive)',
    'gowind-enterprise-features',
    'GoWind Enterprise Edition exclusive features: multi-tenancy management, advanced permissions, data backup, dedicated support, pricing ¥9999-19999/year.',
    '# GoWind Content Hub Enterprise Edition Feature Details

## Exclusive Features
1. Multi-tenancy Management: Manage multiple sites within single system with complete data isolation
2. Advanced Permissions: Fine-grained permission control by role/department
3. Data Backup: Automatic scheduled backups with cross-region disaster recovery support
4. Dedicated Support: 7*24 technical support
5. Custom Development: Custom feature modules on demand

## Pricing Plans
- Basic: ¥9,999/year (up to 10 tenants)
- Standard: ¥14,999/year (up to 50 tenants)
- Premium: ¥19,999/year (unlimited tenants)

## Trial Application
Contact support: 400-123-4567, 15-day free trial available.',
    '# GoWind Content Hub Enterprise Edition Feature Details

## Exclusive Features
1. Multi-tenancy Management: Manage multiple sites within single system with complete data isolation
2. Advanced Permissions: Fine-grained permission control by role/department
3. Data Backup: Automatic scheduled backups with cross-region disaster recovery support
4. Dedicated Support: 7*24 technical support
5. Custom Development: Custom feature modules on demand

## Pricing Plans
- Basic: ¥9,999/year (up to 10 tenants)
- Standard: ¥14,999/year (up to 50 tenants)
- Premium: ¥19,999/year (unlimited tenants)

## Trial Application
Contact support: 400-123-4567, 15-day free trial available.',
    'https://picsum.photos/800/450?random=14',
    'post-enterprise', '/en/blog/gowind-enterprise-features',
    2200,
    'GoWind,Content Hub,Enterprise Edition,Paid Features,Multi-tenancy,Advanced Permissions',
    'GoWind Enterprise Edition exclusive features: multi-tenancy management, advanced permissions, data backup, dedicated support, pricing ¥9999-19999/year.',
    'GoWind Content Hub Enterprise Edition Feature Details (Paid Exclusive) | GoWind Enterprise Services'
),
-- ========== 文章7：常见问题解答 - 英文（草稿） ==========
(
    NOW() - INTERVAL '10 days', NOW(), 7, 'en-US',
    'GoWind Content Hub FAQ (Frequently Asked Questions)',
    'gowind-faq',
    'GoWind Content Hub common questions summary covering installation, configuration, performance, upgrade topics, updated monthly.',
    '# GoWind Content Hub FAQ

## Installation Related
Q1: Database connection failure during installation?
A1: Check database address, port, username and password in .env file.

Q2: Cannot access service after startup?
A2: Check if port is occupied, ensure firewall allows port 8080.

## Configuration Related
Q3: How to enable multi-language support?
A3: Enable in backend Settings > Multi-language, upload translation files.

## To Be Added
- Performance optimization questions
- Upgrade related questions',
    '# GoWind Content Hub FAQ

## Installation Related
Q1: Database connection failure during installation?
A1: Check database address, port, username and password in .env file.

Q2: Cannot access service after startup?
A2: Check if port is occupied, ensure firewall allows port 8080.

## Configuration Related
Q3: How to enable multi-language support?
A3: Enable in backend Settings > Multi-language, upload translation files.

## To Be Added
- Performance optimization questions
- Upgrade related questions',
    'https://picsum.photos/800/450?random=15',
    'post-faq', '/en/blog/gowind-faq',
    1500,
    'GoWind,Content Hub,FAQ,Frequently Asked Questions,Installation,Configuration,Performance',
    'GoWind Content Hub common questions summary covering installation, configuration, performance, upgrade topics, updated monthly.',
    'GoWind Content Hub FAQ (Frequently Asked Questions) | GoWind Help Center'
),
-- ========== 文章8：GoWind Content Hub 性能优化指南 - 英文 ==========
(
    NOW() - INTERVAL '8 days', NOW(), 8, 'en-US',
    'GoWind Content Hub Performance Optimization Guide: Practical Experience from 50K to 100K QPS',
    'gowind-cms-performance-optimization',
    'Share GoWind Content Hub performance optimization practical experience covering database optimization, caching strategies, code-level improvements, achieving 100% QPS increase.',
    '# GoWind Content Hub Performance Optimization Guide

## Optimization Background
v1.9 version had only 50K QPS with 200ms response time, unable to meet high-concurrency requirements.

## Optimization Metrics
| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| QPS | 50,000 | 100,000 | 100% |
| Response Time | 200ms | 80ms | 60% |
| Database Load | 80% | 30% | 62.5% |

## Core Optimization Points
1. Database: New indexes, slow query optimization, read-write separation
2. Caching: Redis cache for categories/articles with dynamic expiration strategy
3. Code: Goroutine optimization, JSON serialization improvement, static resource compression',
    '# GoWind Content Hub Performance Optimization Guide

## Optimization Background
v1.9 version had only 50K QPS with 200ms response time, unable to meet high-concurrency requirements.

## Optimization Metrics
| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| QPS | 50,000 | 100,000 | 100% |
| Response Time | 200ms | 80ms | 60% |
| Database Load | 80% | 30% | 62.5% |

## Core Optimization Points
1. Database: New indexes, slow query optimization, read-write separation
2. Caching: Redis cache for categories/articles with dynamic expiration strategy
3. Code: Goroutine optimization, JSON serialization improvement, static resource compression',
    'https://picsum.photos/800/450?random=16',
    'post-tech', '/en/blog/gowind-cms-performance-optimization',
    3000,
    'GoWind,Content Hub,Performance Optimization,QPS,Database Optimization,Caching Strategy,Go Language',
    'Share GoWind Content Hub performance optimization practical experience covering database optimization, caching strategies, code-level improvements, achieving 100% QPS increase.',
    'GoWind Content Hub Performance Optimization Guide: Practical Experience from 50K to 100K QPS | GoWind Technical Blog'
);
SELECT setval('post_translations_id_seq', (SELECT MAX(id) FROM post_translations));

-- ----------------------------
-- 插入 post_categories 表（文章-分类关联）
-- ----------------------------
INSERT INTO public.post_categories (post_id, category_id, created_at)
VALUES
-- 文章1: [2,4] → 技术文档 + 安装部署
(1, 2, NOW() - INTERVAL '30 days'),
(1, 4, NOW() - INTERVAL '30 days'),
-- 文章2: [1,3] → 博客 + 产品中心
(2, 1, NOW() - INTERVAL '25 days'),
(2, 3, NOW() - INTERVAL '25 days'),
-- 文章3: [2,4,5] → 技术文档 + 安装部署 + 环境配置
(3, 2, NOW() - INTERVAL '22 days'),
(3, 4, NOW() - INTERVAL '22 days'),
(3, 5, NOW() - INTERVAL '22 days'),
-- 文章4: [1,9] → 博客 + 行业资讯
(4, 1, NOW() - INTERVAL '20 days'),
(4, 9, NOW() - INTERVAL '20 days'),
-- 文章5: [2,6] → 技术文档 + 功能教程
(5, 2, NOW() - INTERVAL '15 days'),
(5, 6, NOW() - INTERVAL '15 days'),
-- 文章6: [3] → 产品中心
(6, 3, NOW() - INTERVAL '12 days'),
-- 文章7: [2,7] → 技术文档 + 常见问题
(7, 2, NOW() - INTERVAL '10 days'),
(7, 7, NOW() - INTERVAL '10 days'),
-- 文章8: [1,2] → 博客 + 技术文档
(8, 1, NOW() - INTERVAL '8 days'),
(8, 2, NOW() - INTERVAL '8 days')
;
SELECT setval('post_categories_id_seq', (SELECT MAX(id) FROM post_categories));

-- ----------------------------
-- 插入 post_tags 表（文章-标签关联）
-- ----------------------------
INSERT INTO public.post_tags (post_id, tag_id, created_at)
VALUES
-- 文章1: [1,2,3] → Go语言 + Content Hub + 快速上手
(1, 1, NOW() - INTERVAL '30 days'),
(1, 2, NOW() - INTERVAL '30 days'),
(1, 3, NOW() - INTERVAL '30 days'),
-- 文章2: [4,5,6] → 版本更新 + 新功能 + 升级指南
(2, 4, NOW() - INTERVAL '25 days'),
(2, 5, NOW() - INTERVAL '25 days'),
(2, 6, NOW() - INTERVAL '25 days'),
-- 文章3: [1,7,8] → Go语言 + Linux + 部署教程
(3, 1, NOW() - INTERVAL '22 days'),
(3, 7, NOW() - INTERVAL '22 days'),
(3, 8, NOW() - INTERVAL '22 days'),
-- 文章4: [9,10,11] → 行业分析 + 2024趋势 + 轻量化
(4, 9, NOW() - INTERVAL '20 days'),
(4, 10, NOW() - INTERVAL '20 days'),
(4, 11, NOW() - INTERVAL '20 days'),
-- 文章5: [1,12,13] → Go语言 + 自定义模板 + 开发教程
(5, 1, NOW() - INTERVAL '15 days'),
(5, 12, NOW() - INTERVAL '15 days'),
(5, 13, NOW() - INTERVAL '15 days'),
-- 文章6: [5,14,15] → 新功能 + 企业版 + 付费功能
(6, 5, NOW() - INTERVAL '12 days'),
(6, 14, NOW() - INTERVAL '12 days'),
(6, 15, NOW() - INTERVAL '12 days'),
-- 文章7: [16,17,18] → FAQ + 问题解答 + 故障排除
(7, 16, NOW() - INTERVAL '10 days'),
(7, 17, NOW() - INTERVAL '10 days'),
(7, 18, NOW() - INTERVAL '10 days'),
-- 文章8: [1,19,20] → Go语言 + 性能优化 + 高并发
(8, 1, NOW() - INTERVAL '8 days'),
(8, 19, NOW() - INTERVAL '8 days'),
(8, 20, NOW() - INTERVAL '8 days')
;
SELECT setval('post_tags_id_seq', (SELECT MAX(id) FROM post_tags));

COMMIT;
