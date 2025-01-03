package ssainspect

import (
	"iter"
	"slices"

	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/ssa"
)

// Cursor holds the current function, block and instruction in iteration by [All].
type Cursor struct {
	Func       *ssa.Function
	Block      *ssa.BasicBlock
	Instr      ssa.Instruction
	InstrIndex int
}

// FirstInstr returns whether cur.Instruction is a first instruction in the block.
func (cur *Cursor) FirstInstr() bool {
	return cur.InstrIndex == 0
}

// FirstBlock returns whether cur.Block is a first block in the function.
func (cur *Cursor) FirstBlock() bool {
	return cur.Block.Index == 0
}

// InCycle returns whether the block is in cycle (loop).
func (cur *Cursor) InCycle() bool {
	done := make(map[*ssa.BasicBlock]struct{})
	blocks := []*ssa.BasicBlock{cur.Block}
	for len(blocks) > 0 {
		b := blocks[0]
		blocks = blocks[1:]
		if _, ok := done[b]; ok {
			if b == cur.Block {
				return true
			}
			continue
		}
		done[b] = struct{}{}
		blocks = append(slices.Clone(b.Succs), blocks...)
	}
	return false
}

// Inspector provides an iterator which iterats all SSA functions, basic blocks and instructions.
type Inspector struct {
	cursors []*Cursor
}

// New creates [Inspector].
func New(funcs []*ssa.Function) *Inspector {
	return &Inspector{
		cursors: slices.Collect(All(funcs)),
	}
}

// All returns an iterator which inspects all SSA functions, basic blocks and instructions.
// The iteration result is cached in [Inspector].
func (i *Inspector) All() iter.Seq[*Cursor] {
	return slices.Values(i.cursors)
}

// All returns an iterator which inspects all SSA functions, basic blocks and instructions.
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
