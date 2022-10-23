package transform

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
	log "github.com/sirupsen/logrus"
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

func TargetRef(jsonStr []byte, rdfFilePath string) {
	jqQuery := `select((.Object.kind | test("EndpointsList")) and (.items != null) and (.items[0] != null))
				| [.items[] 
					| .kind as $p1 
					| (["_:", $p1,"-", .metadata.name]|add) as $s 
					| .subsets[]? 
					| .addresses[]? 
					| select((.targetRef != null) and (.targetRef.kind != null) and (.targetRef.name != null)) 
					| {"Subject": $s, "Predicate": ([$p1, "_", .targetRef.kind]|add) , "Object": (["_:", .targetRef.kind, "-", .targetRef.name]|add ) } ]`
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
}

func NodeProp(jsonStr []byte, rdfFilePath string) {
	jqQuery := `[.items[] 
					| [paths(scalars) as $path 
					| {"key": $path 
					| join("_"), "value": getpath($path)}] 
					| from_entries] 
					| [.[] 
					| (["_:", .kind, "-", .metadata_name]|add) as $s 
					| to_entries[] 
					| .key as $p 
					| (.value|tostring) as $o 
					| {"Subject":$s, "Predicate":$p, "Object":$o} ]`
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
}

func OwnRef(jsonStr []byte, rdfFilePath string) {
	jqQuery := `[.items[]? 
					| .kind as $p1 
					| (["_:", $p1, "-", .metadata.name]|add) as $src 
					| [.metadata.ownerReferences[]? 
					| .kind as $p2 
					| (["_:", $p2, "-" ,.name]|add) as $dst 
					| {"Subject":$src, "Predicate":([$p1,"_", $p2]|add), "Object":$dst}]] 
					| flatten`
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
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
