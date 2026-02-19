package driver

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dtm-labs/dtmdriver"

	"github.com/go-kratos/kratos/v2/registry"
	"google.golang.org/grpc/resolver"

	bEtcd "github.com/tx7do/kratos-bootstrap/registry/etcd"

	_ "github.com/go-kratos/kratos/v2/transport/grpc/resolver/direct"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"

	etcdAPI "go.etcd.io/etcd/client/v3"
)

const (
	Name = "dtm-driver-kratos"

	DefaultScheme = "discovery"
	EtcdScheme    = "etcd"
	ConsulScheme  = "consul"
)

type kratosDriver struct{}

func (k *kratosDriver) GetName() string {
	return Name
}

func (k *kratosDriver) RegisterAddrResolver() {

}

func (k *kratosDriver) RegisterService(target string, endpoint string) error {
	if target == "" {
		return nil
	}

	u, err := url.Parse(target)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case DefaultScheme:
		fallthrough

	case EtcdScheme:
		registerInstance := &registry.ServiceInstance{
			Name:      strings.TrimPrefix(u.Path, "/"),
			Endpoints: strings.Split(endpoint, ","),
		}
		client, err := etcdAPI.New(etcdAPI.Config{
			Endpoints: strings.Split(u.Host, ","),
		})
		if err != nil {
			return err
		}
		reg := bEtcd.New(client)
		//add resolver so that dtm can handle discovery://
		resolver.Register(discovery.NewBuilder(reg, discovery.WithInsecure(true)))
		return reg.Register(context.Background(), registerInstance)

	case ConsulScheme:
		return fmt.Errorf("consul scheme not implemented yet")

	default:
		return fmt.Errorf("unknown scheme: %s", u.Scheme)
	}
}

func (k *kratosDriver) ParseServerMethod(uri string) (server string, method string, err error) {
	if !strings.Contains(uri, "//") {
		sep := strings.IndexByte(uri, '/')
		if sep == -1 {
			return "", "", fmt.Errorf("bad url: '%s'. no '/' found", uri)
		}
		return uri[:sep], uri[sep:], nil

	}
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", nil
	}
	index := strings.IndexByte(u.Path[1:], '/') + 1
	return u.Scheme + "://" + u.Host + u.Path[:index], u.Path[index:], nil
}

func init() {
	dtmdriver.Register(&kratosDriver{})
}
