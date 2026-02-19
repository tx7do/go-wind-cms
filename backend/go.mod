module go-wind-cms

go 1.25.3

require (
	entgo.io/ent v0.14.5
	github.com/alicebob/miniredis/v2 v2.36.1
	github.com/dtm-labs/dtm v1.19.0
	github.com/dtm-labs/dtmdriver v0.0.6
	github.com/envoyproxy/protoc-gen-validate v1.3.0
	github.com/getkin/kin-openapi v0.133.0
	github.com/go-kratos/kratos/v2 v2.9.2
	github.com/go-sql-driver/mysql v1.9.3
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/gnostic v0.7.1
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.7.0
	github.com/hibiken/asynq v0.26.0
	github.com/jackc/pgx/v4 v4.18.3
	github.com/jackc/pgx/v5 v5.8.0
	github.com/jinzhu/copier v0.4.0
	github.com/jinzhu/inflection v1.0.0
	github.com/lib/pq v1.10.9
	github.com/menta2k/protoc-gen-redact/v3 v3.0.0-20260213125431-7688a38967d4
	github.com/mileusna/useragent v1.3.5
	github.com/minio/minio-go/v7 v7.0.98
	github.com/redis/go-redis/v9 v9.18.0
	github.com/stretchr/testify v1.11.1
	github.com/tx7do/go-crud/api v0.0.7
	github.com/tx7do/go-crud/entgo v0.0.40
	github.com/tx7do/go-crud/pagination v0.0.11
	github.com/tx7do/go-crud/viewer v0.0.6
	github.com/tx7do/go-utils v1.1.34
	github.com/tx7do/go-utils/aggregator v0.0.4
	github.com/tx7do/go-utils/copierutil v0.0.6
	github.com/tx7do/go-utils/geoip v1.1.8
	github.com/tx7do/go-utils/id v0.0.3
	github.com/tx7do/go-utils/jwtutil v0.0.3
	github.com/tx7do/go-utils/mapper v0.0.3
	github.com/tx7do/go-utils/password v0.0.1
	github.com/tx7do/kratos-authn v1.1.9
	github.com/tx7do/kratos-authn/engine/jwt v1.1.9
	github.com/tx7do/kratos-authz v1.1.7
	github.com/tx7do/kratos-authz/engine/casbin v1.1.11
	github.com/tx7do/kratos-authz/engine/opa v1.1.14
	github.com/tx7do/kratos-authz/middleware v1.1.12
	github.com/tx7do/kratos-bootstrap/api v0.0.35
	github.com/tx7do/kratos-bootstrap/bootstrap v0.1.16
	github.com/tx7do/kratos-bootstrap/cache/redis v0.1.1
	github.com/tx7do/kratos-bootstrap/database/ent v0.1.4
	github.com/tx7do/kratos-bootstrap/oss/minio v0.1.1
	github.com/tx7do/kratos-bootstrap/registry v0.2.2
	github.com/tx7do/kratos-bootstrap/registry/etcd v0.2.2
	github.com/tx7do/kratos-bootstrap/rpc v0.1.1
	github.com/tx7do/kratos-bootstrap/tracer v0.1.3
	github.com/tx7do/kratos-bootstrap/transport/sse v0.0.4
	github.com/tx7do/kratos-swagger-ui v0.0.1
	github.com/tx7do/kratos-transport/broker v1.3.2
	github.com/tx7do/kratos-transport/transport/sse v1.3.2
	github.com/yuin/gopher-lua v1.1.1
	go.etcd.io/etcd/client/v3 v3.6.7
	go.opentelemetry.io/otel/trace v1.40.0
	google.golang.org/genproto v0.0.0-20240711142825-46eb208f015d
	google.golang.org/genproto/googleapis/api v0.0.0-20260217215200-42d3e9bedb6d
	google.golang.org/grpc v1.79.1
	google.golang.org/protobuf v1.36.11
)

require (
	ariga.io/atlas v1.0.0 // indirect
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260209202127-80ab13bee0bf.1 // indirect
	buf.build/go/protovalidate v1.1.2 // indirect
	cel.dev/expr v0.25.1 // indirect
	dario.cat/mergo v1.0.2 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/XSAM/otelsql v0.41.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/bmatcuk/doublestar/v4 v4.9.1 // indirect
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/casbin/casbin/v2 v2.135.0 // indirect
	github.com/casbin/govaluate v1.10.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/clipperhouse/displaywidth v0.6.2 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.6.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dtm-labs/logger v0.0.2 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/inflect v0.21.5 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/go-resty/resty/v2 v2.7.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/cel-go v0.27.0 // indirect
	github.com/google/gnostic-models v0.7.1 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/subcommands v1.2.0 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/graph-gophers/dataloader/v7 v7.1.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.8 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/klauspost/crc32 v1.3.0 // indirect
	github.com/lestrrat-go/blackmagic v1.0.4 // indirect
	github.com/lestrrat-go/dsig v1.0.0 // indirect
	github.com/lestrrat-go/dsig-secp256k1 v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc/v3 v3.0.4 // indirect
	github.com/lestrrat-go/jwx/v3 v3.0.13 // indirect
	github.com/lestrrat-go/option/v2 v2.0.0 // indirect
	github.com/lithammer/shortuuid/v4 v4.2.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20260216142805-b3301c5f2a88 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/minio/crc64nvme v1.1.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/oasdiff/yaml v0.0.0-20250309154309-f31be36b4037 // indirect
	github.com/oasdiff/yaml3 v0.0.0-20250309153720-d2182401db90 // indirect
	github.com/olekukonko/cat v0.0.0-20250911104152-50322a0618f6 // indirect
	github.com/olekukonko/errors v1.1.0 // indirect
	github.com/olekukonko/ll v0.1.3 // indirect
	github.com/olekukonko/tablewriter v1.1.2 // indirect
	github.com/open-policy-agent/opa v1.13.1 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/oschwald/geoip2-golang v1.13.0 // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	github.com/redis/go-redis/extra/rediscmd/v9 v9.17.2 // indirect
	github.com/redis/go-redis/extra/redisotel/v9 v9.17.2 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/segmentio/ksuid v1.0.4 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shoenig/go-m1cpu v0.1.7 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/sony/sonyflake v1.3.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/cobra v1.10.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/swaggest/swgui v1.8.5 // indirect
	github.com/tchap/go-patricia/v2 v2.3.3 // indirect
	github.com/tinylib/msgp v1.6.1 // indirect
	github.com/tklauser/go-sysconf v0.3.16 // indirect
	github.com/tklauser/numcpus v0.11.0 // indirect
	github.com/tx7do/go-crud/audit v0.0.2 // indirect
	github.com/tx7do/kratos-bootstrap/config v0.2.2 // indirect
	github.com/tx7do/kratos-bootstrap/logger v0.1.2 // indirect
	github.com/tx7do/kratos-transport/tracing v1.1.1 // indirect
	github.com/tx7do/kratos-transport/transport v1.3.3 // indirect
	github.com/valyala/fastjson v1.6.7 // indirect
	github.com/vearutop/statigz v1.5.0 // indirect
	github.com/vektah/gqlparser/v2 v2.5.31 // indirect
	github.com/woodsbury/decimal128 v1.3.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.2.0 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zclconf/go-cty v1.17.0 // indirect
	github.com/zclconf/go-cty-yaml v1.2.0 // indirect
	go.einride.tech/aip v0.80.0 // indirect
	go.etcd.io/etcd/api/v3 v3.6.7 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.6.7 // indirect
	go.mongodb.org/mongo-driver v1.17.7 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.40.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.40.0 // indirect
	go.opentelemetry.io/otel/metric v1.40.0 // indirect
	go.opentelemetry.io/otel/sdk v1.40.0 // indirect
	go.opentelemetry.io/proto/otlp v1.9.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/exp v0.0.0-20260212183809-81e46e3db34a // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	golang.org/x/tools v0.42.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260217215200-42d3e9bedb6d // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	gopkg.in/ini.v1 v1.67.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)
