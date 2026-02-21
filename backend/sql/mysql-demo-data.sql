-- Description: 初始化测试数据（租户/用户/组织架构/字典/站内信分类等）- MySQL版
-- Note: 表结构创建完成后执行，执行前备份数据；MySQL需5.7+（支持JSON/COALESCE）
-- 关闭外键检查，避免TRUNCATE关联表报错
SET FOREIGN_KEY_CHECKS = 0;
-- 开启事务，保证数据原子性
START TRANSACTION;

-- 一次性清理相关表并重置自增（MySQL无外键级联TRUNCATE，先关外键检查）
TRUNCATE TABLE sys_org_units AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_positions AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_tasks AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_login_policies AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_dict_types AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_dict_entries AUTO_INCREMENT = 1;
TRUNCATE TABLE internal_message_categories AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_tenants AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_users AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_user_credentials AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_memberships AUTO_INCREMENT = 1;
TRUNCATE TABLE sys_membership_roles AUTO_INCREMENT = 1;

-- 测试租户
INSERT INTO sys_tenants(id, name, code, type, audit_status, status, admin_user_id, created_at)
VALUES (1, '测试租户', 'super', 'PAID', 'APPROVED', 'ON', 2, NOW());
-- 重置自增，后续新增从MAX(id)+1开始（COALESCE处理空表情况）
ALTER TABLE sys_tenants AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_tenants);

-- 插入租户管理员用户
INSERT INTO sys_users (id, tenant_id, username, nickname, realname, email, gender, created_at)
VALUES (2, 1, 'tenant_admin', '租户管理', '张管理员', 'tenant@company.com', 'MALE', NOW());
ALTER TABLE sys_users AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_users);

-- 插入租户管理员的登录凭证（密码：admin，与平台管理员哈希值一致）
INSERT INTO sys_user_credentials (user_id, identity_type, identifier, credential_type, credential, status,
                                  is_primary, created_at)
VALUES
    -- 租户管理员（对应users表id=2）
    (2, 'USERNAME', 'tenant_admin', 'PASSWORD_HASH', '$2a$10$yajZDX20Y40FkG0Bu4N19eXNqRizez/S9fK63.JxGkfLq.RoNKR/a',
     'ENABLED', true, NOW()),
    (2, 'EMAIL', 'tenant@company.com', 'PASSWORD_HASH',
     '$2a$10$yajZDX20Y40FkG0Bu4N19eXNqRizez/S9fK63.JxGkfLq.RoNKR/a', 'ENABLED', false, NOW());
ALTER TABLE sys_user_credentials AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_user_credentials);

-- 组织架构单元
INSERT INTO sys_org_units (id, tenant_id, parent_id, type, name, code, description, path, sort_order, leader_id, status, created_at)
VALUES
    (1, 1, NULL, 'COMPANY', 'XX集团总部', 'HEADQUARTERS', '集团核心管理机构，统筹全集团战略规划、业务管控及资源调配', '/1', 1, 1, 'ON', NOW()),
    (2, 1, 1, 'DIVISION', '技术部', 'TECH', '负责集团整体技术架构规划、研发管理、系统运维及技术创新', '/1/2', 2, 5, 'ON', NOW()),
    (3, 1, 1, 'DIVISION', '财务部', 'FIN', '负责集团财务核算、资金管理、税务筹划、预算编制及财务风控', '/1/3', 3, 8, 'ON', NOW()),
    (4, 1, 1, 'DIVISION', '人事部', 'HR', '负责人力资源规划、招聘配置、薪酬绩效、员工培训及组织发展', '/1/4', 4, 9, 'ON', NOW()),
    (5, 1, 2, 'DEPARTMENT', '研发一部', 'DEV-1', '聚焦新能源领域产品研发、技术迭代及核心模块开发', '/1/2/5', 1, 6, 'ON', NOW()),
    (6, 1, 1, 'REGION', '华北大区', 'NORTH', '负责华北区域市场运营、客户维护、销售管理及本地化服务落地', '/1/6', 3, 12, 'ON', NOW()),
    (7, 1, 1, 'SUBSIDIARY', '广州分公司', 'GZ', '负责华南区域（广州及周边）业务拓展、客户服务及本地化运营', '/1/7', 5, 2, 'ON', NOW()),
    (8, 1, 1, 'SUBSIDIARY', '深圳子公司', 'SZ', '负责深圳区域市场开拓、科技创新业务落地及高端客户对接', '/1/8', 6, 4, 'ON', NOW()),
    (9, 1, 1, 'DIVISION', '销售部', 'SALES', '统筹集团整体销售策略制定、销售团队管理及业绩目标达成', '/1/9', 7, 16, 'ON', NOW()),
    (10, 1, 9, 'DEPARTMENT', '海外事业部', 'INTL', '负责海外市场拓展、国际客户合作、跨境业务管理及本地化运营', '/1/9/10', 1, 17, 'ON', NOW()),
    (11, 1, 10, 'TEAM', '海外销售组', 'INTL-SALES-1', '具体执行海外市场销售任务，跟进客户需求及订单落地', '/1/9/10/11', 1, 18, 'ON', NOW()),
    (12, 1, 5, 'PROJECT', '新能源项目组', 'NEO-PROJ', '专项负责新能源项目的研发、落地、运营及成果转化', '/1/2/5/12', 1, 6, 'ON', NOW()),
    (13, 1, 1, 'COMMITTEE', '审计委员会', 'AUDIT', '独立开展集团内部审计、风控检查、合规监督及问题整改跟进', '/1/13', 8, 12, 'ON', NOW()),
    (14, 1, 1, 'DEPARTMENT', '客服部', 'CS', '负责全集团客户咨询、投诉处理、售后服务及客户满意度提升', '/1/14', 9, 11, 'ON', NOW()),
    (15, 1, 14, 'TEAM', '客服一组', 'CS-1', '承接华南区域客户服务、售后问题处理及客户关系维护', '/1/14/15', 1, 20, 'ON', NOW());
ALTER TABLE sys_org_units AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_org_units);

-- 岗位数据
INSERT INTO sys_positions (id, tenant_id, type, name, code, org_unit_id, reports_to_position_id, description, job_family, job_grade, level, headcount, is_key_position, status, sort_order, created_at)
VALUES
    (1, 1, 'LEADER', '技术总监', 'TECH-DIRECTOR-001', 2, NULL, '负责公司整体技术战略规划、团队管理及核心技术决策', 'TECH', 1, 1, 1, true, 'ON', 1, NOW()),
    (2, 1, 'MANAGER', '技术部经理', 'TECH-MANAGER-001', 2, 1, '负责技术部日常管理、项目排期及团队协作', 'TECH', 2, 2, 1, true, 'ON', 2, NOW()),
    (3, 1, 'MANAGER', '前端主管', 'TECH-FE-LEADER-001', 2, 2, '负责前端团队开发管理、技术方案评审及需求落地', 'TECH', 3, 3, 3, false, 'ON', 3, NOW()),
    (4, 1, 'MANAGER', '后端主管', 'TECH-BE-LEADER-001', 2, 2, '负责后端服务架构设计、数据库优化及接口开发管理', 'TECH', 4, 3, 3, false, 'ON', 4, NOW()),
    (5, 1, 'REGULAR', '前端开发专员', 'TECH-FE-DEV-001', 2, 3, '负责Web/移动端前端页面开发、交互实现及兼容性优化', 'TECH', 5, 4, 5, false, 'ON', 5, NOW()),
    (6, 1, 'REGULAR', '后端开发专员', 'TECH-BE-DEV-001', 2, 4, '负责后端接口开发、业务逻辑实现及系统稳定性维护', 'TECH', 6, 4, 5, false, 'ON', 6, NOW()),
    (7, 1, 'REGULAR', '测试工程师', 'TECH-TEST-001', 2, 2, '负责项目功能测试、性能测试及自动化测试脚本开发', 'TECH', 3, 4, 3, false, 'ON', 7, NOW()),
    (8, 1, 'LEADER', '人力总监', 'HR-DIRECTOR-001', 2, NULL, '负责人力资源战略规划、组织架构设计及人才梯队建设', 'HR', 1, 1, 1, true, 'ON', 1, NOW()),
    (9, 1, 'MANAGER', '招聘主管', 'HR-RECRUIT-LEADER-001', 2, 8, '负责公司各部门招聘需求对接、简历筛选及面试安排', 'HR', 2, 2, 1, false, 'ON', 2, NOW()),
    (10, 1, 'REGULAR', '薪酬绩效专员', 'HR-C&P-001', 2, 8, '负责员工薪酬核算、绩效考核制度落地及社保公积金管理', 'HR', 3, 2, 1, false, 'ON', 3, NOW()),
    (11, 1, 'REGULAR', 'HRBP', 'HR-BP-001', 2, 8, '对接业务部门，提供人力资源支持（入离职、员工关系等）', 'HR', 4, 2, 1, false, 'ON', 4, NOW()),
    (12, 1, 'LEADER', '财务总监', 'FIN-DIRECTOR-001', 2, NULL, '负责公司财务战略、预算管理及财务风险控制', 'FIN', 1, 1, 1, true, 'ON', 1, NOW()),
    (13, 1, 'MANAGER', '会计主管', 'FIN-ACCOUNT-LEADER-001', 2, 12, '负责账务处理、财务报表编制及税务申报管理', 'FIN', 2, 2, 1, false, 'ON', 2, NOW()),
    (14, 1, 'REGULAR', '出纳专员', 'FIN-CASHIER-001', 2, 13, '负责日常资金收付、银行对账及票据管理', 'FIN', 3, 3, 1, false, 'ON', 3, NOW()),
    (15, 1, 'REGULAR', '成本会计', 'FIN-COST-001', 2, 13, '负责成本核算、成本分析及成本控制方案制定', 'FIN', 4, 3, 1, false, 'ON', 4, NOW()),
    (16, 1, 'LEADER', '市场总监', 'MKT-DIRECTOR-001', 4, NULL, '负责市场战略规划、品牌建设及营销活动策划', 'MKT', 1, 1, 1, true, 'ON', 1, NOW()),
    (17, 1, 'MANAGER', '新媒体运营主管', 'MKT-NEWS-LEADER-001', 4, 16, '负责新媒体平台内容运营及用户增长', 'MKT', 2, 2, 1, false, 'ON', 2, NOW()),
    (18, 1, 'REGULAR', '活动策划专员', 'MKT-EVENT-001', 4, 16, '负责线下活动策划、执行及效果复盘', 'MKT', 3, 3, 1, false, 'ON', 3, NOW()),
    (19, 1, 'REGULAR', '市场调研专员', 'MKT-RESEARCH-001', 4, 16, '负责行业动态调研、竞品分析及市场趋势报告撰写', 'MKT', 4, 3, 1, false, 'ON', 4, NOW()),
    (20, 1, 'REGULAR', '行政助理', 'ADMIN-ASSIST-001', 2, 8, '负责办公用品采购、会议安排等行政工作（已合并至HRBP）', 'ADMIN', 5, 5, 1, false, 'OFF', 5, NOW());
ALTER TABLE sys_positions AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_positions);

-- 用户-租户关联关系
INSERT INTO sys_memberships (id, tenant_id, user_id, org_unit_id, position_id, role_id, is_primary, status)
VALUES (2, 1, 2, null, null, 2, true, 'ACTIVE');
ALTER TABLE sys_memberships AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_memberships);

-- 租户成员-角色关联关系（修复原脚本少分号、缺分隔符错误）
INSERT INTO sys_membership_roles (id, membership_id, tenant_id, role_id, is_primary, status)
VALUES (2, 2, 1, 2, true, 'ACTIVE');
ALTER TABLE sys_membership_roles AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_membership_roles);

-- 调度任务（JSON字段直接保留字符串，MySQL5.7+支持）
INSERT INTO sys_tasks(type, type_name, task_payload, cron_spec, enable, created_at)
VALUES ('PERIODIC', 'backup', '{ "name": "test"}', '0 * * * *', true, NOW());
ALTER TABLE sys_tasks AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_tasks);

-- 登录策略
INSERT INTO sys_login_policies(id, target_id, type, method, value, reason, created_at)
VALUES
    (1, 1, 'BLACKLIST', 'IP', '127.0.0.1', '无理由', NOW()),
    (2, 1, 'WHITELIST', 'MAC', '00:1B:44:11:3A:B7 ', '无理由', NOW());
ALTER TABLE sys_login_policies AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_login_policies);

-- 字典类型（移除PG的序列重置指令，MySQL TRUNCATE已重置）
INSERT INTO sys_dict_types(id, type_code, type_name, sort_order, description, is_enabled, created_at)
VALUES
    (1, 'USER_STATUS', '用户状态', 10, '系统用户的状态管理，包括正常、冻结、注销', true, NOW()),
    (2, 'DEVICE_TYPE', '设备类型', 20, 'IoT平台接入的设备品类，新增需同步至设备接入模块', true, NOW()),
    (3, 'ORDER_STATUS', '订单状态', 30, '电商订单的全生命周期状态', true, NOW()),
    (4, 'GENDER', '性别', 40, '用户性别枚举，默认未知', true, NOW()),
    (5, 'PAYMENT_METHOD', '支付方式', 50, '支持的支付渠道，含第三方支付和自有渠道', true, NOW());
ALTER TABLE sys_dict_types AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_dict_types);

-- 字典条目
INSERT INTO sys_dict_entries(id, type_id, entry_value, entry_label, numeric_value, sort_order, description, is_enabled, created_at)
VALUES
    -- 用户状态
    (1, 1, 'NORMAL', '正常', 1, 1, '用户可正常登录和操作', true, NOW()),
    (2, 1, 'FROZEN', '冻结', 2, 2, '因违规被临时冻结，需管理员解冻', true, NOW()),
    (3, 1, 'CANCELED', '注销', 3, 3, '用户主动注销，数据保留但不可登录', true, NOW()),
    -- 设备类型
    (4, 2, 'TEMP_SENSOR', '温湿度传感器', 101, 1, '支持温度（-20~80℃）和湿度（0~100%RH）采集', true, NOW()),
    (5, 2, 'CURRENT_METER', '电流仪表', 102, 2, '交流/直流电流测量，精度0.5级', true, NOW()),
    (6, 2, 'GAS_DETECTOR', '气体探测器', 103, 3, '暂不支持，待硬件适配（2025Q4计划启用）', false, NOW()),
    -- 订单状态
    (7, 3, 'PENDING', '待支付', 1, 1, '下单后未支付，超时自动取消', true, NOW()),
    (8, 3, 'PAID', '已支付', 2, 2, '支付成功，等待发货', true, NOW()),
    (9, 3, 'SHIPPED', '已发货', 3, 3, '商品已出库，物流配送中', true, NOW()),
    (10, 3, 'COMPLETED', '已完成', 4, 4, '用户确认收货，订单结束', true, NOW()),
    (11, 3, 'CANCELED', '已取消', 5, 5, '用户或系统取消订单', true, NOW()),
    -- 性别
    (12, 4, 'MALE', '男', 1, 1, '', true, NOW()),
    (13, 4, 'FEMALE', '女', 2, 2, '', true, NOW()),
    (14, 4, 'UNKNOWN', '未知', 0, 3, '用户未填写时默认值', true, NOW()),
    -- 支付方式
    (15, 5, 'ALIPAY', '支付宝', 1, 1, '支持花呗、余额宝', true, NOW()),
    (16, 5, 'WECHAT', '微信支付', 2, 2, '需绑定微信', true, NOW()),
    (17, 5, 'UNIONPAY', '银联支付', 3, 3, '支持信用卡、储蓄卡', true, NOW()),
    (18, 5, 'CASH', '现金支付', 4, 4, '线下支付，已废弃（2025-01停用）', false, NOW());
ALTER TABLE sys_dict_entries AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM sys_dict_entries);

-- 站内信分类（扁平结构，无树形层级）
INSERT INTO internal_message_categories (id, code, name, remark, sort_order, is_enabled, created_at)
VALUES
    -- 订单相关分类
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
    (403, 'user_permission_changed', '权限变更', '账号角色或功能权限调整通知', 17, true, NOW());
ALTER TABLE internal_message_categories AUTO_INCREMENT = (SELECT COALESCE(MAX(id) + 1, 1) FROM internal_message_categories);

-- 提交事务+恢复外键检查
COMMIT;
SET FOREIGN_KEY_CHECKS = 1;

SELECT '测试数据初始化完成' AS result;