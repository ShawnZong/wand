package util

import (
	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"
)

// given a YAML Node tree and a rego result
// remove prohibited YAML nodes and add messages to a YAML Node as comment
func removeNode(node *yaml.Node, hint map[string]interface{}) {
	for key, msg := range hint {
		elements := FindElements(node, key)
		for _, element := range elements {

			var emptyElement yaml.Node
			emptyElement.Kind = yaml.ScalarNode
			emptyElement.Value = "null"
			emptyElement.HeadComment = AppendComment(element.HeadComment, msg.(string))
			*element = emptyElement
		}
	}
}

// given raw byte data of a YAML, decision results returned by OPA
// remove prohibited YAML nodes and append comments to YAML nodes
// return raw byte data of a updated YAML
func ExecuteProhibitedRule(rawFile *[]byte, queryResult rego.ResultSet) *[]byte {
	// extract prohibited result set
	hints := ExtractRuleResult(queryResult, "prohibited")
	return ExecuteRule(rawFile, hints, removeNode)
}
