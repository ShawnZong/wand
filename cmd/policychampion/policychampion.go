package main

import (
	util "tan/policychampion/internal"
)

func main() {
	filename := "example.yaml"
	rawFile := util.ReadFile(filename)
	resultSet := util.EvalPolicy(rawFile)
	updatedFile := util.ExucutePrihibitedRule(rawFile, resultSet)
	updatedFile = util.AppendOptional2Configuration(updatedFile, resultSet)

	util.WriteFile("updated_"+filename, updatedFile)
}
