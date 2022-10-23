package transform

import (
	log "github.com/sirupsen/logrus"
)

const pod_svcAccount = `select((.Object.kind | test("PodList")) and (.items != null) and (.items[0] != null))
| [.items[] 
    | .kind as $p1 
    | (["_:", $p1, "-", .metadata.name]|add) as $s 
    | select(.spec.serviceAccount != null) 
    | .spec
    | {"Subject": $s, "Predicate": ([$p1, "_ServiceAccount"]|add) , "Object": (["_:ServiceAccount-", .serviceAccount]|add ) } ]`

const rolebinding_role = `select((.Object.kind | test("[RoleBindingList| ClusterRoleBindingList]")) and (.items != null) and (.items[0] != null))
| [.items[] 
    | .kind as $p1 
    | (["_:", $p1, "-", .metadata.name]|add) as $s 
    | select(.roleRef != null)
    | .roleRef
    | {"Subject": $s, "Predicate": ([$p1, "_", .kind]|add) , "Object": (["_:", .kind, "-",.name]|add ) } ]`

const binding_svcAcc = `select((.Object.kind | test("[RoleBindingList| ClusterRoleBindingList]")) and (.items != null) and (.items[0] != null))
| [.items[] 
    | .kind as $p1 
    | (["_:", $p1, "-", .metadata.name]|add) as $s 
    | select(.subjects != null)
    | .subjects[]
    | {"Subject": $s, "Predicate": ([$p1, "_", .kind]|add) , "Object": (["_:", .kind, "-",.name]|add ) } ]`

func RBAC(jsonStr []byte, rdfFilePath string) {
	jqQuery := pod_svcAccount
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
	jqQuery = rolebinding_role
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
	jqQuery = binding_svcAcc
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
}

const endpoint_targetRef = `select((.Object.kind | test("EndpointsList")) and (.items != null) and (.items[0] != null))
| [.items[] 
	| .kind as $p1 
	| (["_:", $p1,"-", .metadata.name]|add) as $s 
	| .subsets[]? 
	| .addresses[]? 
	| select((.targetRef != null) and (.targetRef.kind != null) and (.targetRef.name != null)) 
	| {"Subject": $s, "Predicate": ([$p1, "_", .targetRef.kind]|add) , "Object": (["_:", .targetRef.kind, "-", .targetRef.name]|add ) } ]`

const endpointslice_targetRef = `select((.Object.kind | test("EndpointSliceList")) and (.items != null) and (.items[0] != null))
| [.items[] 
	| .kind as $p1 
	| (["_:", $p1, "-", .metadata.name]|add) as $s 
	| .endpoints[]?
	| select((.targetRef != null) and (.targetRef.kind != null) and (.targetRef.name != null)) 
	| {"Subject": $s, "Predicate": ([$p1, "_", .targetRef.kind]|add) , "Object": (["_:", .targetRef.kind, "-", .targetRef.name]|add ) } ]`

func TargetRef(jsonStr []byte, rdfFilePath string) {
	jqQuery := endpoint_targetRef
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
	jqQuery = endpointslice_targetRef
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
}

const literal = `[.items[] 
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

func NodeProp(jsonStr []byte, rdfFilePath string) {
	jqQuery := literal
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
}

const ownerRef = `[.items[]? 
| .kind as $p1 
| (["_:", $p1, "-", .metadata.name]|add) as $src 
| [.metadata.ownerReferences[]? 
| .kind as $p2 
| (["_:", $p2, "-" ,.name]|add) as $dst 
| {"Subject":$src, "Predicate":([$p1,"_", $p2]|add), "Object":$dst}]] 
| flatten`

func OwnRef(jsonStr []byte, rdfFilePath string) {
	jqQuery := ownerRef
	if err := saveRdf(jqQuery, jsonStr, rdfFilePath); err != nil {
		log.Error(err)
	}
}
