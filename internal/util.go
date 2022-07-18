package util

import (
	"io/ioutil"
	"log"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/loader"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

func ReadFile(path string) *[]byte {
	rawFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error when loading configuration file %v: %v", path, err)
	}
	return &rawFile
}

func WriteFile(filename string, data *[]byte) {
	ioutil.WriteFile(filename, *data, 0644)
}

func GetCompiler(policyPath string) *ast.Compiler {
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

func ParseConfiguration(rawFile *[]byte) map[string]interface{} {

	parsedFile := make(map[string]interface{})

	if err := yaml.Unmarshal(*rawFile, &parsedFile); err != nil {
		log.Fatalf("Error when parsing configuration file: %v", err)
	}

	return parsedFile
}

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
