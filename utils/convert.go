package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/AbdullahWasTaken/kube-miner/collector"
	log "github.com/sirupsen/logrus"
)

func RDF(st collector.ClusterState, path string) {
	uidMap := map[interface{}]string{}
	for _, v := range st {
		for _, d := range v.Items {
			uidMap[d.GetName()] = "_:" + string(d.GetUID())
			uidMap[string(d.GetUID())] = "_:" + string(d.GetUID())
		}
	}
	delete(uidMap, "")
	// creating the output path if it doesn't exists.
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range st {
		f, err := os.Create(filepath.Join(path, fmt.Sprintf("%v.rdf", k)))
		if err != nil {
			log.Error(err)
		}
		for _, res := range v.Items {
			xid := uidMap[res.GetName()]
			flat_res := Flatten(res.Object)
			for k, v := range flat_res {
				if reflect.ValueOf(v).Kind() == reflect.String {
					_, err = f.WriteString(fmt.Sprintf("%v <%v> %q .\n", xid, k, v))
				} else {
					_, err = f.WriteString(fmt.Sprintf("%v <%v> \"%v\" .\n", xid, k, v))
				}
				if err != nil {
					log.Error(err)
				}
				if u, prs := uidMap[v]; prs && (u != xid) {
					_, err = f.WriteString(fmt.Sprintf("\n\n%v <edge> %v(%v) .\n\n", xid, u, v))
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
		f.Sync()
	}
}
