package main

import (
	"flag"
	"path/filepath"

	"github.com/AbdullahWasTaken/kube-miner/collector"
	log "github.com/sirupsen/logrus"
)

func main() {
	outDir := flag.String("out", "./out", "path to the output directory")
	kubeconfig := flag.String("kubeconfig", "", "path to the kubeconfig file")
	jsonPath := flag.String("json", "", "path to state JSON directory")
	flag.Parse()

	// collect fresh state data from running kubernetes cluster
	if *kubeconfig != "" {
		c := collector.NewCollector(*kubeconfig)
		c.Collect(*outDir)
		log.Info("Kubernetes cluster state saved to ", *outDir)
		// if everything was successfull, use the newly created json file in the next step
		*jsonPath = filepath.Join(*outDir, "JSON")
	} else if *jsonPath != "" {

	} else {
		log.Fatal("either kubeconfig or jsonPath must be set")
	}
}
