package transform

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

type triple struct {
	Subject   string
	Predicate string
	Object    string
}

func (r *triple) String() string {
	if strings.HasPrefix(r.Object, "_:") {
		return r.Subject + " <" + r.Predicate + "> " + r.Object + " .\n"
	}
	return r.Subject + " <" + r.Predicate + "> " + strconv.Quote(r.Object) + " .\n"
}

type triples []triple

func (tpl *triples) String() string {
	var sb strings.Builder
	for _, v := range *tpl {
		sb.WriteString(v.String())
	}
	return sb.String()
}

func saveRdf(jqQuery string, jsonStr []byte, rdfFilepath string) error {
	q, err := gojq.Parse(jqQuery)
	if err != nil {
		return err
	}
	var input map[string]interface{}
	err = json.Unmarshal(jsonStr, &input)
	if err != nil {
		return err
	}
	rdf, err := generateRDF(q, input)
	if err != nil {
		return err
	}
	if len(rdf) > 0 {
		// write RDF to file
		err = os.WriteFile(rdfFilepath, []byte(rdf.String()), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateRDF(query *gojq.Query, input map[string]interface{}) (triples, error) {
	iter := query.Run(input)
	v, ok := iter.Next()
	rdf := triples{}
	for ok {
		if err, ok := v.(error); ok {
			return nil, err
		}
		b, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			return nil, err
		}
		var ll triples
		err = json.Unmarshal(b, &ll)
		if err != nil {
			return nil, err
		}
		rdf = append(rdf, ll...)
		v, ok = iter.Next()
	}
	return rdf, nil
}
