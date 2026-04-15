package runner

import (
	"riscv-instruction-encoder/pkg/executor"
	"riscv-instruction-encoder/pkg/isa"
)

func HistoryToPipeline(history []executor.StepResult) []*isa.PipelineInstruction {
	pipelineInstructions := make([]*isa.PipelineInstruction, len(history))
	for i, step := range history {
		pc := int(step.PC)
		pipelineInstructions[i] = &isa.PipelineInstruction{
			Instruction:  step.Instr,
			CurrentStage: 0,
			HasCompleted: false,
			HasStarted:   false,
			Id:           i + 1,
			PC:           pc,
			OriginalPC:   pc,
		}
	}
	return pipelineInstructions
}

func NewPipelineFromExecutionHistory(history []executor.StepResult, forwarding bool, data_hazard bool, control_hazard bool, file_path string) *Pipeline {
	stages := len(isa.Stages)

	return &Pipeline{
		CurrentCycle:   0,
		Instructions:   HistoryToPipeline(history),
		NumStages:      stages,
		forwarding:     forwarding,
		data_hazard:    data_hazard,
		control_hazard: control_hazard,
		file_path:      file_path,
	}
}

func RunFromExecutionHistory(history []executor.StepResult, forwarding bool, data_hazard bool, control_hazard bool, file_path string) {
	p := NewPipelineFromExecutionHistory(history, forwarding, data_hazard, control_hazard, file_path)
	p.Run()
	p.printResult()
	p.writeFile()
}
