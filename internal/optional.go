package util

import (
	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"
)

// given a YAML Node tree and a rego result
// add optional messages to a YAML Node as comment
func appendHint(node *yaml.Node, hint map[string]interface{}) {
	for key, msg := range hint {
		elements := FindElements(node, key)
		for _, element := range elements {
			element.HeadComment = AppendComment(element.HeadComment, msg.(string))
		}
	}
}

// given raw byte data of a YAML, decision results returned by OPA
// append comments to YAML nodes
// return raw byte data of a updated YAML
func ExecuteOptionalRule(rawFile *[]byte, queryResult rego.ResultSet) *[]byte {
	// extract optional result set
	hints := ExtractRuleResult(queryResult, "optional")
	return ExecuteRule(rawFile, hints, appendHint)
}
