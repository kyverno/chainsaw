package main

import (
	"fmt"
	"strings"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	chainsawfunctions "github.com/kyverno/chainsaw/pkg/runner/check/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template/functions"
	kyvernofunctions "github.com/kyverno/kyverno-json/pkg/engine/template/kyverno"
)

func main() {
	fmt.Println("# Functions")
	fmt.Println()
	fmt.Println("## built-in functions")
	fmt.Println()
	printFunctions(jpfunctions.GetDefaultFunctions()...)
	fmt.Println()
	fmt.Println("## kyverno-json functions")
	fmt.Println()
	printFunctions(functions.GetFunctions()...)
	fmt.Println()
	fmt.Println("## kyverno functions")
	fmt.Println()
	printFunctions(kyvernofunctions.GetBareFunctions()...)
	fmt.Println()
	fmt.Println("## chainsaw functions")
	fmt.Println()
	printFunctions(chainsawfunctions.GetFunctions(nil)...)
	fmt.Println()
}

func printFunctions(funcs ...jpfunctions.FunctionEntry) {
	fmt.Println("| Name | Signature |")
	fmt.Println("|---|---|")
	for _, function := range funcs {
		fmt.Println("|", function.Name, "|", "`"+strings.ReplaceAll(functionString(function), "|", `\|`)+"`", "|")
	}
}

func functionString(f jpfunctions.FunctionEntry) string {
	if f.Name == "" {
		return ""
	}
	var args []string
	for _, a := range f.Arguments {
		var aTypes []string
		for _, t := range a.Types {
			aTypes = append(aTypes, string(t))
		}
		args = append(args, strings.Join(aTypes, "|"))
	}
	output := fmt.Sprintf("%s(%s)", f.Name, strings.Join(args, ", "))
	return output
}
