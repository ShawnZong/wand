package util

import (
	"log"

	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"
)

// given OPA query result set, extracts the set of mandatory rules
// return extracted mandatory result set
func ExtractMandatory(queryResult rego.ResultSet) []interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Fail execute mandatory rules: %v", err)
		}
	}()
	results := queryResult[0].Expressions[0].Value.(map[string]interface{})["mandatory"]
	if results == nil {
		log.Println("no results of mandatory rules")
		return []interface{}{}
	}
	return results.([]interface{})
}

// given raw byte data of a YAML, decision results returned by OPA
// append YAML template under specified key and add comments to YAML nodes
// return raw byte data of a updated YAML
func ExecuteMandatoryRule(rawFile *[]byte, queryResult rego.ResultSet) *[]byte {
	// extract mandatory result set
	hints := ExtractMandatory(queryResult)

	var yamlNode yaml.Node
	if err := yaml.Unmarshal(*rawFile, &yamlNode); err != nil {
		log.Fatal(err)
	}

	// func: add optional messages to a YAML Node as comment
	appendMsgTemplate := func(node *yaml.Node, key string, msg string, ref string) {
		var refYAML yaml.Node
		if err := yaml.Unmarshal(*ReadFile("../../" + ref), &refYAML); err != nil {
			log.Fatalf("load ref file error: %v", err)
		}
		elements := FindElements(&yamlNode, key)
		for _, element := range elements {
			// if the YAML Node is empty, replace empty Node directly with template Node
			if len(element.Content) == 0 {
				*element = *refYAML.Content[0]
				element.HeadComment = element.HeadComment + msg
			} else {
				refYAML.Content[0].HeadComment = refYAML.Content[0].HeadComment + msg
				element.Content = append(element.Content, refYAML.Content...)
			}
		}
	}

	// loop each mandatory results and process corresponding nodes
	for _, hint := range hints {
		var key, msg, ref string
		for mapKey, mapValue := range hint.(map[string]interface{}) {
			if mapKey == "templateRef" {
				ref = mapValue.(string)
			} else {
				key = mapKey
				msg = mapValue.(string)
			}
		}
		appendMsgTemplate(&yamlNode, key, msg, ref)
	}
	updatedConfiguration, err := yaml.Marshal(&yamlNode)
	if err != nil {
		log.Fatalf("marshal yaml error: %v", err)
	}
	return &updatedConfiguration
}
