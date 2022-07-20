package util

import (
	"log"

	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"
)

// given raw byte data of a YAML, decision results returned by OPA
// remove prohibited YAML nodes and append comments to YAML nodes
// return raw byte data of a updated YAML
func ExecuteProhibitedRule(rawFile *[]byte, queryResult rego.ResultSet) *[]byte {
	// extract prohibited result set
	hints := ExtractRuleResult(queryResult, "prohibited")

	var yamlNode yaml.Node
	if err := yaml.Unmarshal(*rawFile, &yamlNode); err != nil {
		log.Fatal(err)
	}

	// func: remove prohibited YAML nodes and add messages to a YAML Node as comment
	removeNode := func(node *yaml.Node, key string, msg string) {
		elements := FindElements(&yamlNode, key)
		for _, element := range elements {

			var emptyElement yaml.Node
			emptyElement.Kind = yaml.ScalarNode
			emptyElement.Value = "null"
			emptyElement.HeadComment = AppendComment(element.HeadComment, msg)
			*element = emptyElement
		}

	}

	//loop each result and then execute prohibited rules
	for _, hint := range hints {
		for key, msg := range hint.(map[string]interface{}) {
			removeNode(&yamlNode, key, msg.(string))
		}
	}
	updatedConfiguration, err := yaml.Marshal(&yamlNode)
	if err != nil {
		log.Fatal(err)
	}
	return &updatedConfiguration
}
