package serviceid

const (
	AdminService = "admin-service" // 后台服务
	AppService   = "app-service"   // 前台服务

	CoreService = "core-service" // 核心服务

	DTMService = "dtm-service" // DTM服务
)

var (
	DtmServiceAddress = MakeDiscoveryAddress(DTMService)
)
