package main

import (
	"flag"
	"path/filepath"

	"github.com/AbdullahWasTaken/kube-miner/collector"
	"github.com/AbdullahWasTaken/kube-miner/config"
	"github.com/AbdullahWasTaken/kube-miner/utils"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	outputPath := flag.String("outputPath", "State", "relative path to where the collected data will be stored")
	flag.Parse()
	var c *config.Config = config.BuildConfig(kubeconfig, *outputPath)
	st := collector.GetState(c.DisClient, c.DynClient)
	utils.RDF(st, c.OutputPath)
}
