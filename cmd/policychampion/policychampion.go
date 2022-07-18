package main

import (
	"context"
	"fmt"
	"reflect"

	util "tan/policychampion/internal"

	"github.com/open-policy-agent/opa/rego"
)

func main() {
	filename := "example.yaml"
	rawFile := util.ReadFile(filename)
	yamlfile := util.ParseConfiguration(rawFile)
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

	compiler := util.GetCompiler("../../policies")
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
	hints := util.ExtractOptional(rs)
	updatedFile := util.AppendOptional2Configuration(rawFile, hints)
	util.WriteFile("updated_"+filename, updatedFile)

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
