package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"

	"github.com/AbdullahWasTaken/kube-miner/collector"
	"github.com/AbdullahWasTaken/kube-miner/config"
	"github.com/AbdullahWasTaken/kube-miner/utils"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	outputPath := flag.String("outputPath", "state", "relative path to where the collected data will be stored")
	flag.Parse()
	var c *config.Config = config.BuildConfig(kubeconfig, *outputPath)
	st := collector.GetState(c.DisClient, c.DynClient)
	b, err := json.MarshalIndent(st, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(c.OutputPath, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	f, err := os.Create(filepath.Join(c.OutputPath, "state.json"))
	if err != nil {
		log.Error(err)
	}
	_, err = f.Write(b)
	if err != nil {
		log.Error(err)
	}
	utils.RDF(st, c.OutputPath)
}
