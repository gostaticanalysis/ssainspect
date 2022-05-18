package ssainspect

import (
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

type Inspector struct {
	cursors []*Cursor
	idx     int
}

func New(funcs []*ssa.Function) *Inspector {
	in := &Inspector{
		idx: -1,
	}
	analysisutil.InspectFuncs(funcs, func(i int, instr ssa.Instruction) bool {
		c := &Cursor{
			Func:       instr.Parent(),
			Block:      instr.Block(),
			InstrIndex: i,
			Instr:      instr,
		}

		in.cursors = append(in.cursors, c)
		return true
	})

	return in
}

func (in *Inspector) Next() bool {
	in.idx++
	if in.idx >= len(in.cursors) {
		return false
	}
	return true
}

func (in *Inspector) Cursor() *Cursor {
	if in.idx < 0 || in.idx >= len(in.cursors) {
		return nil
	}
	return in.cursors[in.idx]
}
