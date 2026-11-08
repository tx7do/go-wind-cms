package rpc

import (
	"crypto/tls"

	tlsUtils "github.com/tx7do/go-utils/tls"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func loadServerTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var tlsCfg *tls.Config
	var err error

	if cfg.File != nil {
		if tlsCfg, err = tlsUtils.LoadServerTlsConfigFile(
			cfg.File.GetKeyPath(),
			cfg.File.GetCertPath(),
			cfg.File.GetCaPath(),
			cfg.InsecureSkipVerify,
		); err != nil {
			return nil, err
		}
	} else if cfg.Config != nil {
		if tlsCfg, err = tlsUtils.LoadServerTlsConfigString(
			cfg.Config.GetKeyPem(),
			cfg.Config.GetCertPem(),
			cfg.Config.GetCaPem(),
			cfg.InsecureSkipVerify,
		); err != nil {
			return nil, err
		}
	}

	return tlsCfg, err
}

func loadClientTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var tlsCfg *tls.Config
	var err error

	if cfg.File != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigFile(
			cfg.File.GetKeyPath(),
			cfg.File.GetCertPath(),
			cfg.File.GetCaPath(),
		); err != nil {
			return nil, err
		}
	} else if cfg.Config != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigString(
			cfg.Config.GetKeyPem(),
			cfg.Config.GetCertPem(),
			cfg.Config.GetCaPem(),
		); err != nil {
			return nil, err
		}
	}

	return tlsCfg, err
}
