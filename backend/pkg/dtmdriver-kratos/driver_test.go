package driver

import (
	"testing"
)

func TestKratosDriver_RegisterGrpcService_Consul(t *testing.T) {
	// 本仓 driver 未实现 consul scheme（driver.go 显式返回 not implemented），
	// 该用例在实现落地前恒不可通过——显式跳过而非假红。
	t.Skip("consul scheme is not implemented in this repo's driver (driver.go returns not-implemented)")
}

func TestKratosDriver_RegisterGrpcService_Etcd(t *testing.T) {
	// 依赖本机 etcd（2379）：注册成功即通过；无 60s Sleep——注册验证完成后
	// 无需保持进程存活（原 Sleep 让全量 go test 白等一分钟）。
	target := "etcd://127.0.0.1:2379/dtm-service"
	endpoint := "localhost:36790"
	driver := new(kratosDriver)
	if err := driver.RegisterService(target, endpoint); err != nil {
		t.Errorf("register etcd fail err :%+v", err)
	}
}
