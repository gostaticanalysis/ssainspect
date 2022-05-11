package ssainspect

import (
	"github.com/gostaticanalysis/analysisutil"
	"golang.org/x/tools/go/ssa"
)

type Cursor struct {
	Path       []*Cursor
	Func       *ssa.Function
	Block      *ssa.BasicBlock
	Instr      ssa.Instruction
	InstrIndex int
}

type Inspector struct {
	cursors []*Cursor
	idx     int
}

func New(funcs []*ssa.Function) *Inspector {
	in := &Inspector{
		idx: -1,
	}
	//var path []*Cursor
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

/*
type Inspector struct {
	cycleBlocks map[*ssa.BasicBlock]struct{}
	funcs       []*ssa.Function
	funcIndex   int
	block       *ssa.BasicBlock
	pre         []*ssa.BasicBlock
	done        map[*ssa.BasicBlock]struct{}
	instrIndex  int
}

func New(funcs []*ssa.Function) *Inspector {
	return &Inspector{
		funcIndex:   -1,
		instrIndex:  -1,
		cycleBlocks: make(map[*ssa.BasicBlock]struct{}),
		done:        make(map[*ssa.BasicBlock]struct{}),
		funcs:       funcs,
	}
}

func (in *Inspector) Next() bool {

	// first call
	if in.funcIndex < 0 && in.instrIndex < 0 {
		in.block = in.firstBlock()
		return in.block != nil
	}

	in.instrIndex++
	if in.instrIndex >= len(in.block.Instrs) {
		in.instrIndex = 0
		in.nextBlock()
		if in.block == nil {
			return false
		}
	}

	return true
}

func (in *Inspector) firstBlock() *ssa.BasicBlock {
	in.funcIndex = 0
	for len(in.Function().Blocks) == 0 {
		in.funcIndex++
	}
	if in.funcIndex >= len(in.funcs) {
		return nil
	}
	in.instrIndex = 0
	return in.Function().Blocks[0]
}

func (in *Inspector) nextBlock() {

	if in.block == nil {
		return
	}

	if len(in.block.Succs) > 0 {
		next := in.block.Succs[0]
		if !in.hasDone(next) {
			in.addPre(in.block)
			in.block = next
			in.markDone(in.block)
			return
		}
	}

	if len(in.pre) == 0 {
		in.block = nil
		return
	}

	pre := in.pre
	for _, p := range pre {
		for i := in.block.Index + 1; i < len(p.Succs); i++ {
			next := p.Succs[i]
			if !in.hasDone(next) {
				in.addPre(in.block)
				in.block = next
				in.markDone(in.block)
				return
			}
		}
		in.pre = in.pre[1:]
	}

	for {
		in.funcIndex++
		if in.funcIndex >= len(in.funcs) {
			in.block = nil
			return
		}

		if len(in.Function().Blocks) != 0 {
			next := in.Function().Blocks[0]
			if !in.hasDone(next) {
				in.addPre(in.block)
				in.block = next
				in.markDone(in.block)
				return
			}
		}
	}

	in.block = nil
}

func (in *Inspector) markDone(block *ssa.BasicBlock) {
	if in.done == nil {
		in.done = make(map[*ssa.BasicBlock]struct{})
	}
	in.done[block] = struct{}{}
}

func (in *Inspector) hasDone(b *ssa.BasicBlock) bool {
	_, ok := in.done[b]
	return ok
}

func (in *Inspector) addPre(b *ssa.BasicBlock) {
	in.pre = append(in.pre, nil)
	copy(in.pre[1:], in.pre)
	in.pre[0] = b
}

func (in *Inspector) Function() *ssa.Function {
	return in.funcs[in.funcIndex]
}

func (in *Inspector) Block() *ssa.BasicBlock {
	return in.block
}

func (in *Inspector) Instr() ssa.Instruction {
	if in.instrIndex < len(in.block.Instrs) {
		return in.block.Instrs[in.instrIndex]
	}
	return nil
}

func (in *Inspector) InstrIndex() int {
	return in.instrIndex
}

func (in *Inspector) InCycle() bool {
	if _, ok := in.cycleBlocks[in.Block()]; ok {
		return ok
	}

	done := make(map[*ssa.BasicBlock]bool)
	var walk func(b *ssa.BasicBlock)
	walk = func(b *ssa.BasicBlock) bool {
		if len(b.Succs) == 0 {
			return false
		}
		if done[b] {
			return true
		}
		done[b] = true
	}
}
*/
