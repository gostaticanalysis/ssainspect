package ssainspect_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/tenntenn/golden"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/gostaticanalysis/ssainspect"
)

var flagUpdate bool

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func TestInspectorWithAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)

	rs := analysistest.Run(t, testdata, testAnalyzer, "a")
	buf, ok := rs[0].Result.(*bytes.Buffer)
	if !ok {
		t.Fatal("unexpected error")
	}

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

func run(pass *analysis.Pass) (any, error) {
	inspect, ok := pass.ResultOf[ssainspect.Analyzer].(*ssainspect.Inspector)
	if !ok {
		return nil, errors.New("failed to type assert to *ssainspect.Inspector")
	}

	var buf bytes.Buffer

	for inspect.Next() {
		cur := inspect.Cursor()
		if cur.FirstInstr() {
			fmt.Fprintln(&buf, "Block", cur.Block, "InCycle=", cur.InCycle())
		}
		fmt.Fprintln(&buf, "\t", cur.Instr)
	}

	return &buf, nil
}
