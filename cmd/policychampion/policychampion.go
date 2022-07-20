package main

import (
	"flag"
	util "tan/policychampion/internal"
)

func main() {
	// command-line flag parsing
	inputPath := flag.String("i", "", "set path of input YAML file")
	outputPath := flag.String("o", "updated"+*inputPath, "set path of output file")
	regoNamespace := flag.String("namespace", "main", "set namespace of Rego package")
	policyPath := flag.String("p", "policies", "set path of policy folder")
	flag.Parse()

	rawFile := util.ReadFile(*inputPath)
	resultSet := util.EvalPolicy(rawFile, *regoNamespace, *policyPath)
	updatedFile := util.ExecuteProhibitedRule(rawFile, resultSet)
	updatedFile = util.ExecuteOptionalRule(updatedFile, resultSet)
	updatedFile = util.ExecuteMandatoryRule(updatedFile, resultSet)

	util.WriteFile(*outputPath, updatedFile)
}
