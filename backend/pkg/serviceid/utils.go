package serviceid

// NewDiscoveryName 构建服务发现名称
func NewDiscoveryName(serviceName string) string {
	return ProjectName + "/" + serviceName
}

// MakeDiscoveryAddress 构建服务发现地址
func MakeDiscoveryAddress(serviceName string) string {
	return "discovery:///" + serviceName
}
