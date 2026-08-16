package client

import (
	"errors"

	opensearchCrud "github.com/tx7do/go-crud/opensearch"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-bootstrap/database/opensearch"
)

func NewElasticSearchClient(ctx *bootstrap.Context) (*opensearchCrud.Client, func(), error) {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil, func() {}, nil
	}

	if cfg.Data != nil && cfg.Data.Elasticsearch != nil {
		pw := cfg.Data.Elasticsearch.GetPassword()
		if pw == "" {
			return nil, func() {}, errors.New("OpenSearch password must be set via OPENSEARCH_PASSWORD environment variable; refusing to start with default/empty password")
		}
	}

	cli, err := opensearch.NewClient(ctx.GetLogger(), cfg)
	if err != nil {
		return nil, func() {}, err
	}

	return cli, func() {
	}, nil
}
