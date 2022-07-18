package util

import (
	"fmt"
	"log"

	"github.com/open-policy-agent/opa/rego"
	"gopkg.in/yaml.v3"
)
func ExtractOptional(queryResult rego.ResultSet) []interface{} {
	return queryResult[0].Expressions[0].Value.(map[string]interface{})["optional"].([]interface{})
}


func AppendOptional2Configuration(rawFile *[]byte, hints []interface{}) *[]byte {
	var yamlNode yaml.Node
	if err := yaml.Unmarshal(*rawFile, &yamlNode); err != nil {
		log.Fatal(err)
	}
	// TODO add optional messages to yamlNode
	// for loop hints and then add optional messages to each Node
	appendHint := func(node *yaml.Node, key string, msg string) {
		elements := FindElements(&yamlNode, key)
		for _, element := range elements {
			element.LineComment = element.LineComment + msg
		}

	}
	for _, hint := range hints {
		for key, msg := range hint.(map[string]interface{}) {
			appendHint(&yamlNode, key, msg.(string))
			fmt.Println(key)
			fmt.Println(msg)
			//TODO find the key in yaml node tree and then append hint to head comment
		}
	}
	updatedConfiguration, err := yaml.Marshal(&yamlNode)
	if err != nil {
		log.Fatal(err)
	}
	return &updatedConfiguration
}