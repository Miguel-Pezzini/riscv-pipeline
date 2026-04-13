package hazard

import "riscv-instruction-encoder/pkg/isa"

func HasControlHazard(currentInstruction isa.PipelineInstruction, executing []*isa.PipelineInstruction, forwarding bool) bool {
	for _, prev := range executing {
		if hasUnresolvedBranchHazard(currentInstruction, *prev, forwarding) {
			return true
		}
	}

	return false
}

func hasUnresolvedBranchHazard(currentInstruction isa.PipelineInstruction, previousInstruction isa.PipelineInstruction, forwarding bool) bool {
	prevMeta := previousInstruction.Instruction.GetMeta()
	if !(prevMeta.IsBranch || prevMeta.IsJump) || previousInstruction.HasCompleted {
		return false
	}

	resolveStage := isa.WB
	if forwarding {
		resolveStage = isa.EX
	}

	return previousInstruction.CurrentStage < int(resolveStage)
}
