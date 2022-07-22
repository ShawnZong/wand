package util

import (
	"log"

	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"
)

// given a YAML Node tree and a rego result
// add template YAML to Node tree and append meesage as comment
func appendMsgTemplate(node *yaml.Node, hint map[string]interface{}) {
	// extract ket, msg and ref from a rego result
	var key, msg, ref string
	for mapKey, mapValue := range hint {
		if mapKey == "templateRef" {
			ref = mapValue.(string)
		} else {
			key = mapKey
			msg = mapValue.(string)
		}
	}

	// load reference template, because the template might be different
	// we need to load the template separately each time
	var refYAML yaml.Node
	if err := yaml.Unmarshal(*ReadFile(ref), &refYAML); err != nil {
		log.Fatalf("load ref file error: %v", err)
	}

	// search for specific YAML Nodes according to key
	elements := FindElements(node, key)
	for _, element := range elements {
		// if the YAML Node is empty, replace empty Node directly with template Node
		if len(element.Content) == 0 {
			*element = *refYAML.Content[0]
			element.HeadComment = AppendComment(element.HeadComment, msg)
		} else if element.Content[0].Kind == yaml.ScalarNode {
			// if the parent node is not a Mapping Type, we can't simply append reference Node
			// we need to extract the content of reference node
			refYAML.Content[0].Content[0].HeadComment = AppendComment(refYAML.Content[0].Content[0].HeadComment, msg)
			element.Content = append(element.Content, refYAML.Content[0].Content...)
		} else {
			// if the YAML Node is not empty, append template Node to existing Nodes
			refYAML.Content[0].HeadComment = AppendComment(refYAML.Content[0].HeadComment, msg)
			element.Content = append(element.Content, refYAML.Content...)
		}
	}
}

// given raw byte data of a YAML, decision results returned by OPA
// append YAML template under specified key and add comments to YAML nodes
// return raw byte data of a updated YAML
func ExecuteMandatoryRule(rawFile *[]byte, queryResult rego.ResultSet) *[]byte {
	// extract mandatory result set
	hints := ExtractRuleResult(queryResult, "mandatory")

	return ExecuteRule(rawFile, hints, appendMsgTemplate)
}
