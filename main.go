package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/rego"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
)

func readFile(path string) *[]byte {
	rawFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error when loading configuration file %v: %v", path, err)
	}
	return &rawFile
}

func parseConfiguration(rawFile *[]byte) map[string]interface{} {

	parsedFile := make(map[string]interface{})

	if err := yaml.Unmarshal(*rawFile, &parsedFile); err != nil {
		log.Fatalf("Error when parsing configuration file: %v", err)
	}

	return parsedFile
}

func findElements(yamlNode *yaml.Node, path string) []*yaml.Node {
	pathQuery, err := yamlpath.NewPath(path)
	if err != nil {
		log.Fatalf("cannot create path query: %v", err)
	}
	elements, err := pathQuery.Find(yamlNode)
	if err != nil {
		log.Fatalf("error when finding elements: %v", err)
	}

	return elements
}

func appendOptional2Configuration(rawFile *[]byte, hints []interface{}) *[]byte {
	var yamlNode yaml.Node
	if err := yaml.Unmarshal(*rawFile, &yamlNode); err != nil {
		log.Fatal(err)
	}
	// TODO add optional messages to yamlNode
	// for loop hints and then add optional messages to each Node
	appendHint := func(node *yaml.Node, key string, msg string) {
		elements := findElements(&yamlNode, key)
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

func writeYAML(filename string, data *[]byte) {
	ioutil.WriteFile(filename, *data, 0644)
}

func getCompiler() *ast.Compiler {
	path := []string{"./policies"}
	policies, err := loader.AllRegos(path)
	if err != nil {
		log.Fatalf("Error during loading policies: %v \n", err)
		return nil
	} else if len(policies.Modules) == 0 {
		log.Fatalf("No policy files found in: %v \n", path)
		return nil
	}

	compiler := ast.NewCompiler().WithCapabilities(ast.CapabilitiesForThisVersion())
	compiler.Compile(policies.ParsedModules())
	return compiler
}

func extractOptional(queryResult rego.ResultSet) []interface{} {
	return queryResult[0].Expressions[0].Value.(map[string]interface{})["optional"].([]interface{})
}

func main() {
	filename := "example.yaml"
	rawFile := readFile(filename)
	yamlfile := parseConfiguration(rawFile)
	// fmt.Println(yamlfile)
	ctx := context.Background()

	// input := map[string]interface{}{
	// 	"method": "GET",
	// 	"path":   []interface{}{"salary", "bob"},
	// 	"subject": map[string]interface{}{
	// 		"user":   "bob",
	// 		"groups": []interface{}{"sales", "marketing", "admin"},
	// 	},
	// }

	// rs, err := rego.New(
	// 	rego.Query("data.main"),
	// 	rego.Module("module1", module),
	// 	rego.Module("module2", module2),
	// 	rego.Input(input),
	// ).PrepareForEval(ctx)

	// if err != nil {
	// 	// Handle error.
	// }

	// query, err := rego.New(
	// 	rego.Query("data.main"),
	// 	rego.Module("module1", module),
	// 	rego.Module("module2", module2),
	// 	rego.Input(input),
	// ).PrepareForEval(ctx)

	compiler := getCompiler()
	query, err := rego.New(
		rego.Query("data.main"),
		rego.Compiler(compiler),
	).PrepareForEval(ctx)

	rs, err := query.Eval(ctx, rego.EvalInput(yamlfile))
	// hints := rs[0].Expressions[0].Value.(map[string]interface{})["optional"].([]interface {})[0].(map[string]interface{})
	// for _, value := range rs[0].Expressions[0].Value.([]interface{}) {
	// 	hints := value.(map[string]interface{})
	// 	fmt.Println(hints["key"])
	// }
	// hints := rs[0].Expressions[0].Value.([]interface{})[1].(map[string]interface{})
	hints := extractOptional(rs)
	updatedFile := appendOptional2Configuration(rawFile, hints)
	writeYAML("copy_"+filename, updatedFile)

	fmt.Println(hints)
	fmt.Println(reflect.TypeOf(hints))
	for i, item := range hints {
		fmt.Println(i)
		fmt.Println(item)
		for key, value := range item.(map[string]interface{}) {
			fmt.Println(key)
			fmt.Println(value)
		}
	}
	// fmt.Println(rs)
	// fmt.Println(reflect.TypeOf(rs))

	// fmt.Println(result)
	// fmt.Println(len(rs))
	// fmt.Println(reflect.TypeOf(rs[0]))
	// fmt.Println(rs[0].Expressions[0])
	// fmt.Println(rs[0].Expressions[0].Text)
	// fmt.Println(rs[0].Expressions[0].Location)
	// fmt.Println(reflect.TypeOf(rs[0].Expressions[0]))
	// fmt.Println(rs[0].Expressions[0].Value.([]interface {})[0].(map[string]interface {})["key"])
	// fmt.Println(result["key"])
	// fmt.Println(result["msg"])

	// Inspect result.
	// fmt.Println("value:", rs[0].Bindings)
	fmt.Println("err:", err)

}
