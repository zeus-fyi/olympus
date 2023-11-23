module github.com/zeus-fyi/olympus

go 1.21

require (
	filippo.io/age v1.1.1
	github.com/PagerDuty/go-pagerduty v1.6.0
	github.com/aws/aws-sdk-go-v2 v1.21.0
	github.com/aws/aws-sdk-go-v2/config v1.18.11
	github.com/aws/aws-sdk-go-v2/credentials v1.13.11
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.25
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.39
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.19.7
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.119.0
	github.com/aws/aws-sdk-go-v2/service/eks v1.27.12
	github.com/aws/aws-sdk-go-v2/service/pricing v1.19.6
	github.com/aws/aws-sdk-go-v2/service/s3 v1.29.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.18.2
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.17.1
	github.com/aws/smithy-go v1.14.2
	github.com/caio/go-tdigest/v4 v4.0.1
	github.com/cavaliergopher/grab/v3 v3.0.1
	github.com/cbergoon/merkletree v0.2.0
	github.com/cockroachdb/errors v1.9.1
	github.com/confluentinc/confluent-kafka-go v1.9.2
	github.com/cvcio/twitter v0.0.0-20230702171346-58f04a1d303c
	github.com/dave/jennifer v1.5.1
	github.com/dghubble/go-twitter v0.0.0-20221104224141-912508c3888b
	github.com/dghubble/oauth1 v0.7.2
	github.com/digitalocean/godo v1.98.0
	github.com/ethereum/go-ethereum v1.10.26
	github.com/g8rswimmer/go-twitter/v2 v2.1.5
	github.com/ghodss/yaml v1.0.0
	github.com/go-redis/redis/v9 v9.0.0-beta.2
	github.com/go-resty/resty/v2 v2.7.0
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/google/go-querystring v1.1.0
	github.com/gorilla/schema v1.2.0
	github.com/iancoleman/strcase v0.2.0
	github.com/jackc/pgx/v4 v4.18.1
	github.com/klauspost/compress v1.17.2
	github.com/labstack/echo/v4 v4.10.2
	github.com/labstack/gommon v0.4.0
	github.com/lib/pq v1.10.6
	github.com/nanmu42/etherscan-api v1.10.0
	github.com/oleiade/lane/v2 v2.0.0
	github.com/ovh/go-ovh v1.4.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/phf/go-queue v0.0.0-20170504031614-9abe38d0371d
	github.com/pierrec/lz4 v2.6.1+incompatible
	github.com/pingcap/errors v0.11.4
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.63.0
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.63.0
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/common v0.39.0
	github.com/puzpuzpuz/xsync/v2 v2.4.1
	github.com/rs/zerolog v1.29.0
	github.com/sashabaranov/go-gpt3 v0.0.0-20221216095610-1c20931ead68
	github.com/sashabaranov/go-openai v1.17.8
	github.com/sendgrid/rest v2.6.9+incompatible
	github.com/sendgrid/sendgrid-go v3.12.0+incompatible
	github.com/shirou/gopsutil/v3 v3.22.10
	github.com/shopspring/decimal v1.3.1
	github.com/spf13/viper v1.15.0
	github.com/stretchr/testify v1.8.4
	github.com/stripe/stripe-go/v74 v74.13.0
	github.com/supranational/blst v0.3.11-0.20230406105308-e9dfc5ee724b
	github.com/tidwall/pretty v1.2.1
	github.com/tyler-smith/go-bip32 v1.0.0
	github.com/vartanbeno/go-reddit/v2 v2.0.1
	github.com/wealdtech/go-ed25519hd v0.0.0-20220222130843-fd974f26091e
	github.com/wealdtech/go-eth2-util v1.8.1
	github.com/zeus-fyi/gochain/v4 v4.2.0
	github.com/zeus-fyi/jennifer v0.0.0-20221008222717-be57190219a6
	github.com/zeus-fyi/memoryfs v0.0.0-20221102185216-a3f87725e7fb
	github.com/zeus-fyi/tables-to-go v0.0.0-20221022223825-758445568e31
	github.com/zeus-fyi/tojen v0.0.0-20221009050240-e9617f6e35ef
	go.temporal.io/api v1.11.1-0.20220907050538-6de5285cf463
	go.temporal.io/sdk v1.17.0
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29
	google.golang.org/api v0.143.0
	k8s.io/api v0.26.1
	k8s.io/apimachinery v0.26.1
	k8s.io/cli-runtime v0.26.1
	k8s.io/client-go v0.26.1
	k8s.io/kubectl v0.26.1
	logur.dev/adapter/zerolog v0.6.0
	logur.dev/logur v0.17.0
	nhooyr.io/websocket v1.8.10
)

require (
	cloud.google.com/go/compute v1.23.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/AnimeKaizoku/cacher v1.0.1 // indirect
	github.com/Arman92/go-tdlib v1.0.0 // indirect
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/FactomProject/btcutilecc v0.0.0-20130527213604-d3a63a5752ec // indirect
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6 // indirect
	github.com/aloder/tojen v0.0.0-20200427151940-c77d77b33821 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.9 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.41 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.14.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.27 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.18.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/bwmarrin/discordgo v0.27.1 // indirect
	github.com/celestix/gotgproto v1.0.0-beta13 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/dghubble/sling v1.4.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emicklei/go-restful/v3 v3.10.1 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/ferranbt/fastssz v0.1.3 // indirect
	github.com/getsentry/sentry-go v0.18.0 // indirect
	github.com/go-faster/errors v0.6.1 // indirect
	github.com/go-faster/jx v1.1.0 // indirect
	github.com/go-faster/xor v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/gobwas/ws v1.3.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/status v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.1 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gotd/ige v0.2.2 // indirect
	github.com/gotd/neo v0.1.5 // indirect
	github.com/gotd/td v0.89.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/herumi/bls-eth-go-binary v1.29.1 // indirect
	github.com/holiman/uint256 v1.2.2-0.20230321075855-87b91420868c // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mrjones/oauth v0.0.0-20190623134757-126b35219450 // indirect
	github.com/onsi/ginkgo/v2 v2.8.1 // indirect
	github.com/onsi/gomega v1.27.1 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/status-im/keycard-go v0.2.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220614013038-64ee5596c38a // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/urfave/cli/v2 v2.17.2-0.20221006022127-8f469abc00aa // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/wealdtech/go-bytesutil v1.2.1 // indirect
	github.com/wealdtech/go-eth2-types/v2 v2.8.1 // indirect
	github.com/xelaj/errs v0.0.0-20200831133608-d1c11863e019 // indirect
	github.com/xelaj/go-dry v0.0.0-20201104165138-61a25872c05a // indirect
	github.com/xelaj/mtproto v1.0.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	github.com/zelenin/go-tdlib v0.7.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel v1.19.0 // indirect
	go.opentelemetry.io/otel/trace v1.19.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	gonum.org/v1/gonum v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230920204549-e6e6cdab5c13 // indirect
	google.golang.org/grpc v1.57.0 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gorm.io/driver/sqlite v1.4.4 // indirect
	gorm.io/gorm v1.24.6 // indirect
	k8s.io/apiextensions-apiserver v0.26.1 // indirect
	rsc.io/qr v0.2.0 // indirect
	sigs.k8s.io/controller-runtime v0.14.4 // indirect
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/chai2010/gettext-go v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/fvbommel/sortorder v1.0.1 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.3.1
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.0
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/term v0.0.0-20220808134915-39b0c02b01ae // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/r3labs/sse/v2 v2.8.1
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/cobra v1.7.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	go.starlark.net v0.0.0-20221028183056-acb66ad56dd2 // indirect
	golang.org/x/crypto v0.15.0
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/oauth2 v0.14.0
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/term v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/component-base v0.26.1 // indirect
	k8s.io/klog/v2 v2.90.0 // indirect
	k8s.io/kube-openapi v0.0.0-20230202010329-39b3636cbaa3 // indirect
	k8s.io/utils v0.0.0-20230202215443-34013725500c // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/kustomize/api v0.12.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.9 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
