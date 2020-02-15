module sigs.k8s.io/cluster-api-provider-azure

go 1.13

require (
	cloud.google.com/go v0.53.0 // indirect
	github.com/Azure/azure-sdk-for-go v39.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.5
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/Azure/k8s-infra v0.0.1-alpha1.0.20200215002509-ddd7b0bb8cf8
	github.com/blang/semver v3.5.1+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/golang/mock v1.4.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/pelletier/go-toml v1.6.0
	github.com/pkg/errors v0.9.1
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/crypto v0.0.0-20200214034016-1d94cc7ab1c6
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/tools v0.0.0-20200214225126-5916a50871fb // indirect
	k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver v0.17.3 // indirect
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.3
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20200204173128-addea2498afe // indirect
	k8s.io/utils v0.0.0-20200124190032-861946025e34
	sigs.k8s.io/cluster-api v0.2.6-0.20200106222425-660e6b945a27
	sigs.k8s.io/controller-runtime v0.5.0
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.2.0+incompatible
