package main

import (
	"fmt"
	"strings"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	chainsawfunctions "github.com/kyverno/chainsaw/pkg/engine/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template/functions"
	kyvernofunctions "github.com/kyverno/kyverno-json/pkg/engine/template/kyverno"
)

func main() {
	fmt.Println("# Functions")
	fmt.Println()
	fmt.Println(`!!! warning "Experimental functions"`)
	fmt.Println()
	fmt.Println("    Experimental functions are denoted by the `x_` prefix.")
	fmt.Println()
	fmt.Println("    These are functions that are subject to signature change in a future version.")
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
	{
		var functions []jpfunctions.FunctionEntry
		for _, function := range chainsawfunctions.GetFunctions() {
			functions = append(functions, function.FunctionEntry)
		}
		printFunctions(functions...)
	}
	fmt.Println()
	fmt.Println("## examples")
	fmt.Println()
	fmt.Println("- [x_k8s_get](./examples/x_k8s_get.md)")
	fmt.Println()
	{
		for _, function := range kyvernofunctions.GetFunctions() {
			printFunctionExamples(chainsawfunctions.FunctionEntry{
				FunctionEntry: function.FunctionEntry,
				Note:          function.Note,
			})
		}
	}
}

func printFunctions(funcs ...jpfunctions.FunctionEntry) {
	fmt.Println("| Name | Signature |")
	fmt.Println("|---|---|")
	for _, function := range funcs {
		fmt.Println("|", function.Name, "|", "`"+functionString(function)+"`", "|")
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

func printFunctionExamples(funcs ...chainsawfunctions.FunctionEntry) {
	for _, function := range funcs {
		if function.Note != "" {
			fmt.Println("###", function.Name)
			fmt.Println()
			fmt.Println(function.Note)
			fmt.Println()
		}
	}
}
