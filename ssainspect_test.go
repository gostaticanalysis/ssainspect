package ssainspect_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"iter"
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
	Doc:  "test analyzer for ssainspect.All",
	Run:  run,
	Requires: []*analysis.Analyzer{
		ssainspect.Analyzer,
	},
	ResultType: reflect.TypeFor[*bytes.Buffer](),
}

func run(pass *analysis.Pass) (any, error) {
	seq, ok := pass.ResultOf[ssainspect.Analyzer].(iter.Seq[*ssainspect.Cursor])
	if !ok {
		return nil, errors.New("failed to get result of ssainspect.Analyzer")
	}

	var buf bytes.Buffer

	for cur := range seq {
		if cur.FirstInstr() {
			if cur.FirstBlock() {
				fmt.Fprintln(&buf, "Func", cur.Func)
			}
			fmt.Fprintln(&buf, "Block", cur.Block, "InCycle=", cur.InCycle())
		}
		fmt.Fprintln(&buf, "\t", cur.Instr)
	}

	return &buf, nil
}
