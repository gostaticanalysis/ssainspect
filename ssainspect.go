package ssainspect

import (
	"iter"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/ssa"
)

type Cursor struct {
	Func       *ssa.Function
	Block      *ssa.BasicBlock
	Instr      ssa.Instruction
	InstrIndex int
}

func (c *Cursor) FirstInstr() bool {
	return c.InstrIndex == 0
}

func (c *Cursor) InCycle() bool {
	done := make(map[*ssa.BasicBlock]struct{})
	blocks := []*ssa.BasicBlock{c.Block}
	for len(blocks) > 0 {
		b := blocks[0]
		blocks = blocks[1:]
		if _, ok := done[b]; ok {
			if b == c.Block {
				return true
			}
			continue
		}
		done[b] = struct{}{}
		blocks = append(b.Succs, blocks...)
	}
	return false
}

func All(funcs []*ssa.Function) iter.Seq[*Cursor] {
	return func(yield func(*Cursor) bool) {
		analysisutil.InspectFuncs(funcs, func(i int, instr ssa.Instruction) bool {
			c := &Cursor{
				Func:       instr.Parent(),
				Block:      instr.Block(),
				InstrIndex: i,
				Instr:      instr,
			}
			return yield(c)
		})
	}
}
