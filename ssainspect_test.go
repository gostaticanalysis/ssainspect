package ssainspect_test

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/gostaticanalysis/ssainspect"
	"github.com/gostaticanalysis/testutil"
	"github.com/tenntenn/golden"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

var flagUpdate bool

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func TestInspectorWithAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)

	rs := analysistest.Run(t, testdata, testAnalyzer, "a")
	buf := rs[0].Result.(*bytes.Buffer)
	if flagUpdate {
		golden.Update(t, analysistest.TestData(), "a", buf)
		return
	}

	if diff := golden.Diff(t, analysistest.TestData(), "a", buf); diff != "" {
		t.Error(diff)
	}
}

var testAnalyzer = &analysis.Analyzer{
	Name: "testssainspect",
	Doc:  "test analyzer for ssainspect.Inspector",
	Run:  run,
	Requires: []*analysis.Analyzer{
		ssainspect.Analyzer,
	},
	ResultType: reflect.TypeOf((*bytes.Buffer)(nil)),
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[ssainspect.Analyzer].(*ssainspect.Inspector)

	var buf bytes.Buffer

	for inspect.Next() {
		if inspect.InstrIndex() == 0 {
			fmt.Fprintln(&buf, "Block", inspect.Block())
		}
		fmt.Fprintln(&buf, inspect.Instr())
	}

	return &buf, nil
}
