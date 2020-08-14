module github.com/gardener/machine-controller-manager

go 1.14

require (
	github.com/Azure/azure-sdk-for-go v32.6.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.3
	github.com/Azure/go-autorest/autorest/adal v0.8.0
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v0.0.0-20180828111155-cad214d7d71f
	github.com/aws/aws-sdk-go v1.13.54
	github.com/davecgh/go-spew v1.1.1
	github.com/go-ini/ini v1.36.0 // indirect
	github.com/go-openapi/spec v0.19.8
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gophercloud/gophercloud v0.7.0
	github.com/gophercloud/utils v0.0.0-20190527093828-25f1b77b8c03
	github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8 // indirect
	github.com/metal-stack/metal-go v0.8.3
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.7.0
	github.com/packethost/packngo v0.0.0-20181217122008-b3b45f1b4979
	github.com/prometheus/client_golang v1.7.1
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/api v0.13.0
	k8s.io/api v0.17.11
	k8s.io/apimachinery v0.17.11
	k8s.io/apiserver v0.17.11
	k8s.io/client-go v0.17.11
	k8s.io/cluster-bootstrap v0.17.11
	k8s.io/code-generator v0.17.11
	k8s.io/component-base v0.17.11
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20200731180307-f00132d28269
)

replace (
	github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega => github.com/onsi/gomega v1.5.0
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.3
	google.golang.org/grpc => google.golang.org/grpc v1.25.0
	k8s.io/api => k8s.io/api v0.17.11
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.11
	k8s.io/apiserver => k8s.io/apiserver v0.17.11
	k8s.io/client-go => k8s.io/client-go v0.17.11
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.11
	k8s.io/code-generator => k8s.io/code-generator v0.17.11
)
