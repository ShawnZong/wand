package main

import (
	util "tan/policychampion/internal"
)

func main() {
	filename := "example.yaml"
	rawFile := util.ReadFile(filename)
	resultSet := util.EvalPolicy(rawFile)
	updatedFile := util.ExecuteProhibitedRule(rawFile, resultSet)
	updatedFile = util.ExecuteOptionalRule(updatedFile, resultSet)

	util.WriteFile("updated_"+filename, updatedFile)
}
