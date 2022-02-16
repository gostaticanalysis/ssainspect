package ssainspect

import "golang.org/x/tools/go/ssa"

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
		cycleBlocks: make(map[*ssa.BasicBlock]struct{}),
		done:        make(map[*ssa.BasicBlock]struct{}),
		funcs:       funcs,
	}
}

func (in *Inspector) Next() bool {
	if in.skiped {
		return false
	}

	// first call
	if in.funcIndex == 0 && in.instrIndex == 0 {
		in.block = in.firstBlock()
		if in.block == nil {
			return false
		}
	}

	in.instrIndex++
	if in.instrIndex < len(in.block.Instrs) {
		in.instrIndex = 0
		in.nextBlock()
		if in.block == nil {
			return false
		}
	}

	return true
}

func (in *Inspector) firstBlock() *ssa.BasicBlock {
	for len(in.Function().Blocks) == 0 {
		in.funcIndex++
	}
	if in.funcIndex >= len(in.funcs) {
		return nil
	}
	return in.Function().Blocks[0]
}

func (in *Inspector) nextBlock() {

	if in.block == nil {
		return
	}

	if len(in.block.Succs) > 0 {
		next := in.bool.Succs[0]
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

func (in *Inspector) addPre(b *ssa.BasicBlock) bool {
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
