package ssainspect

import (
	"errors"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

var Analyzer = &analysis.Analyzer{
	Name: "ssainspect",
	Doc:  "make an Inspector",
	Run:  runAnalyzer,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
	ResultType: reflect.TypeFor[*Inspector](),
}

func runAnalyzer(pass *analysis.Pass) (any, error) {
	ssa, ok := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	if !ok {
		return nil, errors.New("failed to get result of buildssa.Analyzer")
	}
	return New(ssa.SrcFuncs), nil
}
