package config

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	kubernetesConfig *rest.Config
	DynClient        *dynamic.Interface
	DisClient        *discovery.DiscoveryClient
	Out              string
}

func BuildConfig(kubeconfig *string, outputPath string) *Config {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	// use: collector
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Error(err)
	}
	// use: collector
	typedClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err)
	}
	return &Config{kubernetesConfig: config, DynClient: &dynamicClient, DisClient: typedClient.DiscoveryClient, Out: outputPath}
}
