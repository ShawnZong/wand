package util

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/loader"
	"github.com/open-policy-agent/opa/rego"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

// read a file
// return raw byte data
func ReadFile(path string) *[]byte {
	rawFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error when loading configuration file %v: %v", path, err)
	}
	return &rawFile
}

// given filename and raw byte data of a file
// create a new file
func WriteFile(filename string, data *[]byte) {
	ioutil.WriteFile(filename, *data, 0644)
}

// load policy files and feed them to a rego compiler
// return a rego compiler loaded with polcy files
func getCompiler(policyPath string) *ast.Compiler {
	path := []string{policyPath}
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

func NewRegoObject() (*rego.PreparedEvalQuery, context.Context) {
	ctx := context.Background()
	compiler := getCompiler("../../policies")
	query, err := rego.New(
		rego.Query("data.main"),
		rego.Compiler(compiler),
	).PrepareForEval(ctx)

	if err != nil {
		log.Fatalf("cannot create new rego object: %v", err)
	}
	return &query, ctx
}

// parse raw byte data of a file to Golang variables
// return sturctured data
func ParseConfiguration(rawFile *[]byte) map[string]interface{} {

	parsedFile := make(map[string]interface{})

	if err := yaml.Unmarshal(*rawFile, &parsedFile); err != nil {
		log.Fatalf("Error when parsing configuration file: %v", err)
	}

	return parsedFile
}

// find corresponding YAML Nodes based on a JSONPath query
// return found Nodes in a YAML Node tree
func FindElements(yamlNode *yaml.Node, path string) []*yaml.Node {
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

func EvalPolicy(rawFile *[]byte) rego.ResultSet {

	yamlfile := ParseConfiguration(rawFile)

	query, ctx := NewRegoObject()
	resultSet, err := query.Eval(ctx, rego.EvalInput(yamlfile))
	if err != nil {
		log.Fatalf("error when evaluating Rego query: %v", err)
	}
	return resultSet
}
