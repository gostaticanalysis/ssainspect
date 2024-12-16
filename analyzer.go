package ssainspect

import (
	"errors"
	"iter"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

var Analyzer = &analysis.Analyzer{
	Name: "ssainspect",
	Doc:  "make iter.Seq[*ssainspect.Cursor]",
	Run:  runAnalyzer,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
	ResultType: reflect.TypeFor[iter.Seq[*Cursor]](),
}

func runAnalyzer(pass *analysis.Pass) (any, error) {
	ssa, ok := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	if !ok {
		return nil, errors.New("failed to get result of buildssa.Analyzer")
	}
	return All(ssa.SrcFuncs), nil
}
