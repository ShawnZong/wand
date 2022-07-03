package main

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

func main() {
	module := `
    package example.authz
    
    import future.keywords
    
    default allow := false
    
    allow {
        is_admin
    }

    optional[{"key":key,"msg":msg}]{
        is_admin
        
        key :="example_key"
        msg :="example message"
    }

    optional[{"key":key,"msg":msg}]{
        is_admin
        
        key :="example_key2"
        msg :="example message2"
    }

    is_admin {
        "admin" in input.subject.groups
    }
    


    `
	ctx := context.Background()
	query, err := rego.New(
		rego.Query("data.example.authz.optional"),
		rego.Module("example.rego", module),
	).PrepareForEval(ctx)

	if err != nil {
		// Handle error.
	}

	input := map[string]interface{}{
		"method": "GET",
		"path":   []interface{}{"salary", "bob"},
		"subject": map[string]interface{}{
			"user":   "bob",
			"groups": []interface{}{"sales", "marketing","admin"},
		},
	}

	rs, err := query.Eval(ctx, rego.EvalInput(input))
	// result := rs[0].Expressions[0].Value.(map[string]interface{})["optional"].([]interface {})[0].(map[string]interface{})
	for _,value := range rs[0].Expressions[0].Value.([]interface{}){
        result:=value.(map[string]interface{})
        fmt.Println(result["key"])
    }
    // result := rs[0].Expressions[0].Value.([]interface{})[1].(map[string]interface{})
	// fmt.Println(rs)
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
