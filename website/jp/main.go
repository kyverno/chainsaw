package main

import (
	"embed"
	"fmt"
	"os"
	"strings"

	jpfunctions "github.com/jmespath-community/go-jmespath/pkg/functions"
	chainsawfunctions "github.com/kyverno/chainsaw/pkg/engine/functions"
	"github.com/kyverno/kyverno-json/pkg/engine/template/functions"
	kyvernofunctions "github.com/kyverno/kyverno-json/pkg/engine/template/kyverno"
)

//go:embed examples
var examples embed.FS

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
	printFunctions(chainsawfunctions.GetFunctions()...)
	fmt.Println()
}

func printFunctions(funcs ...jpfunctions.FunctionEntry) {
	fmt.Println("| Name | Signature | Description |")
	fmt.Println("|---|---|---|")
	for _, function := range funcs {
		sig := functionString(function)
		fmt.Println("|", fmt.Sprintf("[%s](./examples/%s.md)", function.Name, function.Name), "|", "`"+sig+"`", "|", function.Description, "|")
		data := fmt.Sprintf("# %s\n\n## Signature\n\n`%s`\n\n## Description\n\n%s\n\n## Examples\n\n", function.Name, sig, function.Description)
		if e, err := examples.ReadFile(fmt.Sprintf("examples/%s.md", function.Name)); err != nil {
			panic(err)
		} else {
			data += string(e)
		}
		if err := os.WriteFile(fmt.Sprintf("./website/docs/reference/jp/examples/%s.md", function.Name), []byte(data), 0o600); err != nil {
			panic(err)
		}
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
