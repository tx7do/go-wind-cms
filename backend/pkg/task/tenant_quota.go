package task

// TenantQuotaWatermarkTaskType 是租户配额水位扫描任务的类型常量。
// 系统级常驻定时任务：每日扫描全部 ON 租户的四类租户配额（USER_LIMIT / STORAGE /
// API_CALL / AI_TOKENS 月度）用量水位，对达到阈值的租户，按其管理员偏好语言渲染
// 站内信告警并投递给租户管理员（sys_tenants.admin_user_id）。
// handler 为 TenantQuotaService.AsyncTenantQuotaWatermarkScan。
const TenantQuotaWatermarkTaskType = "tenant_quota_watermark_scan"

// TenantQuotaWatermarkCronSpec 水位扫描的 cron（每天 09:00）。
const TenantQuotaWatermarkCronSpec = "0 9 * * *"

// TenantQuotaWatermarkTaskData 水位扫描任务载荷（无参数：扫描对象为全部 ON 租户）。
type TenantQuotaWatermarkTaskData struct{}
