module github.com/SynologyOpenSource/synology-csi

go 1.16

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/container-storage-interface/spec v1.8.0
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/kubernetes-csi/csi-lib-utils v0.13.0
	github.com/kubernetes-csi/csi-test/v4 v4.3.0
	github.com/sirupsen/logrus v1.9.2
	github.com/spf13/cobra v1.7.0
	golang.org/x/net v0.10.0 // indirect
	google.golang.org/genproto v0.0.0-20230524185152-1884fd1fac28 // indirect
	google.golang.org/grpc v1.55.0
	google.golang.org/protobuf v1.30.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/mount-utils v0.27.2
	k8s.io/utils v0.0.0-20230505201702-9f6742963106
)

exclude google.golang.org/grpc v1.37.0
