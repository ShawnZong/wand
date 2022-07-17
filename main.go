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
)

// func ExecuteOptional(modules []string,input map[string]interface{}){

//     for i:= range
// }
func readFile(path string) *[]byte {
	rawFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error when loading configuration file %v: %v", path, err)
	}
	return &rawFile
}

func parseConfiguration(path string) map[string]interface{} {
	rawFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error when loading configuration file %v: %v", path, err)
	}
	parsedFile := make(map[string]interface{})

	if err := yaml.Unmarshal(rawFile, &parsedFile); err != nil {
		log.Fatalf("Error when parsing configuration file %v: %v", path, err)
	}
	return parsedFile
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

	yamlfile := parseConfiguration("./example.yaml")
	fmt.Println(yamlfile["taskdefinition"].(map[string]interface{})["Properties"])
	ctx := context.Background()

	input := map[string]interface{}{
		"method": "GET",
		"path":   []interface{}{"salary", "bob"},
		"subject": map[string]interface{}{
			"user":   "bob",
			"groups": []interface{}{"sales", "marketing", "admin"},
		},
	}

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

	rs, err := query.Eval(ctx, rego.EvalInput(input))
	// result := rs[0].Expressions[0].Value.(map[string]interface{})["optional"].([]interface {})[0].(map[string]interface{})
	// for _, value := range rs[0].Expressions[0].Value.([]interface{}) {
	// 	result := value.(map[string]interface{})
	// 	fmt.Println(result["key"])
	// }
	// result := rs[0].Expressions[0].Value.([]interface{})[1].(map[string]interface{})
	result := extractOptional(rs)
	fmt.Println(result)
	fmt.Println(reflect.TypeOf(result))

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
