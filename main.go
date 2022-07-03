package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/open-policy-agent/opa/rego"
)

func main() {

	ctx := context.Background()

	// Create query that produces a single document.
	rego := rego.New(
		rego.Query("data.example.optional"),
		rego.Module("example.rego",
			`package example

            sites := [
                {"name": "prod"},
                {"name": "smoke1"},
                {"name": "dev"}
            ]
            
            q[name] { name := [
                {"name": "prod"},
                {"name": "smoke1"},
                {"name": "dev"}
            ] }

            optional[{"key":key,"msg":msg}]{
            
                key :="example_key"
                msg :="example message"
            }
            `,
		))

	// Run evaluation.
	rs, err := rego.Eval(ctx)
    result:=rs[0].Expressions[0].Value.([]interface {})[0].(map[string]interface {})
	fmt.Println(rs)
	// fmt.Println(len(rs))
    fmt.Println(reflect.TypeOf(rs[0]))
    fmt.Println(rs[0].Expressions[0])
    fmt.Println(rs[0].Expressions[0].Text)
    fmt.Println(rs[0].Expressions[0].Location)
    fmt.Println(reflect.TypeOf(rs[0].Expressions[0]))
    // fmt.Println(rs[0].Expressions[0].Value.([]interface {})[0].(map[string]interface {})["key"])
    fmt.Println(result["key"])
    fmt.Println(result["msg"])

	// Inspect result.
	// fmt.Println("value:", rs[0].Bindings)
	fmt.Println("err:", err)

}
