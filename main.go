package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/AbdullahWasTaken/kube-miner/collector"
	"github.com/AbdullahWasTaken/kube-miner/transform"
	"github.com/AbdullahWasTaken/kube-miner/utils"
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
		cs := utils.LoadState(*jsonPath)
		rdfPath := filepath.Join(*outDir, "RDF")
		err := os.MkdirAll(rdfPath, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
		for k, v := range cs {
			transform.OwnRef(v, filepath.Join(rdfPath, k+"_ownRef"+".rdf"))
			transform.NodeProp(v, filepath.Join(rdfPath, k+"_nodeProp"+".rdf"))
			transform.TargetRef(v, filepath.Join(rdfPath, k+"_tarRef"+".rdf"))
			transform.RBAC(v, filepath.Join(rdfPath, k+"_rbac"+".rdf"))
		}
	} else {
		log.Fatal("either kubeconfig or jsonPath must be set")
	}
}
