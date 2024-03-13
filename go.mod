module github.com/pivotal-cf/winfs-injector

go 1.21

toolchain go1.21.5

require (
	code.cloudfoundry.org/hydrator v0.0.0-20240112153350-a6bf3b70363f
	github.com/cloudfoundry/bosh-cli/v7 v7.5.2
	github.com/cloudfoundry/bosh-utils v0.0.426
	github.com/jhoonb/archivex v0.0.0-20201016144719-6a343cdae81d
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.31.1
	github.com/pivotal-cf/jhanda v0.0.0-20200619200912-8de8eb943a43
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go v0.112.0 // indirect
	cloud.google.com/go/compute v1.23.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.5 // indirect
	cloud.google.com/go/storage v1.36.0 // indirect
	code.cloudfoundry.org/clock v1.1.0 // indirect
	code.cloudfoundry.org/tlsconfig v0.0.0-20240116140718-a2c58c2ff70c // indirect
	code.cloudfoundry.org/workpool v0.0.0-20230612151832-b93da105e0e8 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/aws/aws-sdk-go v1.50.1 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/charlievieth/fs v0.0.3 // indirect
	github.com/cheggaaa/pb/v3 v3.1.4 // indirect
	github.com/cloudfoundry/bosh-agent v0.0.61-0.20230301011448-4cfe06c13ad7 // indirect
	github.com/cloudfoundry/bosh-davcli v0.0.306 // indirect
	github.com/cloudfoundry/bosh-gcscli v0.0.189 // indirect
	github.com/cloudfoundry/bosh-s3cli v0.0.253 // indirect
	github.com/cloudfoundry/config-server v0.1.165 // indirect
	github.com/cloudfoundry/go-socks5 v0.0.0-20180221174514-54f73bdb8a8e // indirect
	github.com/cloudfoundry/socks5-proxy v0.2.108 // indirect
	github.com/cppforlife/go-patch v0.2.0 // indirect
	github.com/cppforlife/go-semi-semantic v0.0.0-20160921010311-576b6af77ae4 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/frankban/quicktest v1.14.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gopacket/gopacket v1.2.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jessevdk/go-flags v1.5.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc6 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pivotal-cf/paraphernalia v0.0.0-20180203224945-a64ae2051c20 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/vito/go-interact v1.0.1 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.47.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.47.0 // indirect
	go.opentelemetry.io/otel v1.22.0 // indirect
	go.opentelemetry.io/otel/metric v1.22.0 // indirect
	go.opentelemetry.io/otel/trace v1.22.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/oauth2 v0.16.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/term v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/api v0.157.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20240122161410-6c6643bf1457 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240122161410-6c6643bf1457 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240122161410-6c6643bf1457 // indirect
	google.golang.org/grpc v1.60.1 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
