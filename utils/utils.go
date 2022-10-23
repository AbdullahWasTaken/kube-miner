package utils

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Load all the resource states store in a json directory `jsonDir` and
// returns a map containing json strings of loaded resource states.
func LoadState(jsonDir string) map[string][]byte {
	cs := map[string][]byte{}
	files, err := os.ReadDir(jsonDir)
	if err != nil {
		log.Panic(err)
	}

	for _, f := range files {
		jsonStr, err := os.ReadFile(filepath.Join(jsonDir, f.Name()))
		if err != nil {
			log.Error(err)
		} else {
			// add loaded object to cluster state `cs`
			cs[strings.Split(f.Name(), ".")[0]] = jsonStr
		}
	}
	return cs
}
