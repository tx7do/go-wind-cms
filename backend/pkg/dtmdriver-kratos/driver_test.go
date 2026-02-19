package driver

import (
	"testing"
	"time"
)

func TestKratosDriver_RegisterGrpcService_Consul(t *testing.T) {
	target := "consul://127.0.0.1:8500/dtm-service"
	endpoint := "localhost:36790"
	driver := new(kratosDriver)
	if err := driver.RegisterService(target, endpoint); err != nil {
		t.Errorf("register consul fail err :%+v", err)
	}

	time.Sleep(60 * time.Second)
}

func TestKratosDriver_RegisterGrpcService_Etcd(t *testing.T) {
	target := "etcd://127.0.0.1:2379/dtm-service"
	endpoint := "localhost:36790"
	driver := new(kratosDriver)
	if err := driver.RegisterService(target, endpoint); err != nil {
		t.Errorf("register etcd fail err :%+v", err)
	}

	time.Sleep(60 * time.Second)
}
