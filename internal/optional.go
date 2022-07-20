package util

import (
	"log"

	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"
)

// given raw byte data of a YAML, decision results returned by OPA
// append comments to YAML nodes
// return raw byte data of a updated YAML
func ExecuteOptionalRule(rawFile *[]byte, queryResult rego.ResultSet) *[]byte {
	// extract optional result set
	hints := ExtractRuleResult(queryResult, "optional")

	var yamlNode yaml.Node
	if err := yaml.Unmarshal(*rawFile, &yamlNode); err != nil {
		log.Fatal(err)
	}

	// func: add optional messages to a YAML Node as comment
	appendHint := func(node *yaml.Node, key string, msg string) {
		elements := FindElements(&yamlNode, key)
		for _, element := range elements {
			element.HeadComment = AppendComment(element.HeadComment, msg)
		}
	}

	// for loop hints and then add optional messages to each Node
	for _, hint := range hints {
		for key, msg := range hint.(map[string]interface{}) {
			appendHint(&yamlNode, key, msg.(string))
		}
	}
	updatedConfiguration, err := yaml.Marshal(&yamlNode)
	if err != nil {
		log.Fatal(err)
	}
	return &updatedConfiguration
}
