package constants

// sys_config 平台参数键与默认值：口令复杂度/有效期/历史口令策略（等保要求）。
// 阈值经参数管理页（sys_configs，is_built_in）可改，修改落库后即时生效
//（ConfigRepo 缓存写路径同步失效）。口令校验消费方（复杂度/有效期/历史检查）
// 为上游等保批能力，CMS 当前未移植——键与默认值先随配置域落位，消费方移植时即插即用。
const (
	// ConfigKeyPasswordMinLen 最小长度，默认 DefaultPasswordMinLen。
	ConfigKeyPasswordMinLen = "sys.password.minLen"
	// ConfigKeyPasswordMaxAgeDays 有效期天数；<=0 表示不启用有效期。
	ConfigKeyPasswordMaxAgeDays = "sys.password.maxAgeDays"
	// ConfigKeyPasswordHistoryCount 历史口令保留条数；<=0 表示不启用历史检查。
	ConfigKeyPasswordHistoryCount = "sys.password.historyCount"
)

const (
	DefaultPasswordMinLen       = 8
	DefaultPasswordMaxAgeDays   = 90
	DefaultPasswordHistoryCount = 3
)
